// Code generated by MockGen. DO NOT EDIT.
// Source: preheat.go
//
// Generated by this command:
//
//	mockgen -destination mocks/preheat_mock.go -source preheat.go -package mocks
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

// MockPreheat is a mock of Preheat interface.
type MockPreheat struct {
	ctrl     *gomock.Controller
	recorder *MockPreheatMockRecorder
}

// MockPreheatMockRecorder is the mock recorder for MockPreheat.
type MockPreheatMockRecorder struct {
	mock *MockPreheat
}

// NewMockPreheat creates a new mock instance.
func NewMockPreheat(ctrl *gomock.Controller) *MockPreheat {
	mock := &MockPreheat{ctrl: ctrl}
	mock.recorder = &MockPreheatMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPreheat) EXPECT() *MockPreheatMockRecorder {
	return m.recorder
}

// CreatePreheat mocks base method.
func (m *MockPreheat) CreatePreheat(arg0 context.Context, arg1 []models.Scheduler, arg2 types.PreheatArgs) (*job.GroupJobState, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreatePreheat", arg0, arg1, arg2)
	ret0, _ := ret[0].(*job.GroupJobState)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreatePreheat indicates an expected call of CreatePreheat.
func (mr *MockPreheatMockRecorder) CreatePreheat(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreatePreheat", reflect.TypeOf((*MockPreheat)(nil).CreatePreheat), arg0, arg1, arg2)
}
