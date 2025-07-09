package test

import (
	"bytes"
	"context"
	"io"
	"math"

	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	commonv1 "d7y.io/api/v2/pkg/apis/common/v1"
	dfdaemonv1mocks "d7y.io/api/v2/pkg/apis/dfdaemon/v1/mocks"
	"d7y.io/dragonfly/v2/pkg/digest"
)

// GenPieceDigests 计算每个piece的md5和总sha256
func GenPieceDigests(content []byte, pieceSize uint32) ([]string, string) {
	r := bytes.NewBuffer(content)
	var pieces = make([]string, int(math.Ceil(float64(len(content))/float64(pieceSize))))
	for i := range pieces {
		pieces[i] = digest.MD5FromReader(io.LimitReader(r, int64(pieceSize)))
	}
	return pieces, digest.SHA256FromStrings(pieces...)
}

// SetupMockDaemonServer 统一mock DaemonServer
func SetupMockDaemonServer(ctrl *gomock.Controller, pieceGen func(*commonv1.PieceTaskRequest) *commonv1.PiecePacket) *dfdaemonv1mocks.MockDaemonServer {
	daemon := dfdaemonv1mocks.NewMockDaemonServer(ctrl)
	daemon.EXPECT().GetPieceTasks(gomock.Any(), gomock.Any()).AnyTimes().
		DoAndReturn(func(ctx context.Context, request *commonv1.PieceTaskRequest) (*commonv1.PiecePacket, error) {
			return nil, status.Error(codes.Unimplemented, "TODO")
		})
	daemon.EXPECT().SyncPieceTasks(gomock.Any()).AnyTimes().DoAndReturn(func(s interface {
		Recv() (*commonv1.PieceTaskRequest, error)
		Send(*commonv1.PiecePacket) error
	}) error {
		request, err := s.Recv()
		if err != nil {
			return err
		}
		if err = s.Send(pieceGen(request)); err != nil {
			return err
		}
		for {
			request, err = s.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				return err
			}
			if err = s.Send(pieceGen(request)); err != nil {
				return err
			}
		}
		return nil
	})
	return daemon
}
