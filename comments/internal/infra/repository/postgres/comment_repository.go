package postgres

import (
	"context"
	"errors"
	"fmt"
	"route256/comments/internal/domain"
	sqlcrepos "route256/comments/internal/infra/repository/postgres/sqlc/generated"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

// NewCommentRepository создает новый OrderRepository.
func NewCommentRepository(pool sqlcrepos.DBTX) *CommentRepository {
	return &CommentRepository{
		sqlcrepos.New(pool),
	}
}

// CommentRepository предоставляет доступ к хранилищу комментариев из postgres.
type CommentRepository struct {
	querier sqlcrepos.Querier
}

// Insert вставляет новую строку с комментарием в postgres.
func (c *CommentRepository) Insert(ctx context.Context, comment *domain.Comment) (int64, error) {
	commentID, err := c.querier.AddComment(ctx, &sqlcrepos.AddCommentParams{
		UserID:    comment.UserID,
		Sku:       comment.Sku,
		Content:   comment.Content,
		CreatedAt: pgtype.Timestamp{Time: comment.CreatedAt, Valid: true},
	})
	if err != nil {
		return 0, fmt.Errorf("querier.AddComment: %w", err)
	}

	return commentID, err
}

// UpdateContent обновляет содержание комментария в postgres.
func (c *CommentRepository) UpdateContent(ctx context.Context, commentID int64, newComment *domain.Comment) error {
	err := c.querier.UpdateContent(ctx, &sqlcrepos.UpdateContentParams{
		ID:      commentID,
		Content: newComment.Content,
	})
	if err != nil {
		return fmt.Errorf("querier.UpdateContent: %w", err)
	}

	return nil
}

// GetByIDForUpdate возвращает комментарий для обновления в postgres.
func (c *CommentRepository) GetByIDForUpdate(ctx context.Context, commentID int64) (*domain.Comment, error) {
	commentDB, err := c.querier.GetCommentByIDForUpdate(ctx, commentID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrCommentNoExist
		}

		return nil, fmt.Errorf("querier.GetCommentByIDForUpdate: %w", err)
	}

	return &domain.Comment{
		Id:        commentID,
		UserID:    commentDB.UserID,
		Sku:       commentDB.Sku,
		Content:   commentDB.Content,
		CreatedAt: commentDB.CreatedAt.Time,
	}, nil
}

// GetByID возвращает комментарий в postgres.
func (c *CommentRepository) GetByID(ctx context.Context, commentID int64) (*domain.Comment, error) {
	commentDB, err := c.querier.GetCommentByID(ctx, commentID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrCommentNoExist
		}

		return nil, fmt.Errorf("querier.GetCommentByID: %w", err)
	}

	return toComment(commentDB), nil
}

func toComment(commentDB *sqlcrepos.Comment) *domain.Comment {
	return &domain.Comment{
		Id:        commentDB.ID,
		UserID:    commentDB.UserID,
		Sku:       commentDB.Sku,
		Content:   commentDB.Content,
		CreatedAt: commentDB.CreatedAt.Time,
	}
}

// GetListBySKU возвращает список комментариев о товара в postgres.
func (c *CommentRepository) GetListBySKU(ctx context.Context, sku int64) ([]*domain.Comment, error) {
	commentsDB, err := c.querier.GetCommentsBySKU(ctx, sku)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []*domain.Comment{}, nil
		}

		return nil, fmt.Errorf("querier.GetCommentsBySKU: %w", err)
	}

	res := make([]*domain.Comment, 0, len(commentsDB))
	for _, commentDB := range commentsDB {
		res = append(res, toComment(commentDB))
	}

	return res, nil
}

// GetListByUser возвращает список комментариев от пользователя в postgres.
func (c *CommentRepository) GetListByUser(ctx context.Context, userID int64) ([]*domain.Comment, error) {
	commentsDB, err := c.querier.GetCommentsByUser(ctx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []*domain.Comment{}, nil
		}

		return nil, fmt.Errorf("querier.GetCommentsByUser: %w", err)
	}

	res := make([]*domain.Comment, 0, len(commentsDB))
	for _, commentDB := range commentsDB {
		res = append(res, toComment(commentDB))
	}

	return res, nil
}
