package timeout

import (
	"context"
	"time"

	"google.golang.org/grpc"
)

// UnaryClientTimeoutInterceptor automatically start a context with timeout if it is not set.
func UnaryClientTimeoutInterceptor(duration time.Duration) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		ctx, cancel := withTimeout(ctx, duration)
		defer cancel()

		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

// UnaryClientSleepInterceptor automatically start a context with timeout if it is not set.
func UnaryClientSleepInterceptor(duration time.Duration) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		time.Sleep(duration)

		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

// WithUnaryClientTimeoutInterceptor appends UnaryClientTimeoutInterceptor to dial option.
func WithUnaryClientTimeoutInterceptor(duration time.Duration) grpc.DialOption {
	return grpc.WithChainUnaryInterceptor(UnaryClientTimeoutInterceptor(duration))
}

// WithUnaryClientSleepInterceptor appends UnaryClientSleepInterceptor to dial option.
func WithUnaryClientSleepInterceptor(duration time.Duration) grpc.DialOption {
	return grpc.WithChainUnaryInterceptor(UnaryClientSleepInterceptor(duration))
}
