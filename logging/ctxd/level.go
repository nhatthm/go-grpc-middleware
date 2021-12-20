package ctxd

const (
	// LogLevelDebug logs are typically voluminous, and are usually disabled in
	// production.
	LogLevelDebug LogLevel = iota - 1
	// LogLevelInfo is the default logging priority.
	LogLevelInfo
	// LogLevelImportant logs have the same level as Info, but won't be discarded in any circumstances.
	LogLevelImportant
	// LogLevelWarn logs are more important than Info, but don't need individual
	// human review.
	LogLevelWarn
	// LogLevelError logs are high-priority. If an application is running smoothly,
	// it shouldn't generate any error-level logs.
	LogLevelError
)

// LogLevel is a logging priority. Higher levels are more important.
type LogLevel int8
