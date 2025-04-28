package usecase

import "backend.com/forum/auth-servise/internal/entity"

type RegisterResponse struct {
	UserID int64
}

type LoginResponse struct {
	Token string
}

type ValidateTokenResponse struct {
	UserID int64
	Role   string
	Valid  bool
}

type GetUserResponse struct {
	User *entity.User
}
