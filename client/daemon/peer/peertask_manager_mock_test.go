/*
 *     Copyright 2020 The Dragonfly Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

// Code generated by MockGen. DO NOT EDIT.
// Source: peertask_manager.go

// Package peer is a generated GoMock package.
package peer

import (
	context "context"
	io "io"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"

	dflog "d7y.io/dragonfly/v2/pkg/dflog"
	base "d7y.io/dragonfly/v2/pkg/rpc/base"
	scheduler "d7y.io/dragonfly/v2/pkg/rpc/scheduler"
)

// MockPeerTaskManager is a mock of PeerTaskManager interface
type MockPeerTaskManager struct {
	ctrl     *gomock.Controller
	recorder *MockPeerTaskManagerMockRecorder
}

// MockPeerTaskManagerMockRecorder is the mock recorder for MockPeerTaskManager
type MockPeerTaskManagerMockRecorder struct {
	mock *MockPeerTaskManager
}

// NewMockPeerTaskManager creates a new mock instance
func NewMockPeerTaskManager(ctrl *gomock.Controller) *MockPeerTaskManager {
	mock := &MockPeerTaskManager{ctrl: ctrl}
	mock.recorder = &MockPeerTaskManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockPeerTaskManager) EXPECT() *MockPeerTaskManagerMockRecorder {
	return m.recorder
}

// StartFilePeerTask mocks base method
func (m *MockPeerTaskManager) StartFilePeerTask(ctx context.Context, req *FilePeerTaskRequest) (chan *FilePeerTaskProgress, error) {
	ret := m.ctrl.Call(m, "StartFilePeerTask", ctx, req)
	ret0, _ := ret[0].(chan *FilePeerTaskProgress)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// StartFilePeerTask indicates an expected call of StartFilePeerTask
func (mr *MockPeerTaskManagerMockRecorder) StartFilePeerTask(ctx, req interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StartFilePeerTask", reflect.TypeOf((*MockPeerTaskManager)(nil).StartFilePeerTask), ctx, req)
}

// StartStreamPeerTask mocks base method
func (m *MockPeerTaskManager) StartStreamPeerTask(ctx context.Context, req *scheduler.PeerTaskRequest) (io.Reader, map[string]string, error) {
	ret := m.ctrl.Call(m, "StartStreamPeerTask", ctx, req)
	ret0, _ := ret[0].(io.Reader)
	ret1, _ := ret[1].(map[string]string)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// StartStreamPeerTask indicates an expected call of StartStreamPeerTask
func (mr *MockPeerTaskManagerMockRecorder) StartStreamPeerTask(ctx, req interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StartStreamPeerTask", reflect.TypeOf((*MockPeerTaskManager)(nil).StartStreamPeerTask), ctx, req)
}

// Stop mocks base method
func (m *MockPeerTaskManager) Stop(ctx context.Context) error {
	ret := m.ctrl.Call(m, "Stop", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// Stop indicates an expected call of Stop
func (mr *MockPeerTaskManagerMockRecorder) Stop(ctx interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stop", reflect.TypeOf((*MockPeerTaskManager)(nil).Stop), ctx)
}

// MockPeerTask is a mock of PeerTask interface
type MockPeerTask struct {
	ctrl     *gomock.Controller
	recorder *MockPeerTaskMockRecorder
}

// MockPeerTaskMockRecorder is the mock recorder for MockPeerTask
type MockPeerTaskMockRecorder struct {
	mock *MockPeerTask
}

// NewMockPeerTask creates a new mock instance
func NewMockPeerTask(ctrl *gomock.Controller) *MockPeerTask {
	mock := &MockPeerTask{ctrl: ctrl}
	mock.recorder = &MockPeerTaskMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockPeerTask) EXPECT() *MockPeerTaskMockRecorder {
	return m.recorder
}

// ReportPieceResult mocks base method
func (m *MockPeerTask) ReportPieceResult(pieceTask *base.PieceInfo, pieceResult *scheduler.PieceResult) error {
	ret := m.ctrl.Call(m, "ReportPieceResult", pieceTask, pieceResult)
	ret0, _ := ret[0].(error)
	return ret0
}

// ReportPieceResult indicates an expected call of ReportPieceResult
func (mr *MockPeerTaskMockRecorder) ReportPieceResult(pieceTask, pieceResult interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReportPieceResult", reflect.TypeOf((*MockPeerTask)(nil).ReportPieceResult), pieceTask, pieceResult)
}

// GetPeerID mocks base method
func (m *MockPeerTask) GetPeerID() string {
	ret := m.ctrl.Call(m, "GetPeerID")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetPeerID indicates an expected call of GetPeerID
func (mr *MockPeerTaskMockRecorder) GetPeerID() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPeerID", reflect.TypeOf((*MockPeerTask)(nil).GetPeerID))
}

// GetTaskID mocks base method
func (m *MockPeerTask) GetTaskID() string {
	ret := m.ctrl.Call(m, "GetTaskID")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetTaskID indicates an expected call of GetTaskID
func (mr *MockPeerTaskMockRecorder) GetTaskID() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTaskID", reflect.TypeOf((*MockPeerTask)(nil).GetTaskID))
}

// GetTotalPieces mocks base method
func (m *MockPeerTask) GetTotalPieces() int32 {
	ret := m.ctrl.Call(m, "GetTotalPieces")
	ret0, _ := ret[0].(int32)
	return ret0
}

// GetTotalPieces indicates an expected call of GetTotalPieces
func (mr *MockPeerTaskMockRecorder) GetTotalPieces() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTotalPieces", reflect.TypeOf((*MockPeerTask)(nil).GetTotalPieces))
}

// GetContentLength mocks base method
func (m *MockPeerTask) GetContentLength() int64 {
	ret := m.ctrl.Call(m, "GetContentLength")
	ret0, _ := ret[0].(int64)
	return ret0
}

// GetContentLength indicates an expected call of GetContentLength
func (mr *MockPeerTaskMockRecorder) GetContentLength() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetContentLength", reflect.TypeOf((*MockPeerTask)(nil).GetContentLength))
}

// SetContentLength mocks base method
func (m *MockPeerTask) SetContentLength(arg0 int64) error {
	ret := m.ctrl.Call(m, "SetContentLength", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetContentLength indicates an expected call of SetContentLength
func (mr *MockPeerTaskMockRecorder) SetContentLength(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetContentLength", reflect.TypeOf((*MockPeerTask)(nil).SetContentLength), arg0)
}

// SetCallback mocks base method
func (m *MockPeerTask) SetCallback(arg0 PeerTaskCallback) {
	m.ctrl.Call(m, "SetCallback", arg0)
}

// SetCallback indicates an expected call of SetCallback
func (mr *MockPeerTaskMockRecorder) SetCallback(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetCallback", reflect.TypeOf((*MockPeerTask)(nil).SetCallback), arg0)
}

// AddTraffic mocks base method
func (m *MockPeerTask) AddTraffic(arg0 int64) {
	m.ctrl.Call(m, "AddTraffic", arg0)
}

// AddTraffic indicates an expected call of AddTraffic
func (mr *MockPeerTaskMockRecorder) AddTraffic(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddTraffic", reflect.TypeOf((*MockPeerTask)(nil).AddTraffic), arg0)
}

// GetTraffic mocks base method
func (m *MockPeerTask) GetTraffic() int64 {
	ret := m.ctrl.Call(m, "GetTraffic")
	ret0, _ := ret[0].(int64)
	return ret0
}

// GetTraffic indicates an expected call of GetTraffic
func (mr *MockPeerTaskMockRecorder) GetTraffic() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTraffic", reflect.TypeOf((*MockPeerTask)(nil).GetTraffic))
}

// GetContext mocks base method
func (m *MockPeerTask) Context() context.Context {
	ret := m.ctrl.Call(m, "Context")
	ret0, _ := ret[0].(context.Context)
	return ret0
}

// GetContext indicates an expected call of GetContext
func (mr *MockPeerTaskMockRecorder) GetContext() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Context", reflect.TypeOf((*MockPeerTask)(nil).Context))
}

// Log mocks base method
func (m *MockPeerTask) Log() *dflog.SugaredLoggerOnWith {
	ret := m.ctrl.Call(m, "Log")
	ret0, _ := ret[0].(*dflog.SugaredLoggerOnWith)
	return ret0
}

// Log indicates an expected call of Log
func (mr *MockPeerTaskMockRecorder) Log() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Log", reflect.TypeOf((*MockPeerTask)(nil).Log))
}

// MockPeerTaskCallback is a mock of PeerTaskCallback interface
type MockPeerTaskCallback struct {
	ctrl     *gomock.Controller
	recorder *MockPeerTaskCallbackMockRecorder
}

// MockPeerTaskCallbackMockRecorder is the mock recorder for MockPeerTaskCallback
type MockPeerTaskCallbackMockRecorder struct {
	mock *MockPeerTaskCallback
}

// NewMockPeerTaskCallback creates a new mock instance
func NewMockPeerTaskCallback(ctrl *gomock.Controller) *MockPeerTaskCallback {
	mock := &MockPeerTaskCallback{ctrl: ctrl}
	mock.recorder = &MockPeerTaskCallbackMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockPeerTaskCallback) EXPECT() *MockPeerTaskCallbackMockRecorder {
	return m.recorder
}

// Init mocks base method
func (m *MockPeerTaskCallback) Init(pt PeerTask) error {
	ret := m.ctrl.Call(m, "Init", pt)
	ret0, _ := ret[0].(error)
	return ret0
}

// Init indicates an expected call of Init
func (mr *MockPeerTaskCallbackMockRecorder) Init(pt interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Init", reflect.TypeOf((*MockPeerTaskCallback)(nil).Init), pt)
}

// Done mocks base method
func (m *MockPeerTaskCallback) Done(pt PeerTask) error {
	ret := m.ctrl.Call(m, "Done", pt)
	ret0, _ := ret[0].(error)
	return ret0
}

// Done indicates an expected call of Done
func (mr *MockPeerTaskCallbackMockRecorder) Done(pt interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Done", reflect.TypeOf((*MockPeerTaskCallback)(nil).Done), pt)
}

// Update mocks base method
func (m *MockPeerTaskCallback) Update(pt PeerTask) error {
	ret := m.ctrl.Call(m, "Update", pt)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update
func (mr *MockPeerTaskCallbackMockRecorder) Update(pt interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockPeerTaskCallback)(nil).Update), pt)
}

// Fail mocks base method
func (m *MockPeerTaskCallback) Fail(pt PeerTask, reason string) error {
	ret := m.ctrl.Call(m, "Fail", pt, reason)
	ret0, _ := ret[0].(error)
	return ret0
}

// Fail indicates an expected call of Fail
func (mr *MockPeerTaskCallbackMockRecorder) Fail(pt, reason interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Fail", reflect.TypeOf((*MockPeerTaskCallback)(nil).Fail), pt, reason)
}
