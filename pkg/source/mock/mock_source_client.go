// Code generated by MockGen. DO NOT EDIT.
// Source: d7y.io/dragonfly/v2/pkg/source (interfaces: ResourceClient)

// Package mock is a generated GoMock package.
package mock

import (
	io "io"
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
func (m *MockResourceClient) Download(arg0 *source.Request) (io.ReadCloser, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Download", arg0)
	ret0, _ := ret[0].(io.ReadCloser)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Download indicates an expected call of Download.
func (mr *MockResourceClientMockRecorder) Download(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Download", reflect.TypeOf((*MockResourceClient)(nil).Download), arg0)
}

// DownloadWithExpireInfo mocks base method.
func (m *MockResourceClient) DownloadWithExpireInfo(arg0 *source.Request) (io.ReadCloser, *source.ExpireInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DownloadWithExpireInfo", arg0)
	ret0, _ := ret[0].(io.ReadCloser)
	ret1, _ := ret[1].(*source.ExpireInfo)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// DownloadWithExpireInfo indicates an expected call of DownloadWithExpireInfo.
func (mr *MockResourceClientMockRecorder) DownloadWithExpireInfo(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DownloadWithExpireInfo", reflect.TypeOf((*MockResourceClient)(nil).DownloadWithExpireInfo), arg0)
}

// GetContentLength mocks base method.
func (m *MockResourceClient) GetContentLength(arg0 *source.Request) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetContentLength", arg0)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetContentLength indicates an expected call of GetContentLength.
func (mr *MockResourceClientMockRecorder) GetContentLength(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetContentLength", reflect.TypeOf((*MockResourceClient)(nil).GetContentLength), arg0)
}

// GetLastModifiedMillis mocks base method.
func (m *MockResourceClient) GetLastModifiedMillis(arg0 *source.Request) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLastModifiedMillis", arg0)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLastModifiedMillis indicates an expected call of GetLastModifiedMillis.
func (mr *MockResourceClientMockRecorder) GetLastModifiedMillis(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLastModifiedMillis", reflect.TypeOf((*MockResourceClient)(nil).GetLastModifiedMillis), arg0)
}

// IsExpired mocks base method.
func (m *MockResourceClient) IsExpired(arg0 *source.Request, arg1 *source.ExpireInfo) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsExpired", arg0, arg1)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsExpired indicates an expected call of IsExpired.
func (mr *MockResourceClientMockRecorder) IsExpired(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsExpired", reflect.TypeOf((*MockResourceClient)(nil).IsExpired), arg0, arg1)
}

// IsSupportRange mocks base method.
func (m *MockResourceClient) IsSupportRange(arg0 *source.Request) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsSupportRange", arg0)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsSupportRange indicates an expected call of IsSupportRange.
func (mr *MockResourceClientMockRecorder) IsSupportRange(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsSupportRange", reflect.TypeOf((*MockResourceClient)(nil).IsSupportRange), arg0)
}
