package peer

import (
	"context"
	"time"

	"d7y.io/dragonfly/v2/client/daemon/storage"
	"d7y.io/dragonfly/v2/pkg/dfcodes"
	logger "d7y.io/dragonfly/v2/pkg/dflog"
	"d7y.io/dragonfly/v2/pkg/rpc/scheduler"
)

type filePeerTaskCallback struct {
	ctx   context.Context
	ptm   *peerTaskManager
	req   *FilePeerTaskRequest
	start time.Time
}

func (p *filePeerTaskCallback) Init(pt PeerTask) error {
	// prepare storage
	err := p.ptm.storageManager.RegisterTask(p.ctx,
		storage.RegisterTaskRequest{
			CommonTaskRequest: storage.CommonTaskRequest{
				PeerID:      pt.GetPeerID(),
				TaskID:      pt.GetTaskID(),
				Destination: p.req.Output,
			},
			ContentLength: pt.GetContentLength(),
			TotalPieces:   pt.GetTotalPieces(),
		})
	if err != nil {
		logger.Errorf("register task to storage manager failed: %s", err)
	}
	return err
}

func (p *filePeerTaskCallback) Update(pt PeerTask) error {
	// update storage
	err := p.ptm.storageManager.UpdateTask(p.ctx,
		&storage.UpdateTaskRequest{
			PeerTaskMetaData: storage.PeerTaskMetaData{
				PeerID: pt.GetPeerID(),
				TaskID: pt.GetTaskID(),
			},
			ContentLength: pt.GetContentLength(),
			TotalPieces:   pt.GetTotalPieces(),
		})
	if err != nil {
		logger.Errorf("update task to storage manager failed: %s", err)
	}
	return err
}

func (p *filePeerTaskCallback) Done(pt PeerTask) error {
	e := p.ptm.storageManager.Store(
		context.Background(),
		&storage.StoreRequest{
			CommonTaskRequest: storage.CommonTaskRequest{
				PeerID:      pt.GetPeerID(),
				TaskID:      pt.GetTaskID(),
				Destination: p.req.Output,
			},
			MetadataOnly: false,
			TotalPieces:  pt.GetTotalPieces(),
		})
	if e != nil {
		return e
	}
	p.ptm.PeerTaskDone(p.req.PeerId)
	var end = time.Now()
	state, err := p.ptm.schedulerClient.ReportPeerResult(context.Background(), &scheduler.PeerResult{
		TaskId:         pt.GetTaskID(),
		PeerId:         pt.GetPeerID(),
		SrcIp:          p.ptm.host.Ip,
		SecurityDomain: p.ptm.host.SecurityDomain,
		Idc:            p.ptm.host.Idc,
		Url:            p.req.Url,
		ContentLength:  pt.GetContentLength(),
		Traffic:        pt.GetTraffic(),
		Cost:           uint32(end.Sub(p.start).Milliseconds()),
		Success:        true,
		Code:           dfcodes.Success,
	})
	logger.Debugf("task %s/%s report successful peer result, response state: %#v, error: %v", pt.GetTaskID(), pt.GetPeerID(), state, err)
	return nil
}

func (p *filePeerTaskCallback) Fail(pt PeerTask, reason string) error {
	p.ptm.PeerTaskDone(p.req.PeerId)
	var end = time.Now()
	state, err := p.ptm.schedulerClient.ReportPeerResult(context.Background(), &scheduler.PeerResult{
		TaskId:         pt.GetTaskID(),
		PeerId:         pt.GetPeerID(),
		SrcIp:          p.ptm.host.Ip,
		SecurityDomain: p.ptm.host.SecurityDomain,
		Idc:            p.ptm.host.Idc,
		Url:            p.req.Url,
		ContentLength:  pt.GetContentLength(),
		Traffic:        pt.GetTraffic(),
		Cost:           uint32(end.Sub(p.start).Milliseconds()),
		Success:        false,
		Code:           dfcodes.UnknownError,
	})
	logger.Debugf("task %s/%s report fail peer result, response state: %#v, error: %v", pt.GetTaskID(), pt.GetPeerID(), state, err)
	return nil
}
