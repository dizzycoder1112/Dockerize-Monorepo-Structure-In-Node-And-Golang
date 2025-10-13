package rabbitmq

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

// ConsumeQueue starts consuming messages from a queue
func ConsumeQueue(
	conn *Connection,
	queue string,
	handler MessageHandler,
	options *ConsumeOptions,
) error {
	channel, err := conn.GetChannel()
	if err != nil {
		return err
	}

	logger := conn.GetLogger()

	// Use default options if not provided
	if options == nil {
		options = &ConsumeOptions{
			NoAck:     false,
			Exclusive: false,
		}
	}

	// Use default queue options if not provided
	if options.QueueOptions == nil {
		defaultQueueOpts := DefaultQueueOptions()
		options.QueueOptions = &defaultQueueOpts
	}

	// Setup DLQ if enabled
	if options.EnableDLQ {
		if err := setupDLQ(channel, queue, options.QueueOptions); err != nil {
			logger.Error("Failed to setup DLQ", map[string]interface{}{
				"error": err.Error(),
				"queue": queue,
			})
			return fmt.Errorf("failed to setup DLQ for queue %s: %w", queue, err)
		}
		logger.Info("DLQ setup completed", map[string]interface{}{
			"queue":    queue,
			"dlq":      fmt.Sprintf("%s.failed", queue),
		})
	}

	// Assert queue first (must exist before retry strategy binds it)
	_, err = channel.QueueDeclare(
		queue,
		options.QueueOptions.Durable,
		options.QueueOptions.AutoDelete,
		options.QueueOptions.Exclusive,
		options.QueueOptions.NoWait,
		options.QueueOptions.Args,
	)
	if err != nil {
		logger.Error("Failed to declare queue", map[string]interface{}{
			"error": err.Error(),
			"queue": queue,
		})
		return fmt.Errorf("failed to declare queue %s: %w", queue, err)
	}

	// Setup retry strategy after queue is declared
	if options.RetryStrategy != nil {
		if err := options.RetryStrategy.Setup(channel, queue); err != nil {
			logger.Error("Failed to setup retry strategy", map[string]interface{}{
				"error": err.Error(),
				"queue": queue,
			})
			return fmt.Errorf("failed to setup retry strategy for queue %s: %w", queue, err)
		}
	}

	// Start consuming
	msgs, err := channel.Consume(
		queue,
		options.ConsumerTag,
		options.NoAck,
		options.Exclusive,
		false, // no-local
		options.NoWait,
		options.Args,
	)
	if err != nil {
		logger.Error("Failed to start consuming", map[string]interface{}{
			"error": err.Error(),
			"queue": queue,
		})
		return fmt.Errorf("failed to start consuming queue %s: %w", queue, err)
	}

	logger.Info("Started consuming queue", map[string]interface{}{
		"queue": queue,
	})

	// Process messages
	go func() {
		for msg := range msgs {
			if err := processMessage(conn, msg, handler, options); err != nil {
				logger.Error("Error processing message", map[string]interface{}{
					"error": err.Error(),
					"queue": queue,
				})
			}
		}
	}()

	return nil
}

// processMessage handles a single message with retry logic
func processMessage(
	conn *Connection,
	delivery amqp.Delivery,
	handler MessageHandler,
	options *ConsumeOptions,
) error {
	logger := conn.GetLogger()
	channel, err := conn.GetChannel()
	if err != nil {
		return err
	}

	// Execute handler
	err = handler(delivery.Body, delivery)

	if err != nil {
		// Handler failed, check if we should retry
		if options.RetryStrategy != nil && options.RetryStrategy.ShouldRetry(delivery) {
			logger.Debug("Message failed, applying retry strategy", map[string]interface{}{
				"error": err.Error(),
			})

			// Use retry strategy to handle failure
			if retryErr := options.RetryStrategy.HandleFailure(channel, delivery); retryErr != nil {
				logger.Error("Failed to apply retry strategy", map[string]interface{}{
					"error": retryErr.Error(),
				})
				// Nack without requeue if retry strategy fails
				return delivery.Nack(false, false)
			}

			// Ack the original message (retry strategy will handle redelivery)
			return delivery.Ack(false)
		}

		// No retry strategy or retry limit exceeded, nack without requeue
		logger.Error("Message processing failed, no retry", map[string]interface{}{
			"error": err.Error(),
		})
		return delivery.Nack(false, false)
	}

	// Success, ack the message
	if !options.NoAck {
		return delivery.Ack(false)
	}

	return nil
}

// CancelConsumer cancels a consumer by its tag
func CancelConsumer(conn *Connection, consumerTag string) error {
	channel, err := conn.GetChannel()
	if err != nil {
		return err
	}

	logger := conn.GetLogger()

	if err := channel.Cancel(consumerTag, false); err != nil {
		logger.Error("Failed to cancel consumer", map[string]interface{}{
			"error":       err.Error(),
			"consumerTag": consumerTag,
		})
		return fmt.Errorf("failed to cancel consumer %s: %w", consumerTag, err)
	}

	logger.Info("Consumer cancelled", map[string]interface{}{
		"consumerTag": consumerTag,
	})

	return nil
}

// GetRetryMetadata extracts retry metadata from message headers
func GetRetryMetadata(delivery amqp.Delivery) RetryMetadata {
	metadata := RetryMetadata{
		AttemptCount: 0,
	}

	if delivery.Headers == nil {
		return metadata
	}

	if count, ok := delivery.Headers["x-retry-count"].(int32); ok {
		metadata.AttemptCount = int(count)
	}

	if queue, ok := delivery.Headers["x-original-queue"].(string); ok {
		metadata.OriginalQueue = queue
	}

	if timestamp, ok := delivery.Headers["x-first-failed-at"].(int64); ok {
		metadata.FirstFailedAt = timestamp
	}

	return metadata
}

// setupDLQ sets up Dead Letter Queue infrastructure
func setupDLQ(channel *amqp.Channel, originalQueue string, queueOptions *QueueOptions) error {
	dlxName := fmt.Sprintf("%s.failed.dlx", originalQueue)
	dlqName := fmt.Sprintf("%s.failed", originalQueue)

	// Declare DLX exchange for failed messages
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
		return fmt.Errorf("failed to declare DLX exchange: %w", err)
	}

	// Declare DLQ queue (no retry, messages stay here for manual inspection)
	_, err = channel.QueueDeclare(
		dlqName,
		true,  // durable
		false, // auto-delete
		false, // exclusive
		false, // no-wait
		nil,   // no special args - messages stay here permanently
	)
	if err != nil {
		return fmt.Errorf("failed to declare DLQ: %w", err)
	}

	// Bind DLQ to DLX
	err = channel.QueueBind(
		dlqName,
		dlqName, // routing key
		dlxName,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to bind DLQ to DLX: %w", err)
	}

	// Configure original queue options to use DLX
	if queueOptions.Args == nil {
		queueOptions.Args = amqp.Table{}
	}
	queueOptions.Args["x-dead-letter-exchange"] = dlxName
	queueOptions.Args["x-dead-letter-routing-key"] = dlqName

	return nil
}
