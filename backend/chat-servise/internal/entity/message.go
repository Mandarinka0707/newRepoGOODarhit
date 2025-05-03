// internal/entity/message.go
package entity

type Message struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Message  string `json:"message"`
}
