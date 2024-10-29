// Code generated by MockGen. DO NOT EDIT.
// Source: proxy_manager.go
//
// Generated by this command:
//
//	mockgen -destination mocks/proxy_manager_mock.go -source proxy_manager.go -package mocks
//

// Package mocks is a generated GoMock package.
package mocks

import (
	net "net"
	reflect "reflect"

	config "d7y.io/dragonfly/v2/client/config"
	gomock "go.uber.org/mock/gomock"
)

// MockManager is a mock of Manager interface.
type MockManager struct {
	ctrl     *gomock.Controller
	recorder *MockManagerMockRecorder
	isgomock struct{}
}

// MockManagerMockRecorder is the mock recorder for MockManager.
type MockManagerMockRecorder struct {
	mock *MockManager
}

// NewMockManager creates a new mock instance.
func NewMockManager(ctrl *gomock.Controller) *MockManager {
	mock := &MockManager{ctrl: ctrl}
	mock.recorder = &MockManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockManager) EXPECT() *MockManagerMockRecorder {
	return m.recorder
}

// IsEnabled mocks base method.
func (m *MockManager) IsEnabled() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsEnabled")
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsEnabled indicates an expected call of IsEnabled.
func (mr *MockManagerMockRecorder) IsEnabled() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsEnabled", reflect.TypeOf((*MockManager)(nil).IsEnabled))
}

// Serve mocks base method.
func (m *MockManager) Serve(arg0 net.Listener) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Serve", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Serve indicates an expected call of Serve.
func (mr *MockManagerMockRecorder) Serve(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Serve", reflect.TypeOf((*MockManager)(nil).Serve), arg0)
}

// ServeSNI mocks base method.
func (m *MockManager) ServeSNI(arg0 net.Listener) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ServeSNI", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// ServeSNI indicates an expected call of ServeSNI.
func (mr *MockManagerMockRecorder) ServeSNI(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ServeSNI", reflect.TypeOf((*MockManager)(nil).ServeSNI), arg0)
}

// Stop mocks base method.
func (m *MockManager) Stop() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Stop")
	ret0, _ := ret[0].(error)
	return ret0
}

// Stop indicates an expected call of Stop.
func (mr *MockManagerMockRecorder) Stop() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stop", reflect.TypeOf((*MockManager)(nil).Stop))
}

// Watch mocks base method.
func (m *MockManager) Watch(arg0 *config.ProxyOption) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Watch", arg0)
}

// Watch indicates an expected call of Watch.
func (mr *MockManagerMockRecorder) Watch(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Watch", reflect.TypeOf((*MockManager)(nil).Watch), arg0)
}

// MockConfigWatcher is a mock of ConfigWatcher interface.
type MockConfigWatcher struct {
	ctrl     *gomock.Controller
	recorder *MockConfigWatcherMockRecorder
	isgomock struct{}
}

// MockConfigWatcherMockRecorder is the mock recorder for MockConfigWatcher.
type MockConfigWatcherMockRecorder struct {
	mock *MockConfigWatcher
}

// NewMockConfigWatcher creates a new mock instance.
func NewMockConfigWatcher(ctrl *gomock.Controller) *MockConfigWatcher {
	mock := &MockConfigWatcher{ctrl: ctrl}
	mock.recorder = &MockConfigWatcherMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockConfigWatcher) EXPECT() *MockConfigWatcherMockRecorder {
	return m.recorder
}

// Watch mocks base method.
func (m *MockConfigWatcher) Watch(arg0 *config.ProxyOption) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Watch", arg0)
}

// Watch indicates an expected call of Watch.
func (mr *MockConfigWatcherMockRecorder) Watch(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Watch", reflect.TypeOf((*MockConfigWatcher)(nil).Watch), arg0)
}
