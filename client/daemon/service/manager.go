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

package service

import (
	"context"
	"fmt"
	"net"
	"os"
	"time"

	"google.golang.org/grpc"

	"github.com/dragonflyoss/Dragonfly2/client/daemon/peer"
	"github.com/dragonflyoss/Dragonfly2/client/daemon/storage"
	"github.com/dragonflyoss/Dragonfly2/client/util"
	logger "github.com/dragonflyoss/Dragonfly2/pkg/dflog"
	"github.com/dragonflyoss/Dragonfly2/pkg/rpc"
	"github.com/dragonflyoss/Dragonfly2/pkg/rpc/base"
	dfdaemongrpc "github.com/dragonflyoss/Dragonfly2/pkg/rpc/dfdaemon"
	"github.com/dragonflyoss/Dragonfly2/pkg/rpc/scheduler"
)

type Manager interface {
	util.KeepAlive
	ServeDownload(listener net.Listener) error
	ServePeer(listener net.Listener) error
	ServeProxy(listener net.Listener) error
	Stop()
}

type manager struct {
	util.KeepAlive
	peerHost        *scheduler.PeerHost
	peerTaskManager peer.PeerTaskManager
	storageManager  storage.Manager

	downloadServer rpc.Server
	peerServer     rpc.Server
	uploadAddr     string
}

func NewManager(peerHost *scheduler.PeerHost, peerTaskManager peer.PeerTaskManager, storageManager storage.Manager, downloadOpts []grpc.ServerOption, peerOpts []grpc.ServerOption) (Manager, error) {
	mgr := &manager{
		KeepAlive:       util.NewKeepAlive("service manager"),
		peerHost:        peerHost,
		peerTaskManager: peerTaskManager,
		storageManager:  storageManager,
	}
	mgr.downloadServer = rpc.NewServer(mgr, downloadOpts...)
	mgr.peerServer = rpc.NewServer(mgr, peerOpts...)
	return mgr, nil
}

func (m *manager) ServeDownload(listener net.Listener) error {
	return m.downloadServer.Serve(listener)
}

func (m *manager) ServePeer(listener net.Listener) error {
	m.uploadAddr = fmt.Sprintf("%s:%d", m.peerHost.Ip, m.peerHost.DownPort)
	return m.peerServer.Serve(listener)
}

func (m *manager) ServeProxy(lis net.Listener) error {
	// TODO
	return nil
}

func (m *manager) Stop() {
	m.peerServer.GracefulStop()
	m.downloadServer.GracefulStop()
}

func (m *manager) GetPieceTasks(ctx context.Context, request *base.PieceTaskRequest) (*base.PiecePacket, error) {
	m.Keep()
	p, err := m.storageManager.GetPieces(ctx, request)
	if err != nil {
		return &base.PiecePacket{
			State: &base.ResponseState{
				Success: false,
				Code:    base.Code_CLIENT_ERROR,
				Msg:     fmt.Sprintf("get pieces error: %s", err),
			},
			TaskId: request.TaskId,
		}, nil
	}
	p.DstAddr = m.uploadAddr
	return p, nil
}

func (m *manager) CheckHealth(context.Context) (*base.ResponseState, error) {
	m.Keep()
	return &base.ResponseState{
		Success: true,
		Code:    base.Code_SUCCESS,
		Msg:     "Running",
	}, nil
}

func (m *manager) Download(ctx context.Context,
	req *dfdaemongrpc.DownRequest, results chan<- *dfdaemongrpc.DownResult) error {
	m.Keep()
	// init peer task request, peer download request uses different peer id
	peerTask := &peer.FilePeerTaskRequest{
		PeerTaskRequest: scheduler.PeerTaskRequest{
			Url:      req.Url,
			Filter:   req.Filter,
			BizId:    req.BizId,
			UrlMata:  req.UrlMeta,
			PeerId:   m.GenPeerID(),
			PeerHost: m.peerHost,
		},
		Output: req.Output,
	}

	peerTaskProgress, err := m.peerTaskManager.StartFilePeerTask(ctx, peerTask)
	if err != nil {
		return err
	}
loop:
	for {
		select {
		case p := <-peerTaskProgress:
			results <- &dfdaemongrpc.DownResult{
				State:           p.State,
				TaskId:          p.TaskId,
				PeerId:          p.PeerID,
				CompletedLength: uint64(p.CompletedLength),
				Done:            p.Done,
			}
			if p.Done {
				break loop
			}
		case <-ctx.Done():
			results <- &dfdaemongrpc.DownResult{
				State: &base.ResponseState{
					Success: false,
					Code:    base.Code_REQUEST_TIME_OUT,
					Msg:     fmt.Sprintf("%s", ctx.Err()),
				},
				CompletedLength: 0,
				Done:            true,
			}
			logger.Infof("context done due to %s", ctx.Err())
			return ctx.Err()
		}
	}
	return nil
}

func (m *manager) GenPeerID() string {
	// FIXME review peer id format
	return fmt.Sprintf("%s-%d-%d-%d",
		m.peerHost.Ip, m.peerHost.RpcPort, os.Getpid(), time.Now().UnixNano())
}
