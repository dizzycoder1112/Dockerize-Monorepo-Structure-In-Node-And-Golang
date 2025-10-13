package rabbitmq

import amqp "github.com/rabbitmq/amqp091-go"

// Logger interface for custom logging implementations
type Logger interface {
	Info(msg string, fields map[string]interface{})
	Debug(msg string, fields map[string]interface{})
	Error(msg string, fields map[string]interface{})
	Warn(msg string, fields map[string]interface{})
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

// PublishOptions represents message publishing options
type PublishOptions struct {
	Persistent   bool
	Priority     uint8
	Expiration   string
	Headers      amqp.Table
	QueueOptions *QueueOptions
}

// DefaultPublishOptions returns default publish options
func DefaultPublishOptions() PublishOptions {
	queueOpts := DefaultQueueOptions()
	return PublishOptions{
		Persistent:   true,
		Priority:     0,
		Expiration:   "",
		Headers:      nil,
		QueueOptions: &queueOpts,
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
	EnableDLQ     bool // Enable Dead Letter Queue for failed messages
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
