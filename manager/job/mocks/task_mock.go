// Code generated by MockGen. DO NOT EDIT.
// Source: task.go
//
// Generated by this command:
//
//	mockgen -destination mocks/task_mock.go -source task.go -package mocks
//

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	job "d7y.io/dragonfly/v2/internal/job"
	models "d7y.io/dragonfly/v2/manager/models"
	types "d7y.io/dragonfly/v2/manager/types"
	gomock "go.uber.org/mock/gomock"
)

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

// CreateDeleteTask mocks base method.
func (m *MockTask) CreateDeleteTask(arg0 context.Context, arg1 []models.Scheduler, arg2 types.DeleteTaskArgs) (*job.GroupJobState, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateDeleteTask", arg0, arg1, arg2)
	ret0, _ := ret[0].(*job.GroupJobState)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateDeleteTask indicates an expected call of CreateDeleteTask.
func (mr *MockTaskMockRecorder) CreateDeleteTask(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateDeleteTask", reflect.TypeOf((*MockTask)(nil).CreateDeleteTask), arg0, arg1, arg2)
}

// CreateGetTask mocks base method.
func (m *MockTask) CreateGetTask(arg0 context.Context, arg1 []models.Scheduler, arg2 types.GetTaskArgs) (*job.GroupJobState, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateGetTask", arg0, arg1, arg2)
	ret0, _ := ret[0].(*job.GroupJobState)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateGetTask indicates an expected call of CreateGetTask.
func (mr *MockTaskMockRecorder) CreateGetTask(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateGetTask", reflect.TypeOf((*MockTask)(nil).CreateGetTask), arg0, arg1, arg2)
}
