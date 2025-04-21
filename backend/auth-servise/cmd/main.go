package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"backend.com/forum/auth-servise/internal/app"
	"backend.com/forum/auth-servise/internal/repository"
	"backend.com/forum/auth-servise/pkg/logger"

	pb "backend.com/forum/proto"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	port            = flag.String("port", ":50051", "the port to serve on")
	httpPort        = flag.String("http_port", ":8080", "the HTTP port to serve on")
	dbURL           = flag.String("db_url", "postgres://user:password@localhost:5432/database?sslmode=disable", "PostgreSQL connection string")
	migrationsPath  = flag.String("migrations_path", "/Users/darinautalieva/Desktop/GOProject/backend/auth-servise/migrations", "path to migrations files")
	tokenSecret     = flag.String("token_secret", "your-secret-key", "JWT secret key")
	tokenExpiration = flag.Duration("token_expiration", time.Minute*15, "JWT token expiration") // по умолчанию 1 час
	logLevel        = flag.String("log_level", "info", "Log level (debug, info, warn, error)")
)

func main() {
	flag.Parse()

	// Инициализация логгера
	logger, err := logger.NewLogger(*logLevel)
	if err != nil {
		log.Fatalf("failed to initialize logger: %v", err)
	}
	defer func() {
		err := logger.Sync()
		if err != nil {
			log.Printf("failed to sync logger: %v", err)
		}
	}()

	// Подключение к БД
	db, err := sqlx.Connect("postgres", *dbURL)
	if err != nil {
		logger.Fatalf("failed to connect to database: %v", err)
	}
	defer func() {
		err := db.Close()
		if err != nil {
			logger.Errorf("failed to close database connection: %v", err)
		}
	}()

	// Миграции БД
	err = runMigrations(*dbURL, *migrationsPath, logger)
	if err != nil {
		logger.Fatalf("failed to run migrations: %v", err)
	}

	// Инициализация репозиториев
	userRepo := repository.NewUserRepository(db)
	sessionRepo := repository.NewSessionRepository(db)

	// Настройки для AuthService
	cfg := &app.Config{
		TokenSecret:     *tokenSecret,
		TokenExpiration: *tokenExpiration,
	}

	// Инициализация AuthService
	authService := app.NewAuthService(userRepo, sessionRepo, cfg, logger)

	// Запуск gRPC-сервера
	go func() {
		lis, err := net.Listen("tcp", *port)
		if err != nil {
			logger.Fatalf("failed to listen: %v", err)
		}
		s := grpc.NewServer()
		pb.RegisterAuthServiceServer(s, authService)
		reflection.Register(s)
		logger.Infof("gRPC service listening on %s", *port)
		if err := s.Serve(lis); err != nil {
			logger.Fatalf("failed to serve: %v", err)
		}
	}()

	// Инициализация Gin
	r := gin.Default()

	// Настройка CORS
	// В auth-servise/main.go и forum-servise/main.go
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Регистрация маршрутов
	r.POST("/register", func(c *gin.Context) { handleRegister(c, authService) })
	r.POST("/login", func(c *gin.Context) { handleLogin(c, authService) })

	// Запуск HTTP-сервера
	logger.Infof("HTTP service listening on %s", *httpPort)
	if err := r.Run(*httpPort); err != nil {
		logger.Fatalf("failed to run server: %v", err)
	}
}

func handleRegister(c *gin.Context, authService *app.AuthService) {
	var req pb.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := authService.Register(context.Background(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user_id": resp.UserId})
}

func handleLogin(c *gin.Context, authService *app.AuthService) {
	var req pb.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := authService.Login(context.Background(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": resp.Token})
}

func runMigrations(dbURL, migrationsPath string, logger *logger.Logger) error {
	m, err := migrate.New(
		"file://"+migrationsPath,
		dbURL,
	)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}
	defer func() {
		err, sourceErr := m.Close()
		if err != nil {
			logger.Errorf("failed to close migrate instance: %v", err)
		}
		if sourceErr != nil {
			logger.Errorf("failed to close migrate source: %v", sourceErr)
		}
	}()

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %w", err)
	}
	logger.Info("Migrations completed successfully")
	return nil
}
