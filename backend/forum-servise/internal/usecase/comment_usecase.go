package usecase

import (
	"context"
	"errors"

	pb "backend.com/forum/proto"
	"github.com/Mandarinka0707/newRepoGOODarhit/forum-servise/internal/entity"
	"github.com/Mandarinka0707/newRepoGOODarhit/forum-servise/internal/repository"
)

type CommentUseCase struct {
	CommentRepo repository.CommentRepository
	postRepo    repository.PostRepository
	AuthClient  pb.AuthServiceClient
}

func NewCommentUseCase(
	commentRepo repository.CommentRepository,
	postRepo repository.PostRepository,
	authClient pb.AuthServiceClient,
) *CommentUseCase {
	return &CommentUseCase{
		CommentRepo: commentRepo,
		postRepo:    postRepo,
		AuthClient:  authClient,
	}
}

func (uc *CommentUseCase) CreateComment(ctx context.Context, comment *entity.Comment) error {

	_, err := uc.postRepo.GetPostByID(ctx, comment.PostID)
	if err != nil {
		return err
	}

	userResp, err := uc.AuthClient.GetUser(ctx, &pb.GetUserRequest{Id: comment.AuthorID})
	if err != nil || userResp == nil || userResp.User == nil {
		return errors.New("failed to get user info")
	}

	comment.AuthorName = userResp.User.Username
	return uc.CommentRepo.CreateComment(ctx, comment)
}

func (uc *CommentUseCase) GetCommentsByPostID(ctx context.Context, postID int64) ([]entity.Comment, error) {

	_, err := uc.postRepo.GetPostByID(ctx, postID)
	if err != nil {
		return nil, err
	}

	return uc.CommentRepo.GetCommentsByPostID(ctx, postID)
}

// func (uc *CommentUseCase) DeleteComment(ctx context.Context, id int64) error {
// 	// Получаем комментарий
// 	comments, err := uc.CommentRepo.GetCommentsByPostID(ctx, id)
// 	if err != nil {
// 		return err
// 	}
// 	if len(comments) == 0 {
// 		return sql.ErrNoRows
// 	}
// 	comment := comments[0]

// 	// Проверяем права доступа
// 	userResp, err := uc.AuthClient.GetUser(ctx, &pb.GetUserRequest{Id: comment.AuthorID})
// 	if err != nil || userResp == nil || userResp.User == nil {
// 		return errors.New("failed to verify user")
// 	}

// 	// Проверяем, что пользователь является автором
// 	if userResp.User.Id != comment.AuthorID {
// 		return errors.New("permission denied")
// 	}

// 	return uc.CommentRepo.DeleteComment(ctx, id)
// }
