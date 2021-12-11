package timeout

import (
	"context"
	"time"

	"google.golang.org/grpc"
)

// StreamClientTimeoutInterceptor automatically start a context with timeout if it is not set.
func StreamClientTimeoutInterceptor(duration time.Duration) grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		ctx, cancel := withTimeout(ctx, duration)
		defer cancel()

		return streamer(ctx, desc, cc, method, opts...)
	}
}

// StreamClientSleepInterceptor sleeps for a moment before doing the job.
func StreamClientSleepInterceptor(duration time.Duration) grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		time.Sleep(duration)

		return streamer(ctx, desc, cc, method, opts...)
	}
}

// WithStreamClientTimeoutInterceptor appends StreamClientTimeoutInterceptor to dial option.
func WithStreamClientTimeoutInterceptor(duration time.Duration) grpc.DialOption {
	return grpc.WithChainStreamInterceptor(StreamClientTimeoutInterceptor(duration))
}

// WithStreamClientSleepInterceptor appends StreamClientSleepInterceptor to dial option.
func WithStreamClientSleepInterceptor(duration time.Duration) grpc.DialOption {
	return grpc.WithChainStreamInterceptor(StreamClientSleepInterceptor(duration))
}
