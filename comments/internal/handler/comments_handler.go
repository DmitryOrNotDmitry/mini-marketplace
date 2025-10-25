package handler

import (
	"context"
	"errors"
	"fmt"
	"route256/comments/internal/domain"
	"route256/comments/pkg/api/comments/v1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type commentService interface {
	Add(ctx context.Context, comment *domain.Comment) (int64, error)
	GetInfoByID(ctx context.Context, commentID int64) (*domain.Comment, error)
	Edit(ctx context.Context, commentID int64, newComment *domain.Comment) error
	GetListBySKU(ctx context.Context, sku int64) ([]*domain.Comment, error)
	GetListByUser(ctx context.Context, userID int64) ([]*domain.Comment, error)
}

// CommentServerGRPC обрабатывает gRPC-запросы для операций с комментариями.
type CommentServerGRPC struct {
	comments.UnimplementedCommentsServiceV1Server
	commentService commentService
}

// NewCommentServerGRPC создает новый экземпляр CommentServerGRPC.
func NewCommentServerGRPC(commentService commentService) *CommentServerGRPC {
	return &CommentServerGRPC{
		commentService: commentService,
	}
}

// AddCommentV1 создает новый комментарий.
func (c *CommentServerGRPC) AddCommentV1(ctx context.Context, req *comments.AddCommentRequest) (*comments.AddCommentResponse, error) {
	comment := &domain.Comment{
		UserID:  req.UserId,
		Sku:     req.Sku,
		Content: req.Content,
	}

	commentID, err := c.commentService.Add(ctx, comment)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Errorf("internal server error: %w", err).Error())
	}

	res := &comments.AddCommentResponse{
		Id: commentID,
	}

	return res, nil
}

// CommentInfoV1 возвращает информацию по комментарию.
func (c *CommentServerGRPC) CommentInfoV1(ctx context.Context, req *comments.CommentInfoRequest) (*comments.CommentInfoResponse, error) {
	comment, err := c.commentService.GetInfoByID(ctx, req.Id)
	if err != nil {
		if errors.Is(err, domain.ErrCommentNotExist) {
			return nil, status.Error(codes.NotFound, err.Error())
		}

		return nil, status.Error(codes.Internal, fmt.Errorf("internal server error: %w", err).Error())
	}

	res := &comments.CommentInfoResponse{
		Id:        comment.ID,
		UserId:    comment.UserID,
		Sku:       comment.Sku,
		Content:   comment.Content,
		CreatedAt: timestamppb.New(comment.CreatedAt),
	}

	return res, nil
}

// EditCommentV1 изменяет комментарий.
func (c *CommentServerGRPC) EditCommentV1(ctx context.Context, req *comments.EditCommentRequest) (*comments.EditCommentResponse, error) {
	newComment := &domain.Comment{
		UserID:  req.UserId,
		Content: req.Content,
	}

	err := c.commentService.Edit(ctx, req.CommentId, newComment)
	if err != nil {
		if errors.Is(err, domain.ErrCommentNotExist) {
			return nil, status.Error(codes.NotFound, err.Error())
		}

		if errors.Is(err, domain.ErrEditNotMyComment) {
			return nil, status.Error(codes.PermissionDenied, err.Error())
		}

		if errors.Is(err, domain.ErrEditTimeoutExceed) {
			return nil, status.Error(codes.FailedPrecondition, err.Error())
		}

		return nil, status.Error(codes.Internal, fmt.Errorf("internal server error: %w", err).Error())
	}

	return &comments.EditCommentResponse{}, nil
}

// GetCommentsBySKUV1 возвращает список комментарием на товар.
func (c *CommentServerGRPC) GetCommentsBySKUV1(ctx context.Context, req *comments.GetCommentsBySKURequest) (*comments.GetCommentsBySKUResponse, error) {
	commentsData, err := c.commentService.GetListBySKU(ctx, req.Sku)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Errorf("internal server error: %w", err).Error())
	}

	res := &comments.GetCommentsBySKUResponse{
		Comments: make([]*comments.CommentBySKU, 0, len(commentsData)),
	}
	for _, comment := range commentsData {
		res.Comments = append(res.Comments,
			&comments.CommentBySKU{
				Id:        comment.ID,
				UserId:    comment.UserID,
				Content:   comment.Content,
				CreatedAt: timestamppb.New(comment.CreatedAt),
			},
		)
	}

	return res, nil
}

// GetCommentsBySKUV1 возвращает список комментарием, оставленных пользователем.
func (c *CommentServerGRPC) GetCommentsByUserV1(ctx context.Context, req *comments.GetCommentsByUserRequest) (*comments.GetCommentsByUserResponse, error) {
	commentsData, err := c.commentService.GetListByUser(ctx, req.UserId)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Errorf("internal server error: %w", err).Error())
	}

	res := &comments.GetCommentsByUserResponse{
		Comments: make([]*comments.CommentByUser, 0, len(commentsData)),
	}
	for _, comment := range commentsData {
		res.Comments = append(res.Comments,
			&comments.CommentByUser{
				Id:        comment.ID,
				Sku:       comment.Sku,
				Content:   comment.Content,
				CreatedAt: timestamppb.New(comment.CreatedAt),
			},
		)
	}

	return res, nil
}
