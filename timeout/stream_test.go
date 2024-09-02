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

func TestWithStreamClientInterceptors(t *testing.T) {
	t.Parallel()

	duration := time.Millisecond * 20
	buf := bufconn.Listen(1024 * 1024)

	conn, err := grpc.NewClient("passthrough://",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return buf.Dial()
		}),
		timeout.WithStreamClientTimeoutInterceptor(duration),
		// Timeout is set by the interceptor above, so we need to sleep for an amount of time, just enough to let the
		// context expire.
		//
		// Without this interceptor, the returned stream is not nil and no error.
		timeout.WithStreamClientSleepInterceptor(duration*2),
	)
	require.NoError(t, err)

	s, err := conn.NewStream(context.Background(), nil, "test")
	expected := `rpc error: code = DeadlineExceeded desc = context deadline exceeded`

	assert.Nil(t, s)
	assert.EqualError(t, err, expected)
}
