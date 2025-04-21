package repository

import (
	"context"
	"database/sql"
	"errors"

	"backend.com/forum/forum-servise/internal/domain"
	"github.com/jmoiron/sqlx"
)

type CategoryRepository interface {
	CreateCategory(ctx context.Context, category *domain.Category) (int64, error)
	GetCategory(ctx context.Context, id int64) (*domain.Category, error)
	// ... другие методы
}

type categoryRepository struct {
	db *sqlx.DB
}

func NewCategoryRepository(db *sqlx.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) CreateCategory(ctx context.Context, category *domain.Category) (int64, error) {
	query := `INSERT INTO categories (name, description, created_at) VALUES ($1, $2, $3) RETURNING id`
	var id int64
	err := r.db.QueryRowContext(ctx, query, category.Name, category.Description, category.CreatedAt).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *categoryRepository) GetCategory(ctx context.Context, id int64) (*domain.Category, error) {
	query := `SELECT id, name, description, created_at FROM categories WHERE id = $1`
	category := &domain.Category{}
	err := r.db.GetContext(ctx, category, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows // Важно возвращать sql.ErrNoRows
		}
		return nil, err
	}
	return category, nil
}
