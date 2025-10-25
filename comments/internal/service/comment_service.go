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
	UpdateContentWithCheck(ctx context.Context, commentID int64, newComment *domain.Comment, predicate func(oldComment *domain.Comment) error) (err error)
	GetByID(ctx context.Context, commentID int64) (*domain.Comment, error)
	GetListBySKU(ctx context.Context, sku int64, lastCreatedAt time.Time, lastUserID int64, limit int32) ([]*domain.Comment, error)
	GetListByUser(ctx context.Context, userID int64) ([]*domain.Comment, error)
}

// CommentService реализует бизнес-логику управления комментариями на товары.
type CommentService struct {
	commentRepository commentRepository
	editTimeout       time.Duration
	limitRowsBySku    int32
}

// NewCommentService создает новый сервис CommentService.
func NewCommentService(commentRepository commentRepository, editTimeout time.Duration, limitRowsBySku int32) *CommentService {
	return &CommentService{
		commentRepository: commentRepository,
		editTimeout:       editTimeout,
		limitRowsBySku:    limitRowsBySku,
	}
}

// Add добавляет новый комментарий.
func (c *CommentService) Add(ctx context.Context, comment *domain.Comment) (int64, error) {
	comment.CreatedAt = time.Now()
	commentID, err := c.commentRepository.Insert(ctx, comment)
	if err != nil {
		return 0, fmt.Errorf("commentRepository.Insert: %w", err)
	}

	return commentID, nil
}

// GetInfoByID возвращает информацию о комментарии по его идентификатору.
func (c *CommentService) GetInfoByID(ctx context.Context, commentID int64) (*domain.Comment, error) {
	comment, err := c.commentRepository.GetByID(ctx, commentID)
	if err != nil {
		return nil, fmt.Errorf("commentRepository.GetByID: %w", err)
	}

	return comment, nil
}

// Edit редактирует комментарий, если это разрешено.
func (c *CommentService) Edit(ctx context.Context, commentID int64, newComment *domain.Comment) error {
	err := c.commentRepository.UpdateContentWithCheck(ctx, commentID, newComment, func(oldComment *domain.Comment) error {
		if oldComment.UserID != newComment.UserID {
			return domain.ErrEditNotMyComment
		}

		if time.Since(oldComment.CreatedAt) >= c.editTimeout {
			return domain.ErrEditTimeoutExceed
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("commentRepository.UpdateContentWithCheck: %w", err)
	}

	return nil
}

// GetListBySKU возвращает список комментариев по SKU товара.
func (c *CommentService) GetListBySKU(ctx context.Context, sku int64, lastCreatedAt time.Time, lastUserID int64) ([]*domain.Comment, error) {
	comments, err := c.commentRepository.GetListBySKU(ctx, sku, lastCreatedAt, lastUserID, c.limitRowsBySku)
	if err != nil {
		return nil, fmt.Errorf("commentRepository.GetListBySKU: %w", err)
	}

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

// GetListByUser возвращает список комментариев пользователя.
func (c *CommentService) GetListByUser(ctx context.Context, userID int64) ([]*domain.Comment, error) {
	comments, err := c.commentRepository.GetListByUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("commentRepository.GetListByUser: %w", err)
	}

	sortComments(comments)

	return comments, nil
}
