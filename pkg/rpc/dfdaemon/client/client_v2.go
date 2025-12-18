/*
 *     Copyright 2023 The Dragonfly Authors
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

//go:generate mockgen -destination mocks/client_v2_mock.go -source client_v2.go -package mocks

package client

import (
	"context"
	"math"
	"sync"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"golang.org/x/sync/singleflight"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"

	commonv2 "d7y.io/api/v2/pkg/apis/common/v2"
	dfdaemonv2 "d7y.io/api/v2/pkg/apis/dfdaemon/v2"

	logger "d7y.io/dragonfly/v2/internal/dflog"
	pkgbalancer "d7y.io/dragonfly/v2/pkg/balancer"
)

// Pool is the interface for pooling v2 version of the grpc client.
type Pool interface {
	// Serve starts the manager.
	Serve()

	// Stop stops the manager.
	Stop()

	// Get returns the client by address.
	Get(target string, opts ...grpc.DialOption) (V2, error)
}

// pool is the pool for managing v2 version of the dfdaemon client.
type pool struct {
	// pool is a map of client connections for reusing.
	*sync.Map

	// sf is the singleflight instance for concurrent requests.
	sf *singleflight.Group

	// done is a channel to signal the manager is done.
	done chan struct{}
}

// GetV2Pool creates a new pool instance.
func GetV2Pool() Pool {
	return &pool{
		Map:  &sync.Map{},
		sf:   &singleflight.Group{},
		done: make(chan struct{}),
	}
}

// Serve starts the pool.
func (p *pool) Serve() {
	ticker := time.NewTicker(gcInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			p.runGC()
		case <-p.done:
			return
		}
	}
}

// Stop stops the pool.
func (p *pool) Stop() {
	close(p.done)
}

// Get returns a v2 version of the dfdaemon client by address.
func (p *pool) Get(target string, opts ...grpc.DialOption) (V2, error) {
	if client, ok := p.Load(target); ok {
		return client.(V2), nil
	}

	client, err, _ := p.sf.Do(target, func() (any, error) { return GetV2ByAddr(target, opts...) })
	if err != nil {
		return nil, err
	}

	p.Store(target, client)
	return client.(V2), nil
}

// runGC cleans up unhealthy connections.
func (p *pool) runGC() {
	p.Range(func(k, v any) bool {
		// Cleanup the not connecting and ready connections and remove them from the pool.
		if state := v.(*v2).GetState(); state != connectivity.Connecting && state != connectivity.Ready {
			v.(*v2).Close()
			p.Delete(k)
			return true
		}

		return true
	})
}

// V2 is the interface for v2 version of the grpc client.
type V2 interface {
	// SyncPieces syncs pieces from the other peers.
	SyncPieces(context.Context, *dfdaemonv2.SyncPiecesRequest, ...grpc.CallOption) (dfdaemonv2.DfdaemonUpload_SyncPiecesClient, error)

	// DownloadPiece downloads piece from the other peer.
	DownloadPiece(context.Context, *dfdaemonv2.DownloadPieceRequest, ...grpc.CallOption) (*dfdaemonv2.DownloadPieceResponse, error)

	// DownloadTask downloads task from p2p network.
	DownloadTask(context.Context, string, *dfdaemonv2.DownloadTaskRequest, ...grpc.CallOption) (dfdaemonv2.DfdaemonUpload_DownloadTaskClient, error)

	// StatTask stats task information.
	StatTask(context.Context, *dfdaemonv2.StatTaskRequest, ...grpc.CallOption) (*commonv2.Task, error)

	// StatLocalTask stats local task information.
	StatLocalTask(context.Context, *dfdaemonv2.StatLocalTaskRequest, ...grpc.CallOption) (*dfdaemonv2.StatLocalTaskResponse, error)

	// DeleteTask deletes task from p2p network.
	DeleteTask(context.Context, *dfdaemonv2.DeleteTaskRequest, ...grpc.CallOption) error

	// ListTaskEntries lists task entries.
	ListTaskEntries(context.Context, *dfdaemonv2.ListTaskEntriesRequest, ...grpc.CallOption) (*dfdaemonv2.ListTaskEntriesResponse, error)

	// DownloadPersistentTask downloads persistent task from p2p network.
	DownloadPersistentTask(context.Context, *dfdaemonv2.DownloadPersistentTaskRequest, ...grpc.CallOption) (dfdaemonv2.DfdaemonUpload_DownloadPersistentTaskClient, error)

	// UpdatePersistentTask updates persistent task information.
	UpdatePersistentTask(context.Context, *dfdaemonv2.UpdatePersistentTaskRequest, ...grpc.CallOption) error

	// StatPersistentTask stats persistent task information.
	StatPersistentTask(context.Context, *dfdaemonv2.StatPersistentTaskRequest, ...grpc.CallOption) (*commonv2.PersistentTask, error)

	// DeletePersistentTask deletes persistent task from p2p network.
	DeletePersistentTask(context.Context, *dfdaemonv2.DeletePersistentTaskRequest, ...grpc.CallOption) error

	// DownloadPersistentCacheTask downloads persistent cache task from p2p network.
	DownloadPersistentCacheTask(context.Context, *dfdaemonv2.DownloadPersistentCacheTaskRequest, ...grpc.CallOption) (dfdaemonv2.DfdaemonUpload_DownloadPersistentCacheTaskClient, error)

	// UpdatePersistentCacheTask updates persistent cache task information.
	UpdatePersistentCacheTask(context.Context, *dfdaemonv2.UpdatePersistentCacheTaskRequest, ...grpc.CallOption) error

	// StatPersistentCacheTask stats persistent cache task information.
	StatPersistentCacheTask(context.Context, *dfdaemonv2.StatPersistentCacheTaskRequest, ...grpc.CallOption) (*commonv2.PersistentCacheTask, error)

	// DeletePersistentCacheTask deletes persistent cache task from p2p network.
	DeletePersistentCacheTask(context.Context, *dfdaemonv2.DeletePersistentCacheTaskRequest, ...grpc.CallOption) error

	// Close tears down the ClientConn and all underlying connections.
	Close() error
}

// GetV2ByAddr returns v2 version of the dfdaemon client by address.
func GetV2ByAddr(target string, opts ...grpc.DialOption) (V2, error) {
	conn, err := grpc.NewClient(
		target,
		append([]grpc.DialOption{
			grpc.WithIdleTimeout(idleTimeout),
			grpc.WithDefaultCallOptions(
				grpc.MaxCallRecvMsgSize(math.MaxInt32),
				grpc.MaxCallSendMsgSize(math.MaxInt32),
			),
			grpc.WithUnaryInterceptor(grpc_middleware.ChainUnaryClient(
				grpc_prometheus.UnaryClientInterceptor,
				grpc_zap.UnaryClientInterceptor(logger.GrpcLogger.Desugar()),
				grpc_retry.UnaryClientInterceptor(
					grpc_retry.WithMax(maxRetries),
					grpc_retry.WithBackoff(grpc_retry.BackoffLinear(backoffWaitBetween)),
				),
			)),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithStreamInterceptor(grpc_middleware.ChainStreamClient(
				grpc_prometheus.StreamClientInterceptor,
				grpc_zap.StreamClientInterceptor(logger.GrpcLogger.Desugar()),
			)),
		}, opts...)...,
	)
	if err != nil {
		return nil, err
	}

	return &v2{
		DfdaemonUploadClient: dfdaemonv2.NewDfdaemonUploadClient(conn),
		ClientConn:           conn,
	}, nil
}

// v2 provides v2 version of the dfdaemon grpc function.
type v2 struct {
	dfdaemonv2.DfdaemonUploadClient
	*grpc.ClientConn
	*pkgbalancer.ConsistentHashingPickerBuilder
}

// SyncPieces syncs pieces from the other peers.
func (v *v2) SyncPieces(ctx context.Context, req *dfdaemonv2.SyncPiecesRequest, opts ...grpc.CallOption) (dfdaemonv2.DfdaemonUpload_SyncPiecesClient, error) {
	ctx, cancel := context.WithTimeout(ctx, contextTimeout)
	defer cancel()

	return v.DfdaemonUploadClient.SyncPieces(
		context.WithValue(ctx, pkgbalancer.ContextKey, req.TaskId),
		req,
		opts...,
	)
}

// DownloadPiece downloads piece from the other peer.
func (v *v2) DownloadPiece(ctx context.Context, req *dfdaemonv2.DownloadPieceRequest, opts ...grpc.CallOption) (*dfdaemonv2.DownloadPieceResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, contextTimeout)
	defer cancel()

	return v.DfdaemonUploadClient.DownloadPiece(
		context.WithValue(ctx, pkgbalancer.ContextKey, req.TaskId),
		req,
		opts...,
	)
}

// DownloadTask downloads task from p2p network.
func (v *v2) DownloadTask(ctx context.Context, taskID string, req *dfdaemonv2.DownloadTaskRequest, opts ...grpc.CallOption) (dfdaemonv2.DfdaemonUpload_DownloadTaskClient, error) {
	return v.DfdaemonUploadClient.DownloadTask(
		context.WithValue(ctx, pkgbalancer.ContextKey, taskID),
		req,
		opts...,
	)
}

// StatTask stats task information.
func (v *v2) StatTask(ctx context.Context, req *dfdaemonv2.StatTaskRequest, opts ...grpc.CallOption) (*commonv2.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, contextTimeout)
	defer cancel()

	return v.DfdaemonUploadClient.StatTask(ctx, req, opts...)
}

// StatLocalTask stats local task information.
func (v *v2) StatLocalTask(ctx context.Context, req *dfdaemonv2.StatLocalTaskRequest, opts ...grpc.CallOption) (*dfdaemonv2.StatLocalTaskResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, contextTimeout)
	defer cancel()

	return v.DfdaemonUploadClient.StatLocalTask(ctx, req, opts...)
}

// DeleteTask deletes task from p2p network.
func (v *v2) DeleteTask(ctx context.Context, req *dfdaemonv2.DeleteTaskRequest, opts ...grpc.CallOption) error {
	ctx, cancel := context.WithTimeout(ctx, contextTimeout)
	defer cancel()

	_, err := v.DfdaemonUploadClient.DeleteTask(ctx, req, opts...)
	return err
}

// DownloadPersistentTask downloads persistent task from p2p network.
func (v *v2) DownloadPersistentTask(ctx context.Context, req *dfdaemonv2.DownloadPersistentTaskRequest, opts ...grpc.CallOption) (dfdaemonv2.DfdaemonUpload_DownloadPersistentTaskClient, error) {
	return v.DfdaemonUploadClient.DownloadPersistentTask(ctx, req, opts...)
}

// UpdatePersistentTask updates persistent task information.
func (v *v2) UpdatePersistentTask(ctx context.Context, req *dfdaemonv2.UpdatePersistentTaskRequest, opts ...grpc.CallOption) error {
	ctx, cancel := context.WithTimeout(ctx, contextTimeout)
	defer cancel()

	_, err := v.DfdaemonUploadClient.UpdatePersistentTask(ctx, req, opts...)
	return err
}

// StatPersistentTask stats persistent task information.
func (v *v2) StatPersistentTask(ctx context.Context, req *dfdaemonv2.StatPersistentTaskRequest, opts ...grpc.CallOption) (*commonv2.PersistentTask, error) {
	ctx, cancel := context.WithTimeout(ctx, contextTimeout)
	defer cancel()

	return v.DfdaemonUploadClient.StatPersistentTask(ctx, req, opts...)
}

// DeletePersistentTask deletes persistent task from p2p network.
func (v *v2) DeletePersistentTask(ctx context.Context, req *dfdaemonv2.DeletePersistentTaskRequest, opts ...grpc.CallOption) error {
	ctx, cancel := context.WithTimeout(ctx, contextTimeout)
	defer cancel()

	_, err := v.DfdaemonUploadClient.DeletePersistentTask(ctx, req, opts...)
	return err
}

// DownloadPersistentCacheTask downloads persistent cache task from p2p network.
func (v *v2) DownloadPersistentCacheTask(ctx context.Context, req *dfdaemonv2.DownloadPersistentCacheTaskRequest, opts ...grpc.CallOption) (dfdaemonv2.DfdaemonUpload_DownloadPersistentCacheTaskClient, error) {
	return v.DfdaemonUploadClient.DownloadPersistentCacheTask(ctx, req, opts...)
}

// UpdatePersistentCacheTask updates persistent cache task information.
func (v *v2) UpdatePersistentCacheTask(ctx context.Context, req *dfdaemonv2.UpdatePersistentCacheTaskRequest, opts ...grpc.CallOption) error {
	ctx, cancel := context.WithTimeout(ctx, contextTimeout)
	defer cancel()

	_, err := v.DfdaemonUploadClient.UpdatePersistentCacheTask(ctx, req, opts...)
	return err
}

// StatPersistentCacheTask stats persistent cache task information.
func (v *v2) StatPersistentCacheTask(ctx context.Context, req *dfdaemonv2.StatPersistentCacheTaskRequest, opts ...grpc.CallOption) (*commonv2.PersistentCacheTask, error) {
	ctx, cancel := context.WithTimeout(ctx, contextTimeout)
	defer cancel()

	return v.DfdaemonUploadClient.StatPersistentCacheTask(ctx, req, opts...)
}

// DeletePersistentCacheTask deletes persistent cache task from p2p network.
func (v *v2) DeletePersistentCacheTask(ctx context.Context, req *dfdaemonv2.DeletePersistentCacheTaskRequest, opts ...grpc.CallOption) error {
	ctx, cancel := context.WithTimeout(ctx, contextTimeout)
	defer cancel()

	_, err := v.DfdaemonUploadClient.DeletePersistentCacheTask(ctx, req, opts...)
	return err
}

// ListTaskEntries lists task entries from p2p network.
func (v *v2) ListTaskEntries(ctx context.Context, req *dfdaemonv2.ListTaskEntriesRequest, opts ...grpc.CallOption) (*dfdaemonv2.ListTaskEntriesResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, contextTimeout)
	defer cancel()

	return v.DfdaemonUploadClient.ListTaskEntries(ctx, req, opts...)
}
