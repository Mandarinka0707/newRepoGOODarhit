// internal/repository/comment_repository.go
package repository

import (
	"context"

	"backend.com/forum/forum-servise/internal/entity"
	"github.com/jmoiron/sqlx"
)

type CommentRepository interface {
	CreateComment(ctx context.Context, comment *entity.Comment) error
	GetCommentsByPostID(ctx context.Context, postID int64) ([]entity.Comment, error)
	DeleteComment(ctx context.Context, id int64) error
}

// Пример реализации для PostgreSQL
type CommentRepo struct {
	db *sqlx.DB
}

func NewCommentRepository(db *sqlx.DB) CommentRepository {
	return &CommentRepo{db: db}
}

func (r *CommentRepo) CreateComment(ctx context.Context, comment *entity.Comment) error {
	query := `INSERT INTO comments (content, author_id, post_id, author_name) 
        VALUES ($1, $2, $3, $4) RETURNING id` // Исправлено aurhor_name -> author_name
	return r.db.QueryRowContext(ctx, query,
		comment.Content,
		comment.AuthorID,
		comment.PostID,
		comment.AuthorName,
	).Scan(&comment.ID)
}

func (r *CommentRepo) GetCommentsByPostID(ctx context.Context, postID int64) ([]entity.Comment, error) {
	query := `
        SELECT 
            id,
            content,
            author_id,
            post_id,
            author_name 
        FROM comments 
        WHERE post_id = $1
        ORDER BY id DESC` // Убрана лишняя запятая и JOIN
	var comments []entity.Comment
	err := r.db.SelectContext(ctx, &comments, query, postID)
	return comments, err
}

func (r *CommentRepo) DeleteComment(ctx context.Context, id int64) error {
	query := `DELETE FROM comments WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
