// Code generated by MockGen. DO NOT EDIT.
// Source: task_manager.go
//
// Generated by this command:
//
//	mockgen -destination task_manager_mock.go -source task_manager.go -package persistentcache
//

// Package persistentcache is a generated GoMock package.
package persistentcache

import (
	context "context"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockTaskManager is a mock of TaskManager interface.
type MockTaskManager struct {
	ctrl     *gomock.Controller
	recorder *MockTaskManagerMockRecorder
	isgomock struct{}
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

// Delete mocks base method.
func (m *MockTaskManager) Delete(arg0 context.Context, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockTaskManagerMockRecorder) Delete(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockTaskManager)(nil).Delete), arg0, arg1)
}

// Load mocks base method.
func (m *MockTaskManager) Load(arg0 context.Context, arg1 string) (*Task, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Load", arg0, arg1)
	ret0, _ := ret[0].(*Task)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// Load indicates an expected call of Load.
func (mr *MockTaskManagerMockRecorder) Load(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Load", reflect.TypeOf((*MockTaskManager)(nil).Load), arg0, arg1)
}

// LoadAll mocks base method.
func (m *MockTaskManager) LoadAll(arg0 context.Context) ([]*Task, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LoadAll", arg0)
	ret0, _ := ret[0].([]*Task)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LoadAll indicates an expected call of LoadAll.
func (mr *MockTaskManagerMockRecorder) LoadAll(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LoadAll", reflect.TypeOf((*MockTaskManager)(nil).LoadAll), arg0)
}

// LoadCorrentReplicaCount mocks base method.
func (m *MockTaskManager) LoadCorrentReplicaCount(arg0 context.Context, arg1 string) (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LoadCorrentReplicaCount", arg0, arg1)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LoadCorrentReplicaCount indicates an expected call of LoadCorrentReplicaCount.
func (mr *MockTaskManagerMockRecorder) LoadCorrentReplicaCount(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LoadCorrentReplicaCount", reflect.TypeOf((*MockTaskManager)(nil).LoadCorrentReplicaCount), arg0, arg1)
}

// LoadCurrentPersistentReplicaCount mocks base method.
func (m *MockTaskManager) LoadCurrentPersistentReplicaCount(arg0 context.Context, arg1 string) (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LoadCurrentPersistentReplicaCount", arg0, arg1)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LoadCurrentPersistentReplicaCount indicates an expected call of LoadCurrentPersistentReplicaCount.
func (mr *MockTaskManagerMockRecorder) LoadCurrentPersistentReplicaCount(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LoadCurrentPersistentReplicaCount", reflect.TypeOf((*MockTaskManager)(nil).LoadCurrentPersistentReplicaCount), arg0, arg1)
}

// Store mocks base method.
func (m *MockTaskManager) Store(arg0 context.Context, arg1 *Task) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Store", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Store indicates an expected call of Store.
func (mr *MockTaskManagerMockRecorder) Store(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Store", reflect.TypeOf((*MockTaskManager)(nil).Store), arg0, arg1)
}
