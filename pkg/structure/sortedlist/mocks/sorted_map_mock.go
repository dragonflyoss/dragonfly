// Code generated by MockGen. DO NOT EDIT.
// Source: d7y.io/dragonfly/v2/pkg/structure/sortedlist (interfaces: SortedList,Item)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	sortedlist "d7y.io/dragonfly/v2/pkg/structure/sortedlist"
	gomock "github.com/golang/mock/gomock"
)

// MockSortedList is a mock of SortedList interface.
type MockSortedList struct {
	ctrl     *gomock.Controller
	recorder *MockSortedListMockRecorder
}

// MockSortedListMockRecorder is the mock recorder for MockSortedList.
type MockSortedListMockRecorder struct {
	mock *MockSortedList
}

// NewMockSortedList creates a new mock instance.
func NewMockSortedList(ctrl *gomock.Controller) *MockSortedList {
	mock := &MockSortedList{ctrl: ctrl}
	mock.recorder = &MockSortedListMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSortedList) EXPECT() *MockSortedListMockRecorder {
	return m.recorder
}

// Add mocks base method.
func (m *MockSortedList) Add(arg0 sortedlist.Item) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Add", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Add indicates an expected call of Add.
func (mr *MockSortedListMockRecorder) Add(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Add", reflect.TypeOf((*MockSortedList)(nil).Add), arg0)
}

// Delete mocks base method.
func (m *MockSortedList) Delete(arg0 sortedlist.Item) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockSortedListMockRecorder) Delete(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockSortedList)(nil).Delete), arg0)
}

// Range mocks base method.
func (m *MockSortedList) Range(arg0 func(sortedlist.Item) bool) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Range", arg0)
}

// Range indicates an expected call of Range.
func (mr *MockSortedListMockRecorder) Range(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Range", reflect.TypeOf((*MockSortedList)(nil).Range), arg0)
}

// ReverseRange mocks base method.
func (m *MockSortedList) ReverseRange(arg0 func(sortedlist.Item) bool) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ReverseRange", arg0)
}

// ReverseRange indicates an expected call of ReverseRange.
func (mr *MockSortedListMockRecorder) ReverseRange(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReverseRange", reflect.TypeOf((*MockSortedList)(nil).ReverseRange), arg0)
}

// Size mocks base method.
func (m *MockSortedList) Size() int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Size")
	ret0, _ := ret[0].(int)
	return ret0
}

// Size indicates an expected call of Size.
func (mr *MockSortedListMockRecorder) Size() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Size", reflect.TypeOf((*MockSortedList)(nil).Size))
}

// Update mocks base method.
func (m *MockSortedList) Update(arg0 sortedlist.Item) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockSortedListMockRecorder) Update(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockSortedList)(nil).Update), arg0)
}

// UpdateOrAdd mocks base method.
func (m *MockSortedList) UpdateOrAdd(arg0 sortedlist.Item) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateOrAdd", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateOrAdd indicates an expected call of UpdateOrAdd.
func (mr *MockSortedListMockRecorder) UpdateOrAdd(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateOrAdd", reflect.TypeOf((*MockSortedList)(nil).UpdateOrAdd), arg0)
}

// MockItem is a mock of Item interface.
type MockItem struct {
	ctrl     *gomock.Controller
	recorder *MockItemMockRecorder
}

// MockItemMockRecorder is the mock recorder for MockItem.
type MockItemMockRecorder struct {
	mock *MockItem
}

// NewMockItem creates a new mock instance.
func NewMockItem(ctrl *gomock.Controller) *MockItem {
	mock := &MockItem{ctrl: ctrl}
	mock.recorder = &MockItemMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockItem) EXPECT() *MockItemMockRecorder {
	return m.recorder
}

// GetSortKeys mocks base method.
func (m *MockItem) GetSortKeys() (int, int) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSortKeys")
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(int)
	return ret0, ret1
}

// GetSortKeys indicates an expected call of GetSortKeys.
func (mr *MockItemMockRecorder) GetSortKeys() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSortKeys", reflect.TypeOf((*MockItem)(nil).GetSortKeys))
}
