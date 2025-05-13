// controller/mocks_test.go
package controller

import (
	"context"
	"reflect"

	"github.com/Mandarinka0707/newRepoGOODarhit/internal/entity"
	"github.com/Mandarinka0707/newRepoGOODarhit/internal/usecase"
	"github.com/golang/mock/gomock"
)

// gomock implementation for gRPC tests
type MockAuthUsecase struct {
	ctrl     *gomock.Controller
	recorder *MockAuthUsecaseRecorder
}

var _ usecase.AuthUsecaseInterface = (*MockAuthUsecase)(nil)

type MockAuthUsecaseRecorder struct {
	mock *MockAuthUsecase
}

func NewMockAuthUsecase(ctrl *gomock.Controller) *MockAuthUsecase {
	mock := &MockAuthUsecase{ctrl: ctrl}
	mock.recorder = &MockAuthUsecaseRecorder{mock}
	return mock
}

func (m *MockAuthUsecase) EXPECT() *MockAuthUsecaseRecorder {
	return m.recorder
}

func (m *MockAuthUsecase) Register(ctx context.Context, req *usecase.RegisterRequest) (*usecase.RegisterResponse, error) {
	ret := m.ctrl.Call(m, "Register", ctx, req)
	ret0, _ := ret[0].(*usecase.RegisterResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (m *MockAuthUsecase) Login(ctx context.Context, req *usecase.LoginRequest) (*usecase.LoginResponse, error) {
	ret := m.ctrl.Call(m, "Login", ctx, req)
	ret0, _ := ret[0].(*usecase.LoginResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (m *MockAuthUsecase) GetUser(ctx context.Context, req *usecase.GetUserRequest) (*usecase.GetUserResponse, error) {
	ret := m.ctrl.Call(m, "GetUser", ctx, req)
	ret0, _ := ret[0].(*usecase.GetUserResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (m *MockAuthUsecase) ValidateToken(ctx context.Context, req *usecase.ValidateTokenRequest) (*usecase.ValidateTokenResponse, error) {
	ret := m.ctrl.Call(m, "ValidateToken", ctx, req)
	ret0, _ := ret[0].(*usecase.ValidateTokenResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (m *MockAuthUsecase) GetUserByID(ctx context.Context, userID int64) (*entity.User, error) {
	ret := m.ctrl.Call(m, "GetUserByID", ctx, userID)
	ret0, _ := ret[0].(*entity.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *MockAuthUsecaseRecorder) Register(ctx, req interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(
		mr.mock,
		"Register",
		reflect.TypeOf((*MockAuthUsecase)(nil).Register),
		ctx,
		req,
	)
}

func (mr *MockAuthUsecaseRecorder) Login(ctx, req interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(
		mr.mock,
		"Login",
		reflect.TypeOf((*MockAuthUsecase)(nil).Login),
		ctx,
		req,
	)
}

func (mr *MockAuthUsecaseRecorder) GetUser(ctx, req interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(
		mr.mock,
		"GetUser",
		reflect.TypeOf((*MockAuthUsecase)(nil).GetUser),
		ctx,
		req,
	)
}

func (mr *MockAuthUsecaseRecorder) ValidateToken(ctx, req interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(
		mr.mock,
		"ValidateToken",
		reflect.TypeOf((*MockAuthUsecase)(nil).ValidateToken),
		ctx,
		req,
	)
}

func (mr *MockAuthUsecaseRecorder) GetUserByID(ctx, userID interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(
		mr.mock,
		"GetUserByID",
		reflect.TypeOf((*MockAuthUsecase)(nil).GetUserByID),
		ctx,
		userID,
	)
}
