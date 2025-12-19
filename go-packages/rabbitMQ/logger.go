package rabbitmq

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// simpleLogger is the default fallback logger implementation
type simpleLogger struct {
	prefix string
}

// newSimpleLogger creates a new simple logger instance
func newSimpleLogger() *simpleLogger {
	return &simpleLogger{
		prefix: "RabbitMQ",
	}
}

// parseContext converts variadic context arguments to a map
// Supports two formats:
// 1. Key-value pairs: "key1", value1, "key2", value2
// 2. Single map: map[string]any{"key1": value1, "key2": value2}
func parseContext(context []any) map[string]any {
	if len(context) == 0 {
		return nil
	}

	// Check if first argument is a map
	if len(context) == 1 {
		if m, ok := context[0].(map[string]any); ok {
			return m
		}
		// Also support map[string]interface{} for backward compatibility
		if m, ok := context[0].(map[string]interface{}); ok {
			return m
		}
	}

	// Parse as key-value pairs
	result := make(map[string]any)
	for i := 0; i < len(context); i += 2 {
		if i+1 >= len(context) {
			break
		}

		key, ok := context[i].(string)
		if !ok {
			continue
		}

		result[key] = context[i+1]
	}

	return result
}

func (l *simpleLogger) log(level string, msg string, context ...any) {
	timestamp := time.Now().Format(time.RFC3339)

	// Format: 2025-01-01T12:00:00+08:00 [INFO] [RabbitMQ] message
	fmt.Fprintf(os.Stdout, "%s [%s] [%s] %s", timestamp, level, l.prefix, msg)

	// Print context if present
	if len(context) > 0 {
		contextMap := parseContext(context)
		if len(contextMap) > 0 {
			jsonBytes, err := json.Marshal(contextMap)
			if err == nil {
				fmt.Fprintf(os.Stdout, " %s", string(jsonBytes))
			}
		}
	}

	fmt.Fprintln(os.Stdout)
}

func (l *simpleLogger) Info(msg string, context ...any) {
	l.log("INFO", msg, context...)
}

func (l *simpleLogger) Debug(msg string, context ...any) {
	l.log("DEBUG", msg, context...)
}

func (l *simpleLogger) Error(msg string, context ...any) {
	l.log("ERROR", msg, context...)
}

func (l *simpleLogger) Warn(msg string, context ...any) {
	l.log("WARN", msg, context...)
}

// defaultLogger is the fallback logger used when nil is passed to NewConnection
var defaultLogger Logger = newSimpleLogger()
