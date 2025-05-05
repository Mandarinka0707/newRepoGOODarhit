package repository

import (
	"context"
	"database/sql"
	"errors"

	"backend.com/forum/forum-servise/internal/entity"
	"github.com/jmoiron/sqlx"
)

type PostRepository interface {
	CreatePost(ctx context.Context, post *entity.Post) (int64, error)
	GetPosts(ctx context.Context) ([]*entity.Post, error)
	DeletePost(ctx context.Context, postID, authorID int64, role string) error
	GetPostByID(ctx context.Context, id int64) (*entity.Post, error)
	UpdatePost(
		ctx context.Context,
		postID int64,
		authorID int64,
		role string,
		title, content string,
	) (*entity.Post, error)
}

var (
	ErrPermissionDenied = errors.New("permission denied")
	ErrPostNotFound     = errors.New("post not found")
)

type postRepository struct {
	db *sqlx.DB
}

func NewPostRepository(db *sqlx.DB) PostRepository {
	return &postRepository{db: db}
}

func (r *postRepository) CreatePost(ctx context.Context, post *entity.Post) (int64, error) {
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

func (r *postRepository) GetPosts(ctx context.Context) ([]*entity.Post, error) {
	query := `
        SELECT 
            p.id,
            p.title,
            p.content,
            p.author_id,  
            p.created_at
        FROM posts p
        ORDER BY p.created_at DESC
    `

	var posts []*entity.Post
	err := r.db.SelectContext(ctx, &posts, query)
	return posts, err
}
func (r *postRepository) DeletePost(
	ctx context.Context,
	postID,
	authorID int64,
	role string,
) error {
	// 1. Проверяем существование поста
	var exists bool
	err := r.db.GetContext(
		ctx,
		&exists,
		"SELECT EXISTS(SELECT 1 FROM posts WHERE id = $1)",
		postID,
	)
	if err != nil {
		return err
	}
	if !exists {
		return ErrPostNotFound
	}

	// 2. Выполняем удаление с проверкой прав
	query := `
        DELETE FROM posts 
        WHERE id = $1 
        AND (author_id = $2 OR $3 = 'admin')
    `
	result, err := r.db.ExecContext(
		ctx,
		query,
		postID,   // $1
		authorID, // $2
		role,     // $3 <-- Важно! Проверьте порядок параметров
	)
	if err != nil {
		return err
	}

	// 3. Проверяем результат
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return ErrPermissionDenied
	}

	return nil
}
func (r *postRepository) GetPostByID(ctx context.Context, id int64) (*entity.Post, error) {
	query := `
        SELECT 
            p.id,
            p.title,
            p.content,
            p.author_id,
            p.created_at
        FROM posts p
        WHERE p.id = $1
    `

	var post entity.Post
	err := r.db.GetContext(ctx, &post, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrPostNotFound
		}
		return nil, err
	}

	return &post, nil
}

func (r *postRepository) UpdatePost(
	ctx context.Context,
	postID int64,
	authorID int64,
	role string,
	title,
	content string,
) (*entity.Post, error) {
	// 1. Проверяем существование поста
	var exists bool
	err := r.db.GetContext(
		ctx,
		&exists,
		"SELECT EXISTS(SELECT 1 FROM posts WHERE id = $1)",
		postID,
	)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrPostNotFound
	}

	// 2. Выполняем обновление
	query := `
        UPDATE posts 
        SET title = $1, content = $2 
        WHERE id = $3 
        AND (author_id = $4 OR $5 = 'admin')
        RETURNING id, title, content, author_id, created_at
    `

	var post entity.Post
	err = r.db.QueryRowContext(ctx, query,
		title,
		content,
		postID,
		authorID,
		role,
	).Scan(
		&post.ID,
		&post.Title,
		&post.Content,
		&post.AuthorID,
		&post.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrPermissionDenied
		}
		return nil, err
	}

	return &post, nil
}
