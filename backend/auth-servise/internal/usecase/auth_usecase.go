// auth_usecase.go
package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/Mandarinka0707/newRepoGOODarhit/internal/entity"
	"github.com/Mandarinka0707/newRepoGOODarhit/internal/repository"
	"github.com/Mandarinka0707/newRepoGOODarhit/pkg/auth"
	"github.com/dgrijalva/jwt-go"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type AuthUsecase struct {
	userRepo    repository.UserRepository
	sessionRepo repository.SessionRepository
	cfg         *auth.Config
	logger      *zap.Logger
}

type AuthUsecaseInterface interface {
	Register(ctx context.Context, req *RegisterRequest) (*RegisterResponse, error)
	Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error)
	GetUserByID(ctx context.Context, userID int64) (*entity.User, error)
	ValidateToken(ctx context.Context, req *ValidateTokenRequest) (*ValidateTokenResponse, error)
	GetUser(ctx context.Context, req *GetUserRequest) (*GetUserResponse, error)
}

func NewAuthUsecase(
	userRepo repository.UserRepository,
	sessionRepo repository.SessionRepository,
	cfg *auth.Config,
	logger *zap.Logger,
) *AuthUsecase {
	return &AuthUsecase{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
		cfg:         cfg,
		logger:      logger,
	}
}

func (uc *AuthUsecase) Register(
	ctx context.Context,
	req *RegisterRequest,
) (*RegisterResponse, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		uc.logger.Error("Failed to hash password", zap.Error(err))
		return nil, err
	}

	user := &entity.User{
		Username:  req.Username,
		Password:  string(hashedPassword),
		Role:      entity.RoleUser,
		CreatedAt: time.Now(),
	}

	userID, err := uc.userRepo.CreateUser(ctx, user)
	if err != nil {
		uc.logger.Error("Failed to create user", zap.Error(err))
		return nil, err
	}

	return &RegisterResponse{UserID: userID}, nil
}

func (uc *AuthUsecase) Login(
	ctx context.Context,
	req *LoginRequest,
) (*LoginResponse, error) {
	user, err := uc.userRepo.GetUserByUsername(ctx, req.Username)
	if err != nil {
		return nil, fmt.Errorf("invalid username or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, fmt.Errorf("invalid username or password")
	}

	token, err := auth.GenerateToken(
		user.ID,
		user.Role,
		user.Username,
		uc.cfg.TokenSecret,
		uc.cfg.TokenExpiration,
	)
	if err != nil {
		uc.logger.Error("failed to generate token", zap.Error(err))
		return nil, fmt.Errorf("internal server error")
	}

	session := &entity.Session{
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: time.Now().Add(uc.cfg.TokenExpiration),
	}

	if err := uc.sessionRepo.CreateSession(ctx, session); err != nil {
		uc.logger.Error("failed to create session", zap.Error(err))
		return nil, fmt.Errorf("internal server error")
	}

	return &LoginResponse{
		Token:    token,
		Username: user.Username,
	}, nil
}

func (uc *AuthUsecase) GetUserByID(
	ctx context.Context,
	userID int64,
) (*entity.User, error) {
	return uc.userRepo.GetUserByID(ctx, userID)
}

func (uc *AuthUsecase) ValidateToken(
	ctx context.Context,
	req *ValidateTokenRequest,
) (*ValidateTokenResponse, error) {
	uc.logger.Info("Token validation request")

	token, err := auth.ParseToken(req.Token, uc.cfg.TokenSecret)
	if err != nil || !token.Valid {
		uc.logger.Warn("Invalid token", zap.Error(err))
		return &ValidateTokenResponse{Valid: false}, nil
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		uc.logger.Warn("Invalid token claims")
		return &ValidateTokenResponse{Valid: false}, nil
	}

	userID, ok := claims["user_id"].(float64)
	if !ok {
		uc.logger.Warn("Invalid user_id in token")
		return &ValidateTokenResponse{Valid: false}, nil
	}

	role, ok := claims["role"].(string)
	if !ok {
		uc.logger.Warn("Invalid role in token")
		return &ValidateTokenResponse{Valid: false}, nil
	}

	return &ValidateTokenResponse{
		Valid:  true,
		UserID: int64(userID),
		Role:   role,
	}, nil
}

func (uc *AuthUsecase) GetUser(
	ctx context.Context,
	req *GetUserRequest,
) (*GetUserResponse, error) {
	uc.logger.Info("Get user request", zap.Int64("user_id", req.UserID))

	user, err := uc.userRepo.GetUserByID(ctx, req.UserID)
	if err != nil || user == nil {
		uc.logger.Error("User not found", zap.Error(err))
		return nil, fmt.Errorf("user not found")
	}

	return &GetUserResponse{User: user}, nil
}
