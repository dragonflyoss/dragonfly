package server

import (
	"context"
	"fmt"
	"github.com/dragonflyoss/Dragonfly2/pkg/grpc/base"
	"github.com/dragonflyoss/Dragonfly2/pkg/grpc/scheduler"
	"github.com/dragonflyoss/Dragonfly2/scheduler/service"
	"github.com/dragonflyoss/Dragonfly2/scheduler/service/schedule_worker"
	"github.com/dragonflyoss/Dragonfly2/scheduler/types"
)

var _ scheduler.SchedulerServer = &SchedulerServer{}

type SchedulerServer struct {
	svc *service.SchedulerService
	worker schedule_worker.IWorker
}

func (s *SchedulerServer) RegisterPeerTask(ctx context.Context, request *scheduler.PeerTaskRequest) (pkg *scheduler.PiecePackage, err error) {
	pkg = new(scheduler.PiecePackage)
	defer func(){
		e := recover()
		if e != nil {
			err = fmt.Errorf("%v", e)
			return
		}
		if err != nil {
			pkg.State.Code = base.Code_ERROR
			pkg.State.Msg = err.Error()
			pkg.State.Success = false
			err = nil
		} else {
			pkg.State.Code = base.Code_SUCCESS
			pkg.State.Success = true
		}
		return
	}()

	// get or create task
	taskId := s.svc.GenerateTaskId(request.Url, request.Filter)
	task, _ := s.svc.GetTask(taskId)
	if task == nil {
		task = &types.Task{
			TaskId: taskId,
			Url: request.Url,
			Filter: request.Filter,
			BizId: request.BizId,
			UrlMata: request.UrlMata,
		}
		task, err = s.svc.AddTask(task)
		if err != nil {
			return
		}
	}

	// get or create host
	hostId := request.PeerHost.Uuid
	host, _ := s.svc.GetHost(hostId)
	if host == nil {
		host = &types.Host{
			Uuid:            request.PeerHost.Uuid,
			Ip:              request.PeerHost.Ip,
			Port:            request.PeerHost.Port,
			HostName:        request.PeerHost.HostName,
			SecurityDomain:  request.PeerHost.SecurityDomain,
			Location:        request.PeerHost.Location,
			Idc:             request.PeerHost.Idc,
			Switch:          request.PeerHost.Switch,
		}
		host, err = s.svc.AddHost(host)
		if err != nil {
			return
		}
	}

	// get or creat PeerTask
	pid := request.Pid
	peerTask, _ := s.svc.GetPeerTask(pid)
	if peerTask == nil {
		peerTask, err = s.svc.AddPeerTask(pid, task, host)
		if err != nil {
			return
		}
	}

	// do scheduler piece
	pieceList, err := s.svc.Scheduler(peerTask)
	if err != nil {
		return
	}

	// assemble result
	pkg.TaskId = taskId
	if int(peerTask.GetFinishedNum()) + len(pieceList) >= len(task.PieceList) {
		pkg.Last = true
		pkg.ContentLength = task.ContentLength
	}
	for _, p := range pieceList {
		pkg.PieceTasks = append(pkg.PieceTasks, &scheduler.PiecePackage_PieceTask{
		PieceNum : p.Piece.PieceNum,
			PieceRange : p.Piece.PieceRange,
			PieceMd5   : p.Piece.PieceMd5,
			SrcPid     : p.SrcPid ,
			DstPid     : p.DstPid,
			DstAddr    : p.DstAddr,
			PieceOffset: p.Piece.PieceOffset,
			PieceStyle : p.Piece.PieceStyle,
		})
	}

	return
}



func (s *SchedulerServer) PullPieceTasks(server scheduler.Scheduler_PullPieceTasksServer) (err error) {
	schedule_worker.CreateClient(server, s.worker).Start()
	return
}

func (s *SchedulerServer) ReportPeerResult(ctx context.Context, result *scheduler.PeerResult) (ret *base.ResponseState, err error) {
	ret = new(base.ResponseState)
	defer func(){
		e := recover()
		if e != nil {
			err = fmt.Errorf("%v", e)
			return
		}
		if err != nil {
			ret.Code = base.Code_ERROR
			ret.Msg = err.Error()
			ret.Success = false
			err = nil
		} else {
			ret.Code = base.Code_SUCCESS
			ret.Success = true
		}
		return
	}()

	pid := result.SrcIp
	peerTask, err := s.svc.GetPeerTask(pid)
	if err != nil {
		return
	}
	peerTask.SetStatus(result.Traffic, result.Cost, result.Success, result.ErrorCode)

	return
}

func (s *SchedulerServer) LeaveTask(ctx context.Context, target *scheduler.PeerTarget) (ret *base.ResponseState, err error) {
	ret = new(base.ResponseState)
	defer func(){
		e := recover()
		if e != nil {
			err = fmt.Errorf("%v", e)
			return
		}
		if err != nil {
			ret.Code = base.Code_ERROR
			ret.Msg = err.Error()
			ret.Success = false
			err = nil
		} else {
			ret.Code = base.Code_SUCCESS
			ret.Success = true
		}
		return
	}()

	pid := target.Pid

	err = s.svc.DeletePeerTask(pid)

	return
}
