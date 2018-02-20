// Code generated by mockery v1.0.0

// +build mock

package quux

import mock "github.com/stretchr/testify/mock"

// MockClient is an autogenerated mock type for the Client type
type MockClient struct {
	mock.Mock
}

// Echo provides a mock function with given fields: _a0
func (_m *MockClient) Echo(_a0 string) string {
	ret := _m.Called(_a0)

	var r0 string
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}