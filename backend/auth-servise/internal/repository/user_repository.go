package repository

import (
	"context"
	"database/sql"
	"errors"

	domain "backend.com/forum/auth-servise/internal/entity"
)

// DBTX определяет интерфейс, который реализуют как *sqlx.DB, так и моки
type DBTX interface {
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
}

type UserRepository interface {
	CreateUser(ctx context.Context, user *domain.User) (int64, error)
	GetUserByUsername(ctx context.Context, username string) (*domain.User, error)
	GetUserByID(ctx context.Context, id int64) (*domain.User, error)
}

type userRepository struct {
	db DBTX
}

// NewUserRepository создает новый экземпляр UserRepository
func NewUserRepository(db DBTX) UserRepository {
	return &userRepository{db: db}
}

// CreateUser сохраняет пользователя и возвращает его ID
func (r *userRepository) CreateUser(ctx context.Context, user *domain.User) (int64, error) {
	query := `INSERT INTO users (username, password, role, created_at) VALUES ($1, $2, $3, $4) RETURNING id`
	var id int64
	err := r.db.QueryRowContext(ctx, query, user.Username, user.Password, user.Role, user.CreatedAt).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// GetUserByUsername возвращает пользователя по имени
func (r *userRepository) GetUserByUsername(ctx context.Context, username string) (*domain.User, error) {
	query := `SELECT id, username, password, role, created_at FROM users WHERE username = $1`
	user := &domain.User{}
	err := r.db.GetContext(ctx, user, query, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}
	return user, nil
}

// GetUserByID возвращает пользователя по ID
func (r *userRepository) GetUserByID(ctx context.Context, id int64) (*domain.User, error) {
	query := `SELECT id, username, password, role, created_at FROM users WHERE id = $1`
	user := &domain.User{}
	err := r.db.GetContext(ctx, user, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}
