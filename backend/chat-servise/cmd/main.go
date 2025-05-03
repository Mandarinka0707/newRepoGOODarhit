package main

import (
	"chat-microservice-go/internal/handler"
	"chat-microservice-go/internal/repository"
	"chat-microservice-go/internal/usecase"

	"github.com/gin-contrib/cors"

	"database/sql"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	connStr := "postgres://user:password@localhost:5432/database?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	repo := repository.NewMessageRepository(db)
	uc := usecase.NewMessageUseCase(repo)
	h := handler.NewMessageHandler(uc)

	go h.HandleMessages()

	r := gin.Default()
	r.Use(cors.Default())
	r.GET("/ws", h.HandleConnections)
	r.GET("/messages", h.GetMessages)

	log.Println("Listening on :8082...")
	log.Fatal(r.Run(":8082"))
}
