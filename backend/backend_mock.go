// Code generated by mockery. DO NOT EDIT.

package backend

import mock "github.com/stretchr/testify/mock"

// BackendMock is an autogenerated mock type for the Backend type
type BackendMock struct {
	mock.Mock
}

type BackendMock_Expecter struct {
	mock *mock.Mock
}

func (_m *BackendMock) EXPECT() *BackendMock_Expecter {
	return &BackendMock_Expecter{mock: &_m.Mock}
}

// Exec provides a mock function with given fields: arg
func (_m *BackendMock) Exec(arg ...string) ([]byte, error) {
	_va := make([]interface{}, len(arg))
	for _i := range arg {
		_va[_i] = arg[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for Exec")
	}

	var r0 []byte
	var r1 error
	if rf, ok := ret.Get(0).(func(...string) ([]byte, error)); ok {
		return rf(arg...)
	}
	if rf, ok := ret.Get(0).(func(...string) []byte); ok {
		r0 = rf(arg...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	if rf, ok := ret.Get(1).(func(...string) error); ok {
		r1 = rf(arg...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// BackendMock_Exec_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Exec'
type BackendMock_Exec_Call struct {
	*mock.Call
}

// Exec is a helper method to define mock.On call
//   - arg ...string
func (_e *BackendMock_Expecter) Exec(arg ...interface{}) *BackendMock_Exec_Call {
	return &BackendMock_Exec_Call{Call: _e.mock.On("Exec",
		append([]interface{}{}, arg...)...)}
}

func (_c *BackendMock_Exec_Call) Run(run func(arg ...string)) *BackendMock_Exec_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]string, len(args)-0)
		for i, a := range args[0:] {
			if a != nil {
				variadicArgs[i] = a.(string)
			}
		}
		run(variadicArgs...)
	})
	return _c
}

func (_c *BackendMock_Exec_Call) Return(_a0 []byte, _a1 error) *BackendMock_Exec_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *BackendMock_Exec_Call) RunAndReturn(run func(...string) ([]byte, error)) *BackendMock_Exec_Call {
	_c.Call.Return(run)
	return _c
}

// ExecModem provides a mock function with given fields: modem, arg
func (_m *BackendMock) ExecModem(modem string, arg ...string) ([]byte, error) {
	_va := make([]interface{}, len(arg))
	for _i := range arg {
		_va[_i] = arg[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, modem)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for ExecModem")
	}

	var r0 []byte
	var r1 error
	if rf, ok := ret.Get(0).(func(string, ...string) ([]byte, error)); ok {
		return rf(modem, arg...)
	}
	if rf, ok := ret.Get(0).(func(string, ...string) []byte); ok {
		r0 = rf(modem, arg...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	if rf, ok := ret.Get(1).(func(string, ...string) error); ok {
		r1 = rf(modem, arg...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// BackendMock_ExecModem_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ExecModem'
type BackendMock_ExecModem_Call struct {
	*mock.Call
}

// ExecModem is a helper method to define mock.On call
//   - modem string
//   - arg ...string
func (_e *BackendMock_Expecter) ExecModem(modem interface{}, arg ...interface{}) *BackendMock_ExecModem_Call {
	return &BackendMock_ExecModem_Call{Call: _e.mock.On("ExecModem",
		append([]interface{}{modem}, arg...)...)}
}

func (_c *BackendMock_ExecModem_Call) Run(run func(modem string, arg ...string)) *BackendMock_ExecModem_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]string, len(args)-1)
		for i, a := range args[1:] {
			if a != nil {
				variadicArgs[i] = a.(string)
			}
		}
		run(args[0].(string), variadicArgs...)
	})
	return _c
}

func (_c *BackendMock_ExecModem_Call) Return(_a0 []byte, _a1 error) *BackendMock_ExecModem_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *BackendMock_ExecModem_Call) RunAndReturn(run func(string, ...string) ([]byte, error)) *BackendMock_ExecModem_Call {
	_c.Call.Return(run)
	return _c
}

// NewBackendMock creates a new instance of BackendMock. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewBackendMock(t interface {
	mock.TestingT
	Cleanup(func())
}) *BackendMock {
	mock := &BackendMock{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
