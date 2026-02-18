package logger

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// ConsoleLogger is the default logger implementation that outputs to stdout
// Similar to TypeScript's console object, it serves as the fallback logger
type ConsoleLogger struct {
	serviceName string
	colored     bool
}

// ConsoleOptions configures the console logger
type ConsoleOptions struct {
	ServiceName string
	Colored     bool // Enable colored output
}

// NewConsole creates a new console logger
func NewConsole(opts ConsoleOptions) Logger {
	return &ConsoleLogger{
		serviceName: opts.ServiceName,
		colored:     opts.Colored,
	}
}

// ANSI color codes
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorYellow = "\033[33m"
	colorCyan   = "\033[36m"
	colorGray   = "\033[90m"
)

func (c *ConsoleLogger) getColor(level string) string {
	if !c.colored {
		return ""
	}
	switch level {
	case "DEBUG":
		return colorGray
	case "INFO":
		return colorCyan
	case "WARN":
		return colorYellow
	case "ERROR":
		return colorRed
	default:
		return ""
	}
}

func (c *ConsoleLogger) log(level string, msg string, context ...any) {
	timestamp := time.Now().Format(time.RFC3339)
	color := c.getColor(level)
	reset := ""
	if c.colored {
		reset = colorReset
	}

	// Format: 2025-10-14T13:27:15+08:00 [INFO] [service-name] message
	fmt.Fprintf(os.Stdout, "%s %s[%s]%s [%s] %s",
		timestamp,
		color,
		level,
		reset,
		c.serviceName,
		msg,
	)

	// Print context if present
	if len(context) > 0 {
		contextMap := ParseContext(context)
		if len(contextMap) > 0 {
			jsonBytes, err := json.Marshal(contextMap)
			if err == nil {
				fmt.Fprintf(os.Stdout, " %s", string(jsonBytes))
			}
		}
	}

	fmt.Fprintln(os.Stdout)
}

func (c *ConsoleLogger) Info(msg string, context ...any) {
	c.log("INFO", msg, context...)
}

func (c *ConsoleLogger) Error(msg string, context ...any) {
	c.log("ERROR", msg, context...)
}

func (c *ConsoleLogger) Warn(msg string, context ...any) {
	c.log("WARN", msg, context...)
}

func (c *ConsoleLogger) Debug(msg string, context ...any) {
	c.log("DEBUG", msg, context...)
}

// Flush is a no-op for console logger (stdout is unbuffered)
func (c *ConsoleLogger) Flush() error {
	return nil
}

// Default is the global default logger (similar to TS's console)
// Packages should use this as fallback when no logger is provided
var Default Logger = NewConsole(ConsoleOptions{
	ServiceName: "default",
	Colored:     true,
})
