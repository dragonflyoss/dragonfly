// Code generated by MockGen. DO NOT EDIT.
// Source: client_v1.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	scheduler "d7y.io/api/pkg/apis/scheduler/v1"
	gomock "github.com/golang/mock/gomock"
	grpc "google.golang.org/grpc"
)

// MockV1 is a mock of V1 interface.
type MockV1 struct {
	ctrl     *gomock.Controller
	recorder *MockV1MockRecorder
}

// MockV1MockRecorder is the mock recorder for MockV1.
type MockV1MockRecorder struct {
	mock *MockV1
}

// NewMockV1 creates a new mock instance.
func NewMockV1(ctrl *gomock.Controller) *MockV1 {
	mock := &MockV1{ctrl: ctrl}
	mock.recorder = &MockV1MockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockV1) EXPECT() *MockV1MockRecorder {
	return m.recorder
}

// AnnounceHost mocks base method.
func (m *MockV1) AnnounceHost(arg0 context.Context, arg1 *scheduler.AnnounceHostRequest, arg2 ...grpc.CallOption) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "AnnounceHost", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// AnnounceHost indicates an expected call of AnnounceHost.
func (mr *MockV1MockRecorder) AnnounceHost(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AnnounceHost", reflect.TypeOf((*MockV1)(nil).AnnounceHost), varargs...)
}

// AnnounceTask mocks base method.
func (m *MockV1) AnnounceTask(arg0 context.Context, arg1 *scheduler.AnnounceTaskRequest, arg2 ...grpc.CallOption) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "AnnounceTask", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// AnnounceTask indicates an expected call of AnnounceTask.
func (mr *MockV1MockRecorder) AnnounceTask(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AnnounceTask", reflect.TypeOf((*MockV1)(nil).AnnounceTask), varargs...)
}

// Close mocks base method.
func (m *MockV1) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockV1MockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockV1)(nil).Close))
}

// LeaveHost mocks base method.
func (m *MockV1) LeaveHost(arg0 context.Context, arg1 *scheduler.LeaveHostRequest, arg2 ...grpc.CallOption) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "LeaveHost", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// LeaveHost indicates an expected call of LeaveHost.
func (mr *MockV1MockRecorder) LeaveHost(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LeaveHost", reflect.TypeOf((*MockV1)(nil).LeaveHost), varargs...)
}

// LeaveTask mocks base method.
func (m *MockV1) LeaveTask(arg0 context.Context, arg1 *scheduler.PeerTarget, arg2 ...grpc.CallOption) error {
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
func (mr *MockV1MockRecorder) LeaveTask(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LeaveTask", reflect.TypeOf((*MockV1)(nil).LeaveTask), varargs...)
}

// RegisterPeerTask mocks base method.
func (m *MockV1) RegisterPeerTask(arg0 context.Context, arg1 *scheduler.PeerTaskRequest, arg2 ...grpc.CallOption) (*scheduler.RegisterResult, error) {
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
func (mr *MockV1MockRecorder) RegisterPeerTask(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterPeerTask", reflect.TypeOf((*MockV1)(nil).RegisterPeerTask), varargs...)
}

// ReportPeerResult mocks base method.
func (m *MockV1) ReportPeerResult(arg0 context.Context, arg1 *scheduler.PeerResult, arg2 ...grpc.CallOption) error {
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
func (mr *MockV1MockRecorder) ReportPeerResult(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReportPeerResult", reflect.TypeOf((*MockV1)(nil).ReportPeerResult), varargs...)
}

// ReportPieceResult mocks base method.
func (m *MockV1) ReportPieceResult(arg0 context.Context, arg1 *scheduler.PeerTaskRequest, arg2 ...grpc.CallOption) (scheduler.Scheduler_ReportPieceResultClient, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ReportPieceResult", varargs...)
	ret0, _ := ret[0].(scheduler.Scheduler_ReportPieceResultClient)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReportPieceResult indicates an expected call of ReportPieceResult.
func (mr *MockV1MockRecorder) ReportPieceResult(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReportPieceResult", reflect.TypeOf((*MockV1)(nil).ReportPieceResult), varargs...)
}

// StatTask mocks base method.
func (m *MockV1) StatTask(arg0 context.Context, arg1 *scheduler.StatTaskRequest, arg2 ...grpc.CallOption) (*scheduler.Task, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "StatTask", varargs...)
	ret0, _ := ret[0].(*scheduler.Task)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// StatTask indicates an expected call of StatTask.
func (mr *MockV1MockRecorder) StatTask(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StatTask", reflect.TypeOf((*MockV1)(nil).StatTask), varargs...)
}
