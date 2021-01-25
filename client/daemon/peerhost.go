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

package daemon

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"time"

	"golang.org/x/sync/errgroup"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/dragonflyoss/Dragonfly2/client/daemon/gc"
	"github.com/dragonflyoss/Dragonfly2/client/daemon/peer"
	"github.com/dragonflyoss/Dragonfly2/client/daemon/service"
	"github.com/dragonflyoss/Dragonfly2/client/daemon/storage"
	"github.com/dragonflyoss/Dragonfly2/client/daemon/upload"
	"github.com/dragonflyoss/Dragonfly2/pkg/basic/dfnet"
	logger "github.com/dragonflyoss/Dragonfly2/pkg/dflog"
	"github.com/dragonflyoss/Dragonfly2/pkg/rpc"
	"github.com/dragonflyoss/Dragonfly2/pkg/rpc/scheduler"
	schedulerclient "github.com/dragonflyoss/Dragonfly2/pkg/rpc/scheduler/client"
)

type PeerHost interface {
	Serve() error
	Stop()
}

type peerHost struct {
	started chan bool

	schedPeerHost *scheduler.PeerHost

	Option PeerHostOption

	ServiceManager service.Manager
	UploadManager  upload.Manager
	StorageManager storage.Manager
	GCManager      gc.Manager

	PeerTaskManager peer.PeerTaskManager
	PieceManager    peer.PieceManager
}

type PeerHostOption struct {
	Schedulers []dfnet.NetAddr

	// AliveTime indicates alive duration for which daemon keeps no accessing by any uploading and download requests,
	// after this period daemon will automatically exit
	// when AliveTime == 0, will run infinitely
	// TODO keepalive detect
	AliveTime  time.Duration
	GCInterval time.Duration

	Server  ServerOption
	Upload  UploadOption
	Storage StorageOption
}

type ServerOption struct {
	RateLimit    rate.Limit
	DownloadGRPC ListenOption
	PeerGRPC     ListenOption
	Proxy        *ListenOption
}

type UploadOption struct {
	ListenOption
	RateLimit rate.Limit
}

type ListenOption struct {
	Security   SecurityOption
	TCPListen  *TCPListenOption
	UnixListen *UnixListenOption
}

type TCPListenOption struct {
	Listen    string
	PortRange TCPListenPortRange
}

type TCPListenPortRange struct {
	Start int
	End   int
}

type UnixListenOption struct {
	Socket string
}

type SecurityOption struct {
	Insecure  bool
	CACert    string
	Cert      string
	Key       string
	TLSConfig *tls.Config
}

type StorageOption struct {
	storage.Option
	Driver storage.Driver
}

func NewPeerHost(host *scheduler.PeerHost, opt PeerHostOption) (PeerHost, error) {
	sched, err := schedulerclient.CreateClient(opt.Schedulers)
	if err != nil {
		return nil, err
	}

	storageManager, err := storage.NewStorageManager(opt.Storage.Driver, &opt.Storage.Option, func(request storage.CommonTaskRequest) {
		sched.LeaveTask(context.Background(), &scheduler.PeerTarget{
			TaskId: request.TaskID,
			PeerId: request.PeerID,
		})
	})
	if err != nil {
		return nil, err
	}

	pieceManager, err := peer.NewPieceManager(storageManager, peer.WithLimiter(rate.NewLimiter(opt.Server.RateLimit, int(opt.Server.RateLimit))))
	if err != nil {
		return nil, err
	}
	peerTaskManager, err := peer.NewPeerTaskManager(pieceManager, storageManager, sched)
	if err != nil {
		return nil, err
	}

	// TODO(jim): more server options
	var downloadServerOption []grpc.ServerOption
	if !opt.Server.DownloadGRPC.Security.Insecure {
		tlsCredentials, err := loadGPRCTLSCredentials(opt.Server.DownloadGRPC.Security)
		if err != nil {
			return nil, err
		}
		downloadServerOption = append(downloadServerOption, grpc.Creds(tlsCredentials))
	}
	var peerServerOption []grpc.ServerOption
	if !opt.Server.PeerGRPC.Security.Insecure {
		tlsCredentials, err := loadGPRCTLSCredentials(opt.Server.PeerGRPC.Security)
		if err != nil {
			return nil, err
		}
		peerServerOption = append(peerServerOption, grpc.Creds(tlsCredentials))
	}
	serviceManager, err := service.NewManager(host, peerTaskManager, storageManager, downloadServerOption, peerServerOption)
	if err != nil {
		return nil, err
	}

	uploadManager, err := upload.NewUploadManager(storageManager,
		upload.WithLimiter(rate.NewLimiter(opt.Upload.RateLimit, int(opt.Upload.RateLimit))))
	if err != nil {
		return nil, err
	}

	return &peerHost{
		started:       make(chan bool),
		schedPeerHost: host,
		Option:        opt,

		ServiceManager:  serviceManager,
		PeerTaskManager: peerTaskManager,
		PieceManager:    pieceManager,
		UploadManager:   uploadManager,
		StorageManager:  storageManager,
		GCManager:       gc.NewManager(opt.GCInterval),
	}, nil
}

func loadGPRCTLSCredentials(opt SecurityOption) (credentials.TransportCredentials, error) {
	// Load certificate of the CA who signed client's certificate
	pemClientCA, err := ioutil.ReadFile(opt.CACert)
	if err != nil {
		return nil, err
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pemClientCA) {
		return nil, fmt.Errorf("failed to add client CA's certificate")
	}

	// Load server's certificate and private key
	serverCert, err := tls.LoadX509KeyPair(opt.Cert, opt.Key)
	if err != nil {
		return nil, err
	}

	// Create the credentials and return it
	if opt.TLSConfig == nil {
		opt.TLSConfig = &tls.Config{
			Certificates: []tls.Certificate{serverCert},
			ClientAuth:   tls.RequireAndVerifyClientCert,
			ClientCAs:    certPool,
		}
	} else {
		opt.TLSConfig.Certificates = []tls.Certificate{serverCert}
		opt.TLSConfig.ClientAuth = tls.RequireAndVerifyClientCert
		opt.TLSConfig.ClientCAs = certPool
	}

	return credentials.NewTLS(opt.TLSConfig), nil
}

func (ph *peerHost) prepareTCPListener(opt ListenOption, withTLS bool) (net.Listener, int, error) {
	var (
		ln   net.Listener
		port int
		err  error
	)
	if opt.TCPListen != nil {
		ln, port, err = rpc.ListenWithPortRange(opt.TCPListen.Listen, opt.TCPListen.PortRange.Start, opt.TCPListen.PortRange.End)
	}
	if err != nil {
		return nil, -1, err
	}
	// when use grpc, tls config is in server option
	if !withTLS || opt.Security.Insecure {
		return ln, port, err
	}

	// Create the TLS Config with the CA pool and enable Client certificate validation
	if opt.Security.TLSConfig == nil {
		opt.Security.TLSConfig = &tls.Config{}
	}
	tlsConfig := opt.Security.TLSConfig
	if opt.Security.CACert != "" {
		caCert, err := ioutil.ReadFile(opt.Security.CACert)
		if err != nil {
			return nil, -1, err
		}
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)
		tlsConfig.ClientCAs = caCertPool
		tlsConfig.ClientAuth = tls.RequireAndVerifyClientCert
	}
	tlsConfig.Certificates = make([]tls.Certificate, 1)
	tlsConfig.Certificates[0], err = tls.LoadX509KeyPair(opt.Security.Cert, opt.Security.Key)
	if err != nil {
		return nil, -1, err
	}

	return tls.NewListener(ln, tlsConfig), port, nil
}

func (ph *peerHost) Serve() error {
	ph.GCManager.Start()

	// prepare download service listen
	_ = os.Remove(ph.Option.Server.DownloadGRPC.UnixListen.Socket)
	downloadListener, err := rpc.Listen(dfnet.NetAddr{
		Type: dfnet.UNIX,
		Addr: ph.Option.Server.DownloadGRPC.UnixListen.Socket,
	})
	if err != nil {
		logger.Errorf("failed to listen for download grpc service: %v", err)
		return err
	}

	// prepare peer service listen
	peerListener, peerPort, err := ph.prepareTCPListener(ph.Option.Server.PeerGRPC, false)
	if err != nil {
		logger.Errorf("failed to listen for peer grpc service: %v", err)
		return err
	}
	ph.schedPeerHost.RpcPort = int32(peerPort)

	// prepare upload service listen
	uploadListener, uploadPort, err := ph.prepareTCPListener(ph.Option.Upload.ListenOption, true)
	if err != nil {
		logger.Errorf("failed to listen for upload service: %v", err)
		return err
	}
	ph.schedPeerHost.DownPort = int32(uploadPort)

	g := errgroup.Group{}
	// serve download grpc service
	g.Go(func() error {
		logger.Infof("serve download grpc at unix://%s", ph.Option.Server.DownloadGRPC.UnixListen.Socket)
		if err = ph.ServiceManager.ServeDownload(downloadListener); err != nil {
			logger.Errorf("failed to serve for download grpc service: %v", err)
			return err
		}
		return nil
	})

	// serve peer grpc service
	g.Go(func() error {
		logger.Infof("serve peer grpc at %s://%s", peerListener.Addr().Network(), peerListener.Addr().String())
		if err = ph.ServiceManager.ServePeer(peerListener); err != nil {
			logger.Errorf("failed to serve for peer grpc service: %v", err)
			return err
		}
		return nil
	})

	// serve proxy service
	if ph.Option.Server.Proxy != nil {
		g.Go(func() error {
			listener, port, err := ph.prepareTCPListener(*ph.Option.Server.Proxy, true)
			if err != nil {
				logger.Errorf("failed to listen for proxy service: %v", err)
				return err
			}

			logger.Infof("serve proxy at tcp://%s:%d", ph.Option.Server.Proxy.TCPListen.Listen, port)
			if err = ph.ServiceManager.ServeProxy(listener); err != nil {
				logger.Errorf("failed to serve for proxy service: %v", err)
				return err
			}
			return nil
		})
	}

	// serve upload service
	g.Go(func() error {
		logger.Infof("serve upload service at %s://%s", uploadListener.Addr().Network(), uploadListener.Addr().String())
		ph.schedPeerHost.DownPort = int32(uploadPort)
		if err = ph.UploadManager.Serve(uploadListener); err != nil {
			logger.Errorf("failed to serve for upload service: %v", err)
			return err
		}
		return nil
	})

	return g.Wait()
}

func (ph *peerHost) Stop() {
	ph.ServiceManager.Stop()
	ph.UploadManager.Stop()
}
