package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	pb "backend.com/forum/proto"
	_ "github.com/Mandarinka0707/newRepoGOODarhit/forum-servise/docs"
	"github.com/Mandarinka0707/newRepoGOODarhit/forum-servise/internal/handler"
	"github.com/Mandarinka0707/newRepoGOODarhit/forum-servise/internal/repository"
	"github.com/Mandarinka0707/newRepoGOODarhit/forum-servise/internal/usecase"
	"github.com/Mandarinka0707/newRepoGOODarhit/forum-servise/pkg/logger"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	log, err := logger.NewLogger("debug")
	if err != nil {
		fmt.Printf("Failed to create logger: %v\n", err)
		os.Exit(1)
	}

	// Подключение к PostgreSQL
	db, err := sqlx.Connect("postgres", "postgres://user:password@localhost:5432/database?sslmode=disable")
	if err != nil {
		log.Error("Failed to connect to database", err)
		os.Exit(1)
	}
	defer db.Close()

	// Инициализация Gin
	router := gin.Default()

	// Настройка CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Подключение к Auth Service
	authConn, err := grpc.Dial(
		"localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Error("Failed to connect to auth service", err)
		os.Exit(1)
	}
	defer authConn.Close()

	// Инициализация репозиториев и usecases
	authClient := pb.NewAuthServiceClient(authConn)
	postRepo := repository.NewPostRepository(db)
	commentRepo := repository.NewCommentRepository(db)
	postUsecase := usecase.NewPostUsecase(postRepo, authClient, log)
	commentUC := usecase.NewCommentUseCase(commentRepo, postRepo, authClient)

	// Регистрация обработчиков
	postHandler := handler.NewPostHandler(postUsecase, log)
	commentHandler := handler.NewCommentHandler(commentUC)

	// Группировка роутов
	api := router.Group("/api/v1")
	{
		// Роуты для постов
		posts := api.Group("/posts")
		{
			posts.POST("", postHandler.CreatePost)
			posts.GET("", postHandler.GetPosts)
			posts.DELETE("/:id", postHandler.DeletePost)
			posts.PUT("/:id", postHandler.UpdatePost)
		}

		// Роуты для комментариев
		comments := api.Group("/posts/:id/comments")
		{
			comments.POST("", commentHandler.CreateComment)
			comments.GET("", commentHandler.GetCommentsByPostID)
		}
	}

	// Запуск сервера
	server := &http.Server{
		Addr:    ":8081",
		Handler: router,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("Server error", err)
			os.Exit(1)
		}
	}()

	log.Info("Server started on :8081")

	// Graceful shutdown
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
