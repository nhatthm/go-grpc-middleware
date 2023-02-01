package ctxd

import (
	"context"
	"time"

	"github.com/bool64/ctxd"
	grpcLogging "github.com/grpc-ecosystem/go-grpc-middleware/logging"
	"google.golang.org/grpc/codes"
)

const (
	// FieldSystem is a context field.
	FieldSystem = "system"
	// FieldKind is a context field.
	FieldKind = "span.kind"
	// FieldService is a context field.
	FieldService = "grpc.service"
	// FieldMethod is a context field.
	FieldMethod = "grpc.method"
	// FieldStartTime is a context field for execution time.
	FieldStartTime = "grpc.start_time"
	// FieldDeadline is a context field for return code.
	FieldDeadline = "grpc.request.deadline"
	// FieldDuration is a context field for execution time.
	FieldDuration = "grpc.duration_ms"
	// FieldCode is a context field for return code.
	FieldCode = "grpc.code"
)

// CodeToLevel function defines the mapping between gRPC return codes and interceptor log level.
type CodeToLevel func(code codes.Code) LogLevel

// MessageProducer produces a user defined log message.
type MessageProducer func(ctx context.Context, msg string, code codes.Code, err error, duration time.Duration) (context.Context, string)

// Option to set up the logger.
type Option func(l *logger)

type logger struct {
	log ctxd.Logger

	shouldLog      grpcLogging.Decider
	errorToCode    grpcLogging.ErrorToCode
	codeToLevel    CodeToLevel
	produceMessage MessageProducer
}

func defaultLogger(log ctxd.Logger) *logger {
	l := &logger{
		log:            log,
		shouldLog:      grpcLogging.DefaultDeciderMethod,
		errorToCode:    grpcLogging.DefaultErrorToCode,
		codeToLevel:    DefaultCodeToLevel,
		produceMessage: DefaultMessageProducer,
	}

	return l
}

func (l *logger) Write(ctx context.Context, level LogLevel, msg string, code codes.Code, err error, duration time.Duration) {
	ctx, msg = l.produceMessage(ctx, msg, code, err, duration)

	switch level {
	case LogLevelDebug:
		l.log.Debug(ctx, msg)
	case LogLevelInfo:
		l.log.Info(ctx, msg)
	case LogLevelImportant:
		l.log.Important(ctx, msg)
	case LogLevelWarn:
		l.log.Warn(ctx, msg)
	case LogLevelError:
		l.log.Error(ctx, msg)
	}
}

// WithDecider customizes the function for deciding if the gRPC interceptor logs should log.
func WithDecider(d grpcLogging.Decider) Option {
	return func(l *logger) {
		l.shouldLog = d
	}
}

// WithLevels customizes the function for mapping gRPC return codes and interceptor log level statements.
func WithLevels(f CodeToLevel) Option {
	return func(l *logger) {
		l.codeToLevel = f
	}
}

// WithCodes customizes the function for mapping errors to error codes.
func WithCodes(f grpcLogging.ErrorToCode) Option {
	return func(l *logger) {
		l.errorToCode = f
	}
}

// WithMessageProducer customizes the function for message formation.
func WithMessageProducer(f MessageProducer) Option {
	return func(l *logger) {
		l.produceMessage = f
	}
}

// DefaultCodeToLevel is the default implementation of gRPC return codes and interceptor log level for server side.
func DefaultCodeToLevel(code codes.Code) LogLevel { //nolint: cyclop,dupl
	switch code {
	case codes.OK:
		return LogLevelInfo
	case codes.Canceled:
		return LogLevelInfo
	case codes.Unknown:
		return LogLevelError
	case codes.InvalidArgument:
		return LogLevelInfo
	case codes.DeadlineExceeded:
		return LogLevelWarn
	case codes.NotFound:
		return LogLevelInfo
	case codes.AlreadyExists:
		return LogLevelInfo
	case codes.PermissionDenied:
		return LogLevelWarn
	case codes.Unauthenticated:
		return LogLevelInfo // unauthenticated requests can happen
	case codes.ResourceExhausted:
		return LogLevelWarn
	case codes.FailedPrecondition:
		return LogLevelWarn
	case codes.Aborted:
		return LogLevelWarn
	case codes.OutOfRange:
		return LogLevelWarn
	case codes.Unimplemented:
		return LogLevelError
	case codes.Internal:
		return LogLevelError
	case codes.Unavailable:
		return LogLevelWarn
	case codes.DataLoss:
		return LogLevelError
	default:
		return LogLevelError
	}
}

// DefaultClientCodeToLevel is the default implementation of gRPC return codes to log levels for client side.
func DefaultClientCodeToLevel(code codes.Code) LogLevel { //nolint: cyclop,dupl
	switch code {
	case codes.OK:
		return LogLevelDebug
	case codes.Canceled:
		return LogLevelDebug
	case codes.Unknown:
		return LogLevelInfo
	case codes.InvalidArgument:
		return LogLevelDebug
	case codes.DeadlineExceeded:
		return LogLevelInfo
	case codes.NotFound:
		return LogLevelDebug
	case codes.AlreadyExists:
		return LogLevelDebug
	case codes.PermissionDenied:
		return LogLevelInfo
	case codes.Unauthenticated:
		return LogLevelInfo // Unauthenticated requests can happen.
	case codes.ResourceExhausted:
		return LogLevelDebug
	case codes.FailedPrecondition:
		return LogLevelDebug
	case codes.Aborted:
		return LogLevelDebug
	case codes.OutOfRange:
		return LogLevelDebug
	case codes.Unimplemented:
		return LogLevelWarn
	case codes.Internal:
		return LogLevelWarn
	case codes.Unavailable:
		return LogLevelWarn
	case codes.DataLoss:
		return LogLevelWarn
	default:
		return LogLevelInfo
	}
}

// DurationInMilliseconds returns duration in ms format.
func DurationInMilliseconds(d time.Duration) float32 {
	return float32(d.Nanoseconds()/1000) / 1000
}

// DefaultMessageProducer sets the log message and fields.
func DefaultMessageProducer(ctx context.Context, msg string, code codes.Code, err error, duration time.Duration) (context.Context, string) {
	ctx = ctxd.AddFields(ctx,
		FieldCode, code,
		FieldDuration, DurationInMilliseconds(duration),
	)

	if err != nil {
		ctx = ctxd.AddFields(ctx, "error", err)
	}

	return ctx, msg
}
