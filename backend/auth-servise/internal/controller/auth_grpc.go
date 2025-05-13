// controller/auth_grpc.go
package controller

import (
	"context"

	pb "backend.com/forum/proto"
	"github.com/Mandarinka0707/newRepoGOODarhit/internal/entity"
	"github.com/Mandarinka0707/newRepoGOODarhit/internal/usecase"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type AuthController struct {
	uc usecase.AuthUsecaseInterface
	pb.UnimplementedAuthServiceServer
}

func NewAuthController(uc usecase.AuthUsecaseInterface) *AuthController {
	return &AuthController{uc: uc}
}

func (c *AuthController) Register(
	ctx context.Context,
	req *pb.RegisterRequest,
) (*pb.RegisterResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be nil")
	}

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
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be nil")
	}

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
		Username: ucResp.Username,
	}, nil
}

func (c *AuthController) GetUser(
	ctx context.Context,
	req *pb.GetUserRequest,
) (*pb.GetUserResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be nil")
	}

	ucReq := &usecase.GetUserRequest{UserID: req.Id}

	ucResp, err := c.uc.GetUser(ctx, ucReq)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.GetUserResponse{
		User: convertUserToProto(ucResp.User),
	}, nil
}

// auth_grpc.go
func convertUserToProto(user *entity.User) *pb.User {
	if user == nil {
		return nil
	}

	var role string
	switch user.Role {
	case entity.RoleAdmin:
		role = "admin"
	case entity.RoleUser:
		role = "user"
	default:
		role = string(user.Role) // для любых других значений
	}

	return &pb.User{
		Id:        user.ID,
		Username:  user.Username,
		Role:      role,
		CreatedAt: timestamppb.New(user.CreatedAt),
	}
}
func (c *AuthController) ValidateToken(
	ctx context.Context,
	req *pb.ValidateTokenRequest,
) (*pb.ValidateTokenResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be nil")
	}

	ucReq := &usecase.ValidateTokenRequest{Token: req.Token}

	ucResp, err := c.uc.ValidateToken(ctx, ucReq)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.ValidateTokenResponse{
		Valid:  ucResp.Valid,
		UserId: ucResp.UserID,
		Role:   ucResp.Role,
	}, nil
}
