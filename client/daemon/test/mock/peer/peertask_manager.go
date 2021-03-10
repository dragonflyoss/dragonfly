// Code generated by MockGen. DO NOT EDIT.
// Source: ../../../peer/peertask_manager.go

// Package mock_peer is a generated GoMock package.
package mock_peer

import (
	context "context"
	io "io"
	reflect "reflect"

	peer "d7y.io/dragonfly/v2/client/daemon/peer"
	base "d7y.io/dragonfly/v2/pkg/rpc/base"
	scheduler "d7y.io/dragonfly/v2/pkg/rpc/scheduler"
	gomock "github.com/golang/mock/gomock"
)

// MockPeerTaskManager is a mock of PeerTaskManager interface.
type MockPeerTaskManager struct {
	ctrl     *gomock.Controller
	recorder *MockPeerTaskManagerMockRecorder
}

// MockPeerTaskManagerMockRecorder is the mock recorder for MockPeerTaskManager.
type MockPeerTaskManagerMockRecorder struct {
	mock *MockPeerTaskManager
}

// NewMockPeerTaskManager creates a new mock instance.
func NewMockPeerTaskManager(ctrl *gomock.Controller) *MockPeerTaskManager {
	mock := &MockPeerTaskManager{ctrl: ctrl}
	mock.recorder = &MockPeerTaskManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPeerTaskManager) EXPECT() *MockPeerTaskManagerMockRecorder {
	return m.recorder
}

// StartFilePeerTask mocks base method.
func (m *MockPeerTaskManager) StartFilePeerTask(ctx context.Context, req *peer.FilePeerTaskRequest) (chan *peer.PeerTaskProgress, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StartFilePeerTask", ctx, req)
	ret0, _ := ret[0].(chan *peer.PeerTaskProgress)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// StartFilePeerTask indicates an expected call of StartFilePeerTask.
func (mr *MockPeerTaskManagerMockRecorder) StartFilePeerTask(ctx, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StartFilePeerTask", reflect.TypeOf((*MockPeerTaskManager)(nil).StartFilePeerTask), ctx, req)
}

// StartStreamPeerTask mocks base method.
func (m *MockPeerTaskManager) StartStreamPeerTask(ctx context.Context, req *scheduler.PeerTaskRequest) (io.Reader, map[string]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StartStreamPeerTask", ctx, req)
	ret0, _ := ret[0].(io.Reader)
	ret1, _ := ret[1].(map[string]string)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// StartStreamPeerTask indicates an expected call of StartStreamPeerTask.
func (mr *MockPeerTaskManagerMockRecorder) StartStreamPeerTask(ctx, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StartStreamPeerTask", reflect.TypeOf((*MockPeerTaskManager)(nil).StartStreamPeerTask), ctx, req)
}

// Stop mocks base method.
func (m *MockPeerTaskManager) Stop(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Stop", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// Stop indicates an expected call of Stop.
func (mr *MockPeerTaskManagerMockRecorder) Stop(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stop", reflect.TypeOf((*MockPeerTaskManager)(nil).Stop), ctx)
}

// MockPeerTask is a mock of PeerTask interface.
type MockPeerTask struct {
	ctrl     *gomock.Controller
	recorder *MockPeerTaskMockRecorder
}

// MockPeerTaskMockRecorder is the mock recorder for MockPeerTask.
type MockPeerTaskMockRecorder struct {
	mock *MockPeerTask
}

// NewMockPeerTask creates a new mock instance.
func NewMockPeerTask(ctrl *gomock.Controller) *MockPeerTask {
	mock := &MockPeerTask{ctrl: ctrl}
	mock.recorder = &MockPeerTaskMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPeerTask) EXPECT() *MockPeerTaskMockRecorder {
	return m.recorder
}

// AddTraffic mocks base method.
func (m *MockPeerTask) AddTraffic(arg0 int64) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "AddTraffic", arg0)
}

// AddTraffic indicates an expected call of AddTraffic.
func (mr *MockPeerTaskMockRecorder) AddTraffic(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddTraffic", reflect.TypeOf((*MockPeerTask)(nil).AddTraffic), arg0)
}

// GetContentLength mocks base method.
func (m *MockPeerTask) GetContentLength() int64 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetContentLength")
	ret0, _ := ret[0].(int64)
	return ret0
}

// GetContentLength indicates an expected call of GetContentLength.
func (mr *MockPeerTaskMockRecorder) GetContentLength() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetContentLength", reflect.TypeOf((*MockPeerTask)(nil).GetContentLength))
}

// GetPeerID mocks base method.
func (m *MockPeerTask) GetPeerID() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPeerID")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetPeerID indicates an expected call of GetPeerID.
func (mr *MockPeerTaskMockRecorder) GetPeerID() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPeerID", reflect.TypeOf((*MockPeerTask)(nil).GetPeerID))
}

// GetTaskID mocks base method.
func (m *MockPeerTask) GetTaskID() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTaskID")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetTaskID indicates an expected call of GetTaskID.
func (mr *MockPeerTaskMockRecorder) GetTaskID() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTaskID", reflect.TypeOf((*MockPeerTask)(nil).GetTaskID))
}

// GetTraffic mocks base method.
func (m *MockPeerTask) GetTraffic() int64 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTraffic")
	ret0, _ := ret[0].(int64)
	return ret0
}

// GetTraffic indicates an expected call of GetTraffic.
func (mr *MockPeerTaskMockRecorder) GetTraffic() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTraffic", reflect.TypeOf((*MockPeerTask)(nil).GetTraffic))
}

// ReportPieceResult mocks base method.
func (m *MockPeerTask) ReportPieceResult(pieceTask *base.PieceInfo, pieceResult *scheduler.PieceResult) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReportPieceResult", pieceTask, pieceResult)
	ret0, _ := ret[0].(error)
	return ret0
}

// ReportPieceResult indicates an expected call of ReportPieceResult.
func (mr *MockPeerTaskMockRecorder) ReportPieceResult(pieceTask, pieceResult interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReportPieceResult", reflect.TypeOf((*MockPeerTask)(nil).ReportPieceResult), pieceTask, pieceResult)
}

// SetCallback mocks base method.
func (m *MockPeerTask) SetCallback(arg0 peer.PeerTaskCallback) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetCallback", arg0)
}

// SetCallback indicates an expected call of SetCallback.
func (mr *MockPeerTaskMockRecorder) SetCallback(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetCallback", reflect.TypeOf((*MockPeerTask)(nil).SetCallback), arg0)
}

// SetContentLength mocks base method.
func (m *MockPeerTask) SetContentLength(arg0 int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetContentLength", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetContentLength indicates an expected call of SetContentLength.
func (mr *MockPeerTaskMockRecorder) SetContentLength(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetContentLength", reflect.TypeOf((*MockPeerTask)(nil).SetContentLength), arg0)
}

// MockPeerTaskCallback is a mock of PeerTaskCallback interface.
type MockPeerTaskCallback struct {
	ctrl     *gomock.Controller
	recorder *MockPeerTaskCallbackMockRecorder
}

// MockPeerTaskCallbackMockRecorder is the mock recorder for MockPeerTaskCallback.
type MockPeerTaskCallbackMockRecorder struct {
	mock *MockPeerTaskCallback
}

// NewMockPeerTaskCallback creates a new mock instance.
func NewMockPeerTaskCallback(ctrl *gomock.Controller) *MockPeerTaskCallback {
	mock := &MockPeerTaskCallback{ctrl: ctrl}
	mock.recorder = &MockPeerTaskCallbackMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPeerTaskCallback) EXPECT() *MockPeerTaskCallbackMockRecorder {
	return m.recorder
}

// Done mocks base method.
func (m *MockPeerTaskCallback) Done(pt peer.PeerTask) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PeerTaskDone", pt)
	ret0, _ := ret[0].(error)
	return ret0
}

// Done indicates an expected call of Done.
func (mr *MockPeerTaskCallbackMockRecorder) Done(pt interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PeerTaskDone", reflect.TypeOf((*MockPeerTaskCallback)(nil).Done), pt)
}

// Fail mocks base method.
func (m *MockPeerTaskCallback) Fail(pt peer.PeerTask, reason string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Fail", pt, reason)
	ret0, _ := ret[0].(error)
	return ret0
}

// Fail indicates an expected call of Fail.
func (mr *MockPeerTaskCallbackMockRecorder) Fail(pt, reason interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Fail", reflect.TypeOf((*MockPeerTaskCallback)(nil).Fail), pt, reason)
}

// Init mocks base method.
func (m *MockPeerTaskCallback) Init(pt peer.PeerTask) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Init", pt)
	ret0, _ := ret[0].(error)
	return ret0
}

// Init indicates an expected call of Init.
func (mr *MockPeerTaskCallbackMockRecorder) Init(pt interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Init", reflect.TypeOf((*MockPeerTaskCallback)(nil).Init), pt)
}
