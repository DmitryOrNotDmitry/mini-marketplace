package service

import (
	"context"
	"errors"
	"route256/comments/internal/domain"
	mock "route256/comments/mocks"
	"testing"
	"time"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testComponentCS struct {
	commentRepoMock *mock.CommentRepositoryMock
	commentService  *CommentService
}

func newTestComponentCS(t *testing.T) *testComponentCS {
	mc := minimock.NewController(t)
	commentRepoMock := mock.NewCommentRepositoryMock(mc)
	commentService := NewCommentService(commentRepoMock, time.Second, 10)

	return &testComponentCS{
		commentRepoMock: commentRepoMock,
		commentService:  commentService,
	}
}

func TestCommentService(t *testing.T) {
	t.Parallel()

	t.Run("create comment success", func(t *testing.T) {
		t.Parallel()

		tc := newTestComponentCS(t)

		ctx := context.Background()
		comment := &domain.Comment{
			UserID:  1,
			Sku:     1,
			Content: "comment 1",
		}
		expectedCommentID := 10

		tc.commentRepoMock.InsertMock.Return(int64(expectedCommentID), nil)

		commentID, err := tc.commentService.Add(ctx, comment)
		require.NoError(t, err)

		assert.EqualValues(t, expectedCommentID, commentID)

		var emptyTime time.Time
		assert.NotEqualValues(t, emptyTime, comment.CreatedAt)
	})

	t.Run("edit comment success", func(t *testing.T) {
		t.Parallel()

		tc := newTestComponentCS(t)

		ctx := context.Background()
		commentID := int64(1)
		newComment := &domain.Comment{
			UserID:  1,
			Content: "comment 1",
		}

		tc.commentRepoMock.UpdateContentWithCheckMock.Return(nil)

		err := tc.commentService.Edit(ctx, commentID, newComment)
		require.NoError(t, err)
	})

	t.Run("edit comment success", func(t *testing.T) {
		t.Parallel()

		tc := newTestComponentCS(t)

		ctx := context.Background()
		commentID := int64(1)
		oldComment := &domain.Comment{
			ID:        commentID,
			Sku:       1,
			UserID:    1,
			Content:   "comment 1",
			CreatedAt: time.Now(),
		}
		newComment := &domain.Comment{
			UserID:  1,
			Content: "comment 2",
		}

		tc.commentRepoMock.UpdateContentWithCheckMock.Set(func(_ context.Context, _ int64, _ *domain.Comment, predicate func(oldComment *domain.Comment) error) (err error) {
			return predicate(oldComment)
		})

		err := tc.commentService.Edit(ctx, commentID, newComment)
		require.NoError(t, err)
	})

	t.Run("edit comment error with timeout", func(t *testing.T) {
		t.Parallel()

		tc := newTestComponentCS(t)

		ctx := context.Background()
		commentID := int64(1)
		oldComment := &domain.Comment{
			ID:        commentID,
			Sku:       1,
			UserID:    1,
			Content:   "comment 1",
			CreatedAt: time.Now().Add(-10 * time.Second),
		}
		newComment := &domain.Comment{
			UserID:  1,
			Content: "comment 2",
		}

		tc.commentRepoMock.UpdateContentWithCheckMock.Set(func(_ context.Context, _ int64, _ *domain.Comment, predicate func(oldComment *domain.Comment) error) (err error) {
			return predicate(oldComment)
		})

		err := tc.commentService.Edit(ctx, commentID, newComment)
		require.Error(t, err)

		assert.True(t, errors.Is(err, domain.ErrEditTimeoutExceed))
	})

	t.Run("edit comment error not my comment", func(t *testing.T) {
		t.Parallel()

		tc := newTestComponentCS(t)

		ctx := context.Background()
		commentID := int64(1)
		oldComment := &domain.Comment{
			ID:        commentID,
			Sku:       1,
			UserID:    1,
			Content:   "comment 1",
			CreatedAt: time.Now(),
		}
		newComment := &domain.Comment{
			UserID:  2,
			Content: "comment 2",
		}

		tc.commentRepoMock.UpdateContentWithCheckMock.Set(func(_ context.Context, _ int64, _ *domain.Comment, predicate func(oldComment *domain.Comment) error) (err error) {
			return predicate(oldComment)
		})

		err := tc.commentService.Edit(ctx, commentID, newComment)
		require.Error(t, err)

		assert.True(t, errors.Is(err, domain.ErrEditNotMyComment))
	})

	t.Run("get comments by user success", func(t *testing.T) {
		t.Parallel()

		tc := newTestComponentCS(t)

		ctx := context.Background()
		count := 10
		comments := make([]*domain.Comment, 0, count)
		for i := 0; i < count; i++ {
			comments = append(comments, &domain.Comment{
				ID:        int64(i),
				UserID:    1,
				Sku:       1,
				Content:   "content",
				CreatedAt: time.Now().Add(time.Duration(i) * time.Second),
			})
		}

		tc.commentRepoMock.GetListByUserMock.Return(comments, nil)

		actualComments, err := tc.commentService.GetListByUser(ctx, 1)
		require.NoError(t, err)

		assert.EqualValues(t, len(comments), len(actualComments))

		for i := 0; i < count-1; i++ {
			assert.GreaterOrEqual(t, actualComments[i].CreatedAt, actualComments[i+1].CreatedAt)
		}
	})

	t.Run("get comments by sku success", func(t *testing.T) {
		t.Parallel()

		tc := newTestComponentCS(t)

		ctx := context.Background()
		count := 10
		comments := make([]*domain.Comment, 0, count)
		for i := 0; i < count; i++ {
			comments = append(comments, &domain.Comment{
				ID:        int64(i),
				UserID:    int64(i),
				Sku:       1,
				Content:   "content",
				CreatedAt: time.Now().Add(-time.Duration(i) * time.Second),
			})
		}

		tc.commentRepoMock.GetListBySKUMock.Return(comments, nil)

		actualComments, err := tc.commentService.GetListBySKU(ctx, 1, time.Now(), 0)
		require.NoError(t, err)

		assert.EqualValues(t, len(comments), len(actualComments))
	})

	t.Run("get comment by id success", func(t *testing.T) {
		t.Parallel()

		tc := newTestComponentCS(t)

		ctx := context.Background()
		comment := &domain.Comment{
			ID:        1,
			UserID:    1,
			Sku:       1,
			Content:   "content",
			CreatedAt: time.Now(),
		}

		tc.commentRepoMock.GetByIDMock.Return(comment, nil)

		_, err := tc.commentService.GetInfoByID(ctx, 1)
		require.NoError(t, err)
	})

}
