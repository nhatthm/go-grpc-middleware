package timeout

import (
	"context"
	"time"
)

func withTimeout(ctx context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	_, deadlineIsSet := ctx.Deadline()
	cancel := func() {
		// Fake cancel function in case there is no timeout.
	}

	if timeout == 0 || deadlineIsSet {
		return ctx, cancel
	}

	return context.WithTimeout(ctx, timeout)
}
