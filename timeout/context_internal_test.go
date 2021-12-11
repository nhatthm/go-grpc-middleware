package timeout

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWithTimeout_NoCancel(t *testing.T) {
	t.Parallel()

	timeout := time.Millisecond * 50

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	newCtx, cancelNewCtx := withTimeout(ctx, timeout)

	cancelNewCtx()

	// The 2nd cancel doesn't do anything because the context was set with timeout.
	assert.NoError(t, ctx.Err())
	assert.NoError(t, newCtx.Err())

	time.Sleep(timeout * 2)

	assert.ErrorIs(t, ctx.Err(), context.DeadlineExceeded)
	assert.ErrorIs(t, newCtx.Err(), context.DeadlineExceeded)
}

func TestWithTimeout_NoTimeout(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	newCtx, cancel := withTimeout(ctx, 0)

	cancel()

	// Cancel does not affect the contexts because there is no timeout.
	assert.NoError(t, ctx.Err())
	assert.NoError(t, newCtx.Err())
}

func TestWithTimeout_WithTimeout(t *testing.T) {
	t.Parallel()

	timeout := time.Millisecond * 50
	ctx := context.Background()

	newCtx, cancel := withTimeout(ctx, timeout)
	defer cancel()

	time.Sleep(timeout * 2)

	assert.NoError(t, ctx.Err())
	assert.ErrorIs(t, newCtx.Err(), context.DeadlineExceeded)
}
