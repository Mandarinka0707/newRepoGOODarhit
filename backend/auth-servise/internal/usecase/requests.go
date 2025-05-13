package usecase

import "time"

type RegisterRequest struct {
	Username string `bson:"user_name"b json:"uaer_name"`
	Password string `bson:"password" json:"password"`
}

type LoginRequest struct {
	Username string `bson:"user_name"b json:"uaer_name"`
	Password string `bson:"password" json:"password"`
}

type ValidateTokenRequest struct {
	Token string
}

type GetUserRequest struct {
	UserID int64
}
type Config struct {
	TokenSecret     string
	TokenExpiration time.Duration
}
