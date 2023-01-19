// Code generated by MockGen. DO NOT EDIT.
// Source: escort-book-tracking/repositories (interfaces: IEscortProfileRepository)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	models "escort-book-tracking/models"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockIEscortProfileRepository is a mock of IEscortProfileRepository interface.
type MockIEscortProfileRepository struct {
	ctrl     *gomock.Controller
	recorder *MockIEscortProfileRepositoryMockRecorder
}

// MockIEscortProfileRepositoryMockRecorder is the mock recorder for MockIEscortProfileRepository.
type MockIEscortProfileRepositoryMockRecorder struct {
	mock *MockIEscortProfileRepository
}

// NewMockIEscortProfileRepository creates a new mock instance.
func NewMockIEscortProfileRepository(ctrl *gomock.Controller) *MockIEscortProfileRepository {
	mock := &MockIEscortProfileRepository{ctrl: ctrl}
	mock.recorder = &MockIEscortProfileRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIEscortProfileRepository) EXPECT() *MockIEscortProfileRepositoryMockRecorder {
	return m.recorder
}

// GetEscortProfile mocks base method.
func (m *MockIEscortProfileRepository) GetEscortProfile(arg0 context.Context, arg1 string) (*models.EscortProfile, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetEscortProfile", arg0, arg1)
	ret0, _ := ret[0].(*models.EscortProfile)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetEscortProfile indicates an expected call of GetEscortProfile.
func (mr *MockIEscortProfileRepositoryMockRecorder) GetEscortProfile(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetEscortProfile", reflect.TypeOf((*MockIEscortProfileRepository)(nil).GetEscortProfile), arg0, arg1)
}
