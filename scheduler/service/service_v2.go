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

package service

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/bits-and-blooms/bitset"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	commonv2 "d7y.io/api/v2/pkg/apis/common/v2"
	dfdaemonv2 "d7y.io/api/v2/pkg/apis/dfdaemon/v2"
	schedulerv2 "d7y.io/api/v2/pkg/apis/scheduler/v2"

	logger "d7y.io/dragonfly/v2/internal/dflog"
	"d7y.io/dragonfly/v2/pkg/container/set"
	"d7y.io/dragonfly/v2/pkg/digest"
	"d7y.io/dragonfly/v2/pkg/net/http"
	dfdaemonclient "d7y.io/dragonfly/v2/pkg/rpc/dfdaemon/client"
	"d7y.io/dragonfly/v2/pkg/types"
	"d7y.io/dragonfly/v2/scheduler/config"
	"d7y.io/dragonfly/v2/scheduler/metrics"
	"d7y.io/dragonfly/v2/scheduler/resource/persistentcache"
	"d7y.io/dragonfly/v2/scheduler/resource/standard"
	"d7y.io/dragonfly/v2/scheduler/scheduling"
)

// V2 is the interface for v2 version of the service.
type V2 struct {
	// Resource interface.
	resource standard.Resource

	// Persistent cache resource interface.
	persistentCacheResource persistentcache.Resource

	// Scheduling interface.
	scheduling scheduling.Scheduling

	// Scheduler service config.
	config *config.Config

	// Dynamic config.
	dynconfig config.DynconfigInterface
}

// New v2 version of service instance.
func NewV2(
	cfg *config.Config,
	resource standard.Resource,
	persistentCacheResource persistentcache.Resource,
	scheduling scheduling.Scheduling,
	dynconfig config.DynconfigInterface,
) *V2 {
	return &V2{
		resource:                resource,
		persistentCacheResource: persistentCacheResource,
		scheduling:              scheduling,
		config:                  cfg,
		dynconfig:               dynconfig,
	}
}

// AnnouncePeer announces peer to scheduler.
func (v *V2) AnnouncePeer(stream schedulerv2.Scheduler_AnnouncePeerServer) error {
	ctx, cancel := context.WithCancel(stream.Context())
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			logger.Info("context was done")
			return ctx.Err()
		default:
		}

		req, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				return nil
			}

			logger.Errorf("receive error: %s", err.Error())
			return err
		}

		log := logger.WithPeer(req.GetHostId(), req.GetTaskId(), req.GetPeerId())
		switch announcePeerRequest := req.GetRequest().(type) {
		case *schedulerv2.AnnouncePeerRequest_RegisterPeerRequest:
			registerPeerRequest := announcePeerRequest.RegisterPeerRequest
			log.Infof("receive RegisterPeerRequest, url: %s, range: %#v, header: %#v, need back-to-source: %t",
				registerPeerRequest.Download.GetUrl(), registerPeerRequest.Download.GetRange(), registerPeerRequest.Download.GetRequestHeader(), registerPeerRequest.Download.GetNeedBackToSource())
			if err := v.handleRegisterPeerRequest(ctx, stream, req.GetHostId(), req.GetTaskId(), req.GetPeerId(), registerPeerRequest); err != nil {
				log.Error(err)

				// If the peer register failed, and set the peer state to failed. Peer will not need to report
				// the message of the peer download failed.
				if err := v.handleDownloadPeerFailedRequest(ctx, req.GetPeerId()); err != nil {
					log.Error(err)
					return err
				}

				log.Error(err)
				return err
			}
		case *schedulerv2.AnnouncePeerRequest_DownloadPeerStartedRequest:
			log.Info("receive DownloadPeerStartedRequest")
			if err := v.handleDownloadPeerStartedRequest(ctx, req.GetPeerId()); err != nil {
				log.Error(err)

				// If the peer started failed, and set the peer state to failed. Peer will not need to report
				// the message of the peer download failed.
				if err := v.handleDownloadPeerFailedRequest(ctx, req.GetPeerId()); err != nil {
					log.Error(err)
					return err
				}

				log.Error(err)
				return err
			}
		case *schedulerv2.AnnouncePeerRequest_DownloadPeerBackToSourceStartedRequest:
			log.Info("receive DownloadPeerBackToSourceStartedRequest")
			if err := v.handleDownloadPeerBackToSourceStartedRequest(ctx, req.GetPeerId()); err != nil {
				log.Error(err)

				// If the peer started back-to-source failed, and set the peer state to failed. Peer will not need to report
				// the message of the peer download failed.
				if err := v.handleDownloadPeerBackToSourceFailedRequest(ctx, req.GetPeerId()); err != nil {
					log.Error(err)
					return err
				}

				log.Error(err)
				return err
			}
		case *schedulerv2.AnnouncePeerRequest_ReschedulePeerRequest:
			reschedulePeerRequest := announcePeerRequest.ReschedulePeerRequest
			log.Infof("receive ReschedulePeerRequestescription: %s", reschedulePeerRequest.GetDescription())
			if err := v.handleReschedulePeerRequest(ctx, req.GetPeerId(), reschedulePeerRequest.GetCandidateParents()); err != nil {
				log.Error(err)

				// If the peer started back-to-source failed, and set the peer state to failed. Peer will not need to report
				// the message of the peer download failed.
				if err := v.handleDownloadPeerFailedRequest(ctx, req.GetPeerId()); err != nil {
					log.Error(err)
					return err
				}

				log.Error(err)
				return err
			}
		case *schedulerv2.AnnouncePeerRequest_DownloadPeerFinishedRequest:
			downloadPeerFinishedRequest := announcePeerRequest.DownloadPeerFinishedRequest
			log.Infof("receive DownloadPeerFinishedRequest, content length: %d, piece count: %d", downloadPeerFinishedRequest.GetContentLength(), downloadPeerFinishedRequest.GetPieceCount())
			if err := v.handleDownloadPeerFinishedRequest(ctx, req.GetPeerId()); err != nil {
				log.Error(err)
				return err
			}

			// If the task is succeeded, return nil directly and close the stream.
			return nil
		case *schedulerv2.AnnouncePeerRequest_DownloadPeerBackToSourceFinishedRequest:
			downloadPeerBackToSourceFinishedRequest := announcePeerRequest.DownloadPeerBackToSourceFinishedRequest
			log.Infof("receive DownloadPeerBackToSourceFinishedRequest, content length: %d, piece count: %d", downloadPeerBackToSourceFinishedRequest.GetContentLength(), downloadPeerBackToSourceFinishedRequest.GetPieceCount())
			if err := v.handleDownloadPeerBackToSourceFinishedRequest(ctx, req.GetPeerId(), downloadPeerBackToSourceFinishedRequest); err != nil {
				log.Error(err)

				// If the peer started back-to-source failed, and set the peer state to failed. Peer will not need to report
				// the message of the peer download failed.
				if err := v.handleDownloadPeerBackToSourceFailedRequest(ctx, req.GetPeerId()); err != nil {
					log.Error(err)
					return err
				}

				log.Error(err)
				return err
			}

			// If the task is back-to-source succeeded, return nil directly and close the stream.
			return nil
		case *schedulerv2.AnnouncePeerRequest_DownloadPeerFailedRequest:
			log.Infof("receive DownloadPeerFailedRequest, description: %s", announcePeerRequest.DownloadPeerFailedRequest.GetDescription())
			if err := v.handleDownloadPeerFailedRequest(ctx, req.GetPeerId()); err != nil {
				log.Error(err)
				return err
			}

			// If the task is failed, return nil directly and close the stream.
			return nil
		case *schedulerv2.AnnouncePeerRequest_DownloadPeerBackToSourceFailedRequest:
			log.Infof("receive DownloadPeerBackToSourceFailedRequest, description: %s", announcePeerRequest.DownloadPeerBackToSourceFailedRequest.GetDescription())
			if err := v.handleDownloadPeerBackToSourceFailedRequest(ctx, req.GetPeerId()); err != nil {
				log.Error(err)
				return err
			}

			// If the task is back-to-source failed, return nil directly and close the stream.
			return nil
		case *schedulerv2.AnnouncePeerRequest_DownloadPieceFinishedRequest:
			piece := announcePeerRequest.DownloadPieceFinishedRequest.Piece
			log.Infof("receive DownloadPieceFinishedRequest, piece number: %d, piece length: %d, traffic type: %s, cost: %s, parent id: %s", piece.GetNumber(), piece.GetLength(), piece.GetTrafficType(), piece.GetCost().AsDuration().String(), piece.GetParentId())
			if err := v.handleDownloadPieceFinishedRequest(req.GetPeerId(), announcePeerRequest.DownloadPieceFinishedRequest); err != nil {
				log.Error(err)
			}
		case *schedulerv2.AnnouncePeerRequest_DownloadPieceBackToSourceFinishedRequest:
			piece := announcePeerRequest.DownloadPieceBackToSourceFinishedRequest.Piece
			log.Infof("receive DownloadPieceBackToSourceFinishedRequest, piece number: %d, piece length: %d, traffic type: %s, cost: %s, parent id: %s", piece.GetNumber(), piece.GetLength(), piece.GetTrafficType(), piece.GetCost().AsDuration().String(), piece.GetParentId())
			if err := v.handleDownloadPieceBackToSourceFinishedRequest(ctx, req.GetPeerId(), announcePeerRequest.DownloadPieceBackToSourceFinishedRequest); err != nil {
				log.Error(err)
			}
		case *schedulerv2.AnnouncePeerRequest_DownloadPieceFailedRequest:
			downloadPieceFailedRequest := announcePeerRequest.DownloadPieceFailedRequest
			log.Infof("receive DownloadPieceFailedRequest, piece number: %d, temporary: %t, parent id: %s", downloadPieceFailedRequest.GetPieceNumber(), downloadPieceFailedRequest.GetTemporary(), downloadPieceFailedRequest.GetParentId())
			if err := v.handleDownloadPieceFailedRequest(ctx, req.GetPeerId(), downloadPieceFailedRequest); err != nil {
				log.Error(err)
			}
		case *schedulerv2.AnnouncePeerRequest_DownloadPieceBackToSourceFailedRequest:
			downloadPieceBackToSourceFailedRequest := announcePeerRequest.DownloadPieceBackToSourceFailedRequest
			log.Infof("receive DownloadPieceBackToSourceFailedRequest, piece number: %d", downloadPieceBackToSourceFailedRequest.GetPieceNumber())
			if err := v.handleDownloadPieceBackToSourceFailedRequest(ctx, req.GetPeerId(), downloadPieceBackToSourceFailedRequest); err != nil {
				log.Error(err)
			}
		default:
			msg := fmt.Sprintf("receive unknow request: %#v", announcePeerRequest)
			log.Error(msg)
			return status.Error(codes.FailedPrecondition, msg)
		}
	}
}

// StatPeer checks information of peer.
func (v *V2) StatPeer(ctx context.Context, req *schedulerv2.StatPeerRequest) (*commonv2.Peer, error) {
	logger.WithTaskID(req.GetTaskId()).Infof("stat peer request: %#v", req)

	peer, loaded := v.resource.PeerManager().Load(req.GetPeerId())
	if !loaded {
		return nil, status.Errorf(codes.NotFound, "peer %s not found", req.GetPeerId())
	}

	resp := &commonv2.Peer{
		Id:               peer.ID,
		Priority:         peer.Priority,
		Cost:             durationpb.New(peer.Cost.Load()),
		State:            peer.FSM.Current(),
		NeedBackToSource: peer.NeedBackToSource.Load(),
		CreatedAt:        timestamppb.New(peer.CreatedAt.Load()),
		UpdatedAt:        timestamppb.New(peer.UpdatedAt.Load()),
	}

	// Set range to response.
	if peer.Range != nil {
		resp.Range = &commonv2.Range{
			Start:  uint64(peer.Range.Start),
			Length: uint64(peer.Range.Length),
		}
	}

	// Set pieces to response.
	peer.Pieces.Range(func(key, value any) bool {
		piece, ok := value.(*standard.Piece)
		if !ok {
			peer.Log.Errorf("invalid piece %s %#v", key, value)
			return true
		}

		respPiece := &commonv2.Piece{
			Number:      uint32(piece.Number),
			ParentId:    &piece.ParentID,
			Offset:      piece.Offset,
			Length:      piece.Length,
			TrafficType: &piece.TrafficType,
			Cost:        durationpb.New(piece.Cost),
			CreatedAt:   timestamppb.New(piece.CreatedAt),
		}

		if piece.Digest != nil {
			respPiece.Digest = piece.Digest.String()
		}

		resp.Pieces = append(resp.Pieces, respPiece)
		return true
	})

	// Set task to response.
	resp.Task = &commonv2.Task{
		Id:                  peer.Task.ID,
		Type:                peer.Task.Type,
		Url:                 peer.Task.URL,
		Tag:                 &peer.Task.Tag,
		Application:         &peer.Task.Application,
		FilteredQueryParams: peer.Task.FilteredQueryParams,
		RequestHeader:       peer.Task.Header,
		ContentLength:       uint64(peer.Task.ContentLength.Load()),
		PieceCount:          uint32(peer.Task.TotalPieceCount.Load()),
		SizeScope:           peer.Task.SizeScope(),
		State:               peer.Task.FSM.Current(),
		PeerCount:           uint32(peer.Task.PeerCount()),
		CreatedAt:           timestamppb.New(peer.Task.CreatedAt.Load()),
		UpdatedAt:           timestamppb.New(peer.Task.UpdatedAt.Load()),
	}

	// Set digest to task response.
	if peer.Task.Digest != nil {
		dgst := peer.Task.Digest.String()
		resp.Task.Digest = &dgst
	}

	// Set pieces to task response.
	peer.Task.Pieces.Range(func(key, value any) bool {
		piece, ok := value.(*standard.Piece)
		if !ok {
			peer.Task.Log.Errorf("invalid piece %s %#v", key, value)
			return true
		}

		respPiece := &commonv2.Piece{
			Number:      uint32(piece.Number),
			ParentId:    &piece.ParentID,
			Offset:      piece.Offset,
			Length:      piece.Length,
			TrafficType: &piece.TrafficType,
			Cost:        durationpb.New(piece.Cost),
			CreatedAt:   timestamppb.New(piece.CreatedAt),
		}

		if piece.Digest != nil {
			respPiece.Digest = piece.Digest.String()
		}

		resp.Task.Pieces = append(resp.Task.Pieces, respPiece)
		return true
	})

	// Set host to response.
	resp.Host = &commonv2.Host{
		Id:              peer.Host.ID,
		Type:            uint32(peer.Host.Type),
		Hostname:        peer.Host.Hostname,
		Ip:              peer.Host.IP,
		Port:            peer.Host.Port,
		DownloadPort:    peer.Host.DownloadPort,
		Os:              peer.Host.OS,
		Platform:        peer.Host.Platform,
		PlatformFamily:  peer.Host.PlatformFamily,
		PlatformVersion: peer.Host.PlatformVersion,
		KernelVersion:   peer.Host.KernelVersion,
		Cpu: &commonv2.CPU{
			LogicalCount:   peer.Host.CPU.LogicalCount,
			PhysicalCount:  peer.Host.CPU.PhysicalCount,
			Percent:        peer.Host.CPU.Percent,
			ProcessPercent: peer.Host.CPU.ProcessPercent,
			Times: &commonv2.CPUTimes{
				User:      peer.Host.CPU.Times.User,
				System:    peer.Host.CPU.Times.System,
				Idle:      peer.Host.CPU.Times.Idle,
				Nice:      peer.Host.CPU.Times.Nice,
				Iowait:    peer.Host.CPU.Times.Iowait,
				Irq:       peer.Host.CPU.Times.Irq,
				Softirq:   peer.Host.CPU.Times.Softirq,
				Steal:     peer.Host.CPU.Times.Steal,
				Guest:     peer.Host.CPU.Times.Guest,
				GuestNice: peer.Host.CPU.Times.GuestNice,
			},
		},
		Memory: &commonv2.Memory{
			Total:              peer.Host.Memory.Total,
			Available:          peer.Host.Memory.Available,
			Used:               peer.Host.Memory.Used,
			UsedPercent:        peer.Host.Memory.UsedPercent,
			ProcessUsedPercent: peer.Host.Memory.ProcessUsedPercent,
			Free:               peer.Host.Memory.Free,
		},
		Network: &commonv2.Network{
			TcpConnectionCount:       peer.Host.Network.TCPConnectionCount,
			UploadTcpConnectionCount: peer.Host.Network.UploadTCPConnectionCount,
			Location:                 &peer.Host.Network.Location,
			Idc:                      &peer.Host.Network.IDC,
			DownloadRate:             peer.Host.Network.DownloadRate,
			DownloadRateLimit:        peer.Host.Network.DownloadRateLimit,
			UploadRate:               peer.Host.Network.UploadRate,
			UploadRateLimit:          peer.Host.Network.UploadRateLimit,
		},
		Disk: &commonv2.Disk{
			Total:             peer.Host.Disk.Total,
			Free:              peer.Host.Disk.Free,
			Used:              peer.Host.Disk.Used,
			UsedPercent:       peer.Host.Disk.UsedPercent,
			InodesTotal:       peer.Host.Disk.InodesTotal,
			InodesUsed:        peer.Host.Disk.InodesUsed,
			InodesFree:        peer.Host.Disk.InodesFree,
			InodesUsedPercent: peer.Host.Disk.InodesUsedPercent,
			WriteBandwidth:    peer.Host.Disk.WriteBandwidth,
			ReadBandwidth:     peer.Host.Disk.ReadBandwidth,
		},
		Build: &commonv2.Build{
			GitVersion: peer.Host.Build.GitVersion,
			GitCommit:  &peer.Host.Build.GitCommit,
			GoVersion:  &peer.Host.Build.GoVersion,
			Platform:   &peer.Host.Build.Platform,
		},
	}

	return resp, nil
}

// DeletePeer releases peer in scheduler.
func (v *V2) DeletePeer(_ctx context.Context, req *schedulerv2.DeletePeerRequest) error {
	// Context use background to avoid the context canceled by the client and
	// the peer deletion operation is not completed.
	ctx := context.Background()
	log := logger.WithTaskAndPeerID(req.GetTaskId(), req.GetPeerId())
	log.Infof("delete peer request: %#v", req)

	peer, loaded := v.resource.PeerManager().Load(req.GetPeerId())
	if !loaded {
		msg := fmt.Sprintf("peer %s not found", req.GetPeerId())
		log.Error(msg)
		return status.Error(codes.NotFound, msg)
	}

	if err := peer.FSM.Event(ctx, standard.PeerEventLeave); err != nil {
		err = fmt.Errorf("peer fsm event failed: %w", err)
		peer.Log.Error(err)
		return status.Error(codes.FailedPrecondition, err.Error())
	}

	return nil
}

// StatTask checks information of task.
func (v *V2) StatTask(ctx context.Context, req *schedulerv2.StatTaskRequest) (*commonv2.Task, error) {
	log := logger.WithTaskID(req.GetTaskId())
	log.Infof("stat task request: %#v", req)

	task, loaded := v.resource.TaskManager().Load(req.GetTaskId())
	if !loaded {
		msg := fmt.Sprintf("task %s not found", req.GetTaskId())
		log.Error(msg)
		return nil, status.Error(codes.NotFound, msg)
	}

	resp := &commonv2.Task{
		Id:                  task.ID,
		Type:                task.Type,
		Url:                 task.URL,
		Tag:                 &task.Tag,
		Application:         &task.Application,
		FilteredQueryParams: task.FilteredQueryParams,
		RequestHeader:       task.Header,
		ContentLength:       uint64(task.ContentLength.Load()),
		PieceCount:          uint32(task.TotalPieceCount.Load()),
		SizeScope:           task.SizeScope(),
		State:               task.FSM.Current(),
		PeerCount:           uint32(task.PeerCount()),
		CreatedAt:           timestamppb.New(task.CreatedAt.Load()),
		UpdatedAt:           timestamppb.New(task.UpdatedAt.Load()),
	}

	// Set digest to response.
	if task.Digest != nil {
		dgst := task.Digest.String()
		resp.Digest = &dgst
	}

	// Set pieces to response.
	task.Pieces.Range(func(key, value any) bool {
		piece, ok := value.(*standard.Piece)
		if !ok {
			task.Log.Errorf("invalid piece %s %#v", key, value)
			return true
		}

		respPiece := &commonv2.Piece{
			Number:      uint32(piece.Number),
			ParentId:    &piece.ParentID,
			Offset:      piece.Offset,
			Length:      piece.Length,
			TrafficType: &piece.TrafficType,
			Cost:        durationpb.New(piece.Cost),
			CreatedAt:   timestamppb.New(piece.CreatedAt),
		}

		if piece.Digest != nil {
			respPiece.Digest = piece.Digest.String()
		}

		resp.Pieces = append(resp.Pieces, respPiece)
		return true
	})

	return resp, nil
}

// DeleteTask releases task in scheduler.
func (v *V2) DeleteTask(_ctx context.Context, req *schedulerv2.DeleteTaskRequest) error {
	// Context use background to avoid the context canceled by the client and
	// the task deletion operation is not completed.
	ctx := context.Background()
	log := logger.WithHostAndTaskID(req.GetHostId(), req.GetTaskId())
	log.Infof("delete task request: %#v", req)

	host, loaded := v.resource.HostManager().Load(req.GetHostId())
	if !loaded {
		msg := fmt.Sprintf("host %s not found", req.GetHostId())
		log.Error(msg)
		return status.Error(codes.NotFound, msg)
	}

	host.Peers.Range(func(key, value any) bool {
		peer, ok := value.(*standard.Peer)
		if !ok {
			host.Log.Errorf("invalid peer %s %#v", key, value)
			return true
		}

		if peer.Task.ID == req.GetTaskId() {
			if err := peer.FSM.Event(ctx, standard.PeerEventLeave); err != nil {
				err = fmt.Errorf("peer fsm event failed: %w", err)
				peer.Log.Error(err)
				return true
			}
		}

		return true
	})

	return nil
}

// AnnounceHost announces host to scheduler.
func (v *V2) AnnounceHost(ctx context.Context, req *schedulerv2.AnnounceHostRequest) error {
	logger.WithHostID(req.Host.GetId()).Infof("announce host request: %#v", req.GetHost())

	// Get cluster config by manager.
	var concurrentUploadLimit int32
	switch types.HostType(req.Host.GetType()) {
	case types.HostTypeNormal:
		if clientConfig, err := v.dynconfig.GetSchedulerClusterClientConfig(); err == nil {
			concurrentUploadLimit = int32(clientConfig.LoadLimit)
		}
	case types.HostTypeSuperSeed, types.HostTypeStrongSeed, types.HostTypeWeakSeed:
		if seedPeerConfig, err := v.dynconfig.GetSeedPeerClusterConfig(); err == nil {
			concurrentUploadLimit = int32(seedPeerConfig.LoadLimit)
		}
	}

	// Handle standard host.
	host, loaded := v.resource.HostManager().Load(req.Host.GetId())
	if !loaded {
		options := []standard.HostOption{
			standard.WithDisableShared(req.Host.GetDisableShared()),
			standard.WithOS(req.Host.GetOs()),
			standard.WithPlatform(req.Host.GetPlatform()),
			standard.WithPlatformFamily(req.Host.GetPlatformFamily()),
			standard.WithPlatformVersion(req.Host.GetPlatformVersion()),
			standard.WithKernelVersion(req.Host.GetKernelVersion()),
			standard.WithSchedulerClusterID(uint64(v.config.Manager.SchedulerClusterID)),
		}

		if concurrentUploadLimit > 0 {
			options = append(options, standard.WithConcurrentUploadLimit(concurrentUploadLimit))
		}

		if req.Host.GetCpu() != nil {
			options = append(options, standard.WithCPU(standard.CPU{
				LogicalCount:   req.Host.Cpu.GetLogicalCount(),
				PhysicalCount:  req.Host.Cpu.GetPhysicalCount(),
				Percent:        req.Host.Cpu.GetPercent(),
				ProcessPercent: req.Host.Cpu.GetProcessPercent(),
				Times: standard.CPUTimes{
					User:      req.Host.Cpu.Times.GetUser(),
					System:    req.Host.Cpu.Times.GetSystem(),
					Idle:      req.Host.Cpu.Times.GetIdle(),
					Nice:      req.Host.Cpu.Times.GetNice(),
					Iowait:    req.Host.Cpu.Times.GetIowait(),
					Irq:       req.Host.Cpu.Times.GetIrq(),
					Softirq:   req.Host.Cpu.Times.GetSoftirq(),
					Steal:     req.Host.Cpu.Times.GetSteal(),
					Guest:     req.Host.Cpu.Times.GetGuest(),
					GuestNice: req.Host.Cpu.Times.GetGuest(),
				},
			}))
		}

		if req.Host.GetMemory() != nil {
			options = append(options, standard.WithMemory(standard.Memory{
				Total:              req.Host.Memory.GetTotal(),
				Available:          req.Host.Memory.GetAvailable(),
				Used:               req.Host.Memory.GetUsed(),
				UsedPercent:        req.Host.Memory.GetUsedPercent(),
				ProcessUsedPercent: req.Host.Memory.GetProcessUsedPercent(),
				Free:               req.Host.Memory.GetFree(),
			}))
		}

		if req.Host.GetNetwork() != nil {
			options = append(options, standard.WithNetwork(standard.Network{
				TCPConnectionCount:       req.Host.Network.GetTcpConnectionCount(),
				UploadTCPConnectionCount: req.Host.Network.GetUploadTcpConnectionCount(),
				Location:                 req.Host.Network.GetLocation(),
				IDC:                      req.Host.Network.GetIdc(),
				DownloadRate:             req.Host.Network.GetDownloadRate(),
				DownloadRateLimit:        req.Host.Network.GetDownloadRateLimit(),
				UploadRate:               req.Host.Network.GetUploadRate(),
				UploadRateLimit:          req.Host.Network.GetUploadRateLimit(),
			}))
		}

		if req.Host.GetDisk() != nil {
			options = append(options, standard.WithDisk(standard.Disk{
				Total:             req.Host.Disk.GetTotal(),
				Free:              req.Host.Disk.GetFree(),
				Used:              req.Host.Disk.GetUsed(),
				UsedPercent:       req.Host.Disk.GetUsedPercent(),
				InodesTotal:       req.Host.Disk.GetInodesTotal(),
				InodesUsed:        req.Host.Disk.GetInodesUsed(),
				InodesFree:        req.Host.Disk.GetInodesFree(),
				InodesUsedPercent: req.Host.Disk.GetInodesUsedPercent(),
				WriteBandwidth:    req.Host.Disk.GetWriteBandwidth(),
				ReadBandwidth:     req.Host.Disk.GetReadBandwidth(),
			}))
		}

		if req.Host.GetBuild() != nil {
			options = append(options, standard.WithBuild(standard.Build{
				GitVersion: req.Host.Build.GetGitVersion(),
				GitCommit:  req.Host.Build.GetGitCommit(),
				GoVersion:  req.Host.Build.GetGoVersion(),
				Platform:   req.Host.Build.GetPlatform(),
			}))
		}

		if req.GetInterval() != nil {
			options = append(options, standard.WithAnnounceInterval(req.GetInterval().AsDuration()))
		}

		host = standard.NewHost(
			req.Host.GetId(), req.Host.GetIp(), req.Host.GetHostname(),
			req.Host.GetPort(), req.Host.GetDownloadPort(), types.HostType(req.Host.GetType()),
			options...,
		)

		v.resource.HostManager().Store(host)
		host.Log.Infof("announce new host: %#v", req)
	} else {
		// Host already exists and updates properties.
		host.Port = req.Host.GetPort()
		host.DownloadPort = req.Host.GetDownloadPort()
		host.Type = types.HostType(req.Host.GetType())
		host.DisableShared = req.Host.GetDisableShared()
		host.OS = req.Host.GetOs()
		host.Platform = req.Host.GetPlatform()
		host.PlatformFamily = req.Host.GetPlatformFamily()
		host.PlatformVersion = req.Host.GetPlatformVersion()
		host.KernelVersion = req.Host.GetKernelVersion()
		host.UpdatedAt.Store(time.Now())

		if concurrentUploadLimit > 0 {
			host.ConcurrentUploadLimit.Store(concurrentUploadLimit)
		}

		if req.Host.GetCpu() != nil {
			host.CPU = standard.CPU{
				LogicalCount:   req.Host.Cpu.GetLogicalCount(),
				PhysicalCount:  req.Host.Cpu.GetPhysicalCount(),
				Percent:        req.Host.Cpu.GetPercent(),
				ProcessPercent: req.Host.Cpu.GetProcessPercent(),
				Times: standard.CPUTimes{
					User:      req.Host.Cpu.Times.GetUser(),
					System:    req.Host.Cpu.Times.GetSystem(),
					Idle:      req.Host.Cpu.Times.GetIdle(),
					Nice:      req.Host.Cpu.Times.GetNice(),
					Iowait:    req.Host.Cpu.Times.GetIowait(),
					Irq:       req.Host.Cpu.Times.GetIrq(),
					Softirq:   req.Host.Cpu.Times.GetSoftirq(),
					Steal:     req.Host.Cpu.Times.GetSteal(),
					Guest:     req.Host.Cpu.Times.GetGuest(),
					GuestNice: req.Host.Cpu.Times.GetGuestNice(),
				},
			}
		}

		if req.Host.GetMemory() != nil {
			host.Memory = standard.Memory{
				Total:              req.Host.Memory.GetTotal(),
				Available:          req.Host.Memory.GetAvailable(),
				Used:               req.Host.Memory.GetUsed(),
				UsedPercent:        req.Host.Memory.GetUsedPercent(),
				ProcessUsedPercent: req.Host.Memory.GetProcessUsedPercent(),
				Free:               req.Host.Memory.GetFree(),
			}
		}

		if req.Host.GetNetwork() != nil {
			host.Network = standard.Network{
				TCPConnectionCount:       req.Host.Network.GetTcpConnectionCount(),
				UploadTCPConnectionCount: req.Host.Network.GetUploadTcpConnectionCount(),
				Location:                 req.Host.Network.GetLocation(),
				IDC:                      req.Host.Network.GetIdc(),
				DownloadRate:             req.Host.Network.GetDownloadRate(),
				DownloadRateLimit:        req.Host.Network.GetDownloadRateLimit(),
				UploadRate:               req.Host.Network.GetUploadRate(),
				UploadRateLimit:          req.Host.Network.GetUploadRateLimit(),
			}
		}

		if req.Host.GetDisk() != nil {
			host.Disk = standard.Disk{
				Total:             req.Host.Disk.GetTotal(),
				Free:              req.Host.Disk.GetFree(),
				Used:              req.Host.Disk.GetUsed(),
				UsedPercent:       req.Host.Disk.GetUsedPercent(),
				InodesTotal:       req.Host.Disk.GetInodesTotal(),
				InodesUsed:        req.Host.Disk.GetInodesUsed(),
				InodesFree:        req.Host.Disk.GetInodesFree(),
				InodesUsedPercent: req.Host.Disk.GetInodesUsedPercent(),
				WriteBandwidth:    req.Host.Disk.GetWriteBandwidth(),
				ReadBandwidth:     req.Host.Disk.GetReadBandwidth(),
			}
		}

		if req.Host.GetBuild() != nil {
			host.Build = standard.Build{
				GitVersion: req.Host.Build.GetGitVersion(),
				GitCommit:  req.Host.Build.GetGitCommit(),
				GoVersion:  req.Host.Build.GetGoVersion(),
				Platform:   req.Host.Build.GetPlatform(),
			}
		}

		if req.GetInterval() != nil {
			host.AnnounceInterval = req.GetInterval().AsDuration()
		}
	}

	// Handle the persistent cache host. If redis is not enabled,
	// it will not support the persistent cache feature.
	if v.persistentCacheResource == nil {
		return nil
	}

	persistentCacheHost, loaded := v.persistentCacheResource.HostManager().Load(ctx, req.Host.GetId())
	if !loaded {
		persistentCacheHost = persistentcache.NewHost(req.Host.GetId(), req.Host.GetHostname(), req.Host.GetIp(), req.Host.GetOs(),
			req.Host.GetPlatform(), req.Host.GetPlatformFamily(), req.Host.GetPlatformVersion(), req.Host.GetKernelVersion(), req.Host.GetPort(),
			req.Host.GetDownloadPort(), req.Host.GetSchedulerClusterId(), req.Host.GetDisableShared(), types.HostType(req.Host.GetType()),
			persistentcache.CPU{
				LogicalCount:   req.Host.Cpu.GetLogicalCount(),
				PhysicalCount:  req.Host.Cpu.GetPhysicalCount(),
				Percent:        req.Host.Cpu.GetPercent(),
				ProcessPercent: req.Host.Cpu.GetProcessPercent(),
				Times: persistentcache.CPUTimes{
					User:      req.Host.Cpu.Times.GetUser(),
					System:    req.Host.Cpu.Times.GetSystem(),
					Idle:      req.Host.Cpu.Times.GetIdle(),
					Nice:      req.Host.Cpu.Times.GetNice(),
					Iowait:    req.Host.Cpu.Times.GetIowait(),
					Irq:       req.Host.Cpu.Times.GetIrq(),
					Softirq:   req.Host.Cpu.Times.GetSoftirq(),
					Steal:     req.Host.Cpu.Times.GetSteal(),
					Guest:     req.Host.Cpu.Times.GetGuest(),
					GuestNice: req.Host.Cpu.Times.GetGuestNice(),
				},
			},
			persistentcache.Memory{
				Total:              req.Host.Memory.GetTotal(),
				Available:          req.Host.Memory.GetAvailable(),
				Used:               req.Host.Memory.GetUsed(),
				UsedPercent:        req.Host.Memory.GetUsedPercent(),
				ProcessUsedPercent: req.Host.Memory.GetProcessUsedPercent(),
				Free:               req.Host.Memory.GetFree(),
			},
			persistentcache.Network{
				TCPConnectionCount:       req.Host.Network.GetTcpConnectionCount(),
				UploadTCPConnectionCount: req.Host.Network.GetUploadTcpConnectionCount(),
				Location:                 req.Host.Network.GetLocation(),
				IDC:                      req.Host.Network.GetIdc(),
				DownloadRate:             req.Host.Network.GetDownloadRate(),
				DownloadRateLimit:        req.Host.Network.GetDownloadRateLimit(),
				UploadRate:               req.Host.Network.GetUploadRate(),
				UploadRateLimit:          req.Host.Network.GetUploadRateLimit(),
			},
			persistentcache.Disk{
				Total:             req.Host.Disk.GetTotal(),
				Free:              req.Host.Disk.GetFree(),
				Used:              req.Host.Disk.GetUsed(),
				UsedPercent:       req.Host.Disk.GetUsedPercent(),
				InodesTotal:       req.Host.Disk.GetInodesTotal(),
				InodesUsed:        req.Host.Disk.GetInodesUsed(),
				InodesFree:        req.Host.Disk.GetInodesFree(),
				InodesUsedPercent: req.Host.Disk.GetInodesUsedPercent(),
				WriteBandwidth:    req.Host.Disk.GetWriteBandwidth(),
				ReadBandwidth:     req.Host.Disk.GetReadBandwidth(),
			},
			persistentcache.Build{
				GitVersion: req.Host.Build.GetGitVersion(),
				GitCommit:  req.Host.Build.GetGitCommit(),
				GoVersion:  req.Host.Build.GetGoVersion(),
				Platform:   req.Host.Build.GetPlatform(),
			},
			req.GetInterval().AsDuration(),
			time.Now(),
			time.Now(),
			logger.WithHostID(req.Host.GetId()),
		)

		persistentCacheHost.Log.Infof("announce new persistent cache host: %#v", req)
		if err := v.persistentCacheResource.HostManager().Store(ctx, persistentCacheHost); err != nil {
			persistentCacheHost.Log.Errorf("store persistent cache host failed: %s", err)
			return err
		}
	} else {
		// persistentCacheHost already exists and updates properties.
		persistentCacheHost.Port = req.Host.GetPort()
		persistentCacheHost.DownloadPort = req.Host.GetDownloadPort()
		persistentCacheHost.Type = types.HostType(req.Host.GetType())
		persistentCacheHost.DisableShared = req.Host.GetDisableShared()
		persistentCacheHost.SchedulerClusterID = req.Host.GetSchedulerClusterId()
		persistentCacheHost.OS = req.Host.GetOs()
		persistentCacheHost.Platform = req.Host.GetPlatform()
		persistentCacheHost.PlatformFamily = req.Host.GetPlatformFamily()
		persistentCacheHost.PlatformVersion = req.Host.GetPlatformVersion()
		persistentCacheHost.KernelVersion = req.Host.GetKernelVersion()
		persistentCacheHost.UpdatedAt = time.Now()

		if req.Host.GetCpu() != nil {
			persistentCacheHost.CPU = persistentcache.CPU{
				LogicalCount:   req.Host.Cpu.GetLogicalCount(),
				PhysicalCount:  req.Host.Cpu.GetPhysicalCount(),
				Percent:        req.Host.Cpu.GetPercent(),
				ProcessPercent: req.Host.Cpu.GetProcessPercent(),
				Times: persistentcache.CPUTimes{
					User:      req.Host.Cpu.Times.GetUser(),
					System:    req.Host.Cpu.Times.GetSystem(),
					Idle:      req.Host.Cpu.Times.GetIdle(),
					Nice:      req.Host.Cpu.Times.GetNice(),
					Iowait:    req.Host.Cpu.Times.GetIowait(),
					Irq:       req.Host.Cpu.Times.GetIrq(),
					Softirq:   req.Host.Cpu.Times.GetSoftirq(),
					Steal:     req.Host.Cpu.Times.GetSteal(),
					Guest:     req.Host.Cpu.Times.GetGuest(),
					GuestNice: req.Host.Cpu.Times.GetGuestNice(),
				},
			}
		}

		if req.Host.GetMemory() != nil {
			persistentCacheHost.Memory = persistentcache.Memory{
				Total:              req.Host.Memory.GetTotal(),
				Available:          req.Host.Memory.GetAvailable(),
				Used:               req.Host.Memory.GetUsed(),
				UsedPercent:        req.Host.Memory.GetUsedPercent(),
				ProcessUsedPercent: req.Host.Memory.GetProcessUsedPercent(),
				Free:               req.Host.Memory.GetFree(),
			}
		}

		if req.Host.GetNetwork() != nil {
			persistentCacheHost.Network = persistentcache.Network{
				TCPConnectionCount:       req.Host.Network.GetTcpConnectionCount(),
				UploadTCPConnectionCount: req.Host.Network.GetUploadTcpConnectionCount(),
				Location:                 req.Host.Network.GetLocation(),
				IDC:                      req.Host.Network.GetIdc(),
				DownloadRate:             req.Host.Network.GetDownloadRate(),
				DownloadRateLimit:        req.Host.Network.GetDownloadRateLimit(),
				UploadRate:               req.Host.Network.GetUploadRate(),
				UploadRateLimit:          req.Host.Network.GetUploadRateLimit(),
			}
		}

		if req.Host.GetDisk() != nil {
			persistentCacheHost.Disk = persistentcache.Disk{
				Total:             req.Host.Disk.GetTotal(),
				Free:              req.Host.Disk.GetFree(),
				Used:              req.Host.Disk.GetUsed(),
				UsedPercent:       req.Host.Disk.GetUsedPercent(),
				InodesTotal:       req.Host.Disk.GetInodesTotal(),
				InodesUsed:        req.Host.Disk.GetInodesUsed(),
				InodesFree:        req.Host.Disk.GetInodesFree(),
				InodesUsedPercent: req.Host.Disk.GetInodesUsedPercent(),
				WriteBandwidth:    req.Host.Disk.GetWriteBandwidth(),
				ReadBandwidth:     req.Host.Disk.GetReadBandwidth(),
			}
		}

		if req.Host.GetBuild() != nil {
			persistentCacheHost.Build = persistentcache.Build{
				GitVersion: req.Host.Build.GetGitVersion(),
				GitCommit:  req.Host.Build.GetGitCommit(),
				GoVersion:  req.Host.Build.GetGoVersion(),
				Platform:   req.Host.Build.GetPlatform(),
			}
		}

		if req.GetInterval() != nil {
			persistentCacheHost.AnnounceInterval = req.GetInterval().AsDuration()
		}

		persistentCacheHost.Log.Infof("update persistent cache host: %#v", req)
		if err := v.persistentCacheResource.HostManager().Store(ctx, persistentCacheHost); err != nil {
			persistentCacheHost.Log.Errorf("store persistent cache host failed: %s", err)
			return err
		}
	}

	return nil
}

// ListHosts lists hosts in scheduler.
func (v *V2) ListHosts(ctx context.Context) (*schedulerv2.ListHostsResponse, error) {
	hosts := []*commonv2.Host{}
	v.resource.HostManager().Range(func(_ any, value any) bool {
		host, ok := value.(*standard.Host)
		if !ok {
			// Continue to next host.
			logger.Warnf("invalid host %#v", value)
			return true
		}

		hosts = append(hosts, &commonv2.Host{
			Id:              host.ID,
			Type:            uint32(host.Type),
			Hostname:        host.Hostname,
			Ip:              host.IP,
			Port:            host.Port,
			DownloadPort:    host.DownloadPort,
			Os:              host.OS,
			Platform:        host.Platform,
			PlatformFamily:  host.PlatformFamily,
			PlatformVersion: host.PlatformVersion,
			KernelVersion:   host.KernelVersion,
			Cpu: &commonv2.CPU{
				LogicalCount:   host.CPU.LogicalCount,
				PhysicalCount:  host.CPU.PhysicalCount,
				Percent:        host.CPU.Percent,
				ProcessPercent: host.CPU.ProcessPercent,
				Times: &commonv2.CPUTimes{
					User:      host.CPU.Times.User,
					System:    host.CPU.Times.System,
					Idle:      host.CPU.Times.Idle,
					Nice:      host.CPU.Times.Nice,
					Iowait:    host.CPU.Times.Iowait,
					Irq:       host.CPU.Times.Irq,
					Softirq:   host.CPU.Times.Softirq,
					Steal:     host.CPU.Times.Steal,
					Guest:     host.CPU.Times.Guest,
					GuestNice: host.CPU.Times.GuestNice,
				},
			},
			Memory: &commonv2.Memory{
				Total:              host.Memory.Total,
				Available:          host.Memory.Available,
				Used:               host.Memory.Used,
				UsedPercent:        host.Memory.UsedPercent,
				ProcessUsedPercent: host.Memory.ProcessUsedPercent,
				Free:               host.Memory.Free,
			},
			Network: &commonv2.Network{
				TcpConnectionCount:       host.Network.TCPConnectionCount,
				UploadTcpConnectionCount: host.Network.UploadTCPConnectionCount,
				Location:                 &host.Network.Location,
				Idc:                      &host.Network.IDC,
			},
			Disk: &commonv2.Disk{
				Total:             host.Disk.Total,
				Free:              host.Disk.Free,
				Used:              host.Disk.Used,
				UsedPercent:       host.Disk.UsedPercent,
				InodesTotal:       host.Disk.InodesTotal,
				InodesUsed:        host.Disk.InodesUsed,
				InodesFree:        host.Disk.InodesFree,
				InodesUsedPercent: host.Disk.InodesUsedPercent,
				WriteBandwidth:    host.Disk.WriteBandwidth,
				ReadBandwidth:     host.Disk.ReadBandwidth,
			},
			Build: &commonv2.Build{
				GitVersion: host.Build.GitVersion,
				GitCommit:  &host.Build.GitCommit,
				GoVersion:  &host.Build.GoVersion,
				Platform:   &host.Build.Platform,
			},
			SchedulerClusterId: host.SchedulerClusterID,
			DisableShared:      host.DisableShared,
		})

		return true
	})

	return &schedulerv2.ListHostsResponse{
		Hosts: hosts,
	}, nil
}

// DeleteHost releases host in scheduler.
func (v *V2) DeleteHost(_ctx context.Context, req *schedulerv2.DeleteHostRequest) error {
	// Context use background to avoid the context canceled by the client and
	// the host deletion operation is not completed.
	ctx := context.Background()
	log := logger.WithHostID(req.GetHostId())
	log.Infof("delete host request: %#v", req)

	host, loaded := v.resource.HostManager().Load(req.GetHostId())
	if !loaded {
		msg := fmt.Sprintf("host %s not found", req.GetHostId())
		log.Error(msg)
		return status.Error(codes.NotFound, msg)
	}

	// Leave peers in host.
	host.LeavePeers()

	// Delete host in scheduler.
	v.resource.HostManager().Delete(req.GetHostId())

	// Handle the persistent cache host for deletion.
	if v.persistentCacheResource != nil {
		peers, err := v.persistentCacheResource.PeerManager().LoadAllByHostID(ctx, req.GetHostId())
		if err != nil {
			log.Errorf("load persistent cache peers failed: %s", err)
		}

		for _, peer := range peers {
			if err := v.persistentCacheResource.PeerManager().Delete(ctx, peer.ID); err != nil {
				log.Errorf("delete persistent cache peer failed: %s", err)
			}

			// If peer is persistent, replicate the task to another peer.
			if peer.FSM.Is(persistentcache.PeerStateSucceeded) && peer.Persistent {
				blocklist := set.NewSafeSet[string]()
				blocklist.Add(peer.Host.ID)
				go func(peer *persistentcache.Peer, blocklist set.SafeSet[string]) {
					log.Infof("replicate persistent cache task %s", peer.Task.ID)
					if err := v.replicatePersistentCacheTask(context.Background(), peer, blocklist); err != nil {
						log.Errorf("replicate persistent cache task failed %s", err)
					}

					log.Infof("replicate persistent cache task %s finished", peer.Task.ID)
				}(peer, blocklist)
			}
		}

		if err := v.persistentCacheResource.HostManager().Delete(ctx, req.GetHostId()); err != nil {
			log.Errorf("delete persistent cache host failed: %s", err)
			return err
		}
	}

	return nil
}

// handleRegisterPeerRequest handles RegisterPeerRequest of AnnouncePeerRequest.
func (v *V2) handleRegisterPeerRequest(ctx context.Context, stream schedulerv2.Scheduler_AnnouncePeerServer, hostID, taskID, peerID string, req *schedulerv2.RegisterPeerRequest) error {
	// Handle resource included host, task, and peer.
	host, task, peer, err := v.handleResource(ctx, stream, hostID, taskID, peerID, req.GetDownload())
	if err != nil {
		return err
	}

	// Collect RegisterPeerCount metrics.
	priority := peer.CalculatePriority(v.dynconfig)
	metrics.RegisterPeerCount.WithLabelValues(priority.String(), peer.Task.Type.String(),
		peer.Host.Type.Name()).Inc()

	blocklist := set.NewSafeSet[string]()
	blocklist.Add(peer.ID)
	download := proto.Clone(req.Download).(*commonv2.Download)
	switch {
	// If scheduler trigger seed peer download back-to-source,
	// the needBackToSource flag should be true.
	case download.GetNeedBackToSource():
		peer.Log.Info("peer need back to source")
		peer.NeedBackToSource.Store(true)
	// If task is pending, failed, leave, or succeeded and has no available peer,
	// scheduler trigger seed peer download back-to-source.
	case task.FSM.Is(standard.TaskStatePending) ||
		task.FSM.Is(standard.TaskStateFailed) ||
		task.FSM.Is(standard.TaskStateLeave) ||
		task.FSM.Is(standard.TaskStateSucceeded) &&
			!task.HasAvailablePeer(blocklist):

		// If HostType is normal, trigger seed peer download back-to-source.
		if host.Type == types.HostTypeNormal {
			// If trigger the seed peer download back-to-source,
			// the need back-to-source flag should be true.
			download.NeedBackToSource = true

			// Output path should be empty, prevent the seed peer
			// copy file to output path.
			download.OutputPath = nil
			if err := v.downloadTaskBySeedPeer(ctx, taskID, download, peer); err != nil {
				// Collect RegisterPeerFailureCount metrics.
				metrics.RegisterPeerFailureCount.WithLabelValues(priority.String(), peer.Task.Type.String(),
					peer.Host.Type.Name()).Inc()
				return err
			}
		} else {
			// If HostType is not normal, peer is seed peer, and
			// trigger seed peer download back-to-source directly.
			peer.Log.Info("peer need back to source")
			peer.NeedBackToSource.Store(true)
		}
	}

	// Handle task with peer register request.
	if !peer.Task.FSM.Is(standard.TaskStateRunning) {
		if err := peer.Task.FSM.Event(ctx, standard.TaskEventDownload); err != nil {
			// Collect RegisterPeerFailureCount metrics.
			metrics.RegisterPeerFailureCount.WithLabelValues(priority.String(), peer.Task.Type.String(),
				peer.Host.Type.Name()).Inc()
			return status.Error(codes.Internal, err.Error())
		}
	} else {
		peer.Task.UpdatedAt.Store(time.Now())
	}

	// FSM event state transition by size scope.
	sizeScope := peer.Task.SizeScope()
	switch sizeScope {
	case commonv2.SizeScope_EMPTY:
		// Return an EmptyTaskResponse directly.
		peer.Log.Info("scheduling as SizeScope_EMPTY")
		stream, loaded := peer.LoadAnnouncePeerStream()
		if !loaded {
			return status.Error(codes.NotFound, "AnnouncePeerStream not found")
		}

		if err := peer.FSM.Event(ctx, standard.PeerEventRegisterEmpty); err != nil {
			return status.Error(codes.Internal, err.Error())
		}

		if err := stream.Send(&schedulerv2.AnnouncePeerResponse{
			Response: &schedulerv2.AnnouncePeerResponse_EmptyTaskResponse{
				EmptyTaskResponse: &schedulerv2.EmptyTaskResponse{},
			},
		}); err != nil {
			peer.Log.Error(err)
			return status.Error(codes.Internal, err.Error())
		}

		return nil
	case commonv2.SizeScope_NORMAL, commonv2.SizeScope_TINY, commonv2.SizeScope_SMALL, commonv2.SizeScope_UNKNOW:
		peer.Log.Info("scheduling as SizeScope_NORMAL")
		if err := peer.FSM.Event(ctx, standard.PeerEventRegisterNormal); err != nil {
			return status.Error(codes.Internal, err.Error())
		}

		// Scheduling parent for the peer.
		peer.BlockParents.Add(peer.ID)

		// Record the start time.
		start := time.Now()
		if err := v.scheduling.ScheduleCandidateParents(context.Background(), peer, peer.BlockParents); err != nil {
			// Collect RegisterPeerFailureCount metrics.
			metrics.RegisterPeerFailureCount.WithLabelValues(priority.String(), peer.Task.Type.String(),
				peer.Host.Type.Name()).Inc()
			return status.Error(codes.FailedPrecondition, err.Error())
		}

		// Collect SchedulingDuration metrics.
		metrics.ScheduleDuration.Observe(float64(time.Since(start).Milliseconds()))
		return nil
	default:
		return status.Errorf(codes.FailedPrecondition, "invalid size cope %#v", sizeScope)
	}
}

// handleDownloadPeerStartedRequest handles DownloadPeerStartedRequest of AnnouncePeerRequest.
func (v *V2) handleDownloadPeerStartedRequest(ctx context.Context, peerID string) error {
	peer, loaded := v.resource.PeerManager().Load(peerID)
	if !loaded {
		return status.Errorf(codes.NotFound, "peer %s not found", peerID)
	}

	// Collect DownloadPeerStartedCount metrics.
	priority := peer.CalculatePriority(v.dynconfig)
	metrics.DownloadPeerStartedCount.WithLabelValues(priority.String(), peer.Task.Type.String(),
		peer.Host.Type.Name()).Inc()

	// Handle peer with peer started request.
	if !peer.FSM.Is(standard.PeerStateRunning) {
		if err := peer.FSM.Event(ctx, standard.PeerEventDownload); err != nil {
			// Collect DownloadPeerStartedFailureCount metrics.
			metrics.DownloadPeerStartedFailureCount.WithLabelValues(priority.String(), peer.Task.Type.String(),
				peer.Host.Type.Name()).Inc()
			return status.Error(codes.Internal, err.Error())
		}
	}

	return nil
}

// handleDownloadPeerBackToSourceStartedRequest handles DownloadPeerBackToSourceStartedRequest of AnnouncePeerRequest.
func (v *V2) handleDownloadPeerBackToSourceStartedRequest(ctx context.Context, peerID string) error {
	peer, loaded := v.resource.PeerManager().Load(peerID)
	if !loaded {
		return status.Errorf(codes.NotFound, "peer %s not found", peerID)
	}

	// Collect DownloadPeerBackToSourceStartedCount metrics.
	priority := peer.CalculatePriority(v.dynconfig)
	metrics.DownloadPeerBackToSourceStartedCount.WithLabelValues(priority.String(), peer.Task.Type.String(),
		peer.Host.Type.Name()).Inc()

	// Handle peer with peer back-to-source started request.
	if !peer.FSM.Is(standard.PeerStateRunning) {
		if err := peer.FSM.Event(ctx, standard.PeerEventDownloadBackToSource); err != nil {
			// Collect DownloadPeerBackToSourceStartedFailureCount metrics.
			metrics.DownloadPeerBackToSourceStartedFailureCount.WithLabelValues(priority.String(), peer.Task.Type.String(),
				peer.Host.Type.Name()).Inc()
			return status.Error(codes.Internal, err.Error())
		}
	}

	return nil
}

// handleReschedulePeerRequest handles ReschedulePeerRequest of AnnouncePeerRequest.
func (v *V2) handleReschedulePeerRequest(_ context.Context, peerID string, candidateParents []*commonv2.Peer) error {
	peer, loaded := v.resource.PeerManager().Load(peerID)
	if !loaded {
		return status.Errorf(codes.NotFound, "peer %s not found", peerID)
	}

	// Add candidate parent ids to block parents.
	for _, candidateParent := range candidateParents {
		peer.BlockParents.Add(candidateParent.GetId())
	}

	// Record the start time.
	start := time.Now()
	if err := v.scheduling.ScheduleCandidateParents(context.Background(), peer, peer.BlockParents); err != nil {
		return status.Error(codes.FailedPrecondition, err.Error())
	}

	// Collect SchedulingDuration metrics.
	metrics.ScheduleDuration.Observe(float64(time.Since(start).Milliseconds()))
	return nil
}

// handleDownloadPeerFinishedRequest handles DownloadPeerFinishedRequest of AnnouncePeerRequest.
func (v *V2) handleDownloadPeerFinishedRequest(ctx context.Context, peerID string) error {
	peer, loaded := v.resource.PeerManager().Load(peerID)
	if !loaded {
		return status.Errorf(codes.NotFound, "peer %s not found", peerID)
	}

	// Handle peer with peer finished request.
	peer.Cost.Store(time.Since(peer.CreatedAt.Load()))
	if err := peer.FSM.Event(ctx, standard.PeerEventDownloadSucceeded); err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	// Collect DownloadPeerCount and DownloadPeerDuration metrics.
	priority := peer.CalculatePriority(v.dynconfig)
	metrics.DownloadPeerCount.WithLabelValues(priority.String(), peer.Task.Type.String(),
		peer.Host.Type.Name()).Inc()
	// TODO to be determined which traffic type to use, temporarily use TrafficType_REMOTE_PEER instead
	metrics.DownloadPeerDuration.WithLabelValues(metrics.CalculateSizeLevel(peer.Task.ContentLength.Load()).String()).Observe(float64(peer.Cost.Load().Milliseconds()))

	return nil
}

// handleDownloadPeerBackToSourceFinishedRequest handles DownloadPeerBackToSourceFinishedRequest of AnnouncePeerRequest.
func (v *V2) handleDownloadPeerBackToSourceFinishedRequest(ctx context.Context, peerID string, req *schedulerv2.DownloadPeerBackToSourceFinishedRequest) error {
	peer, loaded := v.resource.PeerManager().Load(peerID)
	if !loaded {
		return status.Errorf(codes.NotFound, "peer %s not found", peerID)
	}

	// Handle peer with peer back-to-source finished request.
	peer.Cost.Store(time.Since(peer.CreatedAt.Load()))
	if err := peer.FSM.Event(ctx, standard.PeerEventDownloadSucceeded); err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	// Handle task with peer back-to-source finished request, peer can only represent
	// a successful task after downloading the complete task.
	if peer.Range == nil && !peer.Task.FSM.Is(standard.TaskStateSucceeded) {
		peer.Task.ContentLength.Store(int64(req.GetContentLength()))
		peer.Task.TotalPieceCount.Store(int32(req.GetPieceCount()))
		if err := peer.Task.FSM.Event(ctx, standard.TaskEventDownloadSucceeded); err != nil {
			return status.Error(codes.Internal, err.Error())
		}
	}

	// Collect DownloadPeerCount and DownloadPeerDuration metrics.
	priority := peer.CalculatePriority(v.dynconfig)
	metrics.DownloadPeerCount.WithLabelValues(priority.String(), peer.Task.Type.String(),
		peer.Host.Type.Name()).Inc()
	// TODO to be determined which traffic type to use, temporarily use TrafficType_REMOTE_PEER instead
	metrics.DownloadPeerDuration.WithLabelValues(metrics.CalculateSizeLevel(peer.Task.ContentLength.Load()).String()).Observe(float64(peer.Cost.Load().Milliseconds()))

	return nil
}

// handleDownloadPeerFailedRequest handles DownloadPeerFailedRequest of AnnouncePeerRequest.
func (v *V2) handleDownloadPeerFailedRequest(ctx context.Context, peerID string) error {
	peer, loaded := v.resource.PeerManager().Load(peerID)
	if !loaded {
		return status.Errorf(codes.NotFound, "peer %s not found", peerID)
	}

	// Handle peer with peer failed request.
	if err := peer.FSM.Event(ctx, standard.PeerEventDownloadFailed); err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	// Collect DownloadPeerCount and DownloadPeerFailureCount metrics.
	priority := peer.CalculatePriority(v.dynconfig)
	metrics.DownloadPeerCount.WithLabelValues(priority.String(), peer.Task.Type.String(),
		peer.Host.Type.Name()).Inc()
	metrics.DownloadPeerFailureCount.WithLabelValues(priority.String(), peer.Task.Type.String(),
		peer.Host.Type.Name()).Inc()

	return nil
}

// handleDownloadPeerBackToSourceFailedRequest handles DownloadPeerBackToSourceFailedRequest of AnnouncePeerRequest.
func (v *V2) handleDownloadPeerBackToSourceFailedRequest(ctx context.Context, peerID string) error {
	peer, loaded := v.resource.PeerManager().Load(peerID)
	if !loaded {
		return status.Errorf(codes.NotFound, "peer %s not found", peerID)
	}

	// Handle peer with peer back-to-source failed request.
	if err := peer.FSM.Event(ctx, standard.PeerEventDownloadFailed); err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	// Handle task with peer back-to-source failed request.
	peer.Task.ContentLength.Store(-1)
	peer.Task.TotalPieceCount.Store(0)
	peer.Task.DirectPiece = []byte{}
	if err := peer.Task.FSM.Event(ctx, standard.TaskEventDownloadFailed); err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	// Collect DownloadPeerCount and DownloadPeerBackToSourceFailureCount metrics.
	priority := peer.CalculatePriority(v.dynconfig)
	metrics.DownloadPeerCount.WithLabelValues(priority.String(), peer.Task.Type.String(),
		peer.Host.Type.Name()).Inc()
	metrics.DownloadPeerBackToSourceFailureCount.WithLabelValues(priority.String(), peer.Task.Type.String(),
		peer.Host.Type.Name()).Inc()

	return nil
}

// handleDownloadPieceFinishedRequest handles DownloadPieceFinishedRequest of AnnouncePeerRequest.
func (v *V2) handleDownloadPieceFinishedRequest(peerID string, req *schedulerv2.DownloadPieceFinishedRequest) error {
	// Construct piece.
	piece := &standard.Piece{
		Number:      int32(req.Piece.GetNumber()),
		ParentID:    req.Piece.GetParentId(),
		Offset:      req.Piece.GetOffset(),
		Length:      req.Piece.GetLength(),
		TrafficType: req.Piece.GetTrafficType(),
		Cost:        req.Piece.GetCost().AsDuration(),
		CreatedAt:   req.Piece.GetCreatedAt().AsTime(),
	}

	if len(req.Piece.GetDigest()) > 0 {
		d, err := digest.Parse(req.Piece.GetDigest())
		if err != nil {
			return status.Error(codes.InvalidArgument, err.Error())
		}

		piece.Digest = d
	}

	peer, loaded := v.resource.PeerManager().Load(peerID)
	if !loaded {
		return status.Errorf(codes.NotFound, "peer %s not found", peerID)
	}

	// Handle peer with piece finished request. When the piece is downloaded successfully, peer.UpdatedAt needs
	// to be updated to prevent the peer from being GC during the download process.
	peer.StorePiece(piece)
	peer.FinishedPieces.Set(uint(piece.Number))
	peer.AppendPieceCost(piece.Cost)
	peer.PieceUpdatedAt.Store(time.Now())
	peer.UpdatedAt.Store(time.Now())

	// When the piece is downloaded successfully, parent.UpdatedAt needs to be updated
	// to prevent the parent from being GC during the download process.
	parent, loadedParent := v.resource.PeerManager().Load(piece.ParentID)
	if loadedParent {
		parent.UpdatedAt.Store(time.Now())
		parent.Host.UpdatedAt.Store(time.Now())
	}

	// Handle task with piece finished request.
	peer.Task.UpdatedAt.Store(time.Now())

	// Collect piece and traffic metrics.
	metrics.DownloadPieceCount.WithLabelValues(piece.TrafficType.String(), peer.Task.Type.String(),
		peer.Host.Type.Name()).Inc()
	metrics.Traffic.WithLabelValues(piece.TrafficType.String(), peer.Task.Type.String(),
		peer.Host.Type.Name()).Add(float64(piece.Length))
	if v.config.Metrics.EnableHost {
		metrics.HostTraffic.WithLabelValues(metrics.HostTrafficDownloadType, peer.Task.Type.String(),
			peer.Host.Type.Name(), peer.Host.ID, peer.Host.IP, peer.Host.Hostname).Add(float64(piece.Length))
		if loadedParent {
			metrics.HostTraffic.WithLabelValues(metrics.HostTrafficUploadType, peer.Task.Type.String(),
				parent.Host.Type.Name(), parent.Host.ID, parent.Host.IP, parent.Host.Hostname).Add(float64(piece.Length))
		}
	}

	return nil
}

// handleDownloadPieceBackToSourceFinishedRequest handles DownloadPieceBackToSourceFinishedRequest of AnnouncePeerRequest.
func (v *V2) handleDownloadPieceBackToSourceFinishedRequest(_ context.Context, peerID string, req *schedulerv2.DownloadPieceBackToSourceFinishedRequest) error {
	// Construct piece.
	piece := &standard.Piece{
		Number:      int32(req.Piece.GetNumber()),
		ParentID:    req.Piece.GetParentId(),
		Offset:      req.Piece.GetOffset(),
		Length:      req.Piece.GetLength(),
		TrafficType: req.Piece.GetTrafficType(),
		Cost:        req.Piece.GetCost().AsDuration(),
		CreatedAt:   req.Piece.GetCreatedAt().AsTime(),
	}

	if len(req.Piece.GetDigest()) > 0 {
		d, err := digest.Parse(req.Piece.GetDigest())
		if err != nil {
			return status.Error(codes.InvalidArgument, err.Error())
		}

		piece.Digest = d
	}

	peer, loaded := v.resource.PeerManager().Load(peerID)
	if !loaded {
		return status.Errorf(codes.NotFound, "peer %s not found", peerID)
	}

	// Handle peer with piece back-to-source finished request. When the piece is downloaded successfully, peer.UpdatedAt
	// needs to be updated to prevent the peer from being GC during the download process.
	peer.StorePiece(piece)
	peer.FinishedPieces.Set(uint(piece.Number))
	peer.AppendPieceCost(piece.Cost)
	peer.PieceUpdatedAt.Store(time.Now())
	peer.UpdatedAt.Store(time.Now())

	// Handle task with piece back-to-source finished request.
	peer.Task.StorePiece(piece)
	peer.Task.UpdatedAt.Store(time.Now())

	// Collect piece and traffic metrics.
	metrics.DownloadPieceCount.WithLabelValues(piece.TrafficType.String(), peer.Task.Type.String(),
		peer.Host.Type.Name()).Inc()
	metrics.Traffic.WithLabelValues(piece.TrafficType.String(), peer.Task.Type.String(),
		peer.Host.Type.Name()).Add(float64(piece.Length))
	if v.config.Metrics.EnableHost {
		metrics.HostTraffic.WithLabelValues(metrics.HostTrafficDownloadType, peer.Task.Type.String(),
			peer.Host.Type.Name(), peer.Host.ID, peer.Host.IP, peer.Host.Hostname).Add(float64(piece.Length))
	}

	return nil
}

// handleDownloadPieceFailedRequest handles DownloadPieceFailedRequest of AnnouncePeerRequest.
func (v *V2) handleDownloadPieceFailedRequest(_ context.Context, peerID string, req *schedulerv2.DownloadPieceFailedRequest) error {
	peer, loaded := v.resource.PeerManager().Load(peerID)
	if !loaded {
		return status.Errorf(codes.NotFound, "peer %s not found", peerID)
	}

	// Collect DownloadPieceCount and DownloadPieceFailureCount metrics.
	metrics.DownloadPieceCount.WithLabelValues(commonv2.TrafficType_REMOTE_PEER.String(), peer.Task.Type.String(),
		peer.Host.Type.Name()).Inc()
	metrics.DownloadPieceFailureCount.WithLabelValues(commonv2.TrafficType_REMOTE_PEER.String(), peer.Task.Type.String(),
		peer.Host.Type.Name()).Inc()

	if req.Temporary {
		// Handle peer with piece temporary failed request.
		peer.UpdatedAt.Store(time.Now())
		peer.BlockParents.Add(req.GetParentId())
		if parent, loaded := v.resource.PeerManager().Load(req.GetParentId()); loaded {
			parent.Host.UploadFailedCount.Inc()
		}

		// Handle task with piece temporary failed request.
		peer.Task.UpdatedAt.Store(time.Now())
		return nil
	}

	return status.Error(codes.FailedPrecondition, "download piece failed")
}

// handleDownloadPieceBackToSourceFailedRequest handles DownloadPieceBackToSourceFailedRequest of AnnouncePeerRequest.
func (v *V2) handleDownloadPieceBackToSourceFailedRequest(_ context.Context, peerID string, _ *schedulerv2.DownloadPieceBackToSourceFailedRequest) error {
	peer, loaded := v.resource.PeerManager().Load(peerID)
	if !loaded {
		return status.Errorf(codes.NotFound, "peer %s not found", peerID)
	}

	// Handle peer with piece back-to-source failed request.
	peer.UpdatedAt.Store(time.Now())

	// Handle task with piece back-to-source failed request.
	peer.Task.UpdatedAt.Store(time.Now())

	// Collect DownloadPieceCount and DownloadPieceFailureCount metrics.
	metrics.DownloadPieceCount.WithLabelValues(commonv2.TrafficType_BACK_TO_SOURCE.String(), peer.Task.Type.String(),
		peer.Host.Type.Name()).Inc()
	metrics.DownloadPieceFailureCount.WithLabelValues(commonv2.TrafficType_BACK_TO_SOURCE.String(), peer.Task.Type.String(),
		peer.Host.Type.Name()).Inc()

	return status.Error(codes.Internal, "download piece from source failed")
}

// handleResource handles resource included host, task, and peer.
func (v *V2) handleResource(_ context.Context, stream schedulerv2.Scheduler_AnnouncePeerServer, hostID, taskID, peerID string, download *commonv2.Download) (*standard.Host, *standard.Task, *standard.Peer, error) {
	// If the host does not exist and the host address cannot be found,
	// it may cause an exception.
	host, loaded := v.resource.HostManager().Load(hostID)
	if !loaded {
		return nil, nil, nil, status.Errorf(codes.NotFound, "host %s not found", hostID)
	}

	// Store new task or update task.
	task, loaded := v.resource.TaskManager().Load(taskID)
	if !loaded {
		options := []standard.TaskOption{}
		if download.GetDigest() != "" {
			d, err := digest.Parse(download.GetDigest())
			if err != nil {
				return nil, nil, nil, status.Error(codes.InvalidArgument, err.Error())
			}

			// If request has invalid digest, then new task with the nil digest.
			options = append(options, standard.WithDigest(d))
		}

		task = standard.NewTask(taskID, download.GetUrl(), download.GetTag(), download.GetApplication(), download.GetType(),
			download.GetFilteredQueryParams(), download.GetRequestHeader(), int32(v.config.Scheduler.BackToSourceCount), options...)
		v.resource.TaskManager().Store(task)
	} else {
		task.URL = download.GetUrl()
		task.FilteredQueryParams = download.GetFilteredQueryParams()
		task.Header = download.GetRequestHeader()
	}

	// Store new peer or load peer.
	peer, loaded := v.resource.PeerManager().Load(peerID)
	if !loaded {
		options := []standard.PeerOption{standard.WithPriority(download.GetPriority()), standard.WithAnnouncePeerStream(stream)}
		if download.GetRange() != nil {
			options = append(options, standard.WithRange(http.Range{Start: int64(download.Range.GetStart()), Length: int64(download.Range.GetLength())}))
		}

		peer = standard.NewPeer(peerID, task, host, options...)
		v.resource.PeerManager().Store(peer)
	}

	return host, task, peer, nil
}

// downloadTaskBySeedPeer downloads task by seed peer.
func (v *V2) downloadTaskBySeedPeer(ctx context.Context, taskID string, download *commonv2.Download, peer *standard.Peer) error {
	// Trigger the first download task based on different priority levels,
	// refer to https://github.com/dragonflyoss/api/blob/main/pkg/apis/common/v2/common.proto#L74.
	priority := peer.CalculatePriority(v.dynconfig)
	peer.Log.Infof("peer priority is %s", priority.String())
	switch priority {
	case commonv2.Priority_LEVEL6, commonv2.Priority_LEVEL0:
		// Super peer is first triggered to download back-to-source.
		if v.config.SeedPeer.Enable && !peer.Task.IsSeedPeerFailed() {
			go func(ctx context.Context, taskID string, download *commonv2.Download, hostType types.HostType) {
				peer.Log.Infof("%s seed peer triggers download task", hostType.Name())
				if err := v.resource.SeedPeer().TriggerDownloadTask(context.Background(), taskID, &dfdaemonv2.DownloadTaskRequest{Download: download}); err != nil {
					peer.Log.Errorf("%s seed peer triggers download task failed %s", hostType.Name(), err.Error())
					return
				}

				peer.Log.Infof("%s seed peer triggers download task success", hostType.Name())
			}(ctx, taskID, download, types.HostTypeSuperSeed)

			break
		}

		fallthrough
	case commonv2.Priority_LEVEL5:
		// Strong peer is first triggered to download back-to-source.
		if v.config.SeedPeer.Enable && !peer.Task.IsSeedPeerFailed() {
			go func(ctx context.Context, taskID string, download *commonv2.Download, hostType types.HostType) {
				peer.Log.Infof("%s seed peer triggers download task", hostType.Name())
				if err := v.resource.SeedPeer().TriggerDownloadTask(context.Background(), taskID, &dfdaemonv2.DownloadTaskRequest{Download: download}); err != nil {
					peer.Log.Errorf("%s seed peer triggers download task failed %s", hostType.Name(), err.Error())
					return
				}

				peer.Log.Infof("%s seed peer triggers download task success", hostType.Name())
			}(ctx, taskID, download, types.HostTypeSuperSeed)

			break
		}

		fallthrough
	case commonv2.Priority_LEVEL4:
		// Weak peer is first triggered to download back-to-source.
		if v.config.SeedPeer.Enable && !peer.Task.IsSeedPeerFailed() {
			go func(ctx context.Context, taskID string, download *commonv2.Download, hostType types.HostType) {
				peer.Log.Infof("%s seed peer triggers download task", hostType.Name())
				if err := v.resource.SeedPeer().TriggerDownloadTask(context.Background(), taskID, &dfdaemonv2.DownloadTaskRequest{Download: download}); err != nil {
					peer.Log.Errorf("%s seed peer triggers download task failed %s", hostType.Name(), err.Error())
					return
				}

				peer.Log.Infof("%s seed peer triggers download task success", hostType.Name())
			}(ctx, taskID, download, types.HostTypeSuperSeed)

			break
		}

		fallthrough
	case commonv2.Priority_LEVEL3:
		// When the task has no available peer,
		// the peer is first to download back-to-source.
		peer.NeedBackToSource.Store(true)
	case commonv2.Priority_LEVEL2:
		// Peer is first to download back-to-source.
		return status.Errorf(codes.NotFound, "%s peer not found candidate peers", commonv2.Priority_LEVEL2.String())
	case commonv2.Priority_LEVEL1:
		// Download task is forbidden.
		return status.Errorf(codes.FailedPrecondition, "%s peer is forbidden", commonv2.Priority_LEVEL1.String())
	default:
		return status.Errorf(codes.InvalidArgument, "invalid priority %#v", priority)
	}

	return nil
}

// AnnouncePersistentCachePeer announces persistent cache peer to scheduler.
func (v *V2) AnnouncePersistentCachePeer(stream schedulerv2.Scheduler_AnnouncePersistentCachePeerServer) error {
	if v.persistentCacheResource == nil {
		return status.Error(codes.FailedPrecondition, "redis is not enabled")
	}

	ctx, cancel := context.WithCancel(stream.Context())
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			logger.Info("context was done")
			return ctx.Err()
		default:
		}

		req, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				return nil
			}

			logger.Errorf("receive error: %s", err.Error())
			return err
		}

		log := logger.WithPeer(req.GetHostId(), req.GetTaskId(), req.GetPeerId())
		switch announcePersistentCachePeerRequest := req.GetRequest().(type) {
		case *schedulerv2.AnnouncePersistentCachePeerRequest_RegisterPersistentCachePeerRequest:
			registerPersistentCachePeerRequest := announcePersistentCachePeerRequest.RegisterPersistentCachePeerRequest
			log.Info("receive RegisterPersistentCachePeerRequest")
			if err := v.handleRegisterPersistentCachePeerRequest(ctx, stream, req.GetHostId(), req.GetTaskId(), req.GetPeerId(), registerPersistentCachePeerRequest.GetPersistent()); err != nil {
				// If the peer register failed, and set the peer state to failed. Peer will not need to report
				// the message of the peer download failed.
				if err := v.handleDownloadPersistentCachePeerFailedRequest(ctx, req.GetPeerId()); err != nil {
					log.Error(err)
					return err
				}

				log.Error(err)
				return err
			}
		case *schedulerv2.AnnouncePersistentCachePeerRequest_DownloadPersistentCachePeerStartedRequest:
			log.Info("receive DownloadPersistentCachePeerStartedRequest")
			if err := v.handleDownloadPersistentCachePeerStartedRequest(ctx, req.GetPeerId()); err != nil {
				// If the peer started failed, and set the peer state to failed. Peer will not need to report
				// the message of the peer download failed.
				if err := v.handleDownloadPersistentCachePeerFailedRequest(ctx, req.GetPeerId()); err != nil {
					log.Error(err)
					return err
				}

				log.Error(err)
				return err
			}
		case *schedulerv2.AnnouncePersistentCachePeerRequest_ReschedulePersistentCachePeerRequest:
			reschedulePersistentCachePeerRequest := announcePersistentCachePeerRequest.ReschedulePersistentCachePeerRequest
			log.Info("receive ReschedulePersistentCachePeerRequest")
			if err := v.handleReschedulePersistentCachePeerRequest(ctx, stream, req.GetTaskId(), req.GetPeerId(), reschedulePersistentCachePeerRequest); err != nil {
				// If the peer reschedule failed, and set the peer state to failed. Peer will not need to report
				// the message of the peer download failed.
				if err := v.handleDownloadPersistentCachePeerFailedRequest(ctx, req.GetPeerId()); err != nil {
					log.Error(err)
					return err
				}

				log.Error(err)
				return err
			}
		case *schedulerv2.AnnouncePersistentCachePeerRequest_DownloadPersistentCachePeerFinishedRequest:
			log.Info("receive DownloadPersistentCachePeerFinishedRequest")
			if err := v.handleDownloadPersistentCachePeerFinishedRequest(ctx, req.GetPeerId()); err != nil {
				// If the peer finished failed, and set the peer state to failed. Peer will not need to report
				// the message of the peer download failed.
				if err := v.handleDownloadPersistentCachePeerFailedRequest(ctx, req.GetPeerId()); err != nil {
					log.Error(err)
					return err
				}

				log.Error(err)
				return err
			}

			// If the task is succeeded, return nil directly and close the stream.
			return nil
		case *schedulerv2.AnnouncePersistentCachePeerRequest_DownloadPersistentCachePeerFailedRequest:
			log.Info("receive DownloadPersistentCachePeerFailedRequest")
			if err := v.handleDownloadPersistentCachePeerFailedRequest(ctx, req.GetPeerId()); err != nil {
				log.Error(err)
				return err
			}

			// If the task is failed, return nil directly and close the stream.
			return nil
		case *schedulerv2.AnnouncePersistentCachePeerRequest_DownloadPieceFinishedRequest:
			downloadPieceFinishedRequest := announcePersistentCachePeerRequest.DownloadPieceFinishedRequest

			log.Info("receive DownloadPieceFinishedRequest")
			go func() {
				if err := v.handleDownloadPersistentCachePieceFinishedRequest(context.Background(), req.GetPeerId(), downloadPieceFinishedRequest); err != nil {
					log.Error(err)
				}
			}()
		case *schedulerv2.AnnouncePersistentCachePeerRequest_DownloadPieceFailedRequest:
			downloadPieceFailedRequest := announcePersistentCachePeerRequest.DownloadPieceFailedRequest

			log.Info("receive DownloadPieceFailedRequest")
			go func() {
				if err := v.handleDownloadPersistentCachePieceFailedRequest(context.Background(), req.GetPeerId(), downloadPieceFailedRequest); err != nil {
					log.Error(err)
				}
			}()
		default:
			msg := fmt.Sprintf("receive unknow request: %#v", announcePersistentCachePeerRequest)
			log.Error(msg)
			return status.Error(codes.FailedPrecondition, msg)
		}
	}
}

// handleRegisterPersistentCachePeerRequest handles RegisterPersistentCachePeerRequest of AnnouncePersistentCachePeerRequest.
func (v *V2) handleRegisterPersistentCachePeerRequest(ctx context.Context, stream schedulerv2.Scheduler_AnnouncePersistentCachePeerServer, hostID, taskID, peerID string, persistent bool) error {
	host, loaded := v.persistentCacheResource.HostManager().Load(ctx, hostID)
	if !loaded {
		return status.Errorf(codes.NotFound, "host %s not found", hostID)
	}

	task, loaded := v.persistentCacheResource.TaskManager().Load(ctx, taskID)
	if !loaded {
		return status.Errorf(codes.NotFound, "task %s not found", taskID)
	}

	peer := persistentcache.NewPeer(peerID, persistentcache.PeerStatePending, persistent, &bitset.BitSet{}, []string{}, task, host, 0, time.Now(), time.Now(), logger.WithPeer(hostID, taskID, peerID))

	// Collect RegisterPersistentCachePeerCount metrics.
	metrics.RegisterPersistentCachePeerCount.WithLabelValues(peer.Host.Type.Name()).Inc()

	// Handle task with persistent cache peer register request.
	if !peer.Task.FSM.Is(persistentcache.TaskStateSucceeded) {
		// Collect RegisterPersistentCachePeerFailureCount metrics.
		metrics.RegisterPersistentCachePeerFailureCount.WithLabelValues(peer.Host.Type.Name()).Inc()
		return status.Error(codes.Internal, "persistent cache task is not succeeded")
	} else {
		peer.Task.UpdatedAt = time.Now()
	}

	// FSM event state transition by size scope.
	sizeScope := peer.Task.SizeScope()
	switch sizeScope {
	case commonv2.SizeScope_EMPTY:
		// Return an EmptyTaskResponse directly.
		peer.Log.Info("scheduling as SizeScope_EMPTY")
		if err := peer.FSM.Event(ctx, persistentcache.PeerEventRegisterEmpty); err != nil {
			return status.Error(codes.Internal, err.Error())
		}

		// Update metadata of the persistent cache task.
		if err := v.persistentCacheResource.TaskManager().Store(ctx, peer.Task); err != nil {
			return status.Error(codes.Internal, err.Error())
		}

		// Update metadata of the persistent cache peer.
		if err := v.persistentCacheResource.PeerManager().Store(ctx, peer); err != nil {
			return status.Error(codes.Internal, err.Error())
		}

		if err := stream.Send(&schedulerv2.AnnouncePersistentCachePeerResponse{
			Response: &schedulerv2.AnnouncePersistentCachePeerResponse_EmptyPersistentCacheTaskResponse{
				EmptyPersistentCacheTaskResponse: &schedulerv2.EmptyPersistentCacheTaskResponse{},
			},
		}); err != nil {
			peer.Log.Error(err)
			return status.Error(codes.Internal, err.Error())
		}

		return nil
	case commonv2.SizeScope_NORMAL, commonv2.SizeScope_TINY, commonv2.SizeScope_SMALL, commonv2.SizeScope_UNKNOW:
		peer.Log.Info("scheduling as SizeScope_NORMAL")
		if err := peer.FSM.Event(ctx, persistentcache.PeerEventRegisterNormal); err != nil {
			return status.Error(codes.Internal, err.Error())
		}

		// Scheduling parent for the peer.
		blocklist := set.NewSafeSet[string]()
		blocklist.Add(peer.ID)

		parents, found := v.scheduling.FindCandidatePersistentCacheParents(ctx, peer, blocklist)
		if !found {
			// Collect RegisterPersistentCachePeerFailureCount metrics.
			metrics.RegisterPersistentCachePeerFailureCount.WithLabelValues(peer.Host.Type.Name()).Inc()
			return status.Error(codes.FailedPrecondition, "no candidate parents found")
		}

		currentPersistentReplicaCount, err := v.persistentCacheResource.TaskManager().LoadCurrentPersistentReplicaCount(ctx, taskID)
		if err != nil {
			// Collect RegisterPersistentCachePeerFailureCount metrics.
			metrics.RegisterPersistentCachePeerFailureCount.WithLabelValues(peer.Host.Type.Name()).Inc()
			return status.Error(codes.Internal, err.Error())
		}

		currentReplicaCount, err := v.persistentCacheResource.TaskManager().LoadCurrentReplicaCount(ctx, taskID)
		if err != nil {
			// Collect RegisterPersistentCachePeerFailureCount metrics.
			metrics.RegisterPersistentCachePeerFailureCount.WithLabelValues(peer.Host.Type.Name()).Inc()
			return status.Error(codes.Internal, err.Error())
		}

		candidateParents := make([]*commonv2.PersistentCachePeer, 0, len(parents))
		for _, parent := range parents {
			candidateParents = append(candidateParents, &commonv2.PersistentCachePeer{
				Id:         parent.ID,
				Persistent: parent.Persistent,
				Cost:       durationpb.New(parent.Cost),
				State:      parent.FSM.Current(),
				Task: &commonv2.PersistentCacheTask{
					Id:                            parent.Task.ID,
					PersistentReplicaCount:        parent.Task.PersistentReplicaCount,
					CurrentPersistentReplicaCount: currentPersistentReplicaCount,
					CurrentReplicaCount:           currentReplicaCount,
					Tag:                           &parent.Task.Tag,
					Application:                   &parent.Task.Application,
					PieceLength:                   parent.Task.PieceLength,
					ContentLength:                 parent.Task.ContentLength,
					PieceCount:                    parent.Task.TotalPieceCount,
					State:                         parent.Task.FSM.Current(),
					Ttl:                           durationpb.New(parent.Task.TTL),
					CreatedAt:                     timestamppb.New(parent.Task.CreatedAt),
					UpdatedAt:                     timestamppb.New(parent.Task.UpdatedAt),
				},
				Host: &commonv2.Host{
					Id:              parent.Host.ID,
					Type:            uint32(parent.Host.Type),
					Hostname:        parent.Host.Hostname,
					Ip:              parent.Host.IP,
					Port:            parent.Host.Port,
					DownloadPort:    parent.Host.DownloadPort,
					Os:              parent.Host.OS,
					Platform:        parent.Host.Platform,
					PlatformFamily:  parent.Host.PlatformFamily,
					PlatformVersion: parent.Host.PlatformVersion,
					KernelVersion:   parent.Host.KernelVersion,
					Cpu: &commonv2.CPU{
						LogicalCount:   parent.Host.CPU.LogicalCount,
						PhysicalCount:  parent.Host.CPU.PhysicalCount,
						Percent:        parent.Host.CPU.Percent,
						ProcessPercent: parent.Host.CPU.ProcessPercent,
						Times: &commonv2.CPUTimes{
							User:      parent.Host.CPU.Times.User,
							System:    parent.Host.CPU.Times.System,
							Idle:      parent.Host.CPU.Times.Idle,
							Nice:      parent.Host.CPU.Times.Nice,
							Iowait:    parent.Host.CPU.Times.Iowait,
							Irq:       parent.Host.CPU.Times.Irq,
							Softirq:   parent.Host.CPU.Times.Softirq,
							Steal:     parent.Host.CPU.Times.Steal,
							Guest:     parent.Host.CPU.Times.Guest,
							GuestNice: parent.Host.CPU.Times.GuestNice,
						},
					},
					Memory: &commonv2.Memory{
						Total:              parent.Host.Memory.Total,
						Available:          parent.Host.Memory.Available,
						Used:               parent.Host.Memory.Used,
						UsedPercent:        parent.Host.Memory.UsedPercent,
						ProcessUsedPercent: parent.Host.Memory.ProcessUsedPercent,
						Free:               parent.Host.Memory.Free,
					},
					Network: &commonv2.Network{
						TcpConnectionCount:       parent.Host.Network.TCPConnectionCount,
						UploadTcpConnectionCount: parent.Host.Network.UploadTCPConnectionCount,
						Location:                 &parent.Host.Network.Location,
						Idc:                      &parent.Host.Network.IDC,
						DownloadRate:             parent.Host.Network.DownloadRate,
						DownloadRateLimit:        parent.Host.Network.DownloadRateLimit,
						UploadRate:               parent.Host.Network.UploadRate,
						UploadRateLimit:          parent.Host.Network.UploadRateLimit,
					},
					Disk: &commonv2.Disk{
						Total:             parent.Host.Disk.Total,
						Free:              parent.Host.Disk.Free,
						Used:              parent.Host.Disk.Used,
						UsedPercent:       parent.Host.Disk.UsedPercent,
						InodesTotal:       parent.Host.Disk.InodesTotal,
						InodesUsed:        parent.Host.Disk.InodesUsed,
						InodesFree:        parent.Host.Disk.InodesFree,
						InodesUsedPercent: parent.Host.Disk.InodesUsedPercent,
						ReadBandwidth:     parent.Host.Disk.ReadBandwidth,
						WriteBandwidth:    parent.Host.Disk.WriteBandwidth,
					},
					Build: &commonv2.Build{
						GitVersion:  parent.Host.Build.GitVersion,
						GitCommit:   &parent.Host.Build.GitCommit,
						GoVersion:   &parent.Host.Build.GoVersion,
						RustVersion: &parent.Host.Build.RustVersion,
						Platform:    &parent.Host.Build.Platform,
					},
					SchedulerClusterId: parent.Host.SchedulerClusterID,
					DisableShared:      parent.Host.DisableShared,
				},
				CreatedAt: timestamppb.New(parent.CreatedAt),
				UpdatedAt: timestamppb.New(parent.UpdatedAt),
			})

		}

		// Update metadata of the persistent cache task.
		if err := v.persistentCacheResource.TaskManager().Store(ctx, peer.Task); err != nil {
			return status.Error(codes.Internal, err.Error())
		}

		// Update metadata of the persistent cache peer.
		if err := v.persistentCacheResource.PeerManager().Store(ctx, peer); err != nil {
			return status.Error(codes.Internal, err.Error())
		}

		if err := stream.Send(&schedulerv2.AnnouncePersistentCachePeerResponse{
			Response: &schedulerv2.AnnouncePersistentCachePeerResponse_NormalPersistentCacheTaskResponse{
				NormalPersistentCacheTaskResponse: &schedulerv2.NormalPersistentCacheTaskResponse{
					CandidateParents: candidateParents,
				},
			},
		}); err != nil {
			peer.Log.Error(err)
			return status.Error(codes.Internal, err.Error())
		}

		return nil
	default:
		return status.Errorf(codes.FailedPrecondition, "invalid size cope %#v", sizeScope)
	}
}

// handleDownloadPersistentCachePeerStartedRequest handles DownloadPersistentCachePeerStartedRequest of AnnouncePersistentCachePeerRequest.
func (v *V2) handleDownloadPersistentCachePeerStartedRequest(ctx context.Context, peerID string) error {
	peer, loaded := v.persistentCacheResource.PeerManager().Load(ctx, peerID)
	if !loaded {
		return status.Errorf(codes.NotFound, "peer %s not found", peerID)
	}

	// Collect DownloadPersistentCachePeerStartedCount metrics.
	metrics.DownloadPersistentCachePeerStartedCount.WithLabelValues(peer.Host.Type.Name()).Inc()

	// Handle peer with peer started request.
	if err := peer.FSM.Event(ctx, persistentcache.PeerEventDownload); err != nil {
		// Collect DownloadPersistentCachePeerStartedFailureCount metrics.
		metrics.DownloadPersistentCachePeerStartedFailureCount.WithLabelValues(peer.Host.Type.Name()).Inc()
		return status.Error(codes.Internal, err.Error())
	}

	// Update metadata of the persistent cache peer.
	if err := v.persistentCacheResource.PeerManager().Store(ctx, peer); err != nil {
		// Collect DownloadPersistentCachePeerStartedFailureCount metrics.
		metrics.DownloadPersistentCachePeerStartedFailureCount.WithLabelValues(peer.Host.Type.Name()).Inc()
		return status.Error(codes.Internal, err.Error())
	}

	return nil
}

// handleReschedulePersistentCachePeerRequest handles ReschedulePersistentCachePeerRequest of AnnouncePersistentCachePeerRequest.
func (v *V2) handleReschedulePersistentCachePeerRequest(ctx context.Context, stream schedulerv2.Scheduler_AnnouncePersistentCachePeerServer, taskID, peerID string, req *schedulerv2.ReschedulePersistentCachePeerRequest) error {
	peer, loaded := v.persistentCacheResource.PeerManager().Load(ctx, peerID)
	if !loaded {
		return status.Errorf(codes.NotFound, "peer %s not found", peerID)
	}

	// Scheduling parent for the peer.
	blocklist := set.NewSafeSet[string]()
	blocklist.Add(peer.ID)

	// Add block parents to blocklist.
	for _, parent := range peer.BlockParents {
		blocklist.Add(parent)
	}

	// Add block parents to blocklist by client report.
	for _, parent := range req.GetCandidateParents() {
		blocklist.Add(parent.GetId())
	}
	peer.BlockParents = blocklist.Values()

	parents, found := v.scheduling.FindCandidatePersistentCacheParents(ctx, peer, blocklist)
	if !found {
		// Collect RegisterPersistentCachePeerFailureCount metrics.
		metrics.RegisterPersistentCachePeerFailureCount.WithLabelValues(peer.Host.Type.Name()).Inc()
		return status.Error(codes.FailedPrecondition, "no candidate parents found")
	}

	currentPersistentReplicaCount, err := v.persistentCacheResource.TaskManager().LoadCurrentPersistentReplicaCount(ctx, taskID)
	if err != nil {
		// Collect RegisterPersistentCachePeerFailureCount metrics.
		metrics.RegisterPersistentCachePeerFailureCount.WithLabelValues(peer.Host.Type.Name()).Inc()
		return status.Error(codes.Internal, err.Error())
	}

	currentReplicaCount, err := v.persistentCacheResource.TaskManager().LoadCurrentReplicaCount(ctx, taskID)
	if err != nil {
		// Collect RegisterPersistentCachePeerFailureCount metrics.
		metrics.RegisterPersistentCachePeerFailureCount.WithLabelValues(peer.Host.Type.Name()).Inc()
		return status.Error(codes.Internal, err.Error())
	}

	candidateParents := make([]*commonv2.PersistentCachePeer, 0, len(parents))
	for _, parent := range parents {
		candidateParents = append(candidateParents, &commonv2.PersistentCachePeer{
			Id:         parent.ID,
			Persistent: parent.Persistent,
			Cost:       durationpb.New(parent.Cost),
			State:      parent.FSM.Current(),
			Task: &commonv2.PersistentCacheTask{
				Id:                            parent.Task.ID,
				PersistentReplicaCount:        parent.Task.PersistentReplicaCount,
				CurrentPersistentReplicaCount: currentPersistentReplicaCount,
				CurrentReplicaCount:           currentReplicaCount,
				Tag:                           &parent.Task.Tag,
				Application:                   &parent.Task.Application,
				PieceLength:                   parent.Task.PieceLength,
				ContentLength:                 parent.Task.ContentLength,
				PieceCount:                    parent.Task.TotalPieceCount,
				State:                         parent.Task.FSM.Current(),
				Ttl:                           durationpb.New(parent.Task.TTL),
				CreatedAt:                     timestamppb.New(parent.Task.CreatedAt),
				UpdatedAt:                     timestamppb.New(parent.Task.UpdatedAt),
			},
			Host: &commonv2.Host{
				Id:              parent.Host.ID,
				Type:            uint32(parent.Host.Type),
				Hostname:        parent.Host.Hostname,
				Ip:              parent.Host.IP,
				Port:            parent.Host.Port,
				DownloadPort:    parent.Host.DownloadPort,
				Os:              parent.Host.OS,
				Platform:        parent.Host.Platform,
				PlatformFamily:  parent.Host.PlatformFamily,
				PlatformVersion: parent.Host.PlatformVersion,
				KernelVersion:   parent.Host.KernelVersion,
				Cpu: &commonv2.CPU{
					LogicalCount:   parent.Host.CPU.LogicalCount,
					PhysicalCount:  parent.Host.CPU.PhysicalCount,
					Percent:        parent.Host.CPU.Percent,
					ProcessPercent: parent.Host.CPU.ProcessPercent,
					Times: &commonv2.CPUTimes{
						User:      parent.Host.CPU.Times.User,
						System:    parent.Host.CPU.Times.System,
						Idle:      parent.Host.CPU.Times.Idle,
						Nice:      parent.Host.CPU.Times.Nice,
						Iowait:    parent.Host.CPU.Times.Iowait,
						Irq:       parent.Host.CPU.Times.Irq,
						Softirq:   parent.Host.CPU.Times.Softirq,
						Steal:     parent.Host.CPU.Times.Steal,
						Guest:     parent.Host.CPU.Times.Guest,
						GuestNice: parent.Host.CPU.Times.GuestNice,
					},
				},
				Memory: &commonv2.Memory{
					Total:              parent.Host.Memory.Total,
					Available:          parent.Host.Memory.Available,
					Used:               parent.Host.Memory.Used,
					UsedPercent:        parent.Host.Memory.UsedPercent,
					ProcessUsedPercent: parent.Host.Memory.ProcessUsedPercent,
					Free:               parent.Host.Memory.Free,
				},
				Network: &commonv2.Network{
					TcpConnectionCount:       parent.Host.Network.TCPConnectionCount,
					UploadTcpConnectionCount: parent.Host.Network.UploadTCPConnectionCount,
					Location:                 &parent.Host.Network.Location,
					Idc:                      &parent.Host.Network.IDC,
					DownloadRate:             parent.Host.Network.DownloadRate,
					DownloadRateLimit:        parent.Host.Network.DownloadRateLimit,
					UploadRate:               parent.Host.Network.UploadRate,
					UploadRateLimit:          parent.Host.Network.UploadRateLimit,
				},
				Disk: &commonv2.Disk{
					Total:             parent.Host.Disk.Total,
					Free:              parent.Host.Disk.Free,
					Used:              parent.Host.Disk.Used,
					UsedPercent:       parent.Host.Disk.UsedPercent,
					InodesTotal:       parent.Host.Disk.InodesTotal,
					InodesUsed:        parent.Host.Disk.InodesUsed,
					InodesFree:        parent.Host.Disk.InodesFree,
					InodesUsedPercent: parent.Host.Disk.InodesUsedPercent,
					ReadBandwidth:     parent.Host.Disk.ReadBandwidth,
					WriteBandwidth:    parent.Host.Disk.WriteBandwidth,
				},
				Build: &commonv2.Build{
					GitVersion:  parent.Host.Build.GitVersion,
					GitCommit:   &parent.Host.Build.GitCommit,
					GoVersion:   &parent.Host.Build.GoVersion,
					RustVersion: &parent.Host.Build.RustVersion,
					Platform:    &parent.Host.Build.Platform,
				},
				SchedulerClusterId: parent.Host.SchedulerClusterID,
				DisableShared:      parent.Host.DisableShared,
			},
			CreatedAt: timestamppb.New(parent.CreatedAt),
			UpdatedAt: timestamppb.New(parent.UpdatedAt),
		})
	}

	// Update metadata of the persistent cache peer.
	if err := v.persistentCacheResource.PeerManager().Store(ctx, peer); err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	if err := stream.Send(&schedulerv2.AnnouncePersistentCachePeerResponse{
		Response: &schedulerv2.AnnouncePersistentCachePeerResponse_NormalPersistentCacheTaskResponse{
			NormalPersistentCacheTaskResponse: &schedulerv2.NormalPersistentCacheTaskResponse{
				CandidateParents: candidateParents,
			},
		},
	}); err != nil {
		peer.Log.Error(err)
		return status.Error(codes.Internal, err.Error())
	}

	return nil
}

// handleDownloadPersistentCachePeerFinishedRequest handles DownloadPersistentCachePeerFinishedRequest of AnnouncePersistentCachePeerRequest.
func (v *V2) handleDownloadPersistentCachePeerFinishedRequest(ctx context.Context, peerID string) error {
	peer, loaded := v.persistentCacheResource.PeerManager().Load(ctx, peerID)
	if !loaded {
		return status.Errorf(codes.NotFound, "peer %s not found", peerID)
	}

	// Handle peer with peer finished request.
	peer.Cost = time.Since(peer.CreatedAt)
	if err := peer.FSM.Event(ctx, persistentcache.PeerEventSucceeded); err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	// Update metadata of the persistent cache peer.
	if err := v.persistentCacheResource.PeerManager().Store(ctx, peer); err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	// Collect DownloadPersistentCachePeerCount metrics.
	metrics.DownloadPersistentCachePeerCount.WithLabelValues(peer.Host.Type.Name()).Inc()
	return nil
}

// handleDownloadPersistentCachePeerFailedRequest handles DownloadPersistentCachePeerFailedRequest of AnnouncePersistentCachePeerRequest.
func (v *V2) handleDownloadPersistentCachePeerFailedRequest(ctx context.Context, peerID string) error {
	peer, loaded := v.persistentCacheResource.PeerManager().Load(ctx, peerID)
	if !loaded {
		return status.Errorf(codes.NotFound, "peer %s not found", peerID)
	}

	// Handle peer with peer failed request.
	if err := peer.FSM.Event(ctx, persistentcache.PeerEventFailed); err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	// Update metadata of the persistent cache peer.
	if err := v.persistentCacheResource.PeerManager().Store(ctx, peer); err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	// Collect DownloadPersistentCachePeerCount and DownloadPersistentCachePeerFailureCount metrics.
	metrics.DownloadPersistentCachePeerCount.WithLabelValues(peer.Host.Type.Name()).Inc()
	metrics.DownloadPersistentCachePeerFailureCount.WithLabelValues(peer.Host.Type.Name()).Inc()
	return nil
}

// handleDownloadPersistentCachePieceFinishedRequest handles DownloadPersistentCachePieceFinishedRequest of AnnouncePersistentCachePeerRequest.
func (v *V2) handleDownloadPersistentCachePieceFinishedRequest(ctx context.Context, peerID string, req *schedulerv2.DownloadPieceFinishedRequest) error {
	peer, loaded := v.persistentCacheResource.PeerManager().Load(ctx, peerID)
	if !loaded {
		return status.Errorf(codes.NotFound, "peer %s not found", peerID)
	}

	piece := req.GetPiece()
	peer.FinishedPieces.Set(uint(piece.GetNumber()))
	peer.UpdatedAt = time.Now()

	parent, loadedParent := v.persistentCacheResource.PeerManager().Load(ctx, piece.GetParentId())
	if loadedParent {
		parent.UpdatedAt = time.Now()
		parent.Host.UpdatedAt = time.Now()
	}

	// Update metadata of the persistent cache peer.
	if err := v.persistentCacheResource.PeerManager().Store(ctx, peer); err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	// Update metadata of the persistent cache peer's parent.
	if err := v.persistentCacheResource.PeerManager().Store(ctx, parent); err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	// Collect DownloadPersistentCachePieceCount metrics.
	metrics.DownloadPersistentCachePieceCount.WithLabelValues(peer.Host.Type.Name()).Inc()
	return nil
}

// handleDownloadPersistentCachePieceFailedRequest handles DownloadPersistentCachePieceFailedRequest of AnnouncePersistentCachePeerRequest.
func (v *V2) handleDownloadPersistentCachePieceFailedRequest(ctx context.Context, peerID string, req *schedulerv2.DownloadPieceFailedRequest) error {
	peer, loaded := v.persistentCacheResource.PeerManager().Load(ctx, peerID)
	if !loaded {
		return status.Errorf(codes.NotFound, "peer %s not found", peerID)
	}

	// Collect DownloadPersistentCachePieceCount and DownloadPersistentCachePieceFailureCount metrics.
	metrics.DownloadPersistentCachePieceCount.WithLabelValues(peer.Host.Type.Name()).Inc()
	metrics.DownloadPersistentCachePieceFailureCount.WithLabelValues(peer.Host.Type.Name()).Inc()

	if req.Temporary {
		// Handle peer with piece temporary failed request.
		peer.UpdatedAt = time.Now()
		blocklist := set.NewSafeSet[string]()
		blocklist.Add(peer.ID)

		// Add block parents to blocklist.
		for _, parent := range peer.BlockParents {
			blocklist.Add(parent)
		}

		// Add parent to blocklist, because download piece from parent failed.
		blocklist.Add(req.GetParentId())

		peer.BlockParents = blocklist.Values()
		if parent, loaded := v.resource.PeerManager().Load(req.GetParentId()); loaded {
			parent.Host.UploadFailedCount.Inc()
		}

		// Update metadata of the persistent cache peer.
		if err := v.persistentCacheResource.PeerManager().Store(ctx, peer); err != nil {
			return status.Error(codes.Internal, err.Error())
		}

		return nil
	}

	return status.Error(codes.FailedPrecondition, "download piece failed")
}

// StatPersistentCachePeer checks information of persistent cache peer.
func (v *V2) StatPersistentCachePeer(ctx context.Context, req *schedulerv2.StatPersistentCachePeerRequest) (*commonv2.PersistentCachePeer, error) {
	if v.persistentCacheResource == nil {
		return nil, status.Error(codes.FailedPrecondition, "redis is not enabled")
	}

	log := logger.WithPeer(req.HostId, req.TaskId, req.PeerId)
	log.Info("stat persistent cache peer")

	peer, loaded := v.persistentCacheResource.PeerManager().Load(ctx, req.GetPeerId())
	if !loaded {
		log.Errorf("persistent cache peer %s not found", req.GetPeerId())
		return nil, status.Errorf(codes.NotFound, "persistent cache peer %s not found", req.GetPeerId())
	}

	currentPersistentReplicaCount, err := v.persistentCacheResource.TaskManager().LoadCurrentPersistentReplicaCount(ctx, peer.Task.ID)
	if err != nil {
		log.Errorf("load current persistent replica count failed %s", err.Error())
		return nil, status.Error(codes.Internal, err.Error())
	}

	currentReplicaCount, err := v.persistentCacheResource.TaskManager().LoadCurrentReplicaCount(ctx, peer.Task.ID)
	if err != nil {
		log.Errorf("load current replica count failed %s", err.Error())
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &commonv2.PersistentCachePeer{
		Id:         peer.ID,
		Persistent: peer.Persistent,
		State:      peer.FSM.Current(),
		Cost:       durationpb.New(peer.Cost),
		CreatedAt:  timestamppb.New(peer.CreatedAt),
		UpdatedAt:  timestamppb.New(peer.UpdatedAt),
		Task: &commonv2.PersistentCacheTask{
			Id:                            peer.Task.ID,
			PersistentReplicaCount:        peer.Task.PersistentReplicaCount,
			CurrentPersistentReplicaCount: currentPersistentReplicaCount,
			CurrentReplicaCount:           currentReplicaCount,
			Tag:                           &peer.Task.Tag,
			Application:                   &peer.Task.Application,
			PieceLength:                   peer.Task.PieceLength,
			ContentLength:                 peer.Task.ContentLength,
			PieceCount:                    uint32(peer.Task.TotalPieceCount),
			State:                         peer.Task.FSM.Current(),
			Ttl:                           durationpb.New(peer.Task.TTL),
			CreatedAt:                     timestamppb.New(peer.Task.CreatedAt),
			UpdatedAt:                     timestamppb.New(peer.Task.UpdatedAt),
		},
		Host: &commonv2.Host{
			Id:              peer.Host.ID,
			Type:            uint32(peer.Host.Type),
			Hostname:        peer.Host.Hostname,
			Ip:              peer.Host.IP,
			Port:            peer.Host.Port,
			DownloadPort:    peer.Host.DownloadPort,
			Os:              peer.Host.OS,
			Platform:        peer.Host.Platform,
			PlatformFamily:  peer.Host.PlatformFamily,
			PlatformVersion: peer.Host.PlatformVersion,
			KernelVersion:   peer.Host.KernelVersion,
			Cpu: &commonv2.CPU{
				LogicalCount:   peer.Host.CPU.LogicalCount,
				PhysicalCount:  peer.Host.CPU.PhysicalCount,
				Percent:        peer.Host.CPU.Percent,
				ProcessPercent: peer.Host.CPU.ProcessPercent,
				Times: &commonv2.CPUTimes{
					User:      peer.Host.CPU.Times.User,
					System:    peer.Host.CPU.Times.System,
					Idle:      peer.Host.CPU.Times.Idle,
					Nice:      peer.Host.CPU.Times.Nice,
					Iowait:    peer.Host.CPU.Times.Iowait,
					Irq:       peer.Host.CPU.Times.Irq,
					Softirq:   peer.Host.CPU.Times.Softirq,
					Steal:     peer.Host.CPU.Times.Steal,
					Guest:     peer.Host.CPU.Times.Guest,
					GuestNice: peer.Host.CPU.Times.GuestNice,
				},
			},
			Memory: &commonv2.Memory{
				Total:              peer.Host.Memory.Total,
				Available:          peer.Host.Memory.Available,
				Used:               peer.Host.Memory.Used,
				UsedPercent:        peer.Host.Memory.UsedPercent,
				ProcessUsedPercent: peer.Host.Memory.ProcessUsedPercent,
				Free:               peer.Host.Memory.Free,
			},
			Network: &commonv2.Network{
				TcpConnectionCount:       peer.Host.Network.TCPConnectionCount,
				UploadTcpConnectionCount: peer.Host.Network.UploadTCPConnectionCount,
				Location:                 &peer.Host.Network.Location,
				Idc:                      &peer.Host.Network.IDC,
				DownloadRate:             peer.Host.Network.DownloadRate,
				DownloadRateLimit:        peer.Host.Network.DownloadRateLimit,
				UploadRate:               peer.Host.Network.UploadRate,
				UploadRateLimit:          peer.Host.Network.UploadRateLimit,
			},
			Disk: &commonv2.Disk{
				Total:             peer.Host.Disk.Total,
				Free:              peer.Host.Disk.Free,
				Used:              peer.Host.Disk.Used,
				UsedPercent:       peer.Host.Disk.UsedPercent,
				InodesTotal:       peer.Host.Disk.InodesTotal,
				InodesUsed:        peer.Host.Disk.InodesUsed,
				InodesFree:        peer.Host.Disk.InodesFree,
				InodesUsedPercent: peer.Host.Disk.InodesUsedPercent,
				WriteBandwidth:    peer.Host.Disk.WriteBandwidth,
				ReadBandwidth:     peer.Host.Disk.ReadBandwidth,
			},
			Build: &commonv2.Build{
				GitVersion:  peer.Host.Build.GitVersion,
				GitCommit:   &peer.Host.Build.GitCommit,
				GoVersion:   &peer.Host.Build.GoVersion,
				RustVersion: &peer.Host.Build.RustVersion,
				Platform:    &peer.Host.Build.Platform,
			},
			SchedulerClusterId: uint64(v.config.Manager.SchedulerClusterID),
			DisableShared:      peer.Host.DisableShared,
		},
	}, nil
}

// DeletePersistentCachePeer releases persistent cache peer in scheduler.
func (v *V2) DeletePersistentCachePeer(_ctx context.Context, req *schedulerv2.DeletePersistentCachePeerRequest) error {
	if v.persistentCacheResource == nil {
		return status.Error(codes.FailedPrecondition, "redis is not enabled")
	}

	// Context use background to avoid the context canceled by the client and
	// the peer deletion operation is not completed.
	ctx := context.Background()
	log := logger.WithPeer(req.GetHostId(), req.GetTaskId(), req.GetPeerId())
	log.Info("delete persistent cache peer")

	peer, found := v.persistentCacheResource.PeerManager().Load(ctx, req.GetPeerId())
	if !found {
		log.Errorf("persistent cache peer %s not found", req.GetPeerId())
		return status.Errorf(codes.NotFound, "persistent cache peer %s not found", req.GetPeerId())
	}

	if err := v.persistentCacheResource.PeerManager().Delete(ctx, req.GetPeerId()); err != nil {
		log.Errorf("delete persistent cache peer %s error %s", req.GetPeerId(), err)
		return status.Error(codes.Internal, err.Error())
	}

	// Delete the persistent cache task from the peer, if delete failed, skip it.
	addr := fmt.Sprintf("%s:%d", peer.Host.IP, peer.Host.DownloadPort)
	dialOptions := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	dfdaemonClient, err := dfdaemonclient.GetV2ByAddr(ctx, addr, dialOptions...)
	if err != nil {
		peer.Log.Errorf("get dfdaemon client failed %s", err)
		return err
	}

	advertiseIP := v.config.Server.AdvertiseIP.String()
	if err := dfdaemonClient.DeletePersistentCacheTask(ctx, &dfdaemonv2.DeletePersistentCacheTaskRequest{TaskId: peer.Task.ID, RemoteIp: &advertiseIP}); err != nil {
		peer.Log.Errorf("delete persistent cache task %s from peer %s failed %s", peer.Task.ID, peer.ID, err)
	}

	// Select the remote peer to copy the replica and trigger the download task with asynchronous.
	blocklist := set.NewSafeSet[string]()
	blocklist.Add(req.GetHostId())
	go func(peer *persistentcache.Peer, blocklist set.SafeSet[string]) {
		log.Infof("replicate persistent cache task %s", peer.Task.ID)
		if err := v.replicatePersistentCacheTask(context.Background(), peer, blocklist); err != nil {
			log.Errorf("replicate persistent cache task failed %s", err)
		}

		log.Infof("replicate persistent cache task %s finished", peer.Task.ID)
	}(peer, blocklist)

	return nil
}

// UploadPersistentCacheTaskStarted uploads the metadata of the persistent cache task started.
func (v *V2) UploadPersistentCacheTaskStarted(ctx context.Context, req *schedulerv2.UploadPersistentCacheTaskStartedRequest) error {
	if v.persistentCacheResource == nil {
		return status.Error(codes.FailedPrecondition, "redis is not enabled")
	}

	log := logger.WithPeer(req.GetHostId(), req.GetTaskId(), req.GetPeerId())
	log.Info("upload persistent cache task started")

	host, loaded := v.persistentCacheResource.HostManager().Load(ctx, req.GetHostId())
	if !loaded {
		log.Error("host not found")
		return status.Errorf(codes.NotFound, "host %s not found", req.GetHostId())
	}

	// Handle task with task started request, new task and store it.
	task, loaded := v.persistentCacheResource.TaskManager().Load(ctx, req.GetTaskId())
	if loaded && !task.FSM.Can(persistentcache.TaskEventUpload) {
		log.Errorf("persistent cache task %s is %s cannot upload", task.ID, task.FSM.Current())
		return status.Errorf(codes.AlreadyExists, "persistent cache task %s is %s cannot upload", task.ID, task.FSM.Current())
	}

	task = persistentcache.NewTask(req.GetTaskId(), req.GetTag(), req.GetApplication(), persistentcache.TaskStatePending, req.GetPersistentReplicaCount(),
		req.GetPieceLength(), req.GetContentLength(), req.GetPieceCount(), req.GetTtl().AsDuration(), time.Now(), time.Now(), log)

	if err := task.FSM.Event(ctx, persistentcache.TaskEventUpload); err != nil {
		log.Errorf("task fsm event failed: %s", err.Error())
		return status.Error(codes.Internal, err.Error())
	}

	if err := v.persistentCacheResource.TaskManager().Store(ctx, task); err != nil {
		log.Errorf("store persistent cache task %s error %s", task.ID, err)
		return status.Error(codes.Internal, err.Error())
	}

	// Handle peer with task started request, new peer and store it.
	if peer, loaded := v.persistentCacheResource.PeerManager().Load(ctx, req.GetPeerId()); loaded {
		log.Error("persistent cache peer already exists")
		return status.Errorf(codes.AlreadyExists, "persistent cache peer %s already exists", peer.ID)
	}

	peer := persistentcache.NewPeer(req.GetPeerId(), persistentcache.PeerStatePending, true, bitset.New(uint(req.GetPieceCount())), nil, task, host, 0, time.Now(), time.Now(), log)

	if err := peer.FSM.Event(ctx, persistentcache.PeerEventUpload); err != nil {
		log.Errorf("peer fsm event failed: %s", err.Error())
		return status.Error(codes.Internal, err.Error())
	}

	if err := v.persistentCacheResource.PeerManager().Store(ctx, peer); err != nil {
		log.Errorf("store persistent cache peer %s error %s", peer.ID, err)
		return status.Error(codes.Internal, err.Error())
	}

	return nil
}

// UploadPersistentCacheTaskFinished uploads the metadata of the persistent cache task finished.
func (v *V2) UploadPersistentCacheTaskFinished(ctx context.Context, req *schedulerv2.UploadPersistentCacheTaskFinishedRequest) (*commonv2.PersistentCacheTask, error) {
	if v.persistentCacheResource == nil {
		return nil, status.Error(codes.FailedPrecondition, "redis is not enabled")
	}

	log := logger.WithPeer(req.GetHostId(), req.GetTaskId(), req.GetPeerId())
	log.Info("upload persistent cache task finished")

	// Handle peer with task finished request, load peer and update it.
	peer, loaded := v.persistentCacheResource.PeerManager().Load(ctx, req.GetPeerId())
	if !loaded {
		log.Error("persistent cache peer not found")
		return nil, status.Errorf(codes.NotFound, "persistent cache peer %s not found", req.GetPeerId())
	}

	peer.FinishedPieces.SetAll()
	if err := peer.FSM.Event(ctx, persistentcache.PeerEventSucceeded); err != nil {
		log.Errorf("peer fsm event failed: %s", err.Error())
		return nil, status.Error(codes.Internal, err.Error())
	}
	peer.Cost = time.Since(peer.CreatedAt)
	peer.UpdatedAt = time.Now()

	if err := v.persistentCacheResource.PeerManager().Store(ctx, peer); err != nil {
		log.Errorf("store persistent cache peer %s error %s", peer.ID, err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	// Handle task with peer finished request, load task and update it.
	if err := peer.Task.FSM.Event(ctx, persistentcache.TaskEventSucceeded); err != nil {
		log.Errorf("task fsm event failed: %s", err.Error())
		return nil, status.Error(codes.Internal, err.Error())
	}
	peer.Task.UpdatedAt = time.Now()

	if err := v.persistentCacheResource.TaskManager().Store(ctx, peer.Task); err != nil {
		log.Errorf("store persistent cache task %s error %s", peer.Task.ID, err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	currentPersistentReplicaCount, err := v.persistentCacheResource.TaskManager().LoadCurrentPersistentReplicaCount(ctx, peer.Task.ID)
	if err != nil {
		log.Errorf("load current persistent replica count failed %s", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	currentReplicaCount, err := v.persistentCacheResource.TaskManager().LoadCurrentReplicaCount(ctx, peer.Task.ID)
	if err != nil {
		log.Errorf("load current replica count failed %s", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	persistentCacheTask := &commonv2.PersistentCacheTask{
		Id:                            peer.Task.ID,
		PersistentReplicaCount:        peer.Task.PersistentReplicaCount,
		CurrentPersistentReplicaCount: currentPersistentReplicaCount,
		CurrentReplicaCount:           currentReplicaCount,
		Tag:                           &peer.Task.Tag,
		Application:                   &peer.Task.Application,
		PieceLength:                   peer.Task.PieceLength,
		ContentLength:                 peer.Task.ContentLength,
		PieceCount:                    peer.Task.TotalPieceCount,
		State:                         peer.Task.FSM.Current(),
		Ttl:                           durationpb.New(peer.Task.TTL),
		CreatedAt:                     timestamppb.New(peer.Task.CreatedAt),
		UpdatedAt:                     timestamppb.New(peer.Task.UpdatedAt),
	}

	// Select the remote peer to copy the replica and trigger the download task with asynchronous.
	blocklist := set.NewSafeSet[string]()
	blocklist.Add(peer.Host.ID)
	go func(peer *persistentcache.Peer, blocklist set.SafeSet[string]) {
		log.Infof("replicate persistent cache task %s", peer.Task.ID)
		if err := v.replicatePersistentCacheTask(context.Background(), peer, blocklist); err != nil {
			log.Errorf("replicate persistent cache task failed %s", err)
		}

		log.Infof("replicate persistent cache task %s finished", peer.Task.ID)
	}(peer, blocklist)

	return persistentCacheTask, nil
}

// replicatePersistentCacheTask replicates the persistent cache task to the remote peer.
func (v *V2) replicatePersistentCacheTask(ctx context.Context, peer *persistentcache.Peer, blocklist set.SafeSet[string]) error {
	cachedParents, hosts, found := v.scheduling.FindReplicatePersistentCacheHosts(ctx, peer.Task, blocklist)
	if !found {
		peer.Log.Warn("no replicate hosts found")
		return nil
	}

	// Replicate the persistent cache task to the cached parent, set the persistent flag to true
	// because the cached parent has the temporary cache.
	for _, cachedParent := range cachedParents {
		go func(*persistentcache.Task, *persistentcache.Peer) {
			peer.Log.Infof("replicate to cached parent %s", cachedParent.ID)
			if err := v.persistPersistentCacheTaskByPeer(context.Background(), peer, cachedParent); err != nil {
				peer.Log.Errorf("replicate to cached parent %s failed %s", cachedParent.ID, err)
			}

			peer.Log.Infof("replicate to cached parent %s finished", cachedParent.ID)
		}(peer.Task, cachedParent)
	}

	// Replicate the persistent cache task to the host, trigger the download task from the other peer,
	// because the host has no cache.
	for _, host := range hosts {
		go func(*persistentcache.Task, *persistentcache.Host) {
			peer.Log.Infof("replicate to host %s", host.ID)
			if err := v.downloadPersistentCacheTaskByPeer(context.Background(), peer.Task, host); err != nil {
				peer.Log.Errorf("replicate to host %s failed %s", host.ID, err)
			}

			peer.Log.Infof("replicate to host %s finished", host.ID)
		}(peer.Task, host)
	}

	return nil
}

// downloadPersistentCacheTaskByPeer downloads the persistent cache task by peer.
func (v *V2) downloadPersistentCacheTaskByPeer(ctx context.Context, task *persistentcache.Task, host *persistentcache.Host) error {
	addr := fmt.Sprintf("%s:%d", host.IP, host.DownloadPort)
	dialOptions := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	dfdaemonClient, err := dfdaemonclient.GetV2ByAddr(ctx, addr, dialOptions...)
	if err != nil {
		task.Log.Errorf("get dfdaemon client failed %s", err)
		return err
	}

	advertiseIP := v.config.Server.AdvertiseIP.String()
	stream, err := dfdaemonClient.DownloadPersistentCacheTask(ctx, &dfdaemonv2.DownloadPersistentCacheTaskRequest{
		TaskId:      task.ID,
		Persistent:  true,
		Tag:         &task.Tag,
		Application: &task.Application,
		RemoteIp:    &advertiseIP,
	})
	if err != nil {
		task.Log.Errorf("download persistent cache task failed %s", err)
		return err
	}

	// Wait for the download persistent cache task to complete.
	for {
		_, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				task.Log.Info("download persistent cache task finished")
				return nil
			}

			task.Log.Errorf("download persistent cache task failed %s", err)
			return err
		}
	}
}

// persistPersistentCacheTaskByPeer persists the persistent cache task by peer.
func (v *V2) persistPersistentCacheTaskByPeer(ctx context.Context, peer *persistentcache.Peer, cachedParent *persistentcache.Peer) error {
	addr := fmt.Sprintf("%s:%d", cachedParent.Host.IP, cachedParent.Host.DownloadPort)
	dialOptions := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	dfdaemonClient, err := dfdaemonclient.GetV2ByAddr(ctx, addr, dialOptions...)
	if err != nil {
		peer.Log.Errorf("get dfdaemon client failed %s", err)
		return err
	}

	advertiseIP := v.config.Server.AdvertiseIP.String()
	if err := dfdaemonClient.UpdatePersistentCacheTask(ctx, &dfdaemonv2.UpdatePersistentCacheTaskRequest{
		TaskId:     peer.Task.ID,
		Persistent: true,
		RemoteIp:   &advertiseIP,
	}); err != nil {
		peer.Log.Errorf("update persistent cache task failed %s", err)
		return err
	}

	cachedParent.Persistent = true
	if err := v.persistentCacheResource.PeerManager().Store(ctx, cachedParent); err != nil {
		peer.Log.Errorf("store persistent cache peer %s error %s", cachedParent.ID, err)
		return err
	}

	return nil
}

// UploadPersistentCacheTaskFailed uploads the metadata of the persistent cache task failed.
func (v *V2) UploadPersistentCacheTaskFailed(ctx context.Context, req *schedulerv2.UploadPersistentCacheTaskFailedRequest) error {
	if v.persistentCacheResource == nil {
		return status.Error(codes.FailedPrecondition, "redis is not enabled")
	}

	log := logger.WithPeer(req.GetHostId(), req.GetTaskId(), req.GetPeerId())
	log.Info("upload persistent cache task failed")

	// Handle peer with task failed request, load peer and update it.
	peer, loaded := v.persistentCacheResource.PeerManager().Load(ctx, req.GetPeerId())
	if !loaded {
		log.Error("persistent cache peer not found")
		return status.Errorf(codes.NotFound, "persistent cache peer %s not found", req.GetPeerId())
	}

	if err := peer.FSM.Event(ctx, persistentcache.PeerEventFailed); err != nil {
		log.Errorf("peer fsm event failed: %s", err.Error())
		return status.Error(codes.Internal, err.Error())
	}
	peer.UpdatedAt = time.Now()

	if err := v.persistentCacheResource.PeerManager().Store(ctx, peer); err != nil {
		log.Errorf("store persistent cache peer %s error %s", peer.ID, err)
		return status.Error(codes.Internal, err.Error())
	}

	// Handle task with peer failed request, load task and update it.
	if err := peer.Task.FSM.Event(ctx, persistentcache.TaskEventSucceeded); err != nil {
		log.Errorf("task fsm event failed: %s", err.Error())
		return status.Error(codes.Internal, err.Error())
	}
	peer.Task.UpdatedAt = time.Now()

	if err := v.persistentCacheResource.TaskManager().Store(ctx, peer.Task); err != nil {
		log.Errorf("store persistent cache task %s error %s", peer.Task.ID, err)
		return status.Error(codes.Internal, err.Error())
	}

	return nil
}

// StatPersistentCacheTask checks information of persistent cache task.
func (v *V2) StatPersistentCacheTask(ctx context.Context, req *schedulerv2.StatPersistentCacheTaskRequest) (*commonv2.PersistentCacheTask, error) {
	if v.persistentCacheResource == nil {
		return nil, status.Error(codes.FailedPrecondition, "redis is not enabled")
	}

	log := logger.WithHostAndTaskID(req.GetHostId(), req.GetTaskId())
	log.Info("stat persistent cache task")

	task, loaded := v.persistentCacheResource.TaskManager().Load(ctx, req.GetTaskId())
	if !loaded {
		log.Errorf("persistent cache task %s not found", req.GetTaskId())
		return nil, status.Errorf(codes.NotFound, "persistent cache task %s not found", req.GetTaskId())
	}

	currentPersistentReplicaCount, err := v.persistentCacheResource.TaskManager().LoadCurrentPersistentReplicaCount(ctx, task.ID)
	if err != nil {
		log.Errorf("load current persistent replica count failed %s", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	currentReplicaCount, err := v.persistentCacheResource.TaskManager().LoadCurrentReplicaCount(ctx, task.ID)
	if err != nil {
		log.Errorf("load current replica count failed %s", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &commonv2.PersistentCacheTask{
		Id:                            task.ID,
		PersistentReplicaCount:        task.PersistentReplicaCount,
		CurrentPersistentReplicaCount: currentPersistentReplicaCount,
		CurrentReplicaCount:           currentReplicaCount,
		Tag:                           &task.Tag,
		Application:                   &task.Application,
		PieceLength:                   task.PieceLength,
		ContentLength:                 task.ContentLength,
		PieceCount:                    task.TotalPieceCount,
		State:                         task.FSM.Current(),
		Ttl:                           durationpb.New(task.TTL),
		CreatedAt:                     timestamppb.New(task.CreatedAt),
		UpdatedAt:                     timestamppb.New(task.UpdatedAt),
	}, nil
}

// DeletePersistentCacheTask releases persistent cache task in scheduler.
func (v *V2) DeletePersistentCacheTask(_ctx context.Context, req *schedulerv2.DeletePersistentCacheTaskRequest) error {
	if v.persistentCacheResource == nil {
		return status.Error(codes.FailedPrecondition, "redis is not enabled")
	}

	// Context use background to avoid the context canceled by the client and
	// the task deletion operation is not completed.
	ctx := context.Background()
	log := logger.WithHostAndTaskID(req.GetHostId(), req.GetTaskId())
	log.Info("delete persistent cache task")

	// Delete the persistent cache peers in the redis and peer.
	peers, err := v.persistentCacheResource.PeerManager().LoadAllByTaskID(ctx, req.GetTaskId())
	if err != nil {
		log.Errorf("load persistent cache peers by task %s error %s", req.GetTaskId(), err)
		return status.Error(codes.Internal, err.Error())
	}

	for _, peer := range peers {
		if err := v.persistentCacheResource.PeerManager().Delete(ctx, peer.ID); err != nil {
			log.Errorf("delete persistent cache peer %s error %s", peer.ID, err)
			continue
		}

		// Delete the persistent cache task from the peer, if delete failed, skip it.
		addr := fmt.Sprintf("%s:%d", peer.Host.IP, peer.Host.DownloadPort)
		dialOptions := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
		dfdaemonClient, err := dfdaemonclient.GetV2ByAddr(ctx, addr, dialOptions...)
		if err != nil {
			peer.Log.Errorf("get dfdaemon client failed %s", err)
			continue
		}

		advertiseIP := v.config.Server.AdvertiseIP.String()
		if err := dfdaemonClient.DeletePersistentCacheTask(ctx, &dfdaemonv2.DeletePersistentCacheTaskRequest{TaskId: peer.Task.ID, RemoteIp: &advertiseIP}); err != nil {
			peer.Log.Errorf("delete persistent cache task %s from peer %s failed %s", peer.Task.ID, peer.ID, err)
			continue
		}
	}

	// Delete the persistent cache task in the redis.
	if err := v.persistentCacheResource.TaskManager().Delete(ctx, req.GetTaskId()); err != nil {
		log.Errorf("delete persistent cache task %s error %s", req.GetTaskId(), err)
		return status.Error(codes.Internal, err.Error())
	}

	return nil
}
