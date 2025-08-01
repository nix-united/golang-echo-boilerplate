// Code generated by MockGen. DO NOT EDIT.
// Source: oauth_handler.go
//
// Generated by this command:
//
//	mockgen -source=oauth_handler.go -destination=oauth_handler_mock_test.go -package=handlers_test -typed=true
//

// Package handlers_test is a generated GoMock package.
package handlers_test

import (
	context "context"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockUserAuthenticator is a mock of UserAuthenticator interface.
type MockUserAuthenticator struct {
	ctrl     *gomock.Controller
	recorder *MockUserAuthenticatorMockRecorder
	isgomock struct{}
}

// MockUserAuthenticatorMockRecorder is the mock recorder for MockUserAuthenticator.
type MockUserAuthenticatorMockRecorder struct {
	mock *MockUserAuthenticator
}

// NewMockUserAuthenticator creates a new mock instance.
func NewMockUserAuthenticator(ctrl *gomock.Controller) *MockUserAuthenticator {
	mock := &MockUserAuthenticator{ctrl: ctrl}
	mock.recorder = &MockUserAuthenticatorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserAuthenticator) EXPECT() *MockUserAuthenticatorMockRecorder {
	return m.recorder
}

// GoogleOAuth mocks base method.
func (m *MockUserAuthenticator) GoogleOAuth(ctx context.Context, token string) (string, string, int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GoogleOAuth", ctx, token)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(int64)
	ret3, _ := ret[3].(error)
	return ret0, ret1, ret2, ret3
}

// GoogleOAuth indicates an expected call of GoogleOAuth.
func (mr *MockUserAuthenticatorMockRecorder) GoogleOAuth(ctx, token any) *MockUserAuthenticatorGoogleOAuthCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GoogleOAuth", reflect.TypeOf((*MockUserAuthenticator)(nil).GoogleOAuth), ctx, token)
	return &MockUserAuthenticatorGoogleOAuthCall{Call: call}
}

// MockUserAuthenticatorGoogleOAuthCall wrap *gomock.Call
type MockUserAuthenticatorGoogleOAuthCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockUserAuthenticatorGoogleOAuthCall) Return(arg0, arg1 string, arg2 int64, arg3 error) *MockUserAuthenticatorGoogleOAuthCall {
	c.Call = c.Call.Return(arg0, arg1, arg2, arg3)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockUserAuthenticatorGoogleOAuthCall) Do(f func(context.Context, string) (string, string, int64, error)) *MockUserAuthenticatorGoogleOAuthCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockUserAuthenticatorGoogleOAuthCall) DoAndReturn(f func(context.Context, string) (string, string, int64, error)) *MockUserAuthenticatorGoogleOAuthCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}
