// controller/auth_grpc.go
package controller

import (
	"context"

	"backend.com/forum/auth-servise/internal/entity"
	"backend.com/forum/auth-servise/internal/usecase"
	pb "backend.com/forum/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type AuthController struct {
	uc *usecase.AuthUsecase
	pb.UnimplementedAuthServiceServer
}

func NewAuthController(uc *usecase.AuthUsecase) *AuthController {
	return &AuthController{uc: uc}
}

// Register регистрирует пользователя через gRPC
//
// @Summary Регистрация пользователя через gRPC
// @Description Регистрирует нового пользователя, используя gRPC
// @Tags auth
// @Accept json
// @Produce json
// @Failure 500 {object} entity.ErrorResponse
// @Router /auth/register [post]
func (c *AuthController) Register(
	ctx context.Context,
	req *pb.RegisterRequest,
) (*pb.RegisterResponse, error) {
	ucReq := &usecase.RegisterRequest{
		Username: req.Username,
		Password: req.Password,
	}

	ucResp, err := c.uc.Register(ctx, ucReq)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.RegisterResponse{UserId: ucResp.UserID}, nil
}

// Login выполняет аутентификацию пользователя через gRPC
//
// @Summary Логин пользователя через gRPC
// @Description Выполняет аутентификацию пользователя, используя gRPC
// @Tags auth
// @Accept json
// @Produce json
// @Failure 500 {object} entity.ErrorResponse
// @Router /auth/login [post]
func (c *AuthController) Login(
	ctx context.Context,
	req *pb.LoginRequest,
) (*pb.LoginResponse, error) {
	ucReq := &usecase.LoginRequest{
		Username: req.Username,
		Password: req.Password,
	}

	ucResp, err := c.uc.Login(ctx, ucReq)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.LoginResponse{
		Token:    ucResp.Token,
		Username: ucResp.Username, // Добавьте это
	}, nil
}

// GetUser получает информацию о пользователе через gRPC
//
// @Summary Получить информацию о пользователе через gRPC
// @Description Получает информацию о пользователе по ID
// @Tags auth
// @Accept json
// @Produce json
// @Failure 500 {object} entity.ErrorResponse
// @Param id path string true "ID пользователя"
// @Router /auth/user/{id} [get]
func (c *AuthController) GetUser(
	ctx context.Context,
	req *pb.GetUserRequest,
) (*pb.GetUserResponse, error) {
	ucReq := &usecase.GetUserRequest{UserID: req.Id}

	ucResp, err := c.uc.GetUser(ctx, ucReq)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.GetUserResponse{
		User: convertUserToProto(ucResp.User),
	}, nil
}

func convertUserToProto(user *entity.User) *pb.User {
	return &pb.User{
		Id:        user.ID,
		Username:  user.Username,
		Role:      string(user.Role),
		CreatedAt: timestamppb.New(user.CreatedAt),
	}
}

// ValidateToken проверяет валидность токена через gRPC
//
// @Summary Проверка валидности токена через gRPC
// @Description Проверяет валидность токена и возвращает информацию о пользователе
// @Tags auth
// @Accept json
// @Produce json
// @Failure 500 {object} entity.ErrorResponse
// @Router /auth/validate-token [post]
func (c *AuthController) ValidateToken(
	ctx context.Context,
	req *pb.ValidateTokenRequest,
) (*pb.ValidateTokenResponse, error) {
	ucReq := &usecase.ValidateTokenRequest{Token: req.Token}

	ucResp, err := c.uc.ValidateToken(ctx, ucReq) // Добавлен контекст
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.ValidateTokenResponse{
		Valid:  ucResp.Valid,
		UserId: ucResp.UserID,
		Role:   ucResp.Role,
	}, nil
}
