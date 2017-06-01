// Automatically generated by MockGen. DO NOT EDIT!
// Source: ../interfaces/cache.go

package mock_interfaces

import (
	gomock "github.com/golang/mock/gomock"
)

// Mock of CacheServiceProvider interface
type MockCacheServiceProvider struct {
	ctrl     *gomock.Controller
	recorder *_MockCacheServiceProviderRecorder
}

// Recorder for MockCacheServiceProvider (not exported)
type _MockCacheServiceProviderRecorder struct {
	mock *MockCacheServiceProvider
}

func NewMockCacheServiceProvider(ctrl *gomock.Controller) *MockCacheServiceProvider {
	mock := &MockCacheServiceProvider{ctrl: ctrl}
	mock.recorder = &_MockCacheServiceProviderRecorder{mock}
	return mock
}

func (_m *MockCacheServiceProvider) EXPECT() *_MockCacheServiceProviderRecorder {
	return _m.recorder
}

func (_m *MockCacheServiceProvider) Get(key string) (string, error) {
	ret := _m.ctrl.Call(_m, "Get", key)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockCacheServiceProviderRecorder) Get(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Get", arg0)
}

func (_m *MockCacheServiceProvider) SetWithExpiry(key string, value string, expiryInSeconds int) error {
	ret := _m.ctrl.Call(_m, "SetWithExpiry", key, value, expiryInSeconds)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockCacheServiceProviderRecorder) SetWithExpiry(arg0, arg1, arg2 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "SetWithExpiry", arg0, arg1, arg2)
}

func (_m *MockCacheServiceProvider) Delete(key string) error {
	ret := _m.ctrl.Call(_m, "Delete", key)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockCacheServiceProviderRecorder) Delete(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Delete", arg0)
}

func (_m *MockCacheServiceProvider) Exists(key string) (bool, error) {
	ret := _m.ctrl.Call(_m, "Exists", key)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockCacheServiceProviderRecorder) Exists(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Exists", arg0)
}

func (_m *MockCacheServiceProvider) Flush() error {
	ret := _m.ctrl.Call(_m, "Flush")
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockCacheServiceProviderRecorder) Flush() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Flush")
}