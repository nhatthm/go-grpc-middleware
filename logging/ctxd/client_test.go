package ctxd

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestNewClientLogger(t *testing.T) {
	t.Parallel()

	logger, buf := newCtxdLogger(LogLevelWarn)

	l := newClientLogger(logger,
		WithCodes(func(error) codes.Code {
			return codes.AlreadyExists
		}),
		WithLevels(func(codes.Code) LogLevel {
			return LogLevelWarn
		}),
		WithMessageProducer(func(context.Context, string, codes.Code, error, time.Duration) (context.Context, string) {
			return context.Background(), "error message"
		}),
	)

	code := l.errorToCode(nil)
	level := l.codeToLevel(code)

	l.Write(context.Background(), level, "", code, nil, 0)

	expected := `{"level":"warn","time":"<ignore-diff>","msg":"error message"}`

	assertLogMessage(t, expected, buf.String())
}

func TestUnaryClientInterceptor(t *testing.T) {
	t.Parallel()

	const method = "/grpctest.ItemService/GetItem"

	testCases := []struct {
		scenario           string
		context            context.Context
		loggerLevel        LogLevel
		options            []Option
		invoker            grpc.UnaryInvoker
		expectedError      string
		expectedLogMessage string
	}{
		{
			scenario:    "decider is not used in client calls",
			context:     context.Background(),
			loggerLevel: LogLevelDebug,
			options: []Option{
				WithDecider(func(string, error) bool {
					return false
				}),
			},
			invoker: func(context.Context, string, interface{}, interface{}, *grpc.ClientConn, ...grpc.CallOption) error {
				return status.Error(codes.Internal, "internal error")
			},
			expectedError: `rpc error: code = Internal desc = internal error`,
			expectedLogMessage: `{
    "level": "warn",
    "time": "<ignore-diff>",
    "msg": "finished client unary call",
    "system": "grpc",
    "span.kind": "client",
    "grpc.service": "grpctest.ItemService",
    "grpc.method": "GetItem",
    "grpc.start_time": "<ignore-diff>",
    "grpc.code": "Internal",
    "grpc.duration_ms": "<ignore-diff>",
    "error": "rpc error: code = Internal desc = internal error"
}`,
		},
		{
			scenario:    "error is written when logger is set at warn",
			context:     context.Background(),
			loggerLevel: LogLevelWarn,
			invoker: func(context.Context, string, interface{}, interface{}, *grpc.ClientConn, ...grpc.CallOption) error {
				return status.Error(codes.Internal, "internal error")
			},
			expectedError: `rpc error: code = Internal desc = internal error`,
			expectedLogMessage: `{
    "level": "warn",
    "time": "<ignore-diff>",
    "msg": "finished client unary call",
    "system": "grpc",
    "span.kind": "client",
    "grpc.service": "grpctest.ItemService",
    "grpc.method": "GetItem",
    "grpc.start_time": "<ignore-diff>",
    "grpc.code": "Internal",
    "grpc.duration_ms": "<ignore-diff>",
    "error": "rpc error: code = Internal desc = internal error"
}`,
		},
		{
			scenario:    "no error when logger is set at debug",
			context:     context.Background(),
			loggerLevel: LogLevelDebug,
			invoker: func(context.Context, string, interface{}, interface{}, *grpc.ClientConn, ...grpc.CallOption) error {
				return nil
			},
			expectedLogMessage: `{
    "level": "debug",
    "time": "<ignore-diff>",
    "msg": "finished client unary call",
    "system": "grpc",
    "span.kind": "client",
    "grpc.service": "grpctest.ItemService",
    "grpc.method": "GetItem",
    "grpc.start_time": "<ignore-diff>",
    "grpc.code": "OK",
    "grpc.duration_ms": "<ignore-diff>"
}`,
		},
		{
			scenario:    "with deadline",
			context:     contextWithDeadline(time.Now().Add(time.Hour)),
			loggerLevel: LogLevelDebug,
			invoker: func(context.Context, string, interface{}, interface{}, *grpc.ClientConn, ...grpc.CallOption) error {
				return nil
			},
			expectedLogMessage: `{
    "level": "debug",
    "time": "<ignore-diff>",
    "msg": "finished client unary call",
    "system": "grpc",
    "span.kind": "client",
    "grpc.service": "grpctest.ItemService",
    "grpc.method": "GetItem",
    "grpc.request.deadline": "<ignore-diff>",
    "grpc.start_time": "<ignore-diff>",
    "grpc.code": "OK",
    "grpc.duration_ms": "<ignore-diff>"
}`,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			logger, buf := newCtxdLogger(tc.loggerLevel)

			err := UnaryClientInterceptor(logger, tc.options...)(tc.context, method, nil, nil, nil, tc.invoker)

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}

			assertLogMessage(t, tc.expectedLogMessage, buf.String())
		})
	}
}

func TestStreamClientInterceptor(t *testing.T) {
	t.Parallel()

	const method = "/grpctest.ItemService/ListItems"

	testCases := []struct {
		scenario           string
		context            context.Context
		loggerLevel        LogLevel
		options            []Option
		handler            grpc.Streamer
		expectedError      string
		expectedLogMessage string
	}{
		{
			scenario:    "decider is not used in client calls",
			context:     context.Background(),
			loggerLevel: LogLevelDebug,
			options: []Option{
				WithDecider(func(string, error) bool {
					return false
				}),
			},
			handler: func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, opts ...grpc.CallOption) (grpc.ClientStream, error) {
				return nil, status.Error(codes.Internal, "internal error")
			},
			expectedError: `rpc error: code = Internal desc = internal error`,
			expectedLogMessage: `{
    "level": "warn",
    "time": "<ignore-diff>",
    "msg": "finished client streaming call",
    "system": "grpc",
    "span.kind": "client",
    "grpc.service": "grpctest.ItemService",
    "grpc.method": "ListItems",
    "grpc.start_time": "<ignore-diff>",
    "grpc.code": "Internal",
    "grpc.duration_ms": "<ignore-diff>",
    "error": "rpc error: code = Internal desc = internal error"
}`,
		},
		{
			scenario:    "error is written when logger is set at warn",
			context:     context.Background(),
			loggerLevel: LogLevelWarn,
			handler: func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, opts ...grpc.CallOption) (grpc.ClientStream, error) {
				return nil, status.Error(codes.Internal, "internal error")
			},
			expectedError: `rpc error: code = Internal desc = internal error`,
			expectedLogMessage: `{
    "level": "warn",
    "time": "<ignore-diff>",
    "msg": "finished client streaming call",
    "system": "grpc",
    "span.kind": "client",
    "grpc.service": "grpctest.ItemService",
    "grpc.method": "ListItems",
    "grpc.start_time": "<ignore-diff>",
    "grpc.code": "Internal",
    "grpc.duration_ms": "<ignore-diff>",
    "error": "rpc error: code = Internal desc = internal error"
}`,
		},
		{
			scenario:    "no error when logger is set at debug",
			context:     context.Background(),
			loggerLevel: LogLevelDebug,
			handler: func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, opts ...grpc.CallOption) (grpc.ClientStream, error) {
				return nil, nil
			},
			expectedLogMessage: `{
    "level": "debug",
    "time": "<ignore-diff>",
    "msg": "finished client streaming call",
    "system": "grpc",
    "span.kind": "client",
    "grpc.service": "grpctest.ItemService",
    "grpc.method": "ListItems",
    "grpc.start_time": "<ignore-diff>",
    "grpc.code": "OK",
    "grpc.duration_ms": "<ignore-diff>"
}`,
		},
		{
			scenario:    "with deadline",
			context:     contextWithDeadline(time.Now().Add(time.Hour)),
			loggerLevel: LogLevelDebug,
			handler: func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, opts ...grpc.CallOption) (grpc.ClientStream, error) {
				return nil, nil
			},
			expectedLogMessage: `{
    "level": "debug",
    "time": "<ignore-diff>",
    "msg": "finished client streaming call",
    "system": "grpc",
    "span.kind": "client",
    "grpc.service": "grpctest.ItemService",
    "grpc.method": "ListItems",
    "grpc.request.deadline": "<ignore-diff>",
    "grpc.start_time": "<ignore-diff>",
    "grpc.code": "OK",
    "grpc.duration_ms": "<ignore-diff>"
}`,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			logger, buf := newCtxdLogger(tc.loggerLevel)

			_, err := StreamClientInterceptor(logger, tc.options...)(tc.context, nil, nil, method, tc.handler)

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}

			assertLogMessage(t, tc.expectedLogMessage, buf.String())
		})
	}
}
