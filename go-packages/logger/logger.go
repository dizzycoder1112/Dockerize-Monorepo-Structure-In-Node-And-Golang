package logger

// Logger is the unified logging interface used across all services and packages
// Similar to TypeScript's Logger interface, it provides a consistent API for logging
//
// Usage in packages (with fallback):
//
//	func NewClient(config Config, log logger.Logger) *Client {
//	    if log == nil {
//	        log = logger.Default  // fallback to console
//	    }
//	    return &Client{logger: log}
//	}
//
// Usage in applications (with multiple strategies):
//
//	log := logger.NewMulti(
//	    logger.NewZap(zapLogger),
//	    logger.NewSentry(sentryHub),
//	)
type Logger interface {
	// Info logs an informational message
	// context can be key-value pairs: "key1", value1, "key2", value2
	// or a map: map[string]any{"key1": value1, "key2": value2}
	Info(msg string, context ...any)

	// Error logs an error message
	Error(msg string, context ...any)

	// Warn logs a warning message
	Warn(msg string, context ...any)

	// Debug logs a debug message
	Debug(msg string, context ...any)

	// Flush waits for all pending async operations to complete
	// For sync loggers (Console, Zap), this flushes internal buffers
	// For async loggers (Slack), this waits for in-flight requests
	Flush() error
}
