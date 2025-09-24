package myerrgroup

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestErrorGroup(t *testing.T) {
	t.Parallel()

	t.Run("with no errors", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		errGroup, _ := WithContext(ctx)

		for range 2 {
			errGroup.Go(func() error {
				return nil
			})
		}

		err := errGroup.Wait()
		require.NoError(t, err)
	})

	t.Run("with error", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		errGroup, ctx := WithContext(ctx)

		errGroup.Go(func() error {
			return errors.New("error")
		})
		for range 10 {
			errGroup.Go(func() error {
				return nil
			})
		}

		err := errGroup.Wait()
		require.Error(t, err)
		require.Error(t, ctx.Err())
	})
}
