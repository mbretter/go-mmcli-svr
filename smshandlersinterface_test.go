// Code generated by mockery. DO NOT EDIT.

package main

import (
	http "net/http"

	mock "github.com/stretchr/testify/mock"
)

// SmsHandlersInterfaceMock is an autogenerated mock type for the SmsHandlersInterface type
type SmsHandlersInterfaceMock struct {
	mock.Mock
}

type SmsHandlersInterfaceMock_Expecter struct {
	mock *mock.Mock
}

func (_m *SmsHandlersInterfaceMock) EXPECT() *SmsHandlersInterfaceMock_Expecter {
	return &SmsHandlersInterfaceMock_Expecter{mock: &_m.Mock}
}

// SmsDelete provides a mock function with given fields: w, r
func (_m *SmsHandlersInterfaceMock) SmsDelete(w http.ResponseWriter, r *http.Request) {
	_m.Called(w, r)
}

// SmsHandlersInterfaceMock_SmsDelete_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SmsDelete'
type SmsHandlersInterfaceMock_SmsDelete_Call struct {
	*mock.Call
}

// SmsDelete is a helper method to define mock.On call
//   - w http.ResponseWriter
//   - r *http.Request
func (_e *SmsHandlersInterfaceMock_Expecter) SmsDelete(w interface{}, r interface{}) *SmsHandlersInterfaceMock_SmsDelete_Call {
	return &SmsHandlersInterfaceMock_SmsDelete_Call{Call: _e.mock.On("SmsDelete", w, r)}
}

func (_c *SmsHandlersInterfaceMock_SmsDelete_Call) Run(run func(w http.ResponseWriter, r *http.Request)) *SmsHandlersInterfaceMock_SmsDelete_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(http.ResponseWriter), args[1].(*http.Request))
	})
	return _c
}

func (_c *SmsHandlersInterfaceMock_SmsDelete_Call) Return() *SmsHandlersInterfaceMock_SmsDelete_Call {
	_c.Call.Return()
	return _c
}

func (_c *SmsHandlersInterfaceMock_SmsDelete_Call) RunAndReturn(run func(http.ResponseWriter, *http.Request)) *SmsHandlersInterfaceMock_SmsDelete_Call {
	_c.Call.Return(run)
	return _c
}

// SmsGet provides a mock function with given fields: w, r
func (_m *SmsHandlersInterfaceMock) SmsGet(w http.ResponseWriter, r *http.Request) {
	_m.Called(w, r)
}

// SmsHandlersInterfaceMock_SmsGet_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SmsGet'
type SmsHandlersInterfaceMock_SmsGet_Call struct {
	*mock.Call
}

// SmsGet is a helper method to define mock.On call
//   - w http.ResponseWriter
//   - r *http.Request
func (_e *SmsHandlersInterfaceMock_Expecter) SmsGet(w interface{}, r interface{}) *SmsHandlersInterfaceMock_SmsGet_Call {
	return &SmsHandlersInterfaceMock_SmsGet_Call{Call: _e.mock.On("SmsGet", w, r)}
}

func (_c *SmsHandlersInterfaceMock_SmsGet_Call) Run(run func(w http.ResponseWriter, r *http.Request)) *SmsHandlersInterfaceMock_SmsGet_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(http.ResponseWriter), args[1].(*http.Request))
	})
	return _c
}

func (_c *SmsHandlersInterfaceMock_SmsGet_Call) Return() *SmsHandlersInterfaceMock_SmsGet_Call {
	_c.Call.Return()
	return _c
}

func (_c *SmsHandlersInterfaceMock_SmsGet_Call) RunAndReturn(run func(http.ResponseWriter, *http.Request)) *SmsHandlersInterfaceMock_SmsGet_Call {
	_c.Call.Return(run)
	return _c
}

// SmsSend provides a mock function with given fields: w, r
func (_m *SmsHandlersInterfaceMock) SmsSend(w http.ResponseWriter, r *http.Request) {
	_m.Called(w, r)
}

// SmsHandlersInterfaceMock_SmsSend_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SmsSend'
type SmsHandlersInterfaceMock_SmsSend_Call struct {
	*mock.Call
}

// SmsSend is a helper method to define mock.On call
//   - w http.ResponseWriter
//   - r *http.Request
func (_e *SmsHandlersInterfaceMock_Expecter) SmsSend(w interface{}, r interface{}) *SmsHandlersInterfaceMock_SmsSend_Call {
	return &SmsHandlersInterfaceMock_SmsSend_Call{Call: _e.mock.On("SmsSend", w, r)}
}

func (_c *SmsHandlersInterfaceMock_SmsSend_Call) Run(run func(w http.ResponseWriter, r *http.Request)) *SmsHandlersInterfaceMock_SmsSend_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(http.ResponseWriter), args[1].(*http.Request))
	})
	return _c
}

func (_c *SmsHandlersInterfaceMock_SmsSend_Call) Return() *SmsHandlersInterfaceMock_SmsSend_Call {
	_c.Call.Return()
	return _c
}

func (_c *SmsHandlersInterfaceMock_SmsSend_Call) RunAndReturn(run func(http.ResponseWriter, *http.Request)) *SmsHandlersInterfaceMock_SmsSend_Call {
	_c.Call.Return(run)
	return _c
}

// NewSmsHandlersInterfaceMock creates a new instance of SmsHandlersInterfaceMock. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewSmsHandlersInterfaceMock(t interface {
	mock.TestingT
	Cleanup(func())
}) *SmsHandlersInterfaceMock {
	mock := &SmsHandlersInterfaceMock{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
