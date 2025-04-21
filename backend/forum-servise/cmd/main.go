package main

import (
	"flag"
	"log"
	"time"

	"backend.com/forum/forum-servise/internal/app"
	"backend.com/forum/forum-servise/internal/repository"
	"backend.com/forum/forum-servise/pkg/logger"
	pb "backend.com/forum/proto"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
)

var (
	port            = flag.String("port", ":3000", "the HTTP port to serve on")
	dbURL           = flag.String("db_url", "postgres://user:password@localhost:5432/database?sslmode=disable", "PostgreSQL connection string")
	authServiceAddr = flag.String("auth_service_addr", "localhost:50051", "Auth Service gRPC address")
	logLevel        = flag.String("log_level", "info", "Log level (debug, info, warn, error)")
	chatMessageTTL  = flag.Duration("chat_message_ttl", time.Hour, "Time-to-live for chat messages")
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

	// Инициализация репозиториев
	categoryRepo := repository.NewCategoryRepository(db)
	topicRepo := repository.NewTopicRepository(db)
	messageRepo := repository.NewMessageRepository(db)
	chatMessageRepo := repository.NewChatMessageRepository(db)

	// Инициализация gRPC клиента для Auth Service
	conn, err := grpc.Dial(*authServiceAddr, grpc.WithInsecure())
	if err != nil {
		logger.Fatalf("failed to connect to auth service: %v", err)
	}
	defer conn.Close()
	authClient := pb.NewAuthServiceClient(conn)

	// Настройки для ForumService
	cfg := &app.Config{
		AuthServiceAddress: *authServiceAddr,
		ChatMessageTTL:     *chatMessageTTL,
	}

	// Инициализация ForumService
	forumService := app.NewForumService(categoryRepo, topicRepo, messageRepo, chatMessageRepo, authClient, cfg, logger)

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

	r.Static("/", "./static")

	// Регистрация маршрутов
	r.POST("/websocket", func(c *gin.Context) { forumService.WebsocketHandler(c.Writer, c.Request) })
	r.GET("/manifest.json", func(c *gin.Context) {
		c.File("./static/manifest.json") // Путь к вашему файлу
	})

	r.NoRoute(func(c *gin.Context) {
		c.File("./static/index.html") // Для SPA (чтобы React Router работал)
	})
	// Запуск HTTP-сервера
	logger.Infof("HTTP service listening on %s", *port)
	if err := r.Run(*port); err != nil {
		logger.Fatalf("failed to run server: %v", err)
	}
}
