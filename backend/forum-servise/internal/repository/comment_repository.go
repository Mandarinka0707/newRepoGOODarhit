package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Mandarinka0707/newRepoGOODarhit/forum-servise/internal/entity"
	"github.com/jmoiron/sqlx"
)

type CommentRepository interface {
	CreateComment(ctx context.Context, comment *entity.Comment) error
	GetCommentsByPostID(ctx context.Context, postID int64) ([]entity.Comment, error)
}

type CommentRepo struct {
	db *sqlx.DB
}

func NewCommentRepository(db *sqlx.DB) CommentRepository {
	return &CommentRepo{db: db}
}

func (r *CommentRepo) CreateComment(ctx context.Context, comment *entity.Comment) error {
	query := `INSERT INTO comments (content, author_id, post_id, author_name) 
        VALUES ($1, $2, $3, $4) RETURNING id`
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
        ORDER BY id DESC`

	var comments []entity.Comment
	err := r.db.SelectContext(ctx, &comments, query, postID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []entity.Comment{}, nil
		}
		return nil, err
	}
	return comments, nil
}
