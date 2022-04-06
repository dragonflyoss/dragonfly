// Code generated by MockGen. DO NOT EDIT.
// Source: ../../../pkg/rpc/dfdaemon/dfdaemon.pb.go

// Package mock_server_grpc is a generated GoMock package.
package mock_server_grpc

import (
	context "context"
	reflect "reflect"

	base "d7y.io/dragonfly/v2/pkg/rpc/base"
	dfdaemon "d7y.io/dragonfly/v2/pkg/rpc/dfdaemon"
	gomock "github.com/golang/mock/gomock"
	grpc "google.golang.org/grpc"
	metadata "google.golang.org/grpc/metadata"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// MockDaemonClient is a mock of DaemonClient interface.
type MockDaemonClient struct {
	ctrl     *gomock.Controller
	recorder *MockDaemonClientMockRecorder
}

// MockDaemonClientMockRecorder is the mock recorder for MockDaemonClient.
type MockDaemonClientMockRecorder struct {
	mock *MockDaemonClient
}

// NewMockDaemonClient creates a new mock instance.
func NewMockDaemonClient(ctrl *gomock.Controller) *MockDaemonClient {
	mock := &MockDaemonClient{ctrl: ctrl}
	mock.recorder = &MockDaemonClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDaemonClient) EXPECT() *MockDaemonClientMockRecorder {
	return m.recorder
}

// CheckHealth mocks base method.
func (m *MockDaemonClient) CheckHealth(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "CheckHealth", varargs...)
	ret0, _ := ret[0].(*emptypb.Empty)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckHealth indicates an expected call of CheckHealth.
func (mr *MockDaemonClientMockRecorder) CheckHealth(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckHealth", reflect.TypeOf((*MockDaemonClient)(nil).CheckHealth), varargs...)
}

// Download mocks base method.
func (m *MockDaemonClient) Download(ctx context.Context, in *dfdaemon.DownRequest, opts ...grpc.CallOption) (dfdaemon.Daemon_DownloadClient, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Download", varargs...)
	ret0, _ := ret[0].(dfdaemon.Daemon_DownloadClient)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Download indicates an expected call of Download.
func (mr *MockDaemonClientMockRecorder) Download(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Download", reflect.TypeOf((*MockDaemonClient)(nil).Download), varargs...)
}

// GetPieceTasks mocks base method.
func (m *MockDaemonClient) GetPieceTasks(ctx context.Context, in *base.PieceTaskRequest, opts ...grpc.CallOption) (*base.PiecePacket, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetPieceTasks", varargs...)
	ret0, _ := ret[0].(*base.PiecePacket)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPieceTasks indicates an expected call of GetPieceTasks.
func (mr *MockDaemonClientMockRecorder) GetPieceTasks(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPieceTasks", reflect.TypeOf((*MockDaemonClient)(nil).GetPieceTasks), varargs...)
}

// StatTask mocks base method.
func (m *MockDaemonClient) StatTask(ctx context.Context, in *dfdaemon.StatTaskRequest, opts ...grpc.CallOption) (*base.GrpcDfResult, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "StatTask", varargs...)
	ret0, _ := ret[0].(*base.GrpcDfResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// StatTask indicates an expected call of StatTask.
func (mr *MockDaemonClientMockRecorder) StatTask(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StatTask", reflect.TypeOf((*MockDaemonClient)(nil).StatTask), varargs...)
}

// SyncPieceTasks mocks base method.
func (m *MockDaemonClient) SyncPieceTasks(ctx context.Context, opts ...grpc.CallOption) (dfdaemon.Daemon_SyncPieceTasksClient, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "SyncPieceTasks", varargs...)
	ret0, _ := ret[0].(dfdaemon.Daemon_SyncPieceTasksClient)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SyncPieceTasks indicates an expected call of SyncPieceTasks.
func (mr *MockDaemonClientMockRecorder) SyncPieceTasks(ctx interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SyncPieceTasks", reflect.TypeOf((*MockDaemonClient)(nil).SyncPieceTasks), varargs...)
}

// MockDaemon_DownloadClient is a mock of Daemon_DownloadClient interface.
type MockDaemon_DownloadClient struct {
	ctrl     *gomock.Controller
	recorder *MockDaemon_DownloadClientMockRecorder
}

// MockDaemon_DownloadClientMockRecorder is the mock recorder for MockDaemon_DownloadClient.
type MockDaemon_DownloadClientMockRecorder struct {
	mock *MockDaemon_DownloadClient
}

// NewMockDaemon_DownloadClient creates a new mock instance.
func NewMockDaemon_DownloadClient(ctrl *gomock.Controller) *MockDaemon_DownloadClient {
	mock := &MockDaemon_DownloadClient{ctrl: ctrl}
	mock.recorder = &MockDaemon_DownloadClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDaemon_DownloadClient) EXPECT() *MockDaemon_DownloadClientMockRecorder {
	return m.recorder
}

// CloseSend mocks base method.
func (m *MockDaemon_DownloadClient) CloseSend() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CloseSend")
	ret0, _ := ret[0].(error)
	return ret0
}

// CloseSend indicates an expected call of CloseSend.
func (mr *MockDaemon_DownloadClientMockRecorder) CloseSend() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CloseSend", reflect.TypeOf((*MockDaemon_DownloadClient)(nil).CloseSend))
}

// Context mocks base method.
func (m *MockDaemon_DownloadClient) Context() context.Context {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Context")
	ret0, _ := ret[0].(context.Context)
	return ret0
}

// Context indicates an expected call of Context.
func (mr *MockDaemon_DownloadClientMockRecorder) Context() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Context", reflect.TypeOf((*MockDaemon_DownloadClient)(nil).Context))
}

// Header mocks base method.
func (m *MockDaemon_DownloadClient) Header() (metadata.MD, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Header")
	ret0, _ := ret[0].(metadata.MD)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Header indicates an expected call of Header.
func (mr *MockDaemon_DownloadClientMockRecorder) Header() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Header", reflect.TypeOf((*MockDaemon_DownloadClient)(nil).Header))
}

// Recv mocks base method.
func (m *MockDaemon_DownloadClient) Recv() (*dfdaemon.DownResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Recv")
	ret0, _ := ret[0].(*dfdaemon.DownResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Recv indicates an expected call of Recv.
func (mr *MockDaemon_DownloadClientMockRecorder) Recv() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Recv", reflect.TypeOf((*MockDaemon_DownloadClient)(nil).Recv))
}

// RecvMsg mocks base method.
func (m_2 *MockDaemon_DownloadClient) RecvMsg(m interface{}) error {
	m_2.ctrl.T.Helper()
	ret := m_2.ctrl.Call(m_2, "RecvMsg", m)
	ret0, _ := ret[0].(error)
	return ret0
}

// RecvMsg indicates an expected call of RecvMsg.
func (mr *MockDaemon_DownloadClientMockRecorder) RecvMsg(m interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RecvMsg", reflect.TypeOf((*MockDaemon_DownloadClient)(nil).RecvMsg), m)
}

// SendMsg mocks base method.
func (m_2 *MockDaemon_DownloadClient) SendMsg(m interface{}) error {
	m_2.ctrl.T.Helper()
	ret := m_2.ctrl.Call(m_2, "SendMsg", m)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendMsg indicates an expected call of SendMsg.
func (mr *MockDaemon_DownloadClientMockRecorder) SendMsg(m interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendMsg", reflect.TypeOf((*MockDaemon_DownloadClient)(nil).SendMsg), m)
}

// Trailer mocks base method.
func (m *MockDaemon_DownloadClient) Trailer() metadata.MD {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Trailer")
	ret0, _ := ret[0].(metadata.MD)
	return ret0
}

// Trailer indicates an expected call of Trailer.
func (mr *MockDaemon_DownloadClientMockRecorder) Trailer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Trailer", reflect.TypeOf((*MockDaemon_DownloadClient)(nil).Trailer))
}

// MockDaemon_SyncPieceTasksClient is a mock of Daemon_SyncPieceTasksClient interface.
type MockDaemon_SyncPieceTasksClient struct {
	ctrl     *gomock.Controller
	recorder *MockDaemon_SyncPieceTasksClientMockRecorder
}

// MockDaemon_SyncPieceTasksClientMockRecorder is the mock recorder for MockDaemon_SyncPieceTasksClient.
type MockDaemon_SyncPieceTasksClientMockRecorder struct {
	mock *MockDaemon_SyncPieceTasksClient
}

// NewMockDaemon_SyncPieceTasksClient creates a new mock instance.
func NewMockDaemon_SyncPieceTasksClient(ctrl *gomock.Controller) *MockDaemon_SyncPieceTasksClient {
	mock := &MockDaemon_SyncPieceTasksClient{ctrl: ctrl}
	mock.recorder = &MockDaemon_SyncPieceTasksClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDaemon_SyncPieceTasksClient) EXPECT() *MockDaemon_SyncPieceTasksClientMockRecorder {
	return m.recorder
}

// CloseSend mocks base method.
func (m *MockDaemon_SyncPieceTasksClient) CloseSend() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CloseSend")
	ret0, _ := ret[0].(error)
	return ret0
}

// CloseSend indicates an expected call of CloseSend.
func (mr *MockDaemon_SyncPieceTasksClientMockRecorder) CloseSend() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CloseSend", reflect.TypeOf((*MockDaemon_SyncPieceTasksClient)(nil).CloseSend))
}

// Context mocks base method.
func (m *MockDaemon_SyncPieceTasksClient) Context() context.Context {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Context")
	ret0, _ := ret[0].(context.Context)
	return ret0
}

// Context indicates an expected call of Context.
func (mr *MockDaemon_SyncPieceTasksClientMockRecorder) Context() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Context", reflect.TypeOf((*MockDaemon_SyncPieceTasksClient)(nil).Context))
}

// Header mocks base method.
func (m *MockDaemon_SyncPieceTasksClient) Header() (metadata.MD, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Header")
	ret0, _ := ret[0].(metadata.MD)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Header indicates an expected call of Header.
func (mr *MockDaemon_SyncPieceTasksClientMockRecorder) Header() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Header", reflect.TypeOf((*MockDaemon_SyncPieceTasksClient)(nil).Header))
}

// Recv mocks base method.
func (m *MockDaemon_SyncPieceTasksClient) Recv() (*base.PiecePacket, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Recv")
	ret0, _ := ret[0].(*base.PiecePacket)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Recv indicates an expected call of Recv.
func (mr *MockDaemon_SyncPieceTasksClientMockRecorder) Recv() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Recv", reflect.TypeOf((*MockDaemon_SyncPieceTasksClient)(nil).Recv))
}

// RecvMsg mocks base method.
func (m_2 *MockDaemon_SyncPieceTasksClient) RecvMsg(m interface{}) error {
	m_2.ctrl.T.Helper()
	ret := m_2.ctrl.Call(m_2, "RecvMsg", m)
	ret0, _ := ret[0].(error)
	return ret0
}

// RecvMsg indicates an expected call of RecvMsg.
func (mr *MockDaemon_SyncPieceTasksClientMockRecorder) RecvMsg(m interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RecvMsg", reflect.TypeOf((*MockDaemon_SyncPieceTasksClient)(nil).RecvMsg), m)
}

// Send mocks base method.
func (m *MockDaemon_SyncPieceTasksClient) Send(arg0 *base.PieceTaskRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Send", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Send indicates an expected call of Send.
func (mr *MockDaemon_SyncPieceTasksClientMockRecorder) Send(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Send", reflect.TypeOf((*MockDaemon_SyncPieceTasksClient)(nil).Send), arg0)
}

// SendMsg mocks base method.
func (m_2 *MockDaemon_SyncPieceTasksClient) SendMsg(m interface{}) error {
	m_2.ctrl.T.Helper()
	ret := m_2.ctrl.Call(m_2, "SendMsg", m)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendMsg indicates an expected call of SendMsg.
func (mr *MockDaemon_SyncPieceTasksClientMockRecorder) SendMsg(m interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendMsg", reflect.TypeOf((*MockDaemon_SyncPieceTasksClient)(nil).SendMsg), m)
}

// Trailer mocks base method.
func (m *MockDaemon_SyncPieceTasksClient) Trailer() metadata.MD {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Trailer")
	ret0, _ := ret[0].(metadata.MD)
	return ret0
}

// Trailer indicates an expected call of Trailer.
func (mr *MockDaemon_SyncPieceTasksClientMockRecorder) Trailer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Trailer", reflect.TypeOf((*MockDaemon_SyncPieceTasksClient)(nil).Trailer))
}

// MockDaemonServer is a mock of DaemonServer interface.
type MockDaemonServer struct {
	ctrl     *gomock.Controller
	recorder *MockDaemonServerMockRecorder
}

// MockDaemonServerMockRecorder is the mock recorder for MockDaemonServer.
type MockDaemonServerMockRecorder struct {
	mock *MockDaemonServer
}

// NewMockDaemonServer creates a new mock instance.
func NewMockDaemonServer(ctrl *gomock.Controller) *MockDaemonServer {
	mock := &MockDaemonServer{ctrl: ctrl}
	mock.recorder = &MockDaemonServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDaemonServer) EXPECT() *MockDaemonServerMockRecorder {
	return m.recorder
}

// CheckHealth mocks base method.
func (m *MockDaemonServer) CheckHealth(arg0 context.Context, arg1 *emptypb.Empty) (*emptypb.Empty, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckHealth", arg0, arg1)
	ret0, _ := ret[0].(*emptypb.Empty)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckHealth indicates an expected call of CheckHealth.
func (mr *MockDaemonServerMockRecorder) CheckHealth(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckHealth", reflect.TypeOf((*MockDaemonServer)(nil).CheckHealth), arg0, arg1)
}

// Download mocks base method.
func (m *MockDaemonServer) Download(arg0 *dfdaemon.DownRequest, arg1 dfdaemon.Daemon_DownloadServer) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Download", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Download indicates an expected call of Download.
func (mr *MockDaemonServerMockRecorder) Download(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Download", reflect.TypeOf((*MockDaemonServer)(nil).Download), arg0, arg1)
}

// GetPieceTasks mocks base method.
func (m *MockDaemonServer) GetPieceTasks(arg0 context.Context, arg1 *base.PieceTaskRequest) (*base.PiecePacket, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPieceTasks", arg0, arg1)
	ret0, _ := ret[0].(*base.PiecePacket)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPieceTasks indicates an expected call of GetPieceTasks.
func (mr *MockDaemonServerMockRecorder) GetPieceTasks(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPieceTasks", reflect.TypeOf((*MockDaemonServer)(nil).GetPieceTasks), arg0, arg1)
}

// StatTask mocks base method.
func (m *MockDaemonServer) StatTask(arg0 context.Context, arg1 *dfdaemon.StatTaskRequest) (*base.GrpcDfResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StatTask", arg0, arg1)
	ret0, _ := ret[0].(*base.GrpcDfResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// StatTask indicates an expected call of StatTask.
func (mr *MockDaemonServerMockRecorder) StatTask(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StatTask", reflect.TypeOf((*MockDaemonServer)(nil).StatTask), arg0, arg1)
}

// SyncPieceTasks mocks base method.
func (m *MockDaemonServer) SyncPieceTasks(arg0 dfdaemon.Daemon_SyncPieceTasksServer) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SyncPieceTasks", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SyncPieceTasks indicates an expected call of SyncPieceTasks.
func (mr *MockDaemonServerMockRecorder) SyncPieceTasks(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SyncPieceTasks", reflect.TypeOf((*MockDaemonServer)(nil).SyncPieceTasks), arg0)
}

// MockDaemon_DownloadServer is a mock of Daemon_DownloadServer interface.
type MockDaemon_DownloadServer struct {
	ctrl     *gomock.Controller
	recorder *MockDaemon_DownloadServerMockRecorder
}

// MockDaemon_DownloadServerMockRecorder is the mock recorder for MockDaemon_DownloadServer.
type MockDaemon_DownloadServerMockRecorder struct {
	mock *MockDaemon_DownloadServer
}

// NewMockDaemon_DownloadServer creates a new mock instance.
func NewMockDaemon_DownloadServer(ctrl *gomock.Controller) *MockDaemon_DownloadServer {
	mock := &MockDaemon_DownloadServer{ctrl: ctrl}
	mock.recorder = &MockDaemon_DownloadServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDaemon_DownloadServer) EXPECT() *MockDaemon_DownloadServerMockRecorder {
	return m.recorder
}

// Context mocks base method.
func (m *MockDaemon_DownloadServer) Context() context.Context {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Context")
	ret0, _ := ret[0].(context.Context)
	return ret0
}

// Context indicates an expected call of Context.
func (mr *MockDaemon_DownloadServerMockRecorder) Context() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Context", reflect.TypeOf((*MockDaemon_DownloadServer)(nil).Context))
}

// RecvMsg mocks base method.
func (m_2 *MockDaemon_DownloadServer) RecvMsg(m interface{}) error {
	m_2.ctrl.T.Helper()
	ret := m_2.ctrl.Call(m_2, "RecvMsg", m)
	ret0, _ := ret[0].(error)
	return ret0
}

// RecvMsg indicates an expected call of RecvMsg.
func (mr *MockDaemon_DownloadServerMockRecorder) RecvMsg(m interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RecvMsg", reflect.TypeOf((*MockDaemon_DownloadServer)(nil).RecvMsg), m)
}

// Send mocks base method.
func (m *MockDaemon_DownloadServer) Send(arg0 *dfdaemon.DownResult) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Send", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Send indicates an expected call of Send.
func (mr *MockDaemon_DownloadServerMockRecorder) Send(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Send", reflect.TypeOf((*MockDaemon_DownloadServer)(nil).Send), arg0)
}

// SendHeader mocks base method.
func (m *MockDaemon_DownloadServer) SendHeader(arg0 metadata.MD) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendHeader", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendHeader indicates an expected call of SendHeader.
func (mr *MockDaemon_DownloadServerMockRecorder) SendHeader(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendHeader", reflect.TypeOf((*MockDaemon_DownloadServer)(nil).SendHeader), arg0)
}

// SendMsg mocks base method.
func (m_2 *MockDaemon_DownloadServer) SendMsg(m interface{}) error {
	m_2.ctrl.T.Helper()
	ret := m_2.ctrl.Call(m_2, "SendMsg", m)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendMsg indicates an expected call of SendMsg.
func (mr *MockDaemon_DownloadServerMockRecorder) SendMsg(m interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendMsg", reflect.TypeOf((*MockDaemon_DownloadServer)(nil).SendMsg), m)
}

// SetHeader mocks base method.
func (m *MockDaemon_DownloadServer) SetHeader(arg0 metadata.MD) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetHeader", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetHeader indicates an expected call of SetHeader.
func (mr *MockDaemon_DownloadServerMockRecorder) SetHeader(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetHeader", reflect.TypeOf((*MockDaemon_DownloadServer)(nil).SetHeader), arg0)
}

// SetTrailer mocks base method.
func (m *MockDaemon_DownloadServer) SetTrailer(arg0 metadata.MD) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetTrailer", arg0)
}

// SetTrailer indicates an expected call of SetTrailer.
func (mr *MockDaemon_DownloadServerMockRecorder) SetTrailer(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetTrailer", reflect.TypeOf((*MockDaemon_DownloadServer)(nil).SetTrailer), arg0)
}

// MockDaemon_SyncPieceTasksServer is a mock of Daemon_SyncPieceTasksServer interface.
type MockDaemon_SyncPieceTasksServer struct {
	ctrl     *gomock.Controller
	recorder *MockDaemon_SyncPieceTasksServerMockRecorder
}

// MockDaemon_SyncPieceTasksServerMockRecorder is the mock recorder for MockDaemon_SyncPieceTasksServer.
type MockDaemon_SyncPieceTasksServerMockRecorder struct {
	mock *MockDaemon_SyncPieceTasksServer
}

// NewMockDaemon_SyncPieceTasksServer creates a new mock instance.
func NewMockDaemon_SyncPieceTasksServer(ctrl *gomock.Controller) *MockDaemon_SyncPieceTasksServer {
	mock := &MockDaemon_SyncPieceTasksServer{ctrl: ctrl}
	mock.recorder = &MockDaemon_SyncPieceTasksServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDaemon_SyncPieceTasksServer) EXPECT() *MockDaemon_SyncPieceTasksServerMockRecorder {
	return m.recorder
}

// Context mocks base method.
func (m *MockDaemon_SyncPieceTasksServer) Context() context.Context {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Context")
	ret0, _ := ret[0].(context.Context)
	return ret0
}

// Context indicates an expected call of Context.
func (mr *MockDaemon_SyncPieceTasksServerMockRecorder) Context() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Context", reflect.TypeOf((*MockDaemon_SyncPieceTasksServer)(nil).Context))
}

// Recv mocks base method.
func (m *MockDaemon_SyncPieceTasksServer) Recv() (*base.PieceTaskRequest, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Recv")
	ret0, _ := ret[0].(*base.PieceTaskRequest)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Recv indicates an expected call of Recv.
func (mr *MockDaemon_SyncPieceTasksServerMockRecorder) Recv() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Recv", reflect.TypeOf((*MockDaemon_SyncPieceTasksServer)(nil).Recv))
}

// RecvMsg mocks base method.
func (m_2 *MockDaemon_SyncPieceTasksServer) RecvMsg(m interface{}) error {
	m_2.ctrl.T.Helper()
	ret := m_2.ctrl.Call(m_2, "RecvMsg", m)
	ret0, _ := ret[0].(error)
	return ret0
}

// RecvMsg indicates an expected call of RecvMsg.
func (mr *MockDaemon_SyncPieceTasksServerMockRecorder) RecvMsg(m interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RecvMsg", reflect.TypeOf((*MockDaemon_SyncPieceTasksServer)(nil).RecvMsg), m)
}

// Send mocks base method.
func (m *MockDaemon_SyncPieceTasksServer) Send(arg0 *base.PiecePacket) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Send", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Send indicates an expected call of Send.
func (mr *MockDaemon_SyncPieceTasksServerMockRecorder) Send(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Send", reflect.TypeOf((*MockDaemon_SyncPieceTasksServer)(nil).Send), arg0)
}

// SendHeader mocks base method.
func (m *MockDaemon_SyncPieceTasksServer) SendHeader(arg0 metadata.MD) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendHeader", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendHeader indicates an expected call of SendHeader.
func (mr *MockDaemon_SyncPieceTasksServerMockRecorder) SendHeader(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendHeader", reflect.TypeOf((*MockDaemon_SyncPieceTasksServer)(nil).SendHeader), arg0)
}

// SendMsg mocks base method.
func (m_2 *MockDaemon_SyncPieceTasksServer) SendMsg(m interface{}) error {
	m_2.ctrl.T.Helper()
	ret := m_2.ctrl.Call(m_2, "SendMsg", m)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendMsg indicates an expected call of SendMsg.
func (mr *MockDaemon_SyncPieceTasksServerMockRecorder) SendMsg(m interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendMsg", reflect.TypeOf((*MockDaemon_SyncPieceTasksServer)(nil).SendMsg), m)
}

// SetHeader mocks base method.
func (m *MockDaemon_SyncPieceTasksServer) SetHeader(arg0 metadata.MD) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetHeader", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetHeader indicates an expected call of SetHeader.
func (mr *MockDaemon_SyncPieceTasksServerMockRecorder) SetHeader(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetHeader", reflect.TypeOf((*MockDaemon_SyncPieceTasksServer)(nil).SetHeader), arg0)
}

// SetTrailer mocks base method.
func (m *MockDaemon_SyncPieceTasksServer) SetTrailer(arg0 metadata.MD) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetTrailer", arg0)
}

// SetTrailer indicates an expected call of SetTrailer.
func (mr *MockDaemon_SyncPieceTasksServerMockRecorder) SetTrailer(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetTrailer", reflect.TypeOf((*MockDaemon_SyncPieceTasksServer)(nil).SetTrailer), arg0)
}
