// Code generated by MockGen. DO NOT EDIT.
// Source: source_client.go

// Package mocks is a generated GoMock package.
package mocks

import (
	url "net/url"
	reflect "reflect"

	source "d7y.io/dragonfly/v2/pkg/source"
	gomock "github.com/golang/mock/gomock"
)

// MockResourceClient is a mock of ResourceClient interface.
type MockResourceClient struct {
	ctrl     *gomock.Controller
	recorder *MockResourceClientMockRecorder
}

// MockResourceClientMockRecorder is the mock recorder for MockResourceClient.
type MockResourceClientMockRecorder struct {
	mock *MockResourceClient
}

// NewMockResourceClient creates a new mock instance.
func NewMockResourceClient(ctrl *gomock.Controller) *MockResourceClient {
	mock := &MockResourceClient{ctrl: ctrl}
	mock.recorder = &MockResourceClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockResourceClient) EXPECT() *MockResourceClientMockRecorder {
	return m.recorder
}

// Download mocks base method.
func (m *MockResourceClient) Download(request *source.Request) (*source.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Download", request)
	ret0, _ := ret[0].(*source.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Download indicates an expected call of Download.
func (mr *MockResourceClientMockRecorder) Download(request interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Download", reflect.TypeOf((*MockResourceClient)(nil).Download), request)
}

// GetContentLength mocks base method.
func (m *MockResourceClient) GetContentLength(request *source.Request) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetContentLength", request)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetContentLength indicates an expected call of GetContentLength.
func (mr *MockResourceClientMockRecorder) GetContentLength(request interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetContentLength", reflect.TypeOf((*MockResourceClient)(nil).GetContentLength), request)
}

// GetLastModified mocks base method.
func (m *MockResourceClient) GetLastModified(request *source.Request) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLastModified", request)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLastModified indicates an expected call of GetLastModified.
func (mr *MockResourceClientMockRecorder) GetLastModified(request interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLastModified", reflect.TypeOf((*MockResourceClient)(nil).GetLastModified), request)
}

// IsExpired mocks base method.
func (m *MockResourceClient) IsExpired(request *source.Request, info *source.ExpireInfo) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsExpired", request, info)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsExpired indicates an expected call of IsExpired.
func (mr *MockResourceClientMockRecorder) IsExpired(request, info interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsExpired", reflect.TypeOf((*MockResourceClient)(nil).IsExpired), request, info)
}

// IsSupportRange mocks base method.
func (m *MockResourceClient) IsSupportRange(request *source.Request) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsSupportRange", request)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsSupportRange indicates an expected call of IsSupportRange.
func (mr *MockResourceClientMockRecorder) IsSupportRange(request interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsSupportRange", reflect.TypeOf((*MockResourceClient)(nil).IsSupportRange), request)
}

// MockResourceMetadataGetter is a mock of ResourceMetadataGetter interface.
type MockResourceMetadataGetter struct {
	ctrl     *gomock.Controller
	recorder *MockResourceMetadataGetterMockRecorder
}

// MockResourceMetadataGetterMockRecorder is the mock recorder for MockResourceMetadataGetter.
type MockResourceMetadataGetterMockRecorder struct {
	mock *MockResourceMetadataGetter
}

// NewMockResourceMetadataGetter creates a new mock instance.
func NewMockResourceMetadataGetter(ctrl *gomock.Controller) *MockResourceMetadataGetter {
	mock := &MockResourceMetadataGetter{ctrl: ctrl}
	mock.recorder = &MockResourceMetadataGetterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockResourceMetadataGetter) EXPECT() *MockResourceMetadataGetterMockRecorder {
	return m.recorder
}

// GetMetadata mocks base method.
func (m *MockResourceMetadataGetter) GetMetadata(request *source.Request) (*source.Metadata, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMetadata", request)
	ret0, _ := ret[0].(*source.Metadata)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMetadata indicates an expected call of GetMetadata.
func (mr *MockResourceMetadataGetterMockRecorder) GetMetadata(request interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMetadata", reflect.TypeOf((*MockResourceMetadataGetter)(nil).GetMetadata), request)
}

// MockResourceLister is a mock of ResourceLister interface.
type MockResourceLister struct {
	ctrl     *gomock.Controller
	recorder *MockResourceListerMockRecorder
}

// MockResourceListerMockRecorder is the mock recorder for MockResourceLister.
type MockResourceListerMockRecorder struct {
	mock *MockResourceLister
}

// NewMockResourceLister creates a new mock instance.
func NewMockResourceLister(ctrl *gomock.Controller) *MockResourceLister {
	mock := &MockResourceLister{ctrl: ctrl}
	mock.recorder = &MockResourceListerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockResourceLister) EXPECT() *MockResourceListerMockRecorder {
	return m.recorder
}

// List mocks base method.
func (m *MockResourceLister) List(request *source.Request) ([]*url.URL, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", request)
	ret0, _ := ret[0].([]*url.URL)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockResourceListerMockRecorder) List(request interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockResourceLister)(nil).List), request)
}

// MockClientManager is a mock of ClientManager interface.
type MockClientManager struct {
	ctrl     *gomock.Controller
	recorder *MockClientManagerMockRecorder
}

// MockClientManagerMockRecorder is the mock recorder for MockClientManager.
type MockClientManagerMockRecorder struct {
	mock *MockClientManager
}

// NewMockClientManager creates a new mock instance.
func NewMockClientManager(ctrl *gomock.Controller) *MockClientManager {
	mock := &MockClientManager{ctrl: ctrl}
	mock.recorder = &MockClientManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockClientManager) EXPECT() *MockClientManagerMockRecorder {
	return m.recorder
}

// GetClient mocks base method.
func (m *MockClientManager) GetClient(scheme string, options ...source.Option) (source.ResourceClient, bool) {
	m.ctrl.T.Helper()
	varargs := []interface{}{scheme}
	for _, a := range options {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetClient", varargs...)
	ret0, _ := ret[0].(source.ResourceClient)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// GetClient indicates an expected call of GetClient.
func (mr *MockClientManagerMockRecorder) GetClient(scheme interface{}, options ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{scheme}, options...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetClient", reflect.TypeOf((*MockClientManager)(nil).GetClient), varargs...)
}

// ListClients mocks base method.
func (m *MockClientManager) ListClients() []string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListClients")
	ret0, _ := ret[0].([]string)
	return ret0
}

// ListClients indicates an expected call of ListClients.
func (mr *MockClientManagerMockRecorder) ListClients() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListClients", reflect.TypeOf((*MockClientManager)(nil).ListClients))
}

// Register mocks base method.
func (m *MockClientManager) Register(scheme string, resourceClient source.ResourceClient, adapter source.RequestAdapter, hook ...source.Hook) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{scheme, resourceClient, adapter}
	for _, a := range hook {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Register", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// Register indicates an expected call of Register.
func (mr *MockClientManagerMockRecorder) Register(scheme, resourceClient, adapter interface{}, hook ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{scheme, resourceClient, adapter}, hook...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Register", reflect.TypeOf((*MockClientManager)(nil).Register), varargs...)
}

// UnRegister mocks base method.
func (m *MockClientManager) UnRegister(scheme string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "UnRegister", scheme)
}

// UnRegister indicates an expected call of UnRegister.
func (mr *MockClientManagerMockRecorder) UnRegister(scheme interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UnRegister", reflect.TypeOf((*MockClientManager)(nil).UnRegister), scheme)
}

// MockHook is a mock of Hook interface.
type MockHook struct {
	ctrl     *gomock.Controller
	recorder *MockHookMockRecorder
}

// MockHookMockRecorder is the mock recorder for MockHook.
type MockHookMockRecorder struct {
	mock *MockHook
}

// NewMockHook creates a new mock instance.
func NewMockHook(ctrl *gomock.Controller) *MockHook {
	mock := &MockHook{ctrl: ctrl}
	mock.recorder = &MockHookMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockHook) EXPECT() *MockHookMockRecorder {
	return m.recorder
}

// AfterResponse mocks base method.
func (m *MockHook) AfterResponse(response *source.Response) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AfterResponse", response)
	ret0, _ := ret[0].(error)
	return ret0
}

// AfterResponse indicates an expected call of AfterResponse.
func (mr *MockHookMockRecorder) AfterResponse(response interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AfterResponse", reflect.TypeOf((*MockHook)(nil).AfterResponse), response)
}

// BeforeRequest mocks base method.
func (m *MockHook) BeforeRequest(request *source.Request) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BeforeRequest", request)
	ret0, _ := ret[0].(error)
	return ret0
}

// BeforeRequest indicates an expected call of BeforeRequest.
func (mr *MockHookMockRecorder) BeforeRequest(request interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BeforeRequest", reflect.TypeOf((*MockHook)(nil).BeforeRequest), request)
}
