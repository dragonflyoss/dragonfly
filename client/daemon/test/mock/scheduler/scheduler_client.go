// Code generated by MockGen. DO NOT EDIT.
// Source: pkg/rpc/scheduler/client/client.go

// Package mock_client is a generated GoMock package.
package mock_client

import (
	context "context"
	reflect "reflect"

	dfnet "d7y.io/dragonfly/v2/internal/dfnet"
	scheduler "d7y.io/dragonfly/v2/pkg/rpc/scheduler"
	client "d7y.io/dragonfly/v2/pkg/rpc/scheduler/client"
	gomock "github.com/golang/mock/gomock"
	grpc "google.golang.org/grpc"
)

// MockSchedulerClient is a mock of SchedulerClient interface.
type MockSchedulerClient struct {
	ctrl     *gomock.Controller
	recorder *MockSchedulerClientMockRecorder
}

// MockSchedulerClientMockRecorder is the mock recorder for MockSchedulerClient.
type MockSchedulerClientMockRecorder struct {
	mock *MockSchedulerClient
}

// NewMockSchedulerClient creates a new mock instance.
func NewMockSchedulerClient(ctrl *gomock.Controller) *MockSchedulerClient {
	mock := &MockSchedulerClient{ctrl: ctrl}
	mock.recorder = &MockSchedulerClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSchedulerClient) EXPECT() *MockSchedulerClientMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *MockSchedulerClient) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockSchedulerClientMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockSchedulerClient)(nil).Close))
}

// GetState mocks base method.
func (m *MockSchedulerClient) GetState() []dfnet.NetAddr {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetState")
	ret0, _ := ret[0].([]dfnet.NetAddr)
	return ret0
}

// GetState indicates an expected call of GetState.
func (mr *MockSchedulerClientMockRecorder) GetState() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetState", reflect.TypeOf((*MockSchedulerClient)(nil).GetState))
}

// LeaveTask mocks base method.
func (m *MockSchedulerClient) LeaveTask(arg0 context.Context, arg1 *scheduler.PeerTarget, arg2 ...grpc.CallOption) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "LeaveTask", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// LeaveTask indicates an expected call of LeaveTask.
func (mr *MockSchedulerClientMockRecorder) LeaveTask(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LeaveTask", reflect.TypeOf((*MockSchedulerClient)(nil).LeaveTask), varargs...)
}

// RegisterPeerTask mocks base method.
func (m *MockSchedulerClient) RegisterPeerTask(arg0 context.Context, arg1 *scheduler.PeerTaskRequest, arg2 ...grpc.CallOption) (*scheduler.RegisterResult, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "RegisterPeerTask", varargs...)
	ret0, _ := ret[0].(*scheduler.RegisterResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RegisterPeerTask indicates an expected call of RegisterPeerTask.
func (mr *MockSchedulerClientMockRecorder) RegisterPeerTask(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterPeerTask", reflect.TypeOf((*MockSchedulerClient)(nil).RegisterPeerTask), varargs...)
}

// ReportPeerResult mocks base method.
func (m *MockSchedulerClient) ReportPeerResult(arg0 context.Context, arg1 *scheduler.PeerResult, arg2 ...grpc.CallOption) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ReportPeerResult", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// ReportPeerResult indicates an expected call of ReportPeerResult.
func (mr *MockSchedulerClientMockRecorder) ReportPeerResult(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReportPeerResult", reflect.TypeOf((*MockSchedulerClient)(nil).ReportPeerResult), varargs...)
}

// ReportPieceResult mocks base method.
func (m *MockSchedulerClient) ReportPieceResult(arg0 context.Context, arg1 string, arg2 *scheduler.PeerTaskRequest, arg3 ...grpc.CallOption) (client.PeerPacketStream, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1, arg2}
	for _, a := range arg3 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ReportPieceResult", varargs...)
	ret0, _ := ret[0].(client.PeerPacketStream)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReportPieceResult indicates an expected call of ReportPieceResult.
func (mr *MockSchedulerClientMockRecorder) ReportPieceResult(arg0, arg1, arg2 interface{}, arg3 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1, arg2}, arg3...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReportPieceResult", reflect.TypeOf((*MockSchedulerClient)(nil).ReportPieceResult), varargs...)
}

// UpdateState mocks base method.
func (m *MockSchedulerClient) UpdateState(arg0 []dfnet.NetAddr) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "UpdateState", arg0)
}

// UpdateState indicates an expected call of UpdateState.
func (mr *MockSchedulerClientMockRecorder) UpdateState(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateState", reflect.TypeOf((*MockSchedulerClient)(nil).UpdateState), arg0)
}
