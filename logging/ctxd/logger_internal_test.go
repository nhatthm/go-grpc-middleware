package ctxd

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/bool64/ctxd"
	"github.com/bool64/zapctxd"
	"github.com/stretchr/testify/assert"
	"github.com/swaggest/assertjson"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc/codes"
)

func TestLogger_Write(t *testing.T) {
	t.Parallel()

	const (
		msg      = "hello"
		code     = codes.OK
		duration = time.Second
	)

	testCases := []struct {
		scenario string
		level    LogLevel
		expected string
	}{
		{
			scenario: "LogLevelDebug",
			level:    LogLevelDebug,
			expected: `{"level":"debug","time":"<ignore-diff>","msg":"hello","grpc.code":"OK","grpc.duration_ms":1000}`,
		},
		{
			scenario: "LogLevelInfo",
			level:    LogLevelInfo,
			expected: `{"level":"info","time":"<ignore-diff>","msg":"hello","grpc.code":"OK","grpc.duration_ms":1000}`,
		},
		{
			scenario: "LogLevelImportant",
			level:    LogLevelImportant,
			expected: `{"level":"info","time":"<ignore-diff>","msg":"hello","grpc.code":"OK","grpc.duration_ms":1000}`,
		},
		{
			scenario: "LogLevelWarn",
			level:    LogLevelWarn,
			expected: `{"level":"warn","time":"<ignore-diff>","msg":"hello","grpc.code":"OK","grpc.duration_ms":1000}`,
		},
		{
			scenario: "LogLevelError",
			level:    LogLevelError,
			expected: `{"level":"error","time":"<ignore-diff>","msg":"hello","grpc.code":"OK","grpc.duration_ms":1000}`,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			log, buf := newCtxdLogger(LogLevelDebug)

			l := defaultLogger(log)

			l.Write(context.Background(), tc.level, msg, code, nil, duration)

			assertjson.Equal(t, []byte(tc.expected), buf.Bytes())
		})
	}
}

func newCtxdLogger(level LogLevel) (ctxd.Logger, *bytes.Buffer) {
	buf := new(bytes.Buffer)
	cfg := zapctxd.Config{
		Level:  levelToZapLevel(level),
		Output: buf,
	}

	return zapctxd.New(cfg), buf
}

func levelToZapLevel(level LogLevel) zapcore.Level {
	switch level {
	case LogLevelDebug:
		return zapcore.DebugLevel
	case LogLevelInfo:
		return zapcore.InfoLevel
	case LogLevelImportant:
		return zapcore.InfoLevel
	case LogLevelWarn:
		return zapcore.WarnLevel
	case LogLevelError:
		return zapcore.ErrorLevel
	}

	return zapcore.InfoLevel
}

func assertLogMessage(t assert.TestingT, expected, actual string) bool { //nolint: unparam
	if expected == "" {
		if actual != "" {
			t.Errorf("no message is expected, got %s", actual)

			return false
		}

		return true
	}

	if actual == "" {
		t.Errorf("message is expected: %s, got nothing", expected)

		return false
	}

	return assertjson.Equal(t, []byte(expected), []byte(actual))
}

func contextWithDeadline(deadline time.Time) context.Context {
	//goland:noinspection GoVetLostCancel
	//nolint: govet
	ctx, _ := context.WithDeadline(context.Background(), deadline)

	return ctx
}
