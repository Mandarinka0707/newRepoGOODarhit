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

	return &pb.LoginResponse{Token: ucResp.Token}, nil
}

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
