// internal/usecase/comment_usecase.go
package usecase

import (
	"context"

	"backend.com/forum/forum-servise/internal/entity"
	"backend.com/forum/forum-servise/internal/repository"
	pb "backend.com/forum/proto"
)

// internal/usecase/comment_usecase.go
type CommentUseCase struct {
	commentRepo repository.CommentRepository
	postRepo    repository.PostRepository
	AuthClient  pb.AuthServiceClient // Делаем поле публичным
}

func NewCommentUseCase(
	commentRepo repository.CommentRepository,
	postRepo repository.PostRepository,
	authClient pb.AuthServiceClient,
) *CommentUseCase {
	return &CommentUseCase{
		commentRepo: commentRepo,
		postRepo:    postRepo,
		AuthClient:  authClient, // Исправляем здесь
	}
}
func (uc *CommentUseCase) CreateComment(ctx context.Context, comment *entity.Comment) error {
	_, err := uc.postRepo.GetPostByID(ctx, comment.PostID)
	if err != nil {
		return err
	}
	return uc.commentRepo.CreateComment(ctx, comment)
}

func (uc *CommentUseCase) GetCommentsByPostID(ctx context.Context, postID int64) ([]entity.Comment, error) {
	comments, err := uc.commentRepo.GetCommentsByPostID(ctx, postID)
	if err != nil {
		return nil, err
	}

	// Получаем имена авторов
	for i := range comments {
		userResponse, err := uc.AuthClient.GetUser(ctx, &pb.GetUserRequest{
			Id: comments[i].AuthorID,
		})

		if err == nil && userResponse != nil && userResponse.User != nil {
			comments[i].AuthorName = userResponse.User.Username
		} else {
			comments[i].AuthorName = "Unknown"
		}
	}

	return comments, nil
}

func (uc *CommentUseCase) DeleteComment(ctx context.Context, id int64) error {
	return uc.commentRepo.DeleteComment(ctx, id)
}
