package logger

// MultiLogger dispatches log calls to multiple logger strategies
// Similar to TypeScript's createLogger function, it allows logging to multiple destinations
// simultaneously (e.g., console + Sentry + ELK)
//
// Example:
//
//	logger := logger.NewMulti(
//	    logger.NewConsole(logger.ConsoleOptions{ServiceName: "api"}),
//	    logger.NewZap(zapLogger),
//	    logger.NewSentry(sentryHub),
//	)
type MultiLogger struct {
	strategies []Logger
}

// NewMulti creates a new multi-logger that dispatches to all provided strategies
func NewMulti(strategies ...Logger) Logger {
	return &MultiLogger{
		strategies: strategies,
	}
}

func (m *MultiLogger) Info(msg string, context ...any) {
	for _, strategy := range m.strategies {
		strategy.Info(msg, context...)
	}
}

func (m *MultiLogger) Error(msg string, context ...any) {
	for _, strategy := range m.strategies {
		strategy.Error(msg, context...)
	}
}

func (m *MultiLogger) Warn(msg string, context ...any) {
	for _, strategy := range m.strategies {
		strategy.Warn(msg, context...)
	}
}

func (m *MultiLogger) Debug(msg string, context ...any) {
	for _, strategy := range m.strategies {
		strategy.Debug(msg, context...)
	}
}

// Flush calls Flush on all strategies and returns the first error encountered
func (m *MultiLogger) Flush() error {
	var firstErr error
	for _, strategy := range m.strategies {
		if err := strategy.Flush(); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}
