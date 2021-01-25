// Code generated by MockGen. DO NOT EDIT.
// Source: ../../../storage/storage_manager.go

// Package mock_storage is a generated GoMock package.
package mock_storage

import (
	context "context"
	storage "github.com/dragonflyoss/Dragonfly2/client/daemon/storage"
	base "github.com/dragonflyoss/Dragonfly2/pkg/rpc/base"
	gomock "github.com/golang/mock/gomock"
	io "io"
	reflect "reflect"
)

// MockTaskStorageDriver is a mock of TaskStorageDriver interface
type MockTaskStorageDriver struct {
	ctrl     *gomock.Controller
	recorder *MockTaskStorageDriverMockRecorder
}

// MockTaskStorageDriverMockRecorder is the mock recorder for MockTaskStorageDriver
type MockTaskStorageDriverMockRecorder struct {
	mock *MockTaskStorageDriver
}

// NewMockTaskStorageDriver creates a new mock instance
func NewMockTaskStorageDriver(ctrl *gomock.Controller) *MockTaskStorageDriver {
	mock := &MockTaskStorageDriver{ctrl: ctrl}
	mock.recorder = &MockTaskStorageDriverMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockTaskStorageDriver) EXPECT() *MockTaskStorageDriverMockRecorder {
	return m.recorder
}

// WritePiece mocks base method
func (m *MockTaskStorageDriver) WritePiece(ctx context.Context, req *storage.WritePieceRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WritePiece", ctx, req)
	ret0, _ := ret[0].(error)
	return ret0
}

// WritePiece indicates an expected call of WritePiece
func (mr *MockTaskStorageDriverMockRecorder) WritePiece(ctx, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WritePiece", reflect.TypeOf((*MockTaskStorageDriver)(nil).WritePiece), ctx, req)
}

// ReadPiece mocks base method
func (m *MockTaskStorageDriver) ReadPiece(ctx context.Context, req *storage.ReadPieceRequest) (io.Reader, io.Closer, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadPiece", ctx, req)
	ret0, _ := ret[0].(io.Reader)
	ret1, _ := ret[1].(io.Closer)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// ReadPiece indicates an expected call of ReadPiece
func (mr *MockTaskStorageDriverMockRecorder) ReadPiece(ctx, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadPiece", reflect.TypeOf((*MockTaskStorageDriver)(nil).ReadPiece), ctx, req)
}

// GetPieces mocks base method
func (m *MockTaskStorageDriver) GetPieces(ctx context.Context, req *base.PieceTaskRequest) (*base.PiecePacket, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPieces", ctx, req)
	ret0, _ := ret[0].(*base.PiecePacket)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPieces indicates an expected call of GetPieces
func (mr *MockTaskStorageDriverMockRecorder) GetPieces(ctx, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPieces", reflect.TypeOf((*MockTaskStorageDriver)(nil).GetPieces), ctx, req)
}

// Store mocks base method
func (m *MockTaskStorageDriver) Store(ctx context.Context, req *storage.StoreRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Store", ctx, req)
	ret0, _ := ret[0].(error)
	return ret0
}

// Store indicates an expected call of Store
func (mr *MockTaskStorageDriverMockRecorder) Store(ctx, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Store", reflect.TypeOf((*MockTaskStorageDriver)(nil).Store), ctx, req)
}

// MockManager is a mock of Manager interface
type MockManager struct {
	ctrl     *gomock.Controller
	recorder *MockManagerMockRecorder
}

// MockManagerMockRecorder is the mock recorder for MockManager
type MockManagerMockRecorder struct {
	mock *MockManager
}

// NewMockManager creates a new mock instance
func NewMockManager(ctrl *gomock.Controller) *MockManager {
	mock := &MockManager{ctrl: ctrl}
	mock.recorder = &MockManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockManager) EXPECT() *MockManagerMockRecorder {
	return m.recorder
}

// WritePiece mocks base method
func (m *MockManager) WritePiece(ctx context.Context, req *storage.WritePieceRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WritePiece", ctx, req)
	ret0, _ := ret[0].(error)
	return ret0
}

// WritePiece indicates an expected call of WritePiece
func (mr *MockManagerMockRecorder) WritePiece(ctx, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WritePiece", reflect.TypeOf((*MockManager)(nil).WritePiece), ctx, req)
}

// ReadPiece mocks base method
func (m *MockManager) ReadPiece(ctx context.Context, req *storage.ReadPieceRequest) (io.Reader, io.Closer, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadPiece", ctx, req)
	ret0, _ := ret[0].(io.Reader)
	ret1, _ := ret[1].(io.Closer)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// ReadPiece indicates an expected call of ReadPiece
func (mr *MockManagerMockRecorder) ReadPiece(ctx, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadPiece", reflect.TypeOf((*MockManager)(nil).ReadPiece), ctx, req)
}

// GetPieces mocks base method
func (m *MockManager) GetPieces(ctx context.Context, req *base.PieceTaskRequest) (*base.PiecePacket, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPieces", ctx, req)
	ret0, _ := ret[0].(*base.PiecePacket)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPieces indicates an expected call of GetPieces
func (mr *MockManagerMockRecorder) GetPieces(ctx, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPieces", reflect.TypeOf((*MockManager)(nil).GetPieces), ctx, req)
}

// Store mocks base method
func (m *MockManager) Store(ctx context.Context, req *storage.StoreRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Store", ctx, req)
	ret0, _ := ret[0].(error)
	return ret0
}

// Store indicates an expected call of Store
func (mr *MockManagerMockRecorder) Store(ctx, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Store", reflect.TypeOf((*MockManager)(nil).Store), ctx, req)
}

// RegisterTask mocks base method
func (m *MockManager) RegisterTask(ctx context.Context, req storage.RegisterTaskRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RegisterTask", ctx, req)
	ret0, _ := ret[0].(error)
	return ret0
}

// RegisterTask indicates an expected call of RegisterTask
func (mr *MockManagerMockRecorder) RegisterTask(ctx, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterTask", reflect.TypeOf((*MockManager)(nil).RegisterTask), ctx, req)
}

// MockTaskStorageExecutor is a mock of TaskStorageExecutor interface
type MockTaskStorageExecutor struct {
	ctrl     *gomock.Controller
	recorder *MockTaskStorageExecutorMockRecorder
}

// MockTaskStorageExecutorMockRecorder is the mock recorder for MockTaskStorageExecutor
type MockTaskStorageExecutorMockRecorder struct {
	mock *MockTaskStorageExecutor
}

// NewMockTaskStorageExecutor creates a new mock instance
func NewMockTaskStorageExecutor(ctrl *gomock.Controller) *MockTaskStorageExecutor {
	mock := &MockTaskStorageExecutor{ctrl: ctrl}
	mock.recorder = &MockTaskStorageExecutorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockTaskStorageExecutor) EXPECT() *MockTaskStorageExecutorMockRecorder {
	return m.recorder
}

// TryGC mocks base method
func (m *MockTaskStorageExecutor) TryGC() (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TryGC")
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// TryGC indicates an expected call of TryGC
func (mr *MockTaskStorageExecutorMockRecorder) TryGC() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TryGC", reflect.TypeOf((*MockTaskStorageExecutor)(nil).TryGC))
}

// LoadTask mocks base method
func (m *MockTaskStorageExecutor) LoadTask(meta storage.PeerTaskMetaData) (storage.TaskStorageDriver, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LoadTask", meta)
	ret0, _ := ret[0].(storage.TaskStorageDriver)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// LoadTask indicates an expected call of LoadTask
func (mr *MockTaskStorageExecutorMockRecorder) LoadTask(meta interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LoadTask", reflect.TypeOf((*MockTaskStorageExecutor)(nil).LoadTask), meta)
}

// CreateTask mocks base method
func (m *MockTaskStorageExecutor) CreateTask(request storage.RegisterTaskRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateTask", request)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateTask indicates an expected call of CreateTask
func (mr *MockTaskStorageExecutorMockRecorder) CreateTask(request interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateTask", reflect.TypeOf((*MockTaskStorageExecutor)(nil).CreateTask), request)
}

// ReloadPersistentTask mocks base method
func (m *MockTaskStorageExecutor) ReloadPersistentTask(gcCallback storage.GCCallback) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReloadPersistentTask", gcCallback)
	ret0, _ := ret[0].(error)
	return ret0
}

// ReloadPersistentTask indicates an expected call of ReloadPersistentTask
func (mr *MockTaskStorageExecutorMockRecorder) ReloadPersistentTask(gcCallback interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReloadPersistentTask", reflect.TypeOf((*MockTaskStorageExecutor)(nil).ReloadPersistentTask), gcCallback)
}
