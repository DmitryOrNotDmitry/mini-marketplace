package postgres

import (
	"context"
	"errors"
	"fmt"
	"route256/cart/pkg/myerrgroup"
	"route256/comments/internal/domain"
	sqlcrepos "route256/comments/internal/infra/repository/postgres/sqlc/generated"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

// NewCommentRepository создает новый OrderRepository.
func NewCommentRepository(shardManager *ShardManager) *CommentRepository {
	return &CommentRepository{
		shardManager: shardManager,
	}
}

// CommentRepository предоставляет доступ к хранилищу комментариев из postgres.
type CommentRepository struct {
	shardManager *ShardManager
}

func getQuerier(pool sqlcrepos.DBTX) sqlcrepos.Querier {
	return sqlcrepos.New(pool)
}

// Insert вставляет новую строку с комментарием в postgres.
func (c *CommentRepository) Insert(ctx context.Context, comment *domain.Comment) (int64, error) {
	pool, bucketIdx := c.shardManager.GetShardPool(comment.Sku)
	querier := getQuerier(pool)
	commentID, err := querier.AddComment(ctx, &sqlcrepos.AddCommentParams{
		Column1:   bucketIdx,
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
	querier := getQuerier(c.shardManager.GetShardPoolByID(commentID))
	err := querier.UpdateContent(ctx, &sqlcrepos.UpdateContentParams{
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
	querier := getQuerier(c.shardManager.GetShardPoolByID(commentID))
	commentDB, err := querier.GetCommentByIDForUpdate(ctx, commentID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrCommentNotExist
		}

		return nil, fmt.Errorf("querier.GetCommentByIDForUpdate: %w", err)
	}

	return &domain.Comment{
		ID:        commentID,
		UserID:    commentDB.UserID,
		Sku:       commentDB.Sku,
		Content:   commentDB.Content,
		CreatedAt: commentDB.CreatedAt.Time,
	}, nil
}

// GetByID возвращает комментарий в postgres.
func (c *CommentRepository) GetByID(ctx context.Context, commentID int64) (*domain.Comment, error) {
	querier := getQuerier(c.shardManager.GetShardPoolByID(commentID))
	commentDB, err := querier.GetCommentByID(ctx, commentID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrCommentNotExist
		}

		return nil, fmt.Errorf("querier.GetCommentByID: %w", err)
	}

	return toComment(commentDB), nil
}

func toComment(commentDB *sqlcrepos.Comment) *domain.Comment {
	return &domain.Comment{
		ID:        commentDB.ID,
		UserID:    commentDB.UserID,
		Sku:       commentDB.Sku,
		Content:   commentDB.Content,
		CreatedAt: commentDB.CreatedAt.Time,
	}
}

// GetListBySKU возвращает список комментариев о товара в postgres.
func (c *CommentRepository) GetListBySKU(ctx context.Context, sku int64, lastCreatedAt time.Time, lastUserID int64, limit int32) ([]*domain.Comment, error) {
	pool, _ := c.shardManager.GetShardPool(sku)
	querier := getQuerier(pool)
	commentsDB, err := querier.GetCommentsBySKU(ctx, &sqlcrepos.GetCommentsBySKUParams{
		Sku:       sku,
		CreatedAt: pgtype.Timestamp{Time: lastCreatedAt, Valid: true},
		UserID:    lastUserID,
		Limit:     limit,
	})
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
	pools := c.shardManager.GetAllPools()

	commentsDBList := make([][]*sqlcrepos.Comment, len(pools))

	errgroup := myerrgroup.New()
	for i, pool := range pools {
		errgroup.Go(func() error {
			querier := getQuerier(pool)

			commentsDB, err := querier.GetCommentsByUser(ctx, userID)
			if err != nil {
				if errors.Is(err, pgx.ErrNoRows) {
					commentsDBList[i] = []*sqlcrepos.Comment{}
					return nil
				}

				return fmt.Errorf("querier.GetCommentsByUser: %w", err)
			}

			commentsDBList[i] = commentsDB
			return nil
		})
	}

	err := errgroup.Wait()
	if err != nil {
		return nil, fmt.Errorf("errgroup.Wait: %w", err)
	}

	return groupToCommentList(commentsDBList), nil
}

func groupToCommentList(commentsDBList [][]*sqlcrepos.Comment) []*domain.Comment {
	sumLen := 0
	for _, commsDB := range commentsDBList {
		sumLen += len(commsDB)
	}

	res := make([]*domain.Comment, 0, sumLen)
	for _, commsDB := range commentsDBList {
		for _, commDB := range commsDB {
			res = append(res, toComment(commDB))
		}
	}

	return res
}
