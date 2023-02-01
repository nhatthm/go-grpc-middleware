package ctxd

import (
	"context"
	"path"
	"time"

	"github.com/bool64/ctxd"
	grpcMiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
)

func newServerLogger(log ctxd.Logger, opts ...Option) *logger {
	l := defaultLogger(log)
	l.codeToLevel = DefaultCodeToLevel

	for _, o := range opts {
		o(l)
	}

	return l
}

// UnaryServerInterceptor returns a new unary server interceptors that adds zap.Logger to the context.
func UnaryServerInterceptor(logger ctxd.Logger, opts ...Option) grpc.UnaryServerInterceptor {
	l := newServerLogger(logger, opts...)

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		startTime := time.Now()

		ctx = serverLoggerContext(ctx, info.FullMethod, startTime)
		resp, err := handler(ctx, req)

		duration := time.Since(startTime)

		if !l.shouldLog(info.FullMethod, err) {
			return resp, err
		}

		code := l.errorToCode(err)
		level := l.codeToLevel(code)

		l.Write(ctx, level, "finished unary call", code, err, duration)

		return resp, err
	}
}

// StreamServerInterceptor returns a new streaming server interceptor that adds zap.Logger to the context.
func StreamServerInterceptor(logger ctxd.Logger, opts ...Option) grpc.StreamServerInterceptor {
	l := newServerLogger(logger, opts...)

	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		startTime := time.Now()

		ctx := serverLoggerContext(stream.Context(), info.FullMethod, startTime)
		wrapped := grpcMiddleware.WrapServerStream(stream)
		wrapped.WrappedContext = ctx

		err := handler(srv, wrapped)

		duration := time.Since(startTime)

		if !l.shouldLog(info.FullMethod, err) {
			return err
		}

		code := l.errorToCode(err)
		level := l.codeToLevel(code)

		l.Write(ctx, level, "finished streaming call", code, err, duration)

		return err
	}
}

func serverLoggerContext(ctx context.Context, fullMethodString string, start time.Time) context.Context {
	service := path.Dir(fullMethodString)[1:]
	method := path.Base(fullMethodString)

	ctx = ctxd.AddFields(ctx,
		FieldSystem, "grpc",
		FieldKind, "server",
		FieldService, service,
		FieldMethod, method,
		FieldStartTime, start,
	)

	if d, ok := ctx.Deadline(); ok {
		ctx = ctxd.AddFields(ctx, FieldDeadline, d)
	}

	return ctx
}
