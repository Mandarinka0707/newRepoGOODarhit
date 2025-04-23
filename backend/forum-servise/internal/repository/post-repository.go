package repository

import (
	"context"

	"backend.com/forum/forum-servise/internal/domain"
	"github.com/jmoiron/sqlx"
)

type PostRepository interface {
	CreatePost(ctx context.Context, post *domain.Post) (int64, error)
	GetPosts(ctx context.Context) ([]*domain.Post, error)
}

type postRepository struct {
	db *sqlx.DB
}

func NewPostRepository(db *sqlx.DB) PostRepository {
	return &postRepository{db: db}
}

func (r *postRepository) CreatePost(ctx context.Context, post *domain.Post) (int64, error) {
	query := `
		INSERT INTO posts (title, content, author_id, created_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`
	var id int64
	err := r.db.QueryRowContext(ctx, query,
		post.Title,
		post.Content,
		post.AuthorID,
		post.CreatedAt,
	).Scan(&id)
	return id, err
}

func (r *postRepository) GetPosts(ctx context.Context) ([]*domain.Post, error) {
	query := `
		SELECT id, title, content, author_id, created_at
		FROM posts
		ORDER BY created_at DESC
	`
	var posts []*domain.Post
	err := r.db.SelectContext(ctx, &posts, query)
	return posts, err
}
