/*
 *     Copyright 2020 The Dragonfly Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package manager

import (
	"context"
	"crypto/rand"
	"embed"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io/fs"
	"net"
	"net/http"
	"time"

	"github.com/gin-contrib/static"
	"google.golang.org/grpc"
	"gorm.io/gorm"

	logger "d7y.io/dragonfly/v2/internal/dflog"
	"d7y.io/dragonfly/v2/internal/ratelimiter"
	"d7y.io/dragonfly/v2/manager/cache"
	"d7y.io/dragonfly/v2/manager/config"
	"d7y.io/dragonfly/v2/manager/database"
	managergc "d7y.io/dragonfly/v2/manager/gc"
	"d7y.io/dragonfly/v2/manager/job"
	"d7y.io/dragonfly/v2/manager/metrics"
	"d7y.io/dragonfly/v2/manager/models"
	"d7y.io/dragonfly/v2/manager/permission/rbac"
	"d7y.io/dragonfly/v2/manager/router"
	"d7y.io/dragonfly/v2/manager/rpcserver"
	"d7y.io/dragonfly/v2/manager/searcher"
	"d7y.io/dragonfly/v2/manager/service"
	"d7y.io/dragonfly/v2/pkg/dfpath"
	pkggc "d7y.io/dragonfly/v2/pkg/gc"
	"d7y.io/dragonfly/v2/pkg/net/ip"
	"d7y.io/dragonfly/v2/pkg/objectstorage"
	"d7y.io/dragonfly/v2/pkg/redis"
	"d7y.io/dragonfly/v2/pkg/rpc"
)

const (
	// gracefulStopTimeout specifies a time limit for
	// grpc server to complete a graceful shutdown.
	gracefulStopTimeout = 10 * time.Minute

	// assetsTargetPath is target path of embed assets.
	assetsTargetPath = "dist"
)

//go:embed dist/*
var assets embed.FS

type embedFileSystem struct {
	http.FileSystem
}

func (e embedFileSystem) Exists(prefix string, path string) bool {
	_, err := e.Open(path)
	if err != nil {
		return false
	}
	return true
}

func EmbedFolder(fsEmbed embed.FS, targetPath string) static.ServeFileSystem {
	fsys, err := fs.Sub(fsEmbed, targetPath)
	if err != nil {
		panic(err)
	}

	return embedFileSystem{
		FileSystem: http.FS(fsys),
	}
}

// Server is the manager server.
type Server struct {
	// Server configuration.
	config *config.Config

	// Job server.
	job *job.Job

	// Job rate limiter.
	jobRateLimiter ratelimiter.JobRateLimiter

	// GC server.
	gc pkggc.GC

	// GRPC server.
	grpcServer *grpc.Server

	// REST server.
	restServer *http.Server

	// Metrics server.
	metricsServer *http.Server

	// Redis proxy.
	redisProxy redis.Proxy
}

// New creates a new manager server.
func New(cfg *config.Config, d dfpath.Dfpath) (*Server, error) {
	s := &Server{config: cfg}

	// Initialize database.
	db, err := database.New(cfg)
	if err != nil {
		return nil, err
	}

	// Initialize encryption key
	if cfg.Encryption.Enable {
		if err := initializeEncryptionKey(cfg, db.DB); err != nil {
			return nil, err
		}
		logger.Infof("encryption enabled")
	} else {
		logger.Infof("encryption disabled")
	}

	// Initialize enforcer.
	enforcer, err := rbac.NewEnforcer(db.DB)
	if err != nil {
		return nil, err
	}

	// Initialize cache.
	cache, err := cache.New(cfg)
	if err != nil {
		return nil, err
	}

	// Initialize searcher.
	searcher := searcher.New(d.PluginDir())

	// Initialize job.
	job, err := job.New(cfg, db.DB)
	if err != nil {
		return nil, err
	}
	s.job = job

	// Initialize object storage.
	var objectStorage objectstorage.ObjectStorage
	if cfg.ObjectStorage.Enable {
		objectStorage, err = objectstorage.New(
			cfg.ObjectStorage.Name,
			cfg.ObjectStorage.Region,
			cfg.ObjectStorage.Endpoint,
			cfg.ObjectStorage.AccessKey,
			cfg.ObjectStorage.SecretKey,
			objectstorage.WithS3ForcePathStyle(cfg.ObjectStorage.S3ForcePathStyle),
		)
		if err != nil {
			return nil, err
		}
	}

	// Initialize job rate limiter.
	s.jobRateLimiter, err = ratelimiter.NewJobRateLimiter(db)
	if err != nil {
		return nil, err
	}

	// Initialize garbage collector.
	gc := pkggc.New()
	if err := registerGCTasks(gc, db.DB); err != nil {
		return nil, err
	}
	s.gc = gc

	// Initialize REST server.
	restService := service.New(cfg, db, cache, job, gc, enforcer, objectStorage)
	router, err := router.Init(cfg, d.LogDir(), restService, db, enforcer, s.jobRateLimiter, EmbedFolder(assets, assetsTargetPath))
	if err != nil {
		return nil, err
	}
	s.restServer = &http.Server{
		Addr:    cfg.Server.REST.Addr,
		Handler: router,
	}

	// Initialize roles and check roles.
	err = rbac.InitRBAC(enforcer, router, db.DB)
	if err != nil {
		return nil, err
	}

	// Initialize signing certificate and tls credentials of grpc server.
	var options []grpc.ServerOption
	if cfg.Server.GRPC.TLS != nil {
		// Initialize grpc server with tls.
		transportCredentials, err := rpc.NewServerCredentials(cfg.Server.GRPC.TLS.CACert, cfg.Server.GRPC.TLS.Cert, cfg.Server.GRPC.TLS.Key)
		if err != nil {
			logger.Errorf("failed to create server credentials: %v", err)
			return nil, err
		}

		options = append(options, grpc.Creds(transportCredentials))
	} else {
		// Initialize grpc server without tls.
		options = append(options, grpc.Creds(rpc.NewInsecureCredentials()))
	}

	// Initialize GRPC server.
	_, grpcServer, err := rpcserver.New(cfg, db, cache, searcher, objectStorage, options...)
	if err != nil {
		return nil, err
	}

	s.grpcServer = grpcServer

	// Initialize prometheus.
	if cfg.Metrics.Enable {
		s.metricsServer = metrics.New(&cfg.Metrics, grpcServer)
	}

	if cfg.Database.Redis.Proxy.Enable {
		s.redisProxy = redis.NewProxy(cfg.Database.Redis.Proxy.Addr, cfg.Database.Redis.Addrs[0])
	}

	return s, nil
}

func registerGCTasks(gc pkggc.GC, db *gorm.DB) error {
	// Register job gc task.
	if err := gc.Add(managergc.NewJobGCTask(db)); err != nil {
		return err
	}

	// Register audit gc task.
	if err := gc.Add(managergc.NewAuditGCTask(db)); err != nil {
		return err
	}

	// Register scheduler gc task.
	if err := gc.Add(managergc.NewSchedulerGCTask(db)); err != nil {
		return err
	}

	// Register seed peer gc task.
	if err := gc.Add(managergc.NewSeedPeerGCTask(db)); err != nil {
		return err
	}

	return nil
}

// Initialize encryption key
func initializeEncryptionKey(cfg *config.Config, db *gorm.DB) error {
	// 1. try get key from db
	var existingKey models.EncryptionKey
	if err := db.First(&existingKey).Error; err == nil {
		logger.Infof("encryption key loaded from database, key(hex): %s, key(base64): %s",
			hex.EncodeToString(existingKey.Key),
			base64.StdEncoding.EncodeToString(existingKey.Key),
		)
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("failed to check encryption key: %v", err)
	}

	// 2. if there is no key in db, generate a new one
	keyBytes := make([]byte, 32)
	if _, err := rand.Read(keyBytes); err != nil {
		return fmt.Errorf("failed to generate random encryption key: %v", err)
	}
	if err := db.Create(&models.EncryptionKey{Key: keyBytes}).Error; err != nil {
		return fmt.Errorf("failed to save random encryption key to database: %v", err)
	}
	logger.Infof(
		"generated random encryption key and saved to database, key(hex): %s, key(base64): %s",
		hex.EncodeToString(keyBytes),
		base64.StdEncoding.EncodeToString(keyBytes),
	)
	return nil
}

// Serve starts the manager server.
func (s *Server) Serve() error {
	// Started REST server.
	go func() {
		logger.Infof("started rest server at %s", s.restServer.Addr)
		if s.config.Server.REST.TLS != nil {
			if err := s.restServer.ListenAndServeTLS(s.config.Server.REST.TLS.Cert, s.config.Server.REST.TLS.Key); err != nil {
				if err == http.ErrServerClosed {
					return
				}

				logger.Fatalf("rest server closed unexpect: %v", err)
			}
		} else {
			if err := s.restServer.ListenAndServe(); err != nil {
				if err == http.ErrServerClosed {
					return
				}

				logger.Fatalf("rest server closed unexpect: %v", err)
			}
		}
	}()

	// Started metrics server.
	if s.metricsServer != nil {
		go func() {
			logger.Infof("started metrics server at %s", s.metricsServer.Addr)
			if err := s.metricsServer.ListenAndServe(); err != nil {
				if err == http.ErrServerClosed {
					return
				}

				logger.Fatalf("metrics server closed unexpect: %v", err)
			}
		}()
	}

	// Started redis proxy server.
	if s.redisProxy != nil {
		go func() {
			logger.Infof("started redis proxy server at %s", s.config.Database.Redis.Proxy.Addr)
			if err := s.redisProxy.Serve(); err != nil {
				logger.Fatalf("redis proxy server closed unexpect: %v", err)
			}
		}()
	}

	// Started job server.
	go func() {
		logger.Info("started job server")
		s.job.Serve()
	}()

	// Started job rate limiter server.
	go func() {
		logger.Infof("started job rate limiter server")
		s.jobRateLimiter.Serve()
	}()

	// Started gc server.
	s.gc.Start(context.Background())
	logger.Info("started gc server")

	// Generate GRPC listener.
	ip, ok := ip.FormatIP(s.config.Server.GRPC.ListenIP.String())
	if !ok {
		return fmt.Errorf("format ip failed: %s", ip)
	}

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", ip, s.config.Server.GRPC.Port.Start))
	if err != nil {
		logger.Fatalf("net listener failed to start: %v", err)
	}
	defer listener.Close()

	// Started GRPC server.
	logger.Infof("started grpc server at %s://%s", listener.Addr().Network(), listener.Addr().String())
	if err := s.grpcServer.Serve(listener); err != nil {
		logger.Errorf("stoped grpc server: %+v", err)
		return err
	}

	return nil
}

// Stop stops the manager server.
func (s *Server) Stop() {
	// Stop REST server.
	if err := s.restServer.Shutdown(context.Background()); err != nil {
		logger.Errorf("rest server failed to stop: %+v", err)
	} else {
		logger.Info("rest server closed under request")
	}

	// Stop metrics server.
	if s.metricsServer != nil {
		if err := s.metricsServer.Shutdown(context.Background()); err != nil {
			logger.Errorf("metrics server failed to stop: %+v", err)
		} else {
			logger.Info("metrics server closed under request")
		}
	}

	// Stop redis proxy server.
	if s.redisProxy != nil {
		s.redisProxy.Stop()
		logger.Info("redis proxy server closed under request")
	}

	// Stop job server.
	s.job.Stop()

	// Stop job rate limiter.
	s.jobRateLimiter.Stop()

	// Stop gc server.
	s.gc.Stop()

	// Stop GRPC server.
	stopped := make(chan struct{})
	go func() {
		s.grpcServer.GracefulStop()
		logger.Info("grpc server closed under request")
		close(stopped)
	}()

	t := time.NewTimer(gracefulStopTimeout)
	select {
	case <-t.C:
		s.grpcServer.Stop()
	case <-stopped:
		t.Stop()
	}
}
