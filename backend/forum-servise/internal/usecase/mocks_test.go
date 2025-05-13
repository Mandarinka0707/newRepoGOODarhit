package usecase

import (
	"context"

	pb "backend.com/forum/proto"
	"github.com/Mandarinka0707/newRepoGOODarhit/forum-servise/internal/entity"
	"google.golang.org/grpc"
)

type MockPostRepository struct {
	CreatePostFunc  func(ctx context.Context, post *entity.Post) (int64, error)
	GetPostsFunc    func(ctx context.Context) ([]*entity.Post, error)
	GetPostByIDFunc func(ctx context.Context, id int64) (*entity.Post, error)
	DeletePostFunc  func(ctx context.Context, postID, authorID int64, role string) error
	UpdatePostFunc  func(ctx context.Context, postID, authorID int64, role, title, content string) (*entity.Post, error)
}

func (m *MockPostRepository) CreatePost(ctx context.Context, post *entity.Post) (int64, error) {
	if m.CreatePostFunc != nil {
		return m.CreatePostFunc(ctx, post)
	}
	return 0, nil
}

func (m *MockPostRepository) GetPosts(ctx context.Context) ([]*entity.Post, error) {
	if m.GetPostsFunc != nil {
		return m.GetPostsFunc(ctx)
	}
	return nil, nil
}

func (m *MockPostRepository) GetPostByID(ctx context.Context, id int64) (*entity.Post, error) {
	if m.GetPostByIDFunc != nil {
		return m.GetPostByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockPostRepository) DeletePost(ctx context.Context, postID, authorID int64, role string) error {
	if m.DeletePostFunc != nil {
		return m.DeletePostFunc(ctx, postID, authorID, role)
	}
	return nil
}

func (m *MockPostRepository) UpdatePost(ctx context.Context, postID, authorID int64, role, title, content string) (*entity.Post, error) {
	if m.UpdatePostFunc != nil {
		return m.UpdatePostFunc(ctx, postID, authorID, role, title, content)
	}
	return nil, nil
}

type MockAuthServiceClient struct {
	ValidateTokenFunc func(ctx context.Context, in *pb.ValidateTokenRequest, opts ...grpc.CallOption) (*pb.ValidateTokenResponse, error)
	GetUserFunc       func(ctx context.Context, in *pb.GetUserRequest, opts ...grpc.CallOption) (*pb.GetUserResponse, error)
	LoginFunc         func(ctx context.Context, in *pb.LoginRequest, opts ...grpc.CallOption) (*pb.LoginResponse, error)
	RegisterFunc      func(ctx context.Context, in *pb.RegisterRequest, opts ...grpc.CallOption) (*pb.RegisterResponse, error)
}

func (m *MockAuthServiceClient) ValidateToken(ctx context.Context, in *pb.ValidateTokenRequest, opts ...grpc.CallOption) (*pb.ValidateTokenResponse, error) {
	if m.ValidateTokenFunc != nil {
		return m.ValidateTokenFunc(ctx, in, opts...)
	}
	return nil, nil
}

func (m *MockAuthServiceClient) GetUser(ctx context.Context, in *pb.GetUserRequest, opts ...grpc.CallOption) (*pb.GetUserResponse, error) {
	if m.GetUserFunc != nil {
		return m.GetUserFunc(ctx, in, opts...)
	}
	return nil, nil
}

func (m *MockAuthServiceClient) Login(ctx context.Context, in *pb.LoginRequest, opts ...grpc.CallOption) (*pb.LoginResponse, error) {
	if m.LoginFunc != nil {
		return m.LoginFunc(ctx, in, opts...)
	}
	return nil, nil
}

func (m *MockAuthServiceClient) Register(ctx context.Context, in *pb.RegisterRequest, opts ...grpc.CallOption) (*pb.RegisterResponse, error) {
	if m.RegisterFunc != nil {
		return m.RegisterFunc(ctx, in, opts...)
	}
	return nil, nil
}
