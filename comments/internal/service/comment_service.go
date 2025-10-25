package service

import (
	"cmp"
	"context"
	"fmt"
	"route256/comments/internal/domain"
	"slices"
	"time"
)

type commentRepository interface {
	Insert(ctx context.Context, comment *domain.Comment) (int64, error)
	UpdateContent(ctx context.Context, commentID int64, newComment *domain.Comment) error
	GetByIDForUpdate(ctx context.Context, commentID int64) (*domain.Comment, error)
	GetByID(ctx context.Context, commentID int64) (*domain.Comment, error)
	GetListBySKU(ctx context.Context, sku int64) ([]*domain.Comment, error)
	GetListByUser(ctx context.Context, userID int64) ([]*domain.Comment, error)
}

// CommentService реализует бизнес-логику управления комментариями на товары.
type CommentService struct {
	commentRepository commentRepository
	editTimeout       time.Duration
}

// NewCommentService создает новый сервис CommentService.
func NewCommentService(commentRepository commentRepository, editTimeout time.Duration) *CommentService {
	return &CommentService{
		commentRepository: commentRepository,
		editTimeout:       editTimeout,
	}
}

func (c *CommentService) Add(ctx context.Context, comment *domain.Comment) (int64, error) {
	comment.CreatedAt = time.Now()
	commentID, err := c.commentRepository.Insert(ctx, comment)
	if err != nil {
		return 0, fmt.Errorf("commentRepository.Insert: %w", err)
	}

	return commentID, nil
}

func (c *CommentService) GetInfoByID(ctx context.Context, commentID int64) (*domain.Comment, error) {
	comment, err := c.commentRepository.GetByID(ctx, commentID)
	if err != nil {
		return nil, fmt.Errorf("commentRepository.GetByID: %w", err)
	}

	return comment, nil
}

func (c *CommentService) Edit(ctx context.Context, commentID int64, newComment *domain.Comment) error {
	comment, err := c.commentRepository.GetByIDForUpdate(ctx, commentID)
	if err != nil {
		return fmt.Errorf("commentRepository.GetByIDForUpdate: %w", err)
	}

	if comment.UserID != newComment.UserID {
		return domain.ErrEditNotMyComment
	}

	if time.Since(comment.CreatedAt) >= c.editTimeout {
		return domain.ErrEditTimeoutExceed
	}

	err = c.commentRepository.UpdateContent(ctx, commentID, newComment)
	if err != nil {
		return fmt.Errorf("commentRepository.UpdateContent: %w", err)
	}

	return nil
}

func (c *CommentService) GetListBySKU(ctx context.Context, sku int64) ([]*domain.Comment, error) {
	comments, err := c.commentRepository.GetListBySKU(ctx, sku)
	if err != nil {
		return nil, fmt.Errorf("commentRepository.GetListBySKU: %w", err)
	}

	sortComments(comments)

	return comments, nil
}

func sortComments(comments []*domain.Comment) {
	slices.SortFunc(comments, func(a, b *domain.Comment) int {
		timeCmp := b.CreatedAt.Compare(a.CreatedAt)
		if timeCmp != 0 {
			return timeCmp
		}

		return cmp.Compare(a.UserID, b.UserID)
	})
}

func (c *CommentService) GetListByUser(ctx context.Context, userID int64) ([]*domain.Comment, error) {
	comments, err := c.commentRepository.GetListByUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("commentRepository.GetListByUser: %w", err)
	}

	sortComments(comments)

	return comments, nil
}
