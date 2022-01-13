// Code generated by MockGen. DO NOT EDIT.
// Source: peertask_manager.go

// Package peer is a generated GoMock package.
package peer

import (
	context "context"
	io "io"
	reflect "reflect"

	dflog "d7y.io/dragonfly/v2/internal/dflog"
	scheduler "d7y.io/dragonfly/v2/pkg/rpc/scheduler"
	gomock "github.com/golang/mock/gomock"
)

// MockTaskManager is a mock of TaskManager interface.
type MockTaskManager struct {
	ctrl     *gomock.Controller
	recorder *MockTaskManagerMockRecorder
}

// MockTaskManagerMockRecorder is the mock recorder for MockTaskManager.
type MockTaskManagerMockRecorder struct {
	mock *MockTaskManager
}

// NewMockTaskManager creates a new mock instance.
func NewMockTaskManager(ctrl *gomock.Controller) *MockTaskManager {
	mock := &MockTaskManager{ctrl: ctrl}
	mock.recorder = &MockTaskManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTaskManager) EXPECT() *MockTaskManagerMockRecorder {
	return m.recorder
}

// IsPeerTaskRunning mocks base method.
func (m *MockTaskManager) IsPeerTaskRunning(pid string) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsPeerTaskRunning", pid)
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsPeerTaskRunning indicates an expected call of IsPeerTaskRunning.
func (mr *MockTaskManagerMockRecorder) IsPeerTaskRunning(pid interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsPeerTaskRunning", reflect.TypeOf((*MockTaskManager)(nil).IsPeerTaskRunning), pid)
}

// StartFilePeerTask mocks base method.
func (m *MockTaskManager) StartFilePeerTask(ctx context.Context, req *FilePeerTaskRequest) (chan *FilePeerTaskProgress, *TinyData, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StartFilePeerTask", ctx, req)
	ret0, _ := ret[0].(chan *FilePeerTaskProgress)
	ret1, _ := ret[1].(*TinyData)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// StartFilePeerTask indicates an expected call of StartFilePeerTask.
func (mr *MockTaskManagerMockRecorder) StartFilePeerTask(ctx, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StartFilePeerTask", reflect.TypeOf((*MockTaskManager)(nil).StartFilePeerTask), ctx, req)
}

// StartStreamPeerTask mocks base method.
func (m *MockTaskManager) StartStreamPeerTask(ctx context.Context, req *scheduler.PeerTaskRequest) (io.ReadCloser, map[string]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StartStreamPeerTask", ctx, req)
	ret0, _ := ret[0].(io.ReadCloser)
	ret1, _ := ret[1].(map[string]string)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// StartStreamPeerTask indicates an expected call of StartStreamPeerTask.
func (mr *MockTaskManagerMockRecorder) StartStreamPeerTask(ctx, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StartStreamPeerTask", reflect.TypeOf((*MockTaskManager)(nil).StartStreamPeerTask), ctx, req)
}

// Stop mocks base method.
func (m *MockTaskManager) Stop(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Stop", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// Stop indicates an expected call of Stop.
func (mr *MockTaskManagerMockRecorder) Stop(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stop", reflect.TypeOf((*MockTaskManager)(nil).Stop), ctx)
}

// MockTask is a mock of Task interface.
type MockTask struct {
	ctrl     *gomock.Controller
	recorder *MockTaskMockRecorder
}

// MockTaskMockRecorder is the mock recorder for MockTask.
type MockTaskMockRecorder struct {
	mock *MockTask
}

// NewMockTask creates a new mock instance.
func NewMockTask(ctrl *gomock.Controller) *MockTask {
	mock := &MockTask{ctrl: ctrl}
	mock.recorder = &MockTaskMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTask) EXPECT() *MockTaskMockRecorder {
	return m.recorder
}

// AddTraffic mocks base method.
func (m *MockTask) AddTraffic(arg0 uint64) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "AddTraffic", arg0)
}

// AddTraffic indicates an expected call of AddTraffic.
func (mr *MockTaskMockRecorder) AddTraffic(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddTraffic", reflect.TypeOf((*MockTask)(nil).AddTraffic), arg0)
}

// Context mocks base method.
func (m *MockTask) Context() context.Context {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Context")
	ret0, _ := ret[0].(context.Context)
	return ret0
}

// Context indicates an expected call of Context.
func (mr *MockTaskMockRecorder) Context() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Context", reflect.TypeOf((*MockTask)(nil).Context))
}

// GetContentLength mocks base method.
func (m *MockTask) GetContentLength() int64 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetContentLength")
	ret0, _ := ret[0].(int64)
	return ret0
}

// GetContentLength indicates an expected call of GetContentLength.
func (mr *MockTaskMockRecorder) GetContentLength() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetContentLength", reflect.TypeOf((*MockTask)(nil).GetContentLength))
}

// GetPeerID mocks base method.
func (m *MockTask) GetPeerID() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPeerID")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetPeerID indicates an expected call of GetPeerID.
func (mr *MockTaskMockRecorder) GetPeerID() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPeerID", reflect.TypeOf((*MockTask)(nil).GetPeerID))
}

// GetPieceMd5Sign mocks base method.
func (m *MockTask) GetPieceMd5Sign() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPieceMd5Sign")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetPieceMd5Sign indicates an expected call of GetPieceMd5Sign.
func (mr *MockTaskMockRecorder) GetPieceMd5Sign() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPieceMd5Sign", reflect.TypeOf((*MockTask)(nil).GetPieceMd5Sign))
}

// GetTaskID mocks base method.
func (m *MockTask) GetTaskID() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTaskID")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetTaskID indicates an expected call of GetTaskID.
func (mr *MockTaskMockRecorder) GetTaskID() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTaskID", reflect.TypeOf((*MockTask)(nil).GetTaskID))
}

// GetTotalPieces mocks base method.
func (m *MockTask) GetTotalPieces() int32 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTotalPieces")
	ret0, _ := ret[0].(int32)
	return ret0
}

// GetTotalPieces indicates an expected call of GetTotalPieces.
func (mr *MockTaskMockRecorder) GetTotalPieces() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTotalPieces", reflect.TypeOf((*MockTask)(nil).GetTotalPieces))
}

// GetTraffic mocks base method.
func (m *MockTask) GetTraffic() uint64 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTraffic")
	ret0, _ := ret[0].(uint64)
	return ret0
}

// GetTraffic indicates an expected call of GetTraffic.
func (mr *MockTaskMockRecorder) GetTraffic() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTraffic", reflect.TypeOf((*MockTask)(nil).GetTraffic))
}

// Log mocks base method.
func (m *MockTask) Log() *dflog.SugaredLoggerOnWith {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Log")
	ret0, _ := ret[0].(*dflog.SugaredLoggerOnWith)
	return ret0
}

// Log indicates an expected call of Log.
func (mr *MockTaskMockRecorder) Log() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Log", reflect.TypeOf((*MockTask)(nil).Log))
}

// PublishPieceInfo mocks base method.
func (m *MockTask) PublishPieceInfo(pieceNum int32, size uint32) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "PublishPieceInfo", pieceNum, size)
}

// PublishPieceInfo indicates an expected call of PublishPieceInfo.
func (mr *MockTaskMockRecorder) PublishPieceInfo(pieceNum, size interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PublishPieceInfo", reflect.TypeOf((*MockTask)(nil).PublishPieceInfo), pieceNum, size)
}

// ReportPieceResult mocks base method.
func (m *MockTask) ReportPieceResult(request *DownloadPieceRequest, result *DownloadPieceResult, err error) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ReportPieceResult", request, result, err)
}

// ReportPieceResult indicates an expected call of ReportPieceResult.
func (mr *MockTaskMockRecorder) ReportPieceResult(request, result, err interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReportPieceResult", reflect.TypeOf((*MockTask)(nil).ReportPieceResult), request, result, err)
}

// SetContentLength mocks base method.
func (m *MockTask) SetContentLength(arg0 int64) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetContentLength", arg0)
}

// SetContentLength indicates an expected call of SetContentLength.
func (mr *MockTaskMockRecorder) SetContentLength(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetContentLength", reflect.TypeOf((*MockTask)(nil).SetContentLength), arg0)
}

// SetPieceMd5Sign mocks base method.
func (m *MockTask) SetPieceMd5Sign(arg0 string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetPieceMd5Sign", arg0)
}

// SetPieceMd5Sign indicates an expected call of SetPieceMd5Sign.
func (mr *MockTaskMockRecorder) SetPieceMd5Sign(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetPieceMd5Sign", reflect.TypeOf((*MockTask)(nil).SetPieceMd5Sign), arg0)
}

// SetTotalPieces mocks base method.
func (m *MockTask) SetTotalPieces(arg0 int32) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetTotalPieces", arg0)
}

// SetTotalPieces indicates an expected call of SetTotalPieces.
func (mr *MockTaskMockRecorder) SetTotalPieces(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetTotalPieces", reflect.TypeOf((*MockTask)(nil).SetTotalPieces), arg0)
}

// MockLogger is a mock of Logger interface.
type MockLogger struct {
	ctrl     *gomock.Controller
	recorder *MockLoggerMockRecorder
}

// MockLoggerMockRecorder is the mock recorder for MockLogger.
type MockLoggerMockRecorder struct {
	mock *MockLogger
}

// NewMockLogger creates a new mock instance.
func NewMockLogger(ctrl *gomock.Controller) *MockLogger {
	mock := &MockLogger{ctrl: ctrl}
	mock.recorder = &MockLoggerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockLogger) EXPECT() *MockLoggerMockRecorder {
	return m.recorder
}

// Log mocks base method.
func (m *MockLogger) Log() *dflog.SugaredLoggerOnWith {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Log")
	ret0, _ := ret[0].(*dflog.SugaredLoggerOnWith)
	return ret0
}

// Log indicates an expected call of Log.
func (mr *MockLoggerMockRecorder) Log() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Log", reflect.TypeOf((*MockLogger)(nil).Log))
}
