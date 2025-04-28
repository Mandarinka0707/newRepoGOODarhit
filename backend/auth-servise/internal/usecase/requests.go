package usecase

import "time"

type RegisterRequest struct {
	Username string
	Password string
}

type LoginRequest struct {
	Username string
	Password string
}

type ValidateTokenRequest struct {
	Token string
}

type GetUserRequest struct {
	UserID int64
}
type Config struct { // Конфигурационная структура
	TokenSecret     string
	TokenExpiration time.Duration
}
