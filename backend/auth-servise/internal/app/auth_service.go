package app

import (
	"context"
	"time"

	"backend.com/forum/auth-servise/internal/domain"
	"backend.com/forum/auth-servise/internal/repository"
	"backend.com/forum/auth-servise/pkg/auth"
	"backend.com/forum/auth-servise/pkg/logger"
	pb "backend.com/forum/proto"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthService struct {
	userRepo    repository.UserRepository
	sessionRepo repository.SessionRepository
	cfg         *Config // Используйте Config для настроек
	logger      *logger.Logger
	pb.UnimplementedAuthServiceServer
}

type Config struct { // Конфигурационная структура
	TokenSecret     string
	TokenExpiration time.Duration
}

func NewAuthService(userRepo repository.UserRepository, sessionRepo repository.SessionRepository, cfg *Config, logger *logger.Logger) *AuthService {
	return &AuthService{userRepo: userRepo, sessionRepo: sessionRepo, cfg: cfg, logger: logger}
}

func (s *AuthService) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	s.logger.Infow("Register request received", "username", req.Username)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Errorw("Failed to hash password", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to hash password")
	}

	user := &domain.User{
		Username:  req.Username,
		Password:  string(hashedPassword),
		Role:      domain.RoleUser, // По умолчанию роль user
		CreatedAt: time.Now(),
	}

	userID, err := s.userRepo.CreateUser(ctx, user)
	if err != nil {
		s.logger.Errorw("Failed to create user", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to create user")
	}

	s.logger.Infow("User registered successfully", "username", req.Username, "user_id", userID)
	return &pb.RegisterResponse{UserId: userID}, nil
}

func (s *AuthService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	s.logger.Infow("Login request received", "username", req.Username)
	user, err := s.userRepo.GetUserByUsername(ctx, req.Username)
	if err != nil {
		s.logger.Errorw("Failed to get user by username", "username", req.Username, "error", err)
		return nil, status.Errorf(codes.NotFound, "invalid username or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		s.logger.Errorw("Invalid password", "username", req.Username, "error", err)
		return nil, status.Errorf(codes.NotFound, "invalid username or password")
	}

	token, err := auth.GenerateToken(user.ID, user.Role, s.cfg.TokenSecret, s.cfg.TokenExpiration)
	if err != nil {
		s.logger.Errorw("Failed to generate token", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to generate token")
	}

	session := &domain.Session{
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: time.Now().Add(s.cfg.TokenExpiration),
	}

	if err := s.sessionRepo.CreateSession(ctx, session); err != nil {
		s.logger.Errorw("Failed to create session", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to create session")
	}

	s.logger.Infow("Login successful", "username", req.Username, "user_id", user.ID)
	return &pb.LoginResponse{Token: token}, nil
}

func (s *AuthService) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
	s.logger.Infow("ValidateToken request received", "token", req.Token)

	token, err := auth.ParseToken(req.Token, s.cfg.TokenSecret)

	if err != nil {
		s.logger.Errorw("Failed to parse token", "error", err)
		return &pb.ValidateTokenResponse{Valid: false}, nil
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		s.logger.Warnw("Invalid token", "error", err)
		return &pb.ValidateTokenResponse{Valid: false}, nil
	}

	userID, ok := claims["user_id"].(float64)
	if !ok {
		s.logger.Warnw("Invalid user_id in token")
		return &pb.ValidateTokenResponse{Valid: false}, nil
	}
	role, ok := claims["role"].(string)
	if !ok {
		s.logger.Warnw("Invalid role in token")
		return &pb.ValidateTokenResponse{Valid: false}, nil
	}
	s.logger.Infow("Token is valid", "user_id", userID)
	return &pb.ValidateTokenResponse{Valid: true, UserId: int64(userID), Role: role}, nil
}
