package ctxd

import (
	"context"
	"path"
	"time"

	"github.com/bool64/ctxd"
	"google.golang.org/grpc"
)

func newClientLogger(log ctxd.Logger, opts ...Option) *logger {
	l := defaultLogger(log)
	l.codeToLevel = DefaultClientCodeToLevel

	for _, o := range opts {
		o(l)
	}

	return l
}

// UnaryClientInterceptor returns a new unary client interceptor that optionally logs the execution of external gRPC calls.
func UnaryClientInterceptor(logger ctxd.Logger, opts ...Option) grpc.UnaryClientInterceptor {
	l := newClientLogger(logger, opts...)

	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		startTime := time.Now()

		ctx = clientLoggerContext(ctx, method, startTime)
		err := invoker(ctx, method, req, reply, cc, opts...)

		duration := time.Since(startTime)

		code := l.errorToCode(err)
		level := l.codeToLevel(code)

		l.Write(ctx, level, "finished client unary call", code, err, duration)

		return err
	}
}

// StreamClientInterceptor returns a new streaming client interceptor that optionally logs the execution of external gRPC calls.
func StreamClientInterceptor(logger ctxd.Logger, opts ...Option) grpc.StreamClientInterceptor {
	l := newClientLogger(logger, opts...)

	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		startTime := time.Now()

		ctx = clientLoggerContext(ctx, method, startTime)
		clientStream, err := streamer(ctx, desc, cc, method, opts...)

		duration := time.Since(startTime)

		code := l.errorToCode(err)
		level := l.codeToLevel(code)

		l.Write(ctx, level, "finished client streaming call", code, err, duration)

		return clientStream, err
	}
}

func clientLoggerContext(ctx context.Context, fullMethodString string, start time.Time) context.Context {
	service := path.Dir(fullMethodString)[1:]
	method := path.Base(fullMethodString)

	ctx = ctxd.AddFields(ctx,
		FieldSystem, "grpc",
		FieldKind, "client",
		FieldService, service,
		FieldMethod, method,
		FieldStartTime, start,
	)

	if d, ok := ctx.Deadline(); ok {
		ctx = ctxd.AddFields(ctx, FieldDeadline, d)
	}

	return ctx
}
