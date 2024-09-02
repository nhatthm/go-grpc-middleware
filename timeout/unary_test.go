package timeout_test

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"

	"github.com/nhatthm/go-grpc-middleware/timeout"
)

func TestWithUnaryClientInterceptor(t *testing.T) {
	t.Parallel()

	duration := time.Millisecond * 20
	buf := bufconn.Listen(1024 * 1024)

	srv := grpc.NewServer()
	defer srv.GracefulStop()

	go func() {
		_ = srv.Serve(buf) //nolint: errcheck
	}()

	conn, err := grpc.NewClient("passthrough://",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return buf.Dial()
		}),
		timeout.WithUnaryClientTimeoutInterceptor(duration),
		// Timeout is set by the interceptor above, so we need to sleep for an amount of time, just enough to let the
		// context expire.
		//
		// Without this interceptor, the error is:
		//     `rpc error: code = Unimplemented desc = malformed method name: "test"`
		timeout.WithUnaryClientSleepInterceptor(duration*2),
	)
	require.NoError(t, err)

	err = conn.Invoke(context.Background(), "test", nil, nil)
	expected := `rpc error: code = DeadlineExceeded desc = context deadline exceeded`

	assert.EqualError(t, err, expected)
}
