package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"backend.com/forum/chat-servise/internal/controller"
	"backend.com/forum/chat-servise/internal/repository"
	"backend.com/forum/chat-servise/internal/usecase"
	"backend.com/forum/chat-servise/pkg/logger"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // PostgreSQL driver
)

func main() {
	// 1. Инициализация логгера
	log, err := logger.NewLogger("debug")
	if err != nil {
		fmt.Printf("Failed to create logger: %v\n", err)
		os.Exit(1)
	}

	// 2. Подключение к базе данных
	db, err := sqlx.Connect("postgres", "postgres://user:password@localhost:5432/database?sslmode=disable")
	if err != nil {
		log.Error("Failed to connect to database", err)
		os.Exit(1)
	}
	defer db.Close()

	// 3. Проверка соединения с базой данных
	if err := db.Ping(); err != nil {
		log.Error("Database ping failed", err)
		os.Exit(1)
	}

	// 4. Инициализация репозиториев
	chatMessageRepo := repository.NewChatMessageRepository(db)

	// 5. Конфигурация юзкейса
	chatConfig := &usecase.ChatConfig{
		MessageTTL: 24 * time.Hour,
	}

	// 6. Инициализация юзкейса
	chatUsecase := usecase.NewChatUsecase(
		chatMessageRepo,
		chatConfig,
		log,
	)

	// 7. Инициализация контроллера
	chatController := controller.NewChatController(
		chatUsecase,
		log,
	)

	// 8. Настройка HTTP сервера
	router := gin.Default()

	// CORS configuration
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// WebSocket endpoint
	router.GET("/ws", chatController.WebSocketHandler)

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// 9. Graceful shutdown
	server := &http.Server{
		Addr:    ":8082",
		Handler: router,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("Server error", err)
			os.Exit(1)
		}
	}()

	log.Info("Server started on :8080")

	// Ожидание сигнала для завершения
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Error("Server shutdown error", err)
	}

	log.Info("Server stopped")
}
