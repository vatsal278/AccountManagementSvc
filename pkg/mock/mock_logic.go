// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/vatsal278/AccountManagmentSvc/internal/logic (interfaces: AccountManagmentSvcLogicIer)

// Package mock is a generated GoMock package.
package mock

import (
	reflect "reflect"

	model "github.com/PereRohit/util/model"
	gomock "github.com/golang/mock/gomock"
	model0 "github.com/vatsal278/AccountManagmentSvc/internal/model"
)

// MockAccountManagmentSvcLogicIer is a mock of AccountManagmentSvcLogicIer interface.
type MockAccountManagmentSvcLogicIer struct {
	ctrl     *gomock.Controller
	recorder *MockAccountManagmentSvcLogicIerMockRecorder
}

// MockAccountManagmentSvcLogicIerMockRecorder is the mock recorder for MockAccountManagmentSvcLogicIer.
type MockAccountManagmentSvcLogicIerMockRecorder struct {
	mock *MockAccountManagmentSvcLogicIer
}

// NewMockAccountManagmentSvcLogicIer creates a new mock instance.
func NewMockAccountManagmentSvcLogicIer(ctrl *gomock.Controller) *MockAccountManagmentSvcLogicIer {
	mock := &MockAccountManagmentSvcLogicIer{ctrl: ctrl}
	mock.recorder = &MockAccountManagmentSvcLogicIerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAccountManagmentSvcLogicIer) EXPECT() *MockAccountManagmentSvcLogicIerMockRecorder {
	return m.recorder
}

// AccountDetails mocks base method.
func (m *MockAccountManagmentSvcLogicIer) AccountDetails(arg0 string) *model.Response {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AccountDetails", arg0)
	ret0, _ := ret[0].(*model.Response)
	return ret0
}

// AccountDetails indicates an expected call of AccountDetails.
func (mr *MockAccountManagmentSvcLogicIerMockRecorder) AccountDetails(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AccountDetails", reflect.TypeOf((*MockAccountManagmentSvcLogicIer)(nil).AccountDetails), arg0)
}

// CreateAccount mocks base method.
func (m *MockAccountManagmentSvcLogicIer) CreateAccount(arg0 model0.NewAccount) *model.Response {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateAccount", arg0)
	ret0, _ := ret[0].(*model.Response)
	return ret0
}

// CreateAccount indicates an expected call of CreateAccount.
func (mr *MockAccountManagmentSvcLogicIerMockRecorder) CreateAccount(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateAccount", reflect.TypeOf((*MockAccountManagmentSvcLogicIer)(nil).CreateAccount), arg0)
}

// HealthCheck mocks base method.
func (m *MockAccountManagmentSvcLogicIer) HealthCheck() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HealthCheck")
	ret0, _ := ret[0].(bool)
	return ret0
}

// HealthCheck indicates an expected call of HealthCheck.
func (mr *MockAccountManagmentSvcLogicIerMockRecorder) HealthCheck() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HealthCheck", reflect.TypeOf((*MockAccountManagmentSvcLogicIer)(nil).HealthCheck))
}

// UpdateServices mocks base method.
func (m *MockAccountManagmentSvcLogicIer) UpdateServices(arg0 int, arg1, arg2 string) *model.Response {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateServices", arg0, arg1, arg2)
	ret0, _ := ret[0].(*model.Response)
	return ret0
}

// UpdateServices indicates an expected call of UpdateServices.
func (mr *MockAccountManagmentSvcLogicIerMockRecorder) UpdateServices(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateServices", reflect.TypeOf((*MockAccountManagmentSvcLogicIer)(nil).UpdateServices), arg0, arg1, arg2)
}

// UpdateTransaction mocks base method.
func (m *MockAccountManagmentSvcLogicIer) UpdateTransaction(arg0 int, arg1, arg2 string) *model.Response {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateTransaction", arg0, arg1, arg2)
	ret0, _ := ret[0].(*model.Response)
	return ret0
}

// UpdateTransaction indicates an expected call of UpdateTransaction.
func (mr *MockAccountManagmentSvcLogicIerMockRecorder) UpdateTransaction(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateTransaction", reflect.TypeOf((*MockAccountManagmentSvcLogicIer)(nil).UpdateTransaction), arg0, arg1, arg2)
}
