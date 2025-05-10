package mocks

import (
	"context"
	"strconv"
	"testing"
	"time"

	"backend.com/forum/auth-servise/internal/controller"
	"backend.com/forum/auth-servise/internal/repository"
	"backend.com/forum/auth-servise/internal/usecase"
	"backend.com/forum/auth-servise/pkg/auth"
	pb "backend.com/forum/proto"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

func TestGRPCAuthIntegration(t *testing.T) {
	logger := zaptest.NewLogger(t)
	userRepo := repository.NewUserRepository(db)
	sessionRepo := repository.NewSessionRepository(db)
	cfg := &auth.Config{
		TokenSecret:     "test-secret",
		TokenExpiration: 24 * time.Hour,
	}

	uc := usecase.NewAuthUsecase(userRepo, sessionRepo, cfg, logger)
	grpcCtrl := controller.NewAuthController(uc)

	t.Run("Complete GRPC Auth Flow", func(t *testing.T) {
		ctx := context.Background()
		username := "grpcuser_" + strconv.FormatInt(time.Now().UnixNano(), 10)

		// 1. Test registration
		registerResp, err := grpcCtrl.Register(ctx, &pb.RegisterRequest{
			Username: username,
			Password: "grpcpass123",
		})
		require.NoError(t, err)
		require.Greater(t, registerResp.UserId, int64(0))

		// 2. Test login
		loginResp, err := grpcCtrl.Login(ctx, &pb.LoginRequest{
			Username: username,
			Password: "grpcpass123",
		})
		require.NoError(t, err)
		require.NotEmpty(t, loginResp.Token)

		// 3. Test get user
		userResp, err := grpcCtrl.GetUser(ctx, &pb.GetUserRequest{
			Id: registerResp.UserId,
		})
		require.NoError(t, err)
		require.Equal(t, username, userResp.User.Username)
	})

	t.Run("GRPC Error Cases", func(t *testing.T) {
		ctx := context.Background()

		// Test invalid registration (too long password)
		_, err := grpcCtrl.Register(ctx, &pb.RegisterRequest{
			Username: "invaliduser",
			Password: string(make([]byte, 100)),
		})
		require.Error(t, err)

		// Test login with invalid credentials
		_, err = grpcCtrl.Login(ctx, &pb.LoginRequest{
			Username: "nonexistent",
			Password: "wrongpass",
		})
		require.Error(t, err)
		require.Contains(t, err.Error(), "invalid username or password")

		// Test get user with invalid ID
		_, err = grpcCtrl.GetUser(ctx, &pb.GetUserRequest{
			Id: 999999, // Non-existent user
		})
		require.Error(t, err)
		require.Contains(t, err.Error(), "user not found")
	})
}
