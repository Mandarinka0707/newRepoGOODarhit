package repository

import (
	"context"
	"database/sql"
	"errors"

	domain "github.com/Mandarinka0707/newRepoGOODarhit/internal/entity"
	"github.com/jmoiron/sqlx"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *domain.User) (int64, error)
	GetUserByUsername(ctx context.Context, username string) (*domain.User, error)
	GetUserByID(ctx context.Context, id int64) (*domain.User, error) // Добавьте этот метод
}

type userRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(ctx context.Context, user *domain.User) (int64, error) {
	query := `INSERT INTO users (username, password, role, created_at) VALUES ($1, $2, $3, $4) RETURNING id`
	var id int64
	err := r.db.QueryRowContext(ctx, query, user.Username, user.Password, user.Role, user.CreatedAt).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

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
