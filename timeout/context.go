package timeout

import (
	"context"
	"time"
)

type skipTimeoutCtxKey struct{}

// IsTimeoutSkipped checks whether the timeout interceptor is bypassed.
func IsTimeoutSkipped(ctx context.Context) bool {
	skipped, found := ctx.Value(skipTimeoutCtxKey{}).(bool)

	return found && skipped
}

// SkipTimeout skips the timeout.
func SkipTimeout(ctx context.Context) context.Context {
	return context.WithValue(ctx, skipTimeoutCtxKey{}, true)
}

func withTimeout(ctx context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	_, deadlineIsSet := ctx.Deadline()
	cancel := func() {
		// Fake cancel function in case there is no timeout.
	}

	if timeout == 0 || deadlineIsSet || IsTimeoutSkipped(ctx) {
		return ctx, cancel
	}

	return context.WithTimeout(ctx, timeout)
}
