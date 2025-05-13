package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"backend.com/forum/auth-servise/internal/controller"
	"backend.com/forum/auth-servise/internal/repository"
	"backend.com/forum/auth-servise/internal/usecase"
	"backend.com/forum/auth-servise/pkg/auth"
	"backend.com/forum/auth-servise/pkg/logger"
	"golang.org/x/crypto/bcrypt"

	_ "backend.com/forum/auth-servise/docs"
	pb "backend.com/forum/proto"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// @title           Auth Service API
// @version         1.0
// @description     Authentication and authorization service for the forum
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080

// @securityDefinitions.basic  BasicAuth
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
var (
	grpcPort        = flag.String("grpc-port", ":50051", "gRPC server port")
	httpPort        = flag.String("http-port", ":8080", "HTTP server port")
	dbURL           = flag.String("db-url", "postgres://user:password@localhost:5432/database?sslmode=disable", "Database connection URL")
	migrationsPath  = flag.String("migrations_path", "/Users/darinautalieva/Desktop/GOProject/backend/auth-servise/migrations", "path to migrations files")
	tokenSecret     = flag.String("token-secret", "secret", "JWT token secret")
	tokenExpiration = flag.Duration("token-expiration", 24*time.Hour, "JWT token expiration")
	logLevel        = flag.String("log-level", "info", "Logging level")
)

func main() {
	flag.Parse()
	password := "admin"
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	fmt.Println("PASSWORD     ")
	fmt.Println(string(hash))

	logger, err := logger.NewLogger(*logLevel)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	db, err := sqlx.Connect("postgres", *dbURL)
	if err != nil {
		logger.Fatal("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := runMigrations(*dbURL, *migrationsPath, logger); err != nil {
		logger.Fatal("Migrations failed: %v", err)
	}

	userRepo := repository.NewUserRepository(db)
	sessionRepo := repository.NewSessionRepository(db)

	authConfig := &auth.Config{
		TokenSecret:     *tokenSecret,
		TokenExpiration: *tokenExpiration,
	}

	authUseCase := usecase.NewAuthUsecase(
		userRepo,
		sessionRepo,
		authConfig,
		logger.ZapLogger(),
	)

	grpcController := controller.NewAuthController(authUseCase)
	httpController := controller.NewHTTPAuthController(authUseCase)

	go startGRPCServer(*grpcPort, grpcController, logger)
	startHTTPServer(*httpPort, httpController, logger)
}

func startGRPCServer(port string, controller *controller.AuthController, logger *logger.Logger) {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		logger.Fatal("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterAuthServiceServer(s, controller)
	reflection.Register(s)

	logger.Info("Starting gRPC server on %s", port)
	if err := s.Serve(lis); err != nil {
		logger.Fatal("Failed to serve gRPC: %v", err)
	}
}

// main.go (исправленная часть)
func startHTTPServer(port string, controller *controller.HTTPAuthController, logger *logger.Logger) {
	router := gin.Default()

	// Настройка CORS и Swagger
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Группировка роутов с префиксом /api/v1
	api := router.Group("/api/v1")
	{
		authGroup := api.Group("/auth")
		{
			authGroup.POST("/register", controller.Register)
			authGroup.POST("/login", controller.Login)
			authGroup.GET("/user/:id", controller.GetUser)
		}
	}

	logger.Info("Starting HTTP server on %s", port)
	if err := http.ListenAndServe(port, router); err != nil {
		logger.Fatal("Failed to start HTTP server: %v", err)
	}
}
func runMigrations(dbURL, migrationsPath string, logger *logger.Logger) error {
	m, err := migrate.New(
		"file://"+migrationsPath,
		dbURL,
	)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}
	defer m.Close()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	logger.Info("Database migrations applied successfully")
	return nil
}
