// controller/auth_grpc_test.go
package controller

import (
	"context"
	"errors"
	"testing"
	"time"

	pb "backend.com/forum/proto"
	"github.com/Mandarinka0707/newRepoGOODarhit/internal/entity"
	"github.com/Mandarinka0707/newRepoGOODarhit/internal/usecase"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestAuthController_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := NewMockAuthUsecase(ctrl) // This now matches the implementation
	controller := NewAuthController(mockUC)

	tests := []struct {
		name        string
		req         *pb.RegisterRequest
		mockSetup   func(*MockAuthUsecase)
		want        *pb.RegisterResponse
		expectedErr *status.Status
	}{
		{
			name: "successful registration",
			req: &pb.RegisterRequest{
				Username: "testuser",
				Password: "testpass",
			},
			mockSetup: func(m *MockAuthUsecase) {
				m.EXPECT().Register(
					gomock.Any(),
					&usecase.RegisterRequest{
						Username: "testuser",
						Password: "testpass",
					},
				).Return(&usecase.RegisterResponse{UserID: 1}, nil)
			},
			want: &pb.RegisterResponse{UserId: 1},
		},
		{
			name: "empty username",
			req: &pb.RegisterRequest{
				Username: "",
				Password: "testpass",
			},
			mockSetup: func(m *MockAuthUsecase) {
				m.EXPECT().Register(
					gomock.Any(),
					&usecase.RegisterRequest{
						Username: "",
						Password: "testpass",
					},
				).Return(nil, errors.New("username cannot be empty"))
			},
			expectedErr: status.New(codes.Internal, "username cannot be empty"),
		},
		{
			name: "empty password",
			req: &pb.RegisterRequest{
				Username: "testuser",
				Password: "",
			},
			mockSetup: func(m *MockAuthUsecase) {
				m.EXPECT().Register(
					gomock.Any(),
					&usecase.RegisterRequest{
						Username: "testuser",
						Password: "",
					},
				).Return(nil, errors.New("password cannot be empty"))
			},
			expectedErr: status.New(codes.Internal, "password cannot be empty"),
		},
		{
			name: "usecase returns error",
			req: &pb.RegisterRequest{
				Username: "testuser",
				Password: "testpass",
			},
			mockSetup: func(m *MockAuthUsecase) {
				m.EXPECT().Register(
					gomock.Any(),
					gomock.Any(),
				).Return(nil, errors.New("database error"))
			},
			expectedErr: status.New(codes.Internal, "database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup(mockUC)
			resp, err := controller.Register(context.Background(), tt.req)

			if tt.expectedErr != nil {
				assert.Error(t, err)
				st, ok := status.FromError(err)
				assert.True(t, ok)
				assert.Equal(t, tt.expectedErr.Code(), st.Code())
				assert.Equal(t, tt.expectedErr.Message(), st.Message())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, resp)
			}
		})
	}
}

func TestAuthController_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := NewMockAuthUsecase(ctrl)
	controller := NewAuthController(mockUC)

	tests := []struct {
		name        string
		req         *pb.LoginRequest
		mockSetup   func(*MockAuthUsecase)
		want        *pb.LoginResponse
		expectedErr *status.Status
	}{
		{
			name: "successful login",
			req: &pb.LoginRequest{
				Username: "testuser",
				Password: "testpass",
			},
			mockSetup: func(m *MockAuthUsecase) {
				m.EXPECT().Login(
					gomock.Any(),
					&usecase.LoginRequest{
						Username: "testuser",
						Password: "testpass",
					},
				).Return(&usecase.LoginResponse{
					Token:    "test_token",
					Username: "testuser",
				}, nil)
			},
			want: &pb.LoginResponse{
				Token:    "test_token",
				Username: "testuser",
			},
		},
		{
			name: "empty username",
			req: &pb.LoginRequest{
				Username: "",
				Password: "testpass",
			},
			mockSetup: func(m *MockAuthUsecase) {
				m.EXPECT().Login(
					gomock.Any(),
					&usecase.LoginRequest{
						Username: "",
						Password: "testpass",
					},
				).Return(nil, errors.New("username cannot be empty"))
			},
			expectedErr: status.New(codes.Internal, "username cannot be empty"),
		},
		{
			name: "empty password",
			req: &pb.LoginRequest{
				Username: "testuser",
				Password: "",
			},
			mockSetup: func(m *MockAuthUsecase) {
				m.EXPECT().Login(
					gomock.Any(),
					&usecase.LoginRequest{
						Username: "testuser",
						Password: "",
					},
				).Return(nil, errors.New("password cannot be empty"))
			},
			expectedErr: status.New(codes.Internal, "password cannot be empty"),
		},
		{
			name: "invalid credentials",
			req: &pb.LoginRequest{
				Username: "wronguser",
				Password: "wrongpass",
			},
			mockSetup: func(m *MockAuthUsecase) {
				m.EXPECT().Login(
					gomock.Any(),
					gomock.Any(),
				).Return(nil, errors.New("invalid credentials"))
			},
			expectedErr: status.New(codes.Internal, "invalid credentials"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup(mockUC)
			resp, err := controller.Login(context.Background(), tt.req)

			if tt.expectedErr != nil {
				assert.Error(t, err)
				st, ok := status.FromError(err)
				assert.True(t, ok)
				assert.Equal(t, tt.expectedErr.Code(), st.Code())
				assert.Equal(t, tt.expectedErr.Message(), st.Message())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, resp)
			}
		})
	}
}

func TestAuthController_GetUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := NewMockAuthUsecase(ctrl)
	controller := NewAuthController(mockUC)

	testTime := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name        string
		req         *pb.GetUserRequest
		mockSetup   func(*MockAuthUsecase)
		want        *pb.GetUserResponse
		expectedErr *status.Status
	}{
		{
			name: "successful get user",
			req:  &pb.GetUserRequest{Id: 1},
			mockSetup: func(m *MockAuthUsecase) {
				m.EXPECT().GetUser(
					gomock.Any(),
					&usecase.GetUserRequest{UserID: 1},
				).Return(&usecase.GetUserResponse{
					User: &entity.User{
						ID:        1,
						Username:  "testuser",
						Role:      entity.RoleAdmin,
						CreatedAt: testTime,
					},
				}, nil)
			},
			want: &pb.GetUserResponse{
				User: &pb.User{
					Id:        1,
					Username:  "testuser",
					Role:      "admin",
					CreatedAt: timestamppb.New(testTime),
				},
			},
		},
		{
			name: "empty user id",
			req:  &pb.GetUserRequest{Id: 0},
			mockSetup: func(m *MockAuthUsecase) {
				m.EXPECT().GetUser(
					gomock.Any(),
					&usecase.GetUserRequest{UserID: 0},
				).Return(nil, errors.New("user id cannot be empty"))
			},
			expectedErr: status.New(codes.Internal, "user id cannot be empty"),
		},
		{
			name: "user not found",
			req:  &pb.GetUserRequest{Id: 999},
			mockSetup: func(m *MockAuthUsecase) {
				m.EXPECT().GetUser(
					gomock.Any(),
					&usecase.GetUserRequest{UserID: 999},
				).Return(nil, errors.New("user not found"))
			},
			expectedErr: status.New(codes.Internal, "user not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup(mockUC)
			resp, err := controller.GetUser(context.Background(), tt.req)

			if tt.expectedErr != nil {
				assert.Error(t, err)
				st, ok := status.FromError(err)
				assert.True(t, ok)
				assert.Equal(t, tt.expectedErr.Code(), st.Code())
				assert.Equal(t, tt.expectedErr.Message(), st.Message())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, resp)
			}
		})
	}
}

func TestAuthController_ValidateToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := NewMockAuthUsecase(ctrl)
	controller := NewAuthController(mockUC)

	tests := []struct {
		name        string
		req         *pb.ValidateTokenRequest
		mockSetup   func(*MockAuthUsecase)
		want        *pb.ValidateTokenResponse
		expectedErr *status.Status
	}{
		{
			name: "valid token",
			req:  &pb.ValidateTokenRequest{Token: "valid_token"},
			mockSetup: func(m *MockAuthUsecase) {
				m.EXPECT().ValidateToken(
					gomock.Any(),
					&usecase.ValidateTokenRequest{Token: "valid_token"},
				).Return(&usecase.ValidateTokenResponse{
					Valid:  true,
					UserID: 1,
					Role:   "admin",
				}, nil)
			},
			want: &pb.ValidateTokenResponse{
				Valid:  true,
				UserId: 1,
				Role:   "admin",
			},
		},
		{
			name: "empty token",
			req:  &pb.ValidateTokenRequest{Token: ""},
			mockSetup: func(m *MockAuthUsecase) {
				m.EXPECT().ValidateToken(
					gomock.Any(),
					&usecase.ValidateTokenRequest{Token: ""},
				).Return(nil, errors.New("token cannot be empty"))
			},
			expectedErr: status.New(codes.Internal, "token cannot be empty"),
		},
		{
			name: "invalid token",
			req:  &pb.ValidateTokenRequest{Token: "invalid_token"},
			mockSetup: func(m *MockAuthUsecase) {
				m.EXPECT().ValidateToken(
					gomock.Any(),
					&usecase.ValidateTokenRequest{Token: "invalid_token"},
				).Return(nil, errors.New("invalid token"))
			},
			expectedErr: status.New(codes.Internal, "invalid token"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup(mockUC)
			resp, err := controller.ValidateToken(context.Background(), tt.req)

			if tt.expectedErr != nil {
				assert.Error(t, err)
				st, ok := status.FromError(err)
				assert.True(t, ok)
				assert.Equal(t, tt.expectedErr.Code(), st.Code())
				assert.Equal(t, tt.expectedErr.Message(), st.Message())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, resp)
			}
		})
	}
}

// [Previous imports and mock implementations remain the same...]

func Test_convertUserToProto(t *testing.T) {
	testTime := time.Now()
	user := &entity.User{
		ID:        1,
		Username:  "testuser",
		Role:      entity.RoleUser,
		CreatedAt: testTime,
	}

	want := &pb.User{
		Id:        1,
		Username:  "testuser",
		Role:      "user",
		CreatedAt: timestamppb.New(testTime),
	}

	got := convertUserToProto(user)
	assert.Equal(t, want, got)
}

func Test_convertUserToProto_NilUser(t *testing.T) {
	assert.Nil(t, convertUserToProto(nil))
}

func Test_convertUserToProto_EmptyUser(t *testing.T) {
	user := &entity.User{}
	expected := &pb.User{
		Id:        0,
		Username:  "",
		Role:      "",
		CreatedAt: timestamppb.New(time.Time{}),
	}
	assert.Equal(t, expected, convertUserToProto(user))
}

func Test_convertUserToProto_ZeroTime(t *testing.T) {
	user := &entity.User{
		ID:        1,
		Username:  "test",
		Role:      entity.RoleUser,
		CreatedAt: time.Time{},
	}
	got := convertUserToProto(user)
	assert.Equal(t, timestamppb.New(time.Time{}), got.CreatedAt)
}

func TestNewAuthController(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := NewMockAuthUsecase(ctrl)
	controller := NewAuthController(mockUC)

	assert.NotNil(t, controller)
	assert.Equal(t, mockUC, controller.uc)
}

func Test_convertUserToProto_DifferentRoles(t *testing.T) {
	tests := []struct {
		name     string
		user     *entity.User
		expected *pb.User
	}{
		{
			name: "admin role",
			user: &entity.User{
				ID:        1,
				Username:  "admin",
				Role:      entity.RoleAdmin,
				CreatedAt: time.Now(),
			},
			expected: &pb.User{
				Id:        1,
				Username:  "admin",
				Role:      "admin",
				CreatedAt: timestamppb.New(time.Now()),
			},
		},
		{
			name: "user role",
			user: &entity.User{
				ID:        2,
				Username:  "user",
				Role:      entity.RoleUser,
				CreatedAt: time.Now(),
			},
			expected: &pb.User{
				Id:        2,
				Username:  "user",
				Role:      "user",
				CreatedAt: timestamppb.New(time.Now()),
			},
		},
		{
			name: "unknown role",
			user: &entity.User{
				ID:        3,
				Username:  "unknown",
				Role:      "custom_role",
				CreatedAt: time.Now(),
			},
			expected: &pb.User{
				Id:        3,
				Username:  "unknown",
				Role:      "custom_role",
				CreatedAt: timestamppb.New(time.Now()),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := convertUserToProto(tt.user)
			assert.Equal(t, tt.expected.Id, got.Id)
			assert.Equal(t, tt.expected.Username, got.Username)
			assert.Equal(t, tt.expected.Role, got.Role)
			assert.True(t, tt.expected.CreatedAt.AsTime().Equal(got.CreatedAt.AsTime()))
		})
	}
}

func TestAuthController_GetUser_UnknownRole(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := NewMockAuthUsecase(ctrl)
	controller := NewAuthController(mockUC)

	testTime := time.Now()
	mockUC.EXPECT().GetUser(gomock.Any(), gomock.Any()).Return(&usecase.GetUserResponse{
		User: &entity.User{
			ID:        1,
			Username:  "special",
			Role:      "special_role", // Нестандартная роль
			CreatedAt: testTime,
		},
	}, nil)

	resp, err := controller.GetUser(context.Background(), &pb.GetUserRequest{Id: 1})
	assert.NoError(t, err)
	assert.Equal(t, "special_role", resp.User.Role)
}

func TestAuthController_ValidateToken_Roles(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := NewMockAuthUsecase(ctrl)
	controller := NewAuthController(mockUC)

	tests := []struct {
		name string
		role string
	}{
		{"admin role", "admin"},
		{"user role", "user"},
		{"unknown role", "guest"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUC.EXPECT().ValidateToken(gomock.Any(), gomock.Any()).Return(&usecase.ValidateTokenResponse{
				Valid:  true,
				UserID: 1,
				Role:   tt.role,
			}, nil)

			resp, err := controller.ValidateToken(context.Background(), &pb.ValidateTokenRequest{Token: "test"})
			assert.NoError(t, err)
			assert.Equal(t, tt.role, resp.Role)
		})
	}
}
func TestConvertUserToProto_MinValues(t *testing.T) {
	user := &entity.User{
		ID:        -1, // Отрицательный ID
		Username:  "",
		Role:      "",
		CreatedAt: time.Time{},
	}

	got := convertUserToProto(user)
	assert.Equal(t, int64(-1), got.Id)
	assert.Equal(t, "", got.Username)
	assert.Equal(t, "", got.Role)
	assert.Equal(t, timestamppb.New(time.Time{}), got.CreatedAt)
}

func TestAuthController_Register_InvalidInput(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := NewMockAuthUsecase(ctrl)
	controller := NewAuthController(mockUC)

	// Test with nil request
	_, err := controller.Register(context.Background(), nil)
	assert.Error(t, err)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
	assert.Equal(t, "request cannot be nil", st.Message())
}
func TestAuthController_Login_NilRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := NewMockAuthUsecase(ctrl)
	controller := NewAuthController(mockUC)

	_, err := controller.Login(context.Background(), nil)
	assert.Error(t, err)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
	assert.Equal(t, "request cannot be nil", st.Message())
}

func TestAuthController_GetUser_NilRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := NewMockAuthUsecase(ctrl)
	controller := NewAuthController(mockUC)

	_, err := controller.GetUser(context.Background(), nil)
	assert.Error(t, err)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
	assert.Equal(t, "request cannot be nil", st.Message())
}

func TestAuthController_ValidateToken_NilRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := NewMockAuthUsecase(ctrl)
	controller := NewAuthController(mockUC)

	_, err := controller.ValidateToken(context.Background(), nil)
	assert.Error(t, err)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
	assert.Equal(t, "request cannot be nil", st.Message())
}

func TestAuthController_Register_InvalidArguments(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := NewMockAuthUsecase(ctrl)
	controller := NewAuthController(mockUC)

	tests := []struct {
		name string
		req  *pb.RegisterRequest
		err  string
	}{
		{
			name: "nil request",
			req:  nil,
			err:  "request cannot be nil",
		},
		{
			name: "empty username and password",
			req:  &pb.RegisterRequest{Username: "", Password: ""},
			err:  "username and password cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.req != nil {
				mockUC.EXPECT().Register(gomock.Any(), gomock.Any()).Return(nil, errors.New(tt.err))
			}

			_, err := controller.Register(context.Background(), tt.req)
			assert.Error(t, err)
			st, ok := status.FromError(err)
			assert.True(t, ok)
			if tt.req == nil {
				assert.Equal(t, codes.InvalidArgument, st.Code())
			} else {
				assert.Equal(t, codes.Internal, st.Code())
			}
			assert.Contains(t, st.Message(), tt.err)
		})
	}
}

func TestAuthController_GetUser_InvalidID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := NewMockAuthUsecase(ctrl)
	controller := NewAuthController(mockUC)

	mockUC.EXPECT().GetUser(gomock.Any(), &usecase.GetUserRequest{UserID: -1}).
		Return(nil, errors.New("invalid user id"))

	_, err := controller.GetUser(context.Background(), &pb.GetUserRequest{Id: -1})
	assert.Error(t, err)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Internal, st.Code())
	assert.Equal(t, "invalid user id", st.Message())
}

func TestAuthController_ValidateToken_InvalidResponse(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := NewMockAuthUsecase(ctrl)
	controller := NewAuthController(mockUC)

	mockUC.EXPECT().ValidateToken(gomock.Any(), &usecase.ValidateTokenRequest{Token: "invalid"}).
		Return(nil, errors.New("token validation failed"))

	_, err := controller.ValidateToken(context.Background(), &pb.ValidateTokenRequest{Token: "invalid"})
	assert.Error(t, err)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Internal, st.Code())
	assert.Equal(t, "token validation failed", st.Message())
}
