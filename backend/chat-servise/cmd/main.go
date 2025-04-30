package main

import (
	"flag"
	"net/http"
	"time"

	"github.com/Mandarinka0707/newRepoGOODarhit/chat-servise/internal/controller"
	"github.com/Mandarinka0707/newRepoGOODarhit/chat-servise/internal/repository"
	"github.com/Mandarinka0707/newRepoGOODarhit/chat-servise/internal/usecase"
	"github.com/Mandarinka0707/newRepoGOODarhit/chat-servise/pkg/auth"
	"github.com/Mandarinka0707/newRepoGOODarhit/chat-servise/pkg/logger"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	port     = flag.String("port", ":8082", "HTTP server port")
	dbURL    = flag.String("db-url", "postgres://user:password@localhost:5432/database?sslmode=disable", "Database connection URL")
	logLevel = flag.String("log-level", "info", "Logging level")
)

func main() {
	flag.Parse()

	// Инициализация логгера
	appLogger := logger.NewWithLevel(*logLevel)

	// Подключение к базе данных через sqlx
	db, err := sqlx.Connect("postgres", *dbURL)
	if err != nil {
		appLogger.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Инициализация GORM
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db.DB,
	}), &gorm.Config{})
	if err != nil {
		appLogger.Fatalf("Failed to initialize GORM: %v", err)
	}

	// Инициализация репозитория
	messageRepo := repository.NewMessageRepository(gormDB)

	// Инициализация usecase
	chatUseCase := usecase.NewChatUsecase(messageRepo)
	authClient, err := auth.NewClient("localhost:50051") // Укажите правильный адрес auth-service
	if err != nil {
		appLogger.Fatalf("Failed to create auth client: %v", err)
	}

	// Исправляем вызов NewWebSocketController
	wsController := controller.NewWebSocketController(
		authClient,
		chatUseCase,
		appLogger,
	)

	// Настройка CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           int(12 * time.Hour),
	})

	// Настройка маршрутизатора
	router := mux.NewRouter()
	router.HandleFunc("/ws", wsController.HandleWebSocket)
	router.HandleFunc("/messages", wsController.GetMessages).Methods("GET")

	// Запуск сервера
	appLogger.Infof("Starting chat service on port %s", *port)
	if err := http.ListenAndServe(*port, c.Handler(router)); err != nil {
		appLogger.Fatalf("Failed to start server: %v", err)
	}
}
