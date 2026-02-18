package strategies

import (
	"dizzycoder1112/Dockerize-Monorepo-Structure-In-Node-And-Golang/logger"
	"encoding/json"
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// ZapLogger wraps uber/zap logger to implement the Logger interface
// This allows using structured logging with production-grade features
type ZapLogger struct {
	zap         *zap.Logger
	serviceName string
	isPretty    bool // ✅ Enable multi-line JSON formatting in development mode
}

// ZapOptions configures the Zap logger
type ZapOptions struct {
	ServiceName string
	IsPretty    bool  // Enable pretty console output (for development)
	Level       Level // Minimum log level
}

// Level represents log levels
type Level int

const (
	DebugLevel Level = iota
	InfoLevel
	WarnLevel
	ErrorLevel
)

// ToZapLevel converts our Level to zap's level
func (l Level) ToZapLevel() zapcore.Level {
	switch l {
	case DebugLevel:
		return zapcore.DebugLevel
	case InfoLevel:
		return zapcore.InfoLevel
	case WarnLevel:
		return zapcore.WarnLevel
	case ErrorLevel:
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}

// NewZap creates a new Zap-based logger
func NewZap(opts ZapOptions) (logger.Logger, error) {
	var zapLogger *zap.Logger
	var err error

	if opts.IsPretty {
		// Development mode: pretty console output
		config := zap.NewDevelopmentConfig()
		config.Level = zap.NewAtomicLevelAt(opts.Level.ToZapLevel())
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

		zapLogger, err = config.Build(
			zap.AddCallerSkip(1),              // Skip one level to show correct caller
			zap.AddStacktrace(zapcore.ErrorLevel), // Only stack trace on Error+, not Warn
		)
	} else {
		// Production mode: JSON output
		config := zap.NewProductionConfig()
		config.Level = zap.NewAtomicLevelAt(opts.Level.ToZapLevel())
		config.EncoderConfig.TimeKey = "timestamp"
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

		zapLogger, err = config.Build(
			zap.AddCallerSkip(1),
		)
	}

	if err != nil {
		return nil, err
	}

	return &ZapLogger{
		zap:         zapLogger,
		serviceName: opts.ServiceName,
		isPretty:    opts.IsPretty, // ✅ Save pretty mode configuration
	}, nil
}

// NewZapMust creates a new Zap logger and panics on error
func NewZapMust(opts ZapOptions) logger.Logger {
	logger, err := NewZap(opts)
	if err != nil {
		panic(err)
	}
	return logger
}

func (z *ZapLogger) convertContext(context []any) []zap.Field {
	contextMap := logger.ParseContext(context)
	if len(contextMap) == 0 {
		return []zap.Field{zap.String("service", z.serviceName)}
	}

	fields := make([]zap.Field, 0, len(contextMap)+1)
	fields = append(fields, zap.String("service", z.serviceName))

	for key, value := range contextMap {
		// Use zap.Any which automatically chooses the best field type
		fields = append(fields, zap.Any(key, value))
	}

	return fields
}

// formatPrettyMessage formats message with multi-line JSON for development mode
func (z *ZapLogger) formatPrettyMessage(msg string, context []any) string {
	if len(context) == 0 {
		return msg
	}

	contextMap := logger.ParseContext(context)
	if len(contextMap) == 0 {
		return msg
	}

	// Add service name to context
	contextMap["service"] = z.serviceName

	// Format as indented JSON
	jsonBytes, err := json.MarshalIndent(contextMap, "", "  ")
	if err != nil {
		// Fallback to original message if JSON formatting fails
		return msg
	}

	return fmt.Sprintf("%s\n%s", msg, string(jsonBytes))
}

func (z *ZapLogger) Info(msg string, context ...any) {
	if z.isPretty && len(context) > 0 {
		// ✅ Pretty mode: format as multi-line JSON
		prettyMsg := z.formatPrettyMessage(msg, context)
		z.zap.Info(prettyMsg)
	} else {
		// Production mode: structured fields
		fields := z.convertContext(context)
		z.zap.Info(msg, fields...)
	}
}

func (z *ZapLogger) Error(msg string, context ...any) {
	if z.isPretty && len(context) > 0 {
		prettyMsg := z.formatPrettyMessage(msg, context)
		z.zap.Error(prettyMsg)
	} else {
		fields := z.convertContext(context)
		z.zap.Error(msg, fields...)
	}
}

func (z *ZapLogger) Warn(msg string, context ...any) {
	if z.isPretty && len(context) > 0 {
		prettyMsg := z.formatPrettyMessage(msg, context)
		z.zap.Warn(prettyMsg)
	} else {
		fields := z.convertContext(context)
		z.zap.Warn(msg, fields...)
	}
}

func (z *ZapLogger) Debug(msg string, context ...any) {
	if z.isPretty && len(context) > 0 {
		prettyMsg := z.formatPrettyMessage(msg, context)
		z.zap.Debug(prettyMsg)
	} else {
		fields := z.convertContext(context)
		z.zap.Debug(msg, fields...)
	}
}

// Flush flushes any buffered log entries
// Should be called before application exits
func (z *ZapLogger) Flush() error {
	return z.zap.Sync()
}
