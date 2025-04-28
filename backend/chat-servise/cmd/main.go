package main

import (
	"flag"
	"net/http"
	"time"

	"backend.com/forum/chat-servise/internal/controller"
	"backend.com/forum/chat-servise/internal/repository"
	"backend.com/forum/chat-servise/internal/usecase"
	"backend.com/forum/chat-servise/pkg/logger"
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

	// Инициализация контроллера
	wsController := controller.NewWebSocketController(chatUseCase, appLogger)

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
