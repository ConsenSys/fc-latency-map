// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/ConsenSys/fc-latency-map/manager/geomgr (interfaces: GeoMgr)

// Package geomgr is a generated GoMock package.
package geomgr

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockGeoMgr is a mock of GeoMgr interface.
type MockGeoMgr struct {
	ctrl     *gomock.Controller
	recorder *MockGeoMgrMockRecorder
}

// MockGeoMgrMockRecorder is the mock recorder for MockGeoMgr.
type MockGeoMgrMockRecorder struct {
	mock *MockGeoMgr
}

// NewMockGeoMgr creates a new mock instance.
func NewMockGeoMgr(ctrl *gomock.Controller) *MockGeoMgr {
	mock := &MockGeoMgr{ctrl: ctrl}
	mock.recorder = &MockGeoMgrMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockGeoMgr) EXPECT() *MockGeoMgrMockRecorder {
	return m.recorder
}

// IPGeolocation mocks base method.
func (m *MockGeoMgr) IPGeolocation(arg0 string) (float64, float64) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IPGeolocation", arg0)
	ret0, _ := ret[0].(float64)
	ret1, _ := ret[1].(float64)
	return ret0, ret1
}

// IPGeolocation indicates an expected call of IPGeolocation.
func (mr *MockGeoMgrMockRecorder) IPGeolocation(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IPGeolocation", reflect.TypeOf((*MockGeoMgr)(nil).IPGeolocation), arg0)
}
