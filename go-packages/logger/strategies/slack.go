package strategies

import (
	"bytes"
	"dizzycoder1112/Dockerize-Monorepo-Structure-In-Node-And-Golang/logger"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"
)

// ============================================================
// Fault Tolerance Interface
// ============================================================

// FaultToleranceState represents the current state of fault tolerance mechanism
type FaultToleranceState struct {
	State    string
	Failures int
}

// FaultToleranceStrategy defines the interface for fault tolerance mechanisms
// Implement this interface to create custom fault tolerance (e.g., Circuit Breaker, Rate Limiter)
//
// Example:
//
//	type CircuitBreaker struct {
//	    state    string
//	    failures int
//	}
//
//	func (cb *CircuitBreaker) CanExecute() bool { return cb.state != "OPEN" }
//	func (cb *CircuitBreaker) OnSuccess()       { cb.failures = 0 }
//	func (cb *CircuitBreaker) OnFailure()       { cb.failures++; if cb.failures >= 5 { cb.state = "OPEN" } }
//	func (cb *CircuitBreaker) GetState() FaultToleranceState { return FaultToleranceState{cb.state, cb.failures} }
type FaultToleranceStrategy interface {
	// CanExecute checks if the operation can be executed
	CanExecute() bool
	// OnSuccess is called when the operation succeeds
	OnSuccess()
	// OnFailure is called when the operation fails
	OnFailure()
	// GetState returns current state info (for debugging/monitoring)
	GetState() FaultToleranceState
}

// ============================================================
// Slack Strategy
// ============================================================

// SlackOptions configures the Slack logger strategy
type SlackOptions struct {
	WebhookURL     string
	ServiceName    string
	Environment    string
	FaultTolerance FaultToleranceStrategy // Optional: circuit breaker or rate limiter
}

// SlackStrategy sends error and warning logs to Slack
// Only logs with level "error" or "warn" are sent to Slack
type SlackStrategy struct {
	webhookURL     string
	serviceName    string
	environment    string
	faultTolerance FaultToleranceStrategy

	// Pending request tracking for graceful shutdown
	wg sync.WaitGroup
}

// NewSlack creates a new Slack logger strategy
func NewSlack(opts SlackOptions) logger.Logger {
	return &SlackStrategy{
		webhookURL:     opts.WebhookURL,
		serviceName:    opts.ServiceName,
		environment:    opts.Environment,
		faultTolerance: opts.FaultTolerance,
	}
}

// slackAttachment represents a Slack message attachment
type slackAttachment struct {
	Color  string       `json:"color"`
	Title  string       `json:"title"`
	Text   string       `json:"text,omitempty"`
	Fields []slackField `json:"fields,omitempty"`
	Ts     int64        `json:"ts"`
}

// slackField represents a field in Slack attachment
type slackField struct {
	Title string `json:"title"`
	Value string `json:"value"`
	Short bool   `json:"short"`
}

// slackPayload represents the Slack webhook payload
type slackPayload struct {
	Attachments []slackAttachment `json:"attachments"`
}

// Color constants for Slack messages
const (
	colorError = "#dc3545" // red
	colorWarn  = "#ffc107" // yellow
)

// Emoji constants for Slack messages
const (
	emojiError = "🔴"
	emojiWarn  = "🟡"
)

func (s *SlackStrategy) Info(msg string, context ...any) {
	// Slack only handles error and warn levels
}

func (s *SlackStrategy) Debug(msg string, context ...any) {
	// Slack only handles error and warn levels
}

func (s *SlackStrategy) Warn(msg string, context ...any) {
	s.sendToSlack("warn", msg, context)
}

func (s *SlackStrategy) Error(msg string, context ...any) {
	s.sendToSlack("error", msg, context)
}

// Flush waits for all pending Slack requests to complete
func (s *SlackStrategy) Flush() error {
	s.wg.Wait()
	return nil
}

func (s *SlackStrategy) sendToSlack(level string, msg string, context []any) {
	// Skip if webhook URL is not configured
	if s.webhookURL == "" {
		return
	}

	// Check fault tolerance before sending (if provided)
	if s.faultTolerance != nil && !s.faultTolerance.CanExecute() {
		return
	}

	// Track pending request
	s.wg.Add(1)

	go func() {
		defer s.wg.Done()

		attachment := s.buildAttachment(level, msg, context)
		payload := slackPayload{Attachments: []slackAttachment{attachment}}

		jsonBytes, err := json.Marshal(payload)
		if err != nil {
			fmt.Fprintf(os.Stderr, "[SlackStrategy] Failed to marshal payload: %v\n", err)
			if s.faultTolerance != nil {
				s.faultTolerance.OnFailure()
			}
			return
		}

		resp, err := http.Post(s.webhookURL, "application/json", bytes.NewBuffer(jsonBytes))
		if err != nil {
			fmt.Fprintf(os.Stderr, "[SlackStrategy] Failed to send message: %v\n", err)
			if s.faultTolerance != nil {
				s.faultTolerance.OnFailure()
			}
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			if s.faultTolerance != nil {
				s.faultTolerance.OnSuccess()
			}
		} else {
			fmt.Fprintf(os.Stderr, "[SlackStrategy] HTTP error: %d\n", resp.StatusCode)
			if s.faultTolerance != nil {
				s.faultTolerance.OnFailure()
			}
		}
	}()
}

func (s *SlackStrategy) buildAttachment(level string, msg string, context []any) slackAttachment {
	var color, emoji string
	if level == "error" {
		color = colorError
		emoji = emojiError
	} else {
		color = colorWarn
		emoji = emojiWarn
	}

	title := fmt.Sprintf("%s %s - %s", emoji, stringToUpper(level), s.serviceName)

	fields := make([]slackField, 0)

	// Add environment if provided
	if s.environment != "" {
		fields = append(fields, slackField{
			Title: "Environment",
			Value: s.environment,
			Short: true,
		})
	}

	// Add context fields
	contextMap := logger.ParseContext(context)
	for key, value := range contextMap {
		// Handle error specially - extract stack trace if available
		if key == "error" {
			if err, ok := value.(error); ok {
				fields = append(fields, slackField{
					Title: "Error",
					Value: truncateString(err.Error(), 500),
					Short: false,
				})
				continue
			}
		}

		fields = append(fields, slackField{
			Title: key,
			Value: formatValue(value),
			Short: true,
		})
	}

	return slackAttachment{
		Color:  color,
		Title:  title,
		Text:   msg,
		Fields: fields,
		Ts:     time.Now().Unix(),
	}
}

// stringToUpper converts string to uppercase (simple implementation for level)
func stringToUpper(s string) string {
	switch s {
	case "error":
		return "ERROR"
	case "warn":
		return "WARN"
	default:
		return s
	}
}

// formatValue converts any value to string for Slack display
func formatValue(value any) string {
	switch v := value.(type) {
	case string:
		return v
	case int, int64, int32, float64, float32, bool:
		return fmt.Sprintf("%v", v)
	default:
		jsonBytes, err := json.Marshal(v)
		if err != nil {
			return fmt.Sprintf("%v", v)
		}
		return string(jsonBytes)
	}
}

// truncateString truncates a string to maxLength and adds "... (truncated)" if needed
func truncateString(s string, maxLength int) string {
	if len(s) <= maxLength {
		return s
	}
	return s[:maxLength] + "\n... (truncated)"
}
