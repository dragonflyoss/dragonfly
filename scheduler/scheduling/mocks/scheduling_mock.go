// Code generated by MockGen. DO NOT EDIT.
// Source: scheduling.go
//
// Generated by this command:
//
//	mockgen -destination mocks/scheduling_mock.go -source scheduling.go -package mocks
//

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	set "d7y.io/dragonfly/v2/pkg/container/set"
	persistentcache "d7y.io/dragonfly/v2/scheduler/resource/persistentcache"
	standard "d7y.io/dragonfly/v2/scheduler/resource/standard"
	gomock "go.uber.org/mock/gomock"
)

// MockScheduling is a mock of Scheduling interface.
type MockScheduling struct {
	ctrl     *gomock.Controller
	recorder *MockSchedulingMockRecorder
	isgomock struct{}
}

// MockSchedulingMockRecorder is the mock recorder for MockScheduling.
type MockSchedulingMockRecorder struct {
	mock *MockScheduling
}

// NewMockScheduling creates a new mock instance.
func NewMockScheduling(ctrl *gomock.Controller) *MockScheduling {
	mock := &MockScheduling{ctrl: ctrl}
	mock.recorder = &MockSchedulingMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockScheduling) EXPECT() *MockSchedulingMockRecorder {
	return m.recorder
}

// FindCandidateParents mocks base method.
func (m *MockScheduling) FindCandidateParents(arg0 context.Context, arg1 *standard.Peer, arg2 set.SafeSet[string]) ([]*standard.Peer, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindCandidateParents", arg0, arg1, arg2)
	ret0, _ := ret[0].([]*standard.Peer)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// FindCandidateParents indicates an expected call of FindCandidateParents.
func (mr *MockSchedulingMockRecorder) FindCandidateParents(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindCandidateParents", reflect.TypeOf((*MockScheduling)(nil).FindCandidateParents), arg0, arg1, arg2)
}

// FindCandidatePersistentCacheParents mocks base method.
func (m *MockScheduling) FindCandidatePersistentCacheParents(arg0 context.Context, arg1 *persistentcache.Peer, arg2 set.SafeSet[string]) ([]*persistentcache.Peer, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindCandidatePersistentCacheParents", arg0, arg1, arg2)
	ret0, _ := ret[0].([]*persistentcache.Peer)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// FindCandidatePersistentCacheParents indicates an expected call of FindCandidatePersistentCacheParents.
func (mr *MockSchedulingMockRecorder) FindCandidatePersistentCacheParents(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindCandidatePersistentCacheParents", reflect.TypeOf((*MockScheduling)(nil).FindCandidatePersistentCacheParents), arg0, arg1, arg2)
}

// FindParentAndCandidateParents mocks base method.
func (m *MockScheduling) FindParentAndCandidateParents(arg0 context.Context, arg1 *standard.Peer, arg2 set.SafeSet[string]) ([]*standard.Peer, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindParentAndCandidateParents", arg0, arg1, arg2)
	ret0, _ := ret[0].([]*standard.Peer)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// FindParentAndCandidateParents indicates an expected call of FindParentAndCandidateParents.
func (mr *MockSchedulingMockRecorder) FindParentAndCandidateParents(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindParentAndCandidateParents", reflect.TypeOf((*MockScheduling)(nil).FindParentAndCandidateParents), arg0, arg1, arg2)
}

// FindReplicatePersistentCacheParents mocks base method.
func (m *MockScheduling) FindReplicatePersistentCacheParents(arg0 context.Context, arg1 *persistentcache.Task, arg2 set.SafeSet[string]) ([]*persistentcache.Peer, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindReplicatePersistentCacheParents", arg0, arg1, arg2)
	ret0, _ := ret[0].([]*persistentcache.Peer)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// FindReplicatePersistentCacheParents indicates an expected call of FindReplicatePersistentCacheParents.
func (mr *MockSchedulingMockRecorder) FindReplicatePersistentCacheParents(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindReplicatePersistentCacheParents", reflect.TypeOf((*MockScheduling)(nil).FindReplicatePersistentCacheParents), arg0, arg1, arg2)
}

// FindSuccessParent mocks base method.
func (m *MockScheduling) FindSuccessParent(arg0 context.Context, arg1 *standard.Peer, arg2 set.SafeSet[string]) (*standard.Peer, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindSuccessParent", arg0, arg1, arg2)
	ret0, _ := ret[0].(*standard.Peer)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// FindSuccessParent indicates an expected call of FindSuccessParent.
func (mr *MockSchedulingMockRecorder) FindSuccessParent(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindSuccessParent", reflect.TypeOf((*MockScheduling)(nil).FindSuccessParent), arg0, arg1, arg2)
}

// ScheduleCandidateParents mocks base method.
func (m *MockScheduling) ScheduleCandidateParents(arg0 context.Context, arg1 *standard.Peer, arg2 set.SafeSet[string]) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ScheduleCandidateParents", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// ScheduleCandidateParents indicates an expected call of ScheduleCandidateParents.
func (mr *MockSchedulingMockRecorder) ScheduleCandidateParents(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ScheduleCandidateParents", reflect.TypeOf((*MockScheduling)(nil).ScheduleCandidateParents), arg0, arg1, arg2)
}

// ScheduleParentAndCandidateParents mocks base method.
func (m *MockScheduling) ScheduleParentAndCandidateParents(arg0 context.Context, arg1 *standard.Peer, arg2 set.SafeSet[string]) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ScheduleParentAndCandidateParents", arg0, arg1, arg2)
}

// ScheduleParentAndCandidateParents indicates an expected call of ScheduleParentAndCandidateParents.
func (mr *MockSchedulingMockRecorder) ScheduleParentAndCandidateParents(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ScheduleParentAndCandidateParents", reflect.TypeOf((*MockScheduling)(nil).ScheduleParentAndCandidateParents), arg0, arg1, arg2)
}
