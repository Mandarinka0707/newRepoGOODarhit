// internal/entity/message.go
package entity

type Message struct {
	ID       int    `json:"id" example:"1"`
	Username string `json:"username" example:"john_doe"`
	Message  string `json:"message" example:"Hello, world!"`
}