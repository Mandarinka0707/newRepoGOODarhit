package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "backend.com/forum/forum-servise/docs"
	"backend.com/forum/forum-servise/internal/handler"
	"backend.com/forum/forum-servise/internal/repository"
	"backend.com/forum/forum-servise/internal/usecase"
	"backend.com/forum/forum-servise/pkg/logger"
	pb "backend.com/forum/proto"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // PostgreSQL driver
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// @title Forum Service API
// @version 1.0
// @description API for managing forum posts and comments
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8081
// @BasePath /api/v1
// @schemes http

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	// 1. Initialize logger
	log, err := logger.NewLogger("debug")
	if err != nil {
		fmt.Printf("Failed to create logger: %v\n", err)
		os.Exit(1)
	}

	// 2. Connect to database
	db, err := sqlx.Connect("postgres", "postgres://user:password@localhost:5432/database?sslmode=disable")
	if err != nil {
		log.Error("Failed to connect to database", err)
		os.Exit(1)
	}
	defer db.Close()

	// Verify tables exist
	_, err = db.ExecContext(context.Background(), "SELECT 1 FROM comments LIMIT 1")
	if err != nil {
		log.Fatal("Comments table does not exist or inaccessible: ", err)
	}

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

	// Swagger documentation route
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 3. Connect to auth gRPC service
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

	// 4. Initialize repositories
	postRepo := repository.NewPostRepository(db)
	commentRepo := repository.NewCommentRepository(db)

	// 5. Initialize use cases
	postUsecase := usecase.NewPostUsecase(
		postRepo,
		authClient,
		log,
	)
	commentUC := usecase.NewCommentUseCase(commentRepo, postRepo, authClient)

	// 6. Initialize handlers
	postHandler := handler.NewPostHandler(postUsecase, log)
	commentHandler := handler.NewCommentHandler(commentUC)

	// 7. Setup API routes
	api := router.Group("/api/v1")
	{
		// Post routes
		api.POST("/posts", postHandler.CreatePost)
		api.GET("/posts", postHandler.GetPosts)
		api.DELETE("/posts/:id", postHandler.DeletePost)
		api.PUT("/posts/:id", postHandler.UpdatePost)

		// Comment routes
		api.POST("/posts/:id/comments", commentHandler.CreateComment)
		api.GET("/posts/:id/comments", commentHandler.GetCommentsByPostID)
	}

	// 8. Configure HTTP server
	server := &http.Server{
		Addr:    ":8081",
		Handler: router,
	}

	// 9. Start server
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("Server error", err)
			os.Exit(1)
		}
	}()

	log.Info("Server started on :8081")

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
