package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"backend.com/forum/forum-servise/internal/handler"
	"backend.com/forum/forum-servise/internal/repository"
	"backend.com/forum/forum-servise/internal/usecase"
	"backend.com/forum/forum-servise/pkg/logger"
	pb "backend.com/forum/proto"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // PostgreSQL driver
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

	// 3. Подключение к gRPC сервису аутентификации
	authConn, err := grpc.Dial(
		"localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Error("Failed to connect to auth service", err)
		os.Exit(1)
	}
	defer authConn.Close()

	authClient := pb.NewAuthServiceClient(authConn)

	// 4. Инициализация репозиториев
	postRepo := repository.NewPostRepository(db)

	// 5. Инициализация юзкейсов
	postUsecase := usecase.NewPostUsecase(
		postRepo,
		authClient,
		log,
	)

	// 6. Инициализация хендлеров
	postHandler := handler.NewPostHandler(postUsecase, log)

	// 7. Настройка маршрутизатора
	router.POST("/api/v1/posts", postHandler.CreatePost)
	router.GET("/api/v1/posts", postHandler.GetPosts)

	// 8. Настройка HTTP сервера
	server := &http.Server{
		Addr:    ":8081",
		Handler: router,
	}

	// 9. Запуск сервера
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("Server error", err)
			os.Exit(1)
		}
	}()

	log.Info("Server started on :8081") // Исправлен порт в сообщении

	// 10. Graceful shutdown
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
