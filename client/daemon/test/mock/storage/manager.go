// Code generated by MockGen. DO NOT EDIT.
// Source: ../../../storage/storage_manager.go

// Package mock_storage is a generated GoMock package.
package mock_storage

import (
	context "context"
	io "io"
	reflect "reflect"
	time "time"

	storage "d7y.io/dragonfly/v2/client/daemon/storage"
	base "d7y.io/dragonfly/v2/pkg/rpc/base"
	gomock "github.com/golang/mock/gomock"
)

// MockTaskStorageDriver is a mock of TaskStorageDriver interface.
type MockTaskStorageDriver struct {
	ctrl     *gomock.Controller
	recorder *MockTaskStorageDriverMockRecorder
}

// MockTaskStorageDriverMockRecorder is the mock recorder for MockTaskStorageDriver.
type MockTaskStorageDriverMockRecorder struct {
	mock *MockTaskStorageDriver
}

// NewMockTaskStorageDriver creates a new mock instance.
func NewMockTaskStorageDriver(ctrl *gomock.Controller) *MockTaskStorageDriver {
	mock := &MockTaskStorageDriver{ctrl: ctrl}
	mock.recorder = &MockTaskStorageDriverMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTaskStorageDriver) EXPECT() *MockTaskStorageDriverMockRecorder {
	return m.recorder
}

// GetPieces mocks base method.
func (m *MockTaskStorageDriver) GetPieces(ctx context.Context, req *base.PieceTaskRequest) (*base.PiecePacket, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPieces", ctx, req)
	ret0, _ := ret[0].(*base.PiecePacket)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPieces indicates an expected call of GetPieces.
func (mr *MockTaskStorageDriverMockRecorder) GetPieces(ctx, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPieces", reflect.TypeOf((*MockTaskStorageDriver)(nil).GetPieces), ctx, req)
}

// ReadPiece mocks base method.
func (m *MockTaskStorageDriver) ReadPiece(ctx context.Context, req *storage.ReadPieceRequest) (io.Reader, io.Closer, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadPiece", ctx, req)
	ret0, _ := ret[0].(io.Reader)
	ret1, _ := ret[1].(io.Closer)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// ReadPiece indicates an expected call of ReadPiece.
func (mr *MockTaskStorageDriverMockRecorder) ReadPiece(ctx, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadPiece", reflect.TypeOf((*MockTaskStorageDriver)(nil).ReadPiece), ctx, req)
}

// Store mocks base method.
func (m *MockTaskStorageDriver) Store(ctx context.Context, req *storage.StoreRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Store", ctx, req)
	ret0, _ := ret[0].(error)
	return ret0
}

// Store indicates an expected call of Store.
func (mr *MockTaskStorageDriverMockRecorder) Store(ctx, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Store", reflect.TypeOf((*MockTaskStorageDriver)(nil).Store), ctx, req)
}

// UpdateTask mocks base method.
func (m *MockTaskStorageDriver) UpdateTask(ctx context.Context, req *storage.UpdateTaskRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateTask", ctx, req)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateTask indicates an expected call of UpdateTask.
func (mr *MockTaskStorageDriverMockRecorder) UpdateTask(ctx, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateTask", reflect.TypeOf((*MockTaskStorageDriver)(nil).UpdateTask), ctx, req)
}

// WritePiece mocks base method.
func (m *MockTaskStorageDriver) WritePiece(ctx context.Context, req *storage.WritePieceRequest) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WritePiece", ctx, req)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// WritePiece indicates an expected call of WritePiece.
func (mr *MockTaskStorageDriverMockRecorder) WritePiece(ctx, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WritePiece", reflect.TypeOf((*MockTaskStorageDriver)(nil).WritePiece), ctx, req)
}

// MockManager is a mock of Manager interface.
type MockManager struct {
	ctrl     *gomock.Controller
	recorder *MockManagerMockRecorder
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

// Alive mocks base method.
func (m *MockManager) Alive(alive time.Duration) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Alive", alive)
	ret0, _ := ret[0].(bool)
	return ret0
}

// Alive indicates an expected call of Alive.
func (mr *MockManagerMockRecorder) Alive(alive interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Alive", reflect.TypeOf((*MockManager)(nil).Alive), alive)
}

// Clean mocks base method.
func (m *MockManager) CleanUp() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "CleanUp")
}

// Clean indicates an expected call of Clean.
func (mr *MockManagerMockRecorder) Clean() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CleanUp", reflect.TypeOf((*MockManager)(nil).CleanUp))
}

// GetPieces mocks base method.
func (m *MockManager) GetPieces(ctx context.Context, req *base.PieceTaskRequest) (*base.PiecePacket, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPieces", ctx, req)
	ret0, _ := ret[0].(*base.PiecePacket)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPieces indicates an expected call of GetPieces.
func (mr *MockManagerMockRecorder) GetPieces(ctx, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPieces", reflect.TypeOf((*MockManager)(nil).GetPieces), ctx, req)
}

// Keep mocks base method.
func (m *MockManager) Keep() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Keep")
}

// Keep indicates an expected call of Keep.
func (mr *MockManagerMockRecorder) Keep() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Keep", reflect.TypeOf((*MockManager)(nil).Keep))
}

// ReadPiece mocks base method.
func (m *MockManager) ReadPiece(ctx context.Context, req *storage.ReadPieceRequest) (io.Reader, io.Closer, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadPiece", ctx, req)
	ret0, _ := ret[0].(io.Reader)
	ret1, _ := ret[1].(io.Closer)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// ReadPiece indicates an expected call of ReadPiece.
func (mr *MockManagerMockRecorder) ReadPiece(ctx, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadPiece", reflect.TypeOf((*MockManager)(nil).ReadPiece), ctx, req)
}

// RegisterTask mocks base method.
func (m *MockManager) RegisterTask(ctx context.Context, req storage.RegisterTaskRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RegisterTask", ctx, req)
	ret0, _ := ret[0].(error)
	return ret0
}

// RegisterTask indicates an expected call of RegisterTask.
func (mr *MockManagerMockRecorder) RegisterTask(ctx, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterTask", reflect.TypeOf((*MockManager)(nil).RegisterTask), ctx, req)
}

// Store mocks base method.
func (m *MockManager) Store(ctx context.Context, req *storage.StoreRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Store", ctx, req)
	ret0, _ := ret[0].(error)
	return ret0
}

// Store indicates an expected call of Store.
func (mr *MockManagerMockRecorder) Store(ctx, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Store", reflect.TypeOf((*MockManager)(nil).Store), ctx, req)
}

// UpdateTask mocks base method.
func (m *MockManager) UpdateTask(ctx context.Context, req *storage.UpdateTaskRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateTask", ctx, req)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateTask indicates an expected call of UpdateTask.
func (mr *MockManagerMockRecorder) UpdateTask(ctx, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateTask", reflect.TypeOf((*MockManager)(nil).UpdateTask), ctx, req)
}

// WritePiece mocks base method.
func (m *MockManager) WritePiece(ctx context.Context, req *storage.WritePieceRequest) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WritePiece", ctx, req)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// WritePiece indicates an expected call of WritePiece.
func (mr *MockManagerMockRecorder) WritePiece(ctx, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WritePiece", reflect.TypeOf((*MockManager)(nil).WritePiece), ctx, req)
}
