package retry

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
	"time"
)

func TestRetry(t *testing.T) {
	t.Run("retry with interval", func(t *testing.T) {
		var count int
		err := Times(5).WithInterval(time.Second).Do(func() error {
			count++
			return nil
		})
		assert.NoError(t, err)
		assert.Equal(t, 1, count)
	})

	t.Run("retry with interval and err", func(t *testing.T) {
		var count int
		err := Times(5).WithInterval(time.Second).Do(func() error {
			count++
			return errors.New("test err")
		})
		assert.Equal(t, "test err", strings.TrimSpace(err.Error()))
		assert.Equal(t, 5, count)
	})

	t.Run("retry with ctx", func(t *testing.T) {
		var count int
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()
		err := Times(5).
			WithInterval(time.Second).
			WithContext(ctx).
			Do(func() error {
				count++
				return errors.New("test err")
			})
		assert.Equal(t, "test err", strings.TrimSpace(err.Error()))
		assert.Equal(t, 3, count)
	})
}
