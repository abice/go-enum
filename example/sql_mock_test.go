// Code generated by MockGen. DO NOT EDIT.
// Source: database/sql/driver (interfaces: Conn,Driver,Stmt,Result,Rows)

// Package example is a generated GoMock package.
package example

import (
	driver "database/sql/driver"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockConn is a mock of Conn interface
type MockConn struct {
	ctrl     *gomock.Controller
	recorder *MockConnMockRecorder
}

// MockConnMockRecorder is the mock recorder for MockConn
type MockConnMockRecorder struct {
	mock *MockConn
}

// NewMockConn creates a new mock instance
func NewMockConn(ctrl *gomock.Controller) *MockConn {
	mock := &MockConn{ctrl: ctrl}
	mock.recorder = &MockConnMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockConn) EXPECT() *MockConnMockRecorder {
	return m.recorder
}

// Begin mocks base method
func (m *MockConn) Begin() (driver.Tx, error) {
	ret := m.ctrl.Call(m, "Begin")
	ret0, _ := ret[0].(driver.Tx)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Begin indicates an expected call of Begin
func (mr *MockConnMockRecorder) Begin() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Begin", reflect.TypeOf((*MockConn)(nil).Begin))
}

// Close mocks base method
func (m *MockConn) Close() error {
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close
func (mr *MockConnMockRecorder) Close() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockConn)(nil).Close))
}

// Prepare mocks base method
func (m *MockConn) Prepare(arg0 string) (driver.Stmt, error) {
	ret := m.ctrl.Call(m, "Prepare", arg0)
	ret0, _ := ret[0].(driver.Stmt)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Prepare indicates an expected call of Prepare
func (mr *MockConnMockRecorder) Prepare(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Prepare", reflect.TypeOf((*MockConn)(nil).Prepare), arg0)
}

// MockDriver is a mock of Driver interface
type MockDriver struct {
	ctrl     *gomock.Controller
	recorder *MockDriverMockRecorder
}

// MockDriverMockRecorder is the mock recorder for MockDriver
type MockDriverMockRecorder struct {
	mock *MockDriver
}

// NewMockDriver creates a new mock instance
func NewMockDriver(ctrl *gomock.Controller) *MockDriver {
	mock := &MockDriver{ctrl: ctrl}
	mock.recorder = &MockDriverMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockDriver) EXPECT() *MockDriverMockRecorder {
	return m.recorder
}

// Open mocks base method
func (m *MockDriver) Open(arg0 string) (driver.Conn, error) {
	ret := m.ctrl.Call(m, "Open", arg0)
	ret0, _ := ret[0].(driver.Conn)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Open indicates an expected call of Open
func (mr *MockDriverMockRecorder) Open(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Open", reflect.TypeOf((*MockDriver)(nil).Open), arg0)
}

// MockStmt is a mock of Stmt interface
type MockStmt struct {
	ctrl     *gomock.Controller
	recorder *MockStmtMockRecorder
}

// MockStmtMockRecorder is the mock recorder for MockStmt
type MockStmtMockRecorder struct {
	mock *MockStmt
}

// NewMockStmt creates a new mock instance
func NewMockStmt(ctrl *gomock.Controller) *MockStmt {
	mock := &MockStmt{ctrl: ctrl}
	mock.recorder = &MockStmtMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockStmt) EXPECT() *MockStmtMockRecorder {
	return m.recorder
}

// Close mocks base method
func (m *MockStmt) Close() error {
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close
func (mr *MockStmtMockRecorder) Close() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockStmt)(nil).Close))
}

// Exec mocks base method
func (m *MockStmt) Exec(arg0 []driver.Value) (driver.Result, error) {
	ret := m.ctrl.Call(m, "Exec", arg0)
	ret0, _ := ret[0].(driver.Result)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Exec indicates an expected call of Exec
func (mr *MockStmtMockRecorder) Exec(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Exec", reflect.TypeOf((*MockStmt)(nil).Exec), arg0)
}

// NumInput mocks base method
func (m *MockStmt) NumInput() int {
	ret := m.ctrl.Call(m, "NumInput")
	ret0, _ := ret[0].(int)
	return ret0
}

// NumInput indicates an expected call of NumInput
func (mr *MockStmtMockRecorder) NumInput() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NumInput", reflect.TypeOf((*MockStmt)(nil).NumInput))
}

// Query mocks base method
func (m *MockStmt) Query(arg0 []driver.Value) (driver.Rows, error) {
	ret := m.ctrl.Call(m, "Query", arg0)
	ret0, _ := ret[0].(driver.Rows)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Query indicates an expected call of Query
func (mr *MockStmtMockRecorder) Query(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Query", reflect.TypeOf((*MockStmt)(nil).Query), arg0)
}

// MockResult is a mock of Result interface
type MockResult struct {
	ctrl     *gomock.Controller
	recorder *MockResultMockRecorder
}

// MockResultMockRecorder is the mock recorder for MockResult
type MockResultMockRecorder struct {
	mock *MockResult
}

// NewMockResult creates a new mock instance
func NewMockResult(ctrl *gomock.Controller) *MockResult {
	mock := &MockResult{ctrl: ctrl}
	mock.recorder = &MockResultMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockResult) EXPECT() *MockResultMockRecorder {
	return m.recorder
}

// LastInsertId mocks base method
func (m *MockResult) LastInsertId() (int64, error) {
	ret := m.ctrl.Call(m, "LastInsertId")
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LastInsertId indicates an expected call of LastInsertId
func (mr *MockResultMockRecorder) LastInsertId() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LastInsertId", reflect.TypeOf((*MockResult)(nil).LastInsertId))
}

// RowsAffected mocks base method
func (m *MockResult) RowsAffected() (int64, error) {
	ret := m.ctrl.Call(m, "RowsAffected")
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RowsAffected indicates an expected call of RowsAffected
func (mr *MockResultMockRecorder) RowsAffected() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RowsAffected", reflect.TypeOf((*MockResult)(nil).RowsAffected))
}

// MockRows is a mock of Rows interface
type MockRows struct {
	ctrl     *gomock.Controller
	recorder *MockRowsMockRecorder
}

// MockRowsMockRecorder is the mock recorder for MockRows
type MockRowsMockRecorder struct {
	mock *MockRows
}

// NewMockRows creates a new mock instance
func NewMockRows(ctrl *gomock.Controller) *MockRows {
	mock := &MockRows{ctrl: ctrl}
	mock.recorder = &MockRowsMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockRows) EXPECT() *MockRowsMockRecorder {
	return m.recorder
}

// Close mocks base method
func (m *MockRows) Close() error {
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close
func (mr *MockRowsMockRecorder) Close() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockRows)(nil).Close))
}

// Columns mocks base method
func (m *MockRows) Columns() []string {
	ret := m.ctrl.Call(m, "Columns")
	ret0, _ := ret[0].([]string)
	return ret0
}

// Columns indicates an expected call of Columns
func (mr *MockRowsMockRecorder) Columns() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Columns", reflect.TypeOf((*MockRows)(nil).Columns))
}

// Next mocks base method
func (m *MockRows) Next(arg0 []driver.Value) error {
	ret := m.ctrl.Call(m, "Next", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Next indicates an expected call of Next
func (mr *MockRowsMockRecorder) Next(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Next", reflect.TypeOf((*MockRows)(nil).Next), arg0)
}
