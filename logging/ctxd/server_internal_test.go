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

func TestNewServerLogger(t *testing.T) {
	t.Parallel()

	logger, buf := newCtxdLogger(LogLevelWarn)

	l := newServerLogger(logger,
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

func TestUnaryServerInterceptor(t *testing.T) {
	t.Parallel()

	info := &grpc.UnaryServerInfo{
		FullMethod: "/grpctest.ItemService/GetItem",
	}

	testCases := []struct {
		scenario           string
		context            context.Context
		loggerLevel        LogLevel
		options            []Option
		handler            grpc.UnaryHandler
		expectedResponse   interface{}
		expectedError      string
		expectedLogMessage string
	}{
		{
			scenario:    "should not log",
			context:     context.Background(),
			loggerLevel: LogLevelDebug,
			options: []Option{
				WithDecider(func(string, error) bool {
					return false
				}),
			},
			handler: func(context.Context, interface{}) (interface{}, error) {
				return nil, status.Error(codes.Internal, "internal error")
			},
			expectedError: `rpc error: code = Internal desc = internal error`,
		},
		{
			scenario:    "error is written when logger is set at warn",
			context:     context.Background(),
			loggerLevel: LogLevelWarn,
			handler: func(context.Context, interface{}) (interface{}, error) {
				return nil, status.Error(codes.Internal, "internal error")
			},
			expectedError: `rpc error: code = Internal desc = internal error`,
			expectedLogMessage: `{
    "level": "error",
    "time": "<ignore-diff>",
    "msg": "finished unary call",
    "system": "grpc",
    "span.kind": "server",
    "grpc.service": "grpctest.ItemService",
    "grpc.method": "GetItem",
    "grpc.start_time": "<ignore-diff>",
    "grpc.code": "Internal",
    "grpc.duration_ms": "<ignore-diff>",
    "error": "rpc error: code = Internal desc = internal error"
}`,
		},
		{
			scenario:    "no error when logger is set at info",
			context:     context.Background(),
			loggerLevel: LogLevelInfo,
			handler: func(context.Context, interface{}) (interface{}, error) {
				return 42, nil
			},
			expectedResponse: 42,
			expectedLogMessage: `{
    "level": "info",
    "time": "<ignore-diff>",
    "msg": "finished unary call",
    "system": "grpc",
    "span.kind": "server",
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
			loggerLevel: LogLevelInfo,
			handler: func(context.Context, interface{}) (interface{}, error) {
				return 42, nil
			},
			expectedResponse: 42,
			expectedLogMessage: `{
    "level": "info",
    "time": "<ignore-diff>",
    "msg": "finished unary call",
    "system": "grpc",
    "span.kind": "server",
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

			resp, err := UnaryServerInterceptor(logger, tc.options...)(tc.context, nil, info, tc.handler)

			assert.Equal(t, tc.expectedResponse, resp)

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}

			assertLogMessage(t, tc.expectedLogMessage, buf.String())
		})
	}
}

func TestStreamServerInterceptor(t *testing.T) {
	t.Parallel()

	info := &grpc.StreamServerInfo{
		FullMethod: "/grpctest.ItemService/ListItems",
	}

	testCases := []struct {
		scenario           string
		context            context.Context
		loggerLevel        LogLevel
		options            []Option
		handler            grpc.StreamHandler
		expectedError      string
		expectedLogMessage string
	}{
		{
			scenario:    "should not log",
			context:     context.Background(),
			loggerLevel: LogLevelDebug,
			options: []Option{
				WithDecider(func(string, error) bool {
					return false
				}),
			},
			handler: func(srv interface{}, stream grpc.ServerStream) error {
				return status.Error(codes.Internal, "internal error")
			},
			expectedError: `rpc error: code = Internal desc = internal error`,
		},
		{
			scenario:    "error is written when logger is set at warn",
			context:     context.Background(),
			loggerLevel: LogLevelWarn,
			handler: func(srv interface{}, stream grpc.ServerStream) error {
				return status.Error(codes.Internal, "internal error")
			},
			expectedError: `rpc error: code = Internal desc = internal error`,
			expectedLogMessage: `{
    "level": "error",
    "time": "<ignore-diff>",
    "msg": "finished streaming call",
    "system": "grpc",
    "span.kind": "server",
    "grpc.service": "grpctest.ItemService",
    "grpc.method": "ListItems",
    "grpc.start_time": "<ignore-diff>",
    "grpc.code": "Internal",
    "grpc.duration_ms": "<ignore-diff>",
    "error": "rpc error: code = Internal desc = internal error"
}`,
		},
		{
			scenario:    "no error when logger is set at info",
			context:     context.Background(),
			loggerLevel: LogLevelInfo,
			handler: func(srv interface{}, stream grpc.ServerStream) error {
				return nil
			},
			expectedLogMessage: `{
    "level": "info",
    "time": "<ignore-diff>",
    "msg": "finished streaming call",
    "system": "grpc",
    "span.kind": "server",
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
			loggerLevel: LogLevelInfo,
			handler: func(srv interface{}, stream grpc.ServerStream) error {
				return nil
			},
			expectedLogMessage: `{
    "level": "info",
    "time": "<ignore-diff>",
    "msg": "finished streaming call",
    "system": "grpc",
    "span.kind": "server",
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
			stream := serverStreamWithContext(tc.context)

			err := StreamServerInterceptor(logger, tc.options...)(nil, stream, info, tc.handler)

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}

			assertLogMessage(t, tc.expectedLogMessage, buf.String())
		})
	}
}

func serverStreamWithContext(ctx context.Context) grpc.ServerStream {
	return &serverStream{context: ctx}
}

type serverStream struct {
	grpc.ServerStream
	context context.Context
}

func (s *serverStream) Context() context.Context {
	return s.context
}
