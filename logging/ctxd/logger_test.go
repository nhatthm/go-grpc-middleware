package ctxd_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"

	"github.com/nhatthm/go-grpc-middleware/logging/ctxd"
)

func TestDefaultCodeToLevel(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		code  codes.Code
		level ctxd.LogLevel
	}{
		{
			code:  codes.OK,
			level: ctxd.LogLevelInfo,
		},
		{
			code:  codes.Canceled,
			level: ctxd.LogLevelInfo,
		},
		{
			code:  codes.Unknown,
			level: ctxd.LogLevelError,
		},
		{
			code:  codes.InvalidArgument,
			level: ctxd.LogLevelInfo,
		},
		{
			code:  codes.DeadlineExceeded,
			level: ctxd.LogLevelWarn,
		},
		{
			code:  codes.NotFound,
			level: ctxd.LogLevelInfo,
		},
		{
			code:  codes.AlreadyExists,
			level: ctxd.LogLevelInfo,
		},
		{
			code:  codes.PermissionDenied,
			level: ctxd.LogLevelWarn,
		},
		{
			code:  codes.Unauthenticated,
			level: ctxd.LogLevelInfo,
		},
		{
			code:  codes.ResourceExhausted,
			level: ctxd.LogLevelWarn,
		},
		{
			code:  codes.FailedPrecondition,
			level: ctxd.LogLevelWarn,
		},
		{
			code:  codes.Aborted,
			level: ctxd.LogLevelWarn,
		},
		{
			code:  codes.OutOfRange,
			level: ctxd.LogLevelWarn,
		},
		{
			code:  codes.Unimplemented,
			level: ctxd.LogLevelError,
		},
		{
			code:  codes.Internal,
			level: ctxd.LogLevelError,
		},
		{
			code:  codes.Unavailable,
			level: ctxd.LogLevelWarn,
		},
		{
			code:  codes.DataLoss,
			level: ctxd.LogLevelError,
		},
		{
			code:  100,
			level: ctxd.LogLevelError,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.code.String(), func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tc.level, ctxd.DefaultCodeToLevel(tc.code))
		})
	}
}

func TestDefaultClientCodeToLevel(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		code  codes.Code
		level ctxd.LogLevel
	}{
		{
			code:  codes.OK,
			level: ctxd.LogLevelDebug,
		},
		{
			code:  codes.Canceled,
			level: ctxd.LogLevelDebug,
		},
		{
			code:  codes.Unknown,
			level: ctxd.LogLevelInfo,
		},
		{
			code:  codes.InvalidArgument,
			level: ctxd.LogLevelDebug,
		},
		{
			code:  codes.DeadlineExceeded,
			level: ctxd.LogLevelInfo,
		},
		{
			code:  codes.NotFound,
			level: ctxd.LogLevelDebug,
		},
		{
			code:  codes.AlreadyExists,
			level: ctxd.LogLevelDebug,
		},
		{
			code:  codes.PermissionDenied,
			level: ctxd.LogLevelInfo,
		},
		{
			code:  codes.Unauthenticated,
			level: ctxd.LogLevelInfo,
		},
		{
			code:  codes.ResourceExhausted,
			level: ctxd.LogLevelDebug,
		},
		{
			code:  codes.FailedPrecondition,
			level: ctxd.LogLevelDebug,
		},
		{
			code:  codes.Aborted,
			level: ctxd.LogLevelDebug,
		},
		{
			code:  codes.OutOfRange,
			level: ctxd.LogLevelDebug,
		},
		{
			code:  codes.Unimplemented,
			level: ctxd.LogLevelWarn,
		},
		{
			code:  codes.Internal,
			level: ctxd.LogLevelWarn,
		},
		{
			code:  codes.Unavailable,
			level: ctxd.LogLevelWarn,
		},
		{
			code:  codes.DataLoss,
			level: ctxd.LogLevelWarn,
		},
		{
			code:  100,
			level: ctxd.LogLevelInfo,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.code.String(), func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tc.level, ctxd.DefaultClientCodeToLevel(tc.code))
		})
	}
}
