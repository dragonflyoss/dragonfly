// Code generated by MockGen. DO NOT EDIT.
// Source: dag.go
//
// Generated by this command:
//
//	mockgen -destination mocks/dag_mock.go -source dag.go -package mocks
//

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	dag "d7y.io/dragonfly/v2/pkg/graph/dag"
	gomock "go.uber.org/mock/gomock"
)

// MockDAG is a mock of DAG interface.
type MockDAG[T comparable] struct {
	ctrl     *gomock.Controller
	recorder *MockDAGMockRecorder[T]
	isgomock struct{}
}

// MockDAGMockRecorder is the mock recorder for MockDAG.
type MockDAGMockRecorder[T comparable] struct {
	mock *MockDAG[T]
}

// NewMockDAG creates a new mock instance.
func NewMockDAG[T comparable](ctrl *gomock.Controller) *MockDAG[T] {
	mock := &MockDAG[T]{ctrl: ctrl}
	mock.recorder = &MockDAGMockRecorder[T]{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDAG[T]) EXPECT() *MockDAGMockRecorder[T] {
	return m.recorder
}

// AddEdge mocks base method.
func (m *MockDAG[T]) AddEdge(fromVertexID, toVertexID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddEdge", fromVertexID, toVertexID)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddEdge indicates an expected call of AddEdge.
func (mr *MockDAGMockRecorder[T]) AddEdge(fromVertexID, toVertexID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddEdge", reflect.TypeOf((*MockDAG[T])(nil).AddEdge), fromVertexID, toVertexID)
}

// AddVertex mocks base method.
func (m *MockDAG[T]) AddVertex(id string, value T) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddVertex", id, value)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddVertex indicates an expected call of AddVertex.
func (mr *MockDAGMockRecorder[T]) AddVertex(id, value any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddVertex", reflect.TypeOf((*MockDAG[T])(nil).AddVertex), id, value)
}

// CanAddEdge mocks base method.
func (m *MockDAG[T]) CanAddEdge(fromVertexID, toVertexID string) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CanAddEdge", fromVertexID, toVertexID)
	ret0, _ := ret[0].(bool)
	return ret0
}

// CanAddEdge indicates an expected call of CanAddEdge.
func (mr *MockDAGMockRecorder[T]) CanAddEdge(fromVertexID, toVertexID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CanAddEdge", reflect.TypeOf((*MockDAG[T])(nil).CanAddEdge), fromVertexID, toVertexID)
}

// DeleteEdge mocks base method.
func (m *MockDAG[T]) DeleteEdge(fromVertexID, toVertexID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteEdge", fromVertexID, toVertexID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteEdge indicates an expected call of DeleteEdge.
func (mr *MockDAGMockRecorder[T]) DeleteEdge(fromVertexID, toVertexID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteEdge", reflect.TypeOf((*MockDAG[T])(nil).DeleteEdge), fromVertexID, toVertexID)
}

// DeleteVertex mocks base method.
func (m *MockDAG[T]) DeleteVertex(id string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "DeleteVertex", id)
}

// DeleteVertex indicates an expected call of DeleteVertex.
func (mr *MockDAGMockRecorder[T]) DeleteVertex(id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteVertex", reflect.TypeOf((*MockDAG[T])(nil).DeleteVertex), id)
}

// DeleteVertexInEdges mocks base method.
func (m *MockDAG[T]) DeleteVertexInEdges(id string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteVertexInEdges", id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteVertexInEdges indicates an expected call of DeleteVertexInEdges.
func (mr *MockDAGMockRecorder[T]) DeleteVertexInEdges(id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteVertexInEdges", reflect.TypeOf((*MockDAG[T])(nil).DeleteVertexInEdges), id)
}

// DeleteVertexOutEdges mocks base method.
func (m *MockDAG[T]) DeleteVertexOutEdges(id string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteVertexOutEdges", id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteVertexOutEdges indicates an expected call of DeleteVertexOutEdges.
func (mr *MockDAGMockRecorder[T]) DeleteVertexOutEdges(id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteVertexOutEdges", reflect.TypeOf((*MockDAG[T])(nil).DeleteVertexOutEdges), id)
}

// GetRandomVertices mocks base method.
func (m *MockDAG[T]) GetRandomVertices(n uint) []*dag.Vertex[T] {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRandomVertices", n)
	ret0, _ := ret[0].([]*dag.Vertex[T])
	return ret0
}

// GetRandomVertices indicates an expected call of GetRandomVertices.
func (mr *MockDAGMockRecorder[T]) GetRandomVertices(n any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRandomVertices", reflect.TypeOf((*MockDAG[T])(nil).GetRandomVertices), n)
}

// GetSinkVertices mocks base method.
func (m *MockDAG[T]) GetSinkVertices() []*dag.Vertex[T] {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSinkVertices")
	ret0, _ := ret[0].([]*dag.Vertex[T])
	return ret0
}

// GetSinkVertices indicates an expected call of GetSinkVertices.
func (mr *MockDAGMockRecorder[T]) GetSinkVertices() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSinkVertices", reflect.TypeOf((*MockDAG[T])(nil).GetSinkVertices))
}

// GetSourceVertices mocks base method.
func (m *MockDAG[T]) GetSourceVertices() []*dag.Vertex[T] {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSourceVertices")
	ret0, _ := ret[0].([]*dag.Vertex[T])
	return ret0
}

// GetSourceVertices indicates an expected call of GetSourceVertices.
func (mr *MockDAGMockRecorder[T]) GetSourceVertices() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSourceVertices", reflect.TypeOf((*MockDAG[T])(nil).GetSourceVertices))
}

// GetVertex mocks base method.
func (m *MockDAG[T]) GetVertex(id string) (*dag.Vertex[T], error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetVertex", id)
	ret0, _ := ret[0].(*dag.Vertex[T])
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetVertex indicates an expected call of GetVertex.
func (mr *MockDAGMockRecorder[T]) GetVertex(id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetVertex", reflect.TypeOf((*MockDAG[T])(nil).GetVertex), id)
}

// GetVertices mocks base method.
func (m *MockDAG[T]) GetVertices() map[string]*dag.Vertex[T] {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetVertices")
	ret0, _ := ret[0].(map[string]*dag.Vertex[T])
	return ret0
}

// GetVertices indicates an expected call of GetVertices.
func (mr *MockDAGMockRecorder[T]) GetVertices() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetVertices", reflect.TypeOf((*MockDAG[T])(nil).GetVertices))
}

// VertexCount mocks base method.
func (m *MockDAG[T]) VertexCount() uint64 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "VertexCount")
	ret0, _ := ret[0].(uint64)
	return ret0
}

// VertexCount indicates an expected call of VertexCount.
func (mr *MockDAGMockRecorder[T]) VertexCount() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "VertexCount", reflect.TypeOf((*MockDAG[T])(nil).VertexCount))
}
