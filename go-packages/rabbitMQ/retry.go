package rabbitmq

import (
	"fmt"
	"math"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// ImmediateRetryStrategy retries immediately using Nack with requeue
type ImmediateRetryStrategy struct {
	MaxAttempts int
}

// NewImmediateRetry creates a new immediate retry strategy
func NewImmediateRetry(maxAttempts int) *ImmediateRetryStrategy {
	return &ImmediateRetryStrategy{
		MaxAttempts: maxAttempts,
	}
}

func (s *ImmediateRetryStrategy) ShouldRetry(delivery amqp.Delivery) bool {
	metadata := GetRetryMetadata(delivery)
	return metadata.AttemptCount < s.MaxAttempts
}

func (s *ImmediateRetryStrategy) GetDelay(attemptCount int) int {
	return 0 // No delay for immediate retry
}

func (s *ImmediateRetryStrategy) Setup(channel *amqp.Channel, originalQueue string) error {
	// No setup needed for immediate retry
	return nil
}

func (s *ImmediateRetryStrategy) HandleFailure(channel *amqp.Channel, delivery amqp.Delivery) error {
	metadata := GetRetryMetadata(delivery)

	// Update retry count
	if delivery.Headers == nil {
		delivery.Headers = amqp.Table{}
	}
	delivery.Headers["x-retry-count"] = int32(metadata.AttemptCount + 1)
	delivery.Headers["x-original-queue"] = delivery.RoutingKey

	if metadata.FirstFailedAt == 0 {
		delivery.Headers["x-first-failed-at"] = time.Now().Unix()
	}

	// Nack with requeue for immediate retry
	return delivery.Nack(false, true)
}

// FixedDelayRetryStrategy retries with a fixed delay using DLX
type FixedDelayRetryStrategy struct {
	MaxAttempts int
	DelayMs     int
}

// NewFixedDelayRetry creates a new fixed delay retry strategy
func NewFixedDelayRetry(maxAttempts int, delayMs int) *FixedDelayRetryStrategy {
	return &FixedDelayRetryStrategy{
		MaxAttempts: maxAttempts,
		DelayMs:     delayMs,
	}
}

func (s *FixedDelayRetryStrategy) ShouldRetry(delivery amqp.Delivery) bool {
	metadata := GetRetryMetadata(delivery)
	return metadata.AttemptCount < s.MaxAttempts
}

func (s *FixedDelayRetryStrategy) GetDelay(attemptCount int) int {
	return s.DelayMs
}

func (s *FixedDelayRetryStrategy) Setup(channel *amqp.Channel, originalQueue string) error {
	waitQueueName := fmt.Sprintf("%s.wait", originalQueue)
	dlxName := fmt.Sprintf("%s.dlx", originalQueue)

	// Declare DLX exchange
	err := channel.ExchangeDeclare(
		dlxName,
		"direct",
		true,  // durable
		false, // auto-delete
		false, // internal
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		return fmt.Errorf("failed to declare DLX: %w", err)
	}

	// Declare wait queue with DLX and TTL
	_, err = channel.QueueDeclare(
		waitQueueName,
		true,  // durable
		false, // auto-delete
		false, // exclusive
		false, // no-wait
		amqp.Table{
			"x-dead-letter-exchange":    dlxName,
			"x-dead-letter-routing-key": originalQueue,
			"x-message-ttl":             int32(s.DelayMs),
		},
	)
	if err != nil {
		return fmt.Errorf("failed to declare wait queue: %w", err)
	}

	// Bind original queue to DLX
	err = channel.QueueBind(
		originalQueue,
		originalQueue,
		dlxName,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to bind original queue to DLX: %w", err)
	}

	return nil
}

func (s *FixedDelayRetryStrategy) HandleFailure(channel *amqp.Channel, delivery amqp.Delivery) error {
	metadata := GetRetryMetadata(delivery)
	waitQueueName := fmt.Sprintf("%s.wait", delivery.RoutingKey)

	// Update headers
	if delivery.Headers == nil {
		delivery.Headers = amqp.Table{}
	}
	delivery.Headers["x-retry-count"] = int32(metadata.AttemptCount + 1)
	delivery.Headers["x-original-queue"] = delivery.RoutingKey

	if metadata.FirstFailedAt == 0 {
		delivery.Headers["x-first-failed-at"] = time.Now().Unix()
	}

	// Publish to wait queue
	err := channel.Publish(
		"",            // exchange
		waitQueueName, // routing key
		false,         // mandatory
		false,         // immediate
		amqp.Publishing{
			ContentType:  delivery.ContentType,
			Body:         delivery.Body,
			DeliveryMode: delivery.DeliveryMode,
			Priority:     delivery.Priority,
			Headers:      delivery.Headers,
		},
	)

	if err != nil {
		return fmt.Errorf("failed to publish to wait queue: %w", err)
	}

	return nil
}

// ExponentialBackoffStrategy retries with exponential backoff using DLX
type ExponentialBackoffStrategy struct {
	MaxAttempts    int
	InitialDelayMs int
	Multiplier     float64
	MaxDelayMs     int
}

// NewExponentialBackoff creates a new exponential backoff retry strategy
func NewExponentialBackoff(maxAttempts int, initialDelayMs int, multiplier float64) *ExponentialBackoffStrategy {
	return &ExponentialBackoffStrategy{
		MaxAttempts:    maxAttempts,
		InitialDelayMs: initialDelayMs,
		Multiplier:     multiplier,
		MaxDelayMs:     300000, // Default max 5 minutes
	}
}

// NewExponentialBackoffWithMaxDelay creates a new exponential backoff retry strategy with custom max delay
func NewExponentialBackoffWithMaxDelay(maxAttempts int, initialDelayMs int, multiplier float64, maxDelayMs int) *ExponentialBackoffStrategy {
	return &ExponentialBackoffStrategy{
		MaxAttempts:    maxAttempts,
		InitialDelayMs: initialDelayMs,
		Multiplier:     multiplier,
		MaxDelayMs:     maxDelayMs,
	}
}

func (s *ExponentialBackoffStrategy) ShouldRetry(delivery amqp.Delivery) bool {
	metadata := GetRetryMetadata(delivery)
	return metadata.AttemptCount < s.MaxAttempts
}

func (s *ExponentialBackoffStrategy) GetDelay(attemptCount int) int {
	delay := float64(s.InitialDelayMs) * math.Pow(s.Multiplier, float64(attemptCount))

	if int(delay) > s.MaxDelayMs {
		return s.MaxDelayMs
	}

	return int(delay)
}

func (s *ExponentialBackoffStrategy) Setup(channel *amqp.Channel, originalQueue string) error {
	dlxName := fmt.Sprintf("%s.dlx", originalQueue)

	// Declare DLX exchange
	err := channel.ExchangeDeclare(
		dlxName,
		"direct",
		true,  // durable
		false, // auto-delete
		false, // internal
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		return fmt.Errorf("failed to declare DLX: %w", err)
	}

	// Create wait queues for each retry attempt with different TTLs
	for i := 0; i < s.MaxAttempts; i++ {
		waitQueueName := fmt.Sprintf("%s.wait.%d", originalQueue, i)
		ttl := s.GetDelay(i)

		_, err = channel.QueueDeclare(
			waitQueueName,
			true,  // durable
			false, // auto-delete
			false, // exclusive
			false, // no-wait
			amqp.Table{
				"x-dead-letter-exchange":    dlxName,
				"x-dead-letter-routing-key": originalQueue,
				"x-message-ttl":             int32(ttl),
			},
		)
		if err != nil {
			return fmt.Errorf("failed to declare wait queue %s: %w", waitQueueName, err)
		}
	}

	// Bind original queue to DLX
	err = channel.QueueBind(
		originalQueue,
		originalQueue,
		dlxName,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to bind original queue to DLX: %w", err)
	}

	return nil
}

func (s *ExponentialBackoffStrategy) HandleFailure(channel *amqp.Channel, delivery amqp.Delivery) error {
	metadata := GetRetryMetadata(delivery)
	waitQueueName := fmt.Sprintf("%s.wait.%d", delivery.RoutingKey, metadata.AttemptCount)

	// Update headers
	if delivery.Headers == nil {
		delivery.Headers = amqp.Table{}
	}
	delivery.Headers["x-retry-count"] = int32(metadata.AttemptCount + 1)
	delivery.Headers["x-original-queue"] = delivery.RoutingKey

	if metadata.FirstFailedAt == 0 {
		delivery.Headers["x-first-failed-at"] = time.Now().Unix()
	}

	// Publish to wait queue
	err := channel.Publish(
		"",            // exchange
		waitQueueName, // routing key
		false,         // mandatory
		false,         // immediate
		amqp.Publishing{
			ContentType:  delivery.ContentType,
			Body:         delivery.Body,
			DeliveryMode: delivery.DeliveryMode,
			Priority:     delivery.Priority,
			Headers:      delivery.Headers,
		},
	)

	if err != nil {
		return fmt.Errorf("failed to publish to wait queue: %w", err)
	}

	return nil
}
