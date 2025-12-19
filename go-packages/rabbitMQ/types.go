package rabbitmq

import amqp "github.com/rabbitmq/amqp091-go"

// Logger interface for custom logging implementations
// Supports variadic context arguments in two formats:
// 1. Key-value pairs: "key1", value1, "key2", value2
// 2. Single map: map[string]any{"key1": value1, "key2": value2}
type Logger interface {
	Info(msg string, context ...any)
	Debug(msg string, context ...any)
	Error(msg string, context ...any)
	Warn(msg string, context ...any)
}

// Config holds RabbitMQ connection configuration
type Config struct {
	URL      string
	Prefetch int
}

// QueueOptions represents queue declaration options
type QueueOptions struct {
	Durable    bool
	AutoDelete bool
	Exclusive  bool
	NoWait     bool
	Args       amqp.Table
}

// DefaultQueueOptions returns default queue options
func DefaultQueueOptions() QueueOptions {
	return QueueOptions{
		Durable:    true,
		AutoDelete: false,
		Exclusive:  false,
		NoWait:     false,
		Args:       nil,
	}
}

// ExchangeOptions represents exchange declaration options
type ExchangeOptions struct {
	Type       string // direct, topic, fanout, headers
	Durable    bool
	AutoDelete bool
	Internal   bool
	NoWait     bool
	Args       amqp.Table
}

// DefaultExchangeOptions returns default exchange options
// Default: topic exchange, durable
func DefaultExchangeOptions() ExchangeOptions {
	return ExchangeOptions{
		Type:       "topic",
		Durable:    true,
		AutoDelete: false,
		Internal:   false,
		NoWait:     false,
		Args:       nil,
	}
}

// PublishOptions represents message publishing options
type PublishOptions struct {
	Persistent         bool
	Priority           uint8
	Expiration         string
	Headers            amqp.Table
	QueueOptions       *QueueOptions
	EnableQueueDeclare bool   // Enable queue declaration (default: false, assume queue already exists)
	ChannelID          string // Optional channel ID for channel isolation. Empty string uses default channel.
}

// DefaultPublishOptions returns default publish options
// By default, queue declaration is disabled (assume queue already exists)
func DefaultPublishOptions() PublishOptions {
	return PublishOptions{
		Persistent:         true,
		Priority:           0,
		Expiration:         "",
		Headers:            nil,
		QueueOptions:       nil,
		EnableQueueDeclare: false, // Default: don't declare, assume queue exists
	}
}

// ConsumeOptions represents consumer configuration options
type ConsumeOptions struct {
	NoAck         bool
	Exclusive     bool
	ConsumerTag   string
	NoWait        bool
	Args          amqp.Table
	QueueOptions  *QueueOptions
	RetryStrategy RetryStrategy
	EnableDLQ     bool   // Enable Dead Letter Queue for failed messages
	ChannelID     string // Optional channel ID for channel isolation. Empty string uses default channel.
}

// MessageHandler is a function type for handling consumed messages
type MessageHandler func(payload []byte, delivery amqp.Delivery) error

// RetryStrategy defines the interface for retry strategies
type RetryStrategy interface {
	// ShouldRetry determines if a message should be retried based on the delivery
	ShouldRetry(delivery amqp.Delivery) bool

	// GetDelay returns the delay in milliseconds before retry
	GetDelay(attemptCount int) int

	// Setup configures the necessary queues and exchanges for retry mechanism
	Setup(channel *amqp.Channel, originalQueue string) error

	// HandleFailure handles a failed message according to the strategy
	HandleFailure(channel *amqp.Channel, delivery amqp.Delivery) error
}

// RetryMetadata holds retry-related metadata from message headers
type RetryMetadata struct {
	AttemptCount int
	OriginalQueue string
	FirstFailedAt int64
}
