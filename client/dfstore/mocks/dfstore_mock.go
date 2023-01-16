// Code generated by MockGen. DO NOT EDIT.
// Source: dfstore.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	io "io"
	http "net/http"
	reflect "reflect"

	dfstore "d7y.io/dragonfly/v2/client/dfstore"
	objectstorage "d7y.io/dragonfly/v2/pkg/objectstorage"
	gomock "github.com/golang/mock/gomock"
)

// MockDfstore is a mock of Dfstore interface.
type MockDfstore struct {
	ctrl     *gomock.Controller
	recorder *MockDfstoreMockRecorder
}

// MockDfstoreMockRecorder is the mock recorder for MockDfstore.
type MockDfstoreMockRecorder struct {
	mock *MockDfstore
}

// NewMockDfstore creates a new mock instance.
func NewMockDfstore(ctrl *gomock.Controller) *MockDfstore {
	mock := &MockDfstore{ctrl: ctrl}
	mock.recorder = &MockDfstoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDfstore) EXPECT() *MockDfstoreMockRecorder {
	return m.recorder
}

// DeleteObjectRequestWithContext mocks base method.
func (m *MockDfstore) DeleteObjectRequestWithContext(ctx context.Context, input *dfstore.DeleteObjectInput) (*http.Request, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteObjectRequestWithContext", ctx, input)
	ret0, _ := ret[0].(*http.Request)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteObjectRequestWithContext indicates an expected call of DeleteObjectRequestWithContext.
func (mr *MockDfstoreMockRecorder) DeleteObjectRequestWithContext(ctx, input interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteObjectRequestWithContext", reflect.TypeOf((*MockDfstore)(nil).DeleteObjectRequestWithContext), ctx, input)
}

// DeleteObjectWithContext mocks base method.
func (m *MockDfstore) DeleteObjectWithContext(ctx context.Context, input *dfstore.DeleteObjectInput) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteObjectWithContext", ctx, input)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteObjectWithContext indicates an expected call of DeleteObjectWithContext.
func (mr *MockDfstoreMockRecorder) DeleteObjectWithContext(ctx, input interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteObjectWithContext", reflect.TypeOf((*MockDfstore)(nil).DeleteObjectWithContext), ctx, input)
}

// GetObjectMetadataRequestWithContext mocks base method.
func (m *MockDfstore) GetObjectMetadataRequestWithContext(ctx context.Context, input *dfstore.GetObjectMetadataInput) (*http.Request, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetObjectMetadataRequestWithContext", ctx, input)
	ret0, _ := ret[0].(*http.Request)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetObjectMetadataRequestWithContext indicates an expected call of GetObjectMetadataRequestWithContext.
func (mr *MockDfstoreMockRecorder) GetObjectMetadataRequestWithContext(ctx, input interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetObjectMetadataRequestWithContext", reflect.TypeOf((*MockDfstore)(nil).GetObjectMetadataRequestWithContext), ctx, input)
}

// GetObjectMetadataWithContext mocks base method.
func (m *MockDfstore) GetObjectMetadataWithContext(ctx context.Context, input *dfstore.GetObjectMetadataInput) (*objectstorage.ObjectMetadata, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetObjectMetadataWithContext", ctx, input)
	ret0, _ := ret[0].(*objectstorage.ObjectMetadata)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetObjectMetadataWithContext indicates an expected call of GetObjectMetadataWithContext.
func (mr *MockDfstoreMockRecorder) GetObjectMetadataWithContext(ctx, input interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetObjectMetadataWithContext", reflect.TypeOf((*MockDfstore)(nil).GetObjectMetadataWithContext), ctx, input)
}

// GetObjectRequestWithContext mocks base method.
func (m *MockDfstore) GetObjectRequestWithContext(ctx context.Context, input *dfstore.GetObjectInput) (*http.Request, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetObjectRequestWithContext", ctx, input)
	ret0, _ := ret[0].(*http.Request)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetObjectRequestWithContext indicates an expected call of GetObjectRequestWithContext.
func (mr *MockDfstoreMockRecorder) GetObjectRequestWithContext(ctx, input interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetObjectRequestWithContext", reflect.TypeOf((*MockDfstore)(nil).GetObjectRequestWithContext), ctx, input)
}

// GetObjectWithContext mocks base method.
func (m *MockDfstore) GetObjectWithContext(ctx context.Context, input *dfstore.GetObjectInput) (io.ReadCloser, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetObjectWithContext", ctx, input)
	ret0, _ := ret[0].(io.ReadCloser)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetObjectWithContext indicates an expected call of GetObjectWithContext.
func (mr *MockDfstoreMockRecorder) GetObjectWithContext(ctx, input interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetObjectWithContext", reflect.TypeOf((*MockDfstore)(nil).GetObjectWithContext), ctx, input)
}

// IsObjectExistRequestWithContext mocks base method.
func (m *MockDfstore) IsObjectExistRequestWithContext(ctx context.Context, input *dfstore.IsObjectExistInput) (*http.Request, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsObjectExistRequestWithContext", ctx, input)
	ret0, _ := ret[0].(*http.Request)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsObjectExistRequestWithContext indicates an expected call of IsObjectExistRequestWithContext.
func (mr *MockDfstoreMockRecorder) IsObjectExistRequestWithContext(ctx, input interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsObjectExistRequestWithContext", reflect.TypeOf((*MockDfstore)(nil).IsObjectExistRequestWithContext), ctx, input)
}

// IsObjectExistWithContext mocks base method.
func (m *MockDfstore) IsObjectExistWithContext(ctx context.Context, input *dfstore.IsObjectExistInput) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsObjectExistWithContext", ctx, input)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsObjectExistWithContext indicates an expected call of IsObjectExistWithContext.
func (mr *MockDfstoreMockRecorder) IsObjectExistWithContext(ctx, input interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsObjectExistWithContext", reflect.TypeOf((*MockDfstore)(nil).IsObjectExistWithContext), ctx, input)
}

// PutObjectRequestWithContext mocks base method.
func (m *MockDfstore) PutObjectRequestWithContext(ctx context.Context, input *dfstore.PutObjectInput) (*http.Request, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PutObjectRequestWithContext", ctx, input)
	ret0, _ := ret[0].(*http.Request)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PutObjectRequestWithContext indicates an expected call of PutObjectRequestWithContext.
func (mr *MockDfstoreMockRecorder) PutObjectRequestWithContext(ctx, input interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PutObjectRequestWithContext", reflect.TypeOf((*MockDfstore)(nil).PutObjectRequestWithContext), ctx, input)
}

// PutObjectWithContext mocks base method.
func (m *MockDfstore) PutObjectWithContext(ctx context.Context, input *dfstore.PutObjectInput) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PutObjectWithContext", ctx, input)
	ret0, _ := ret[0].(error)
	return ret0
}

// PutObjectWithContext indicates an expected call of PutObjectWithContext.
func (mr *MockDfstoreMockRecorder) PutObjectWithContext(ctx, input interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PutObjectWithContext", reflect.TypeOf((*MockDfstore)(nil).PutObjectWithContext), ctx, input)
}
