package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

// PublishToQueue publishes a message to a queue
func PublishToQueue(
	conn *Connection,
	queue string,
	payload interface{},
	options *PublishOptions,
) error {
	channel, err := conn.GetChannel()
	if err != nil {
		return err
	}

	logger := conn.GetLogger()

	// Use default options if not provided
	if options == nil {
		defaultOpts := DefaultPublishOptions()
		options = &defaultOpts
	}

	// Use default queue options if not provided
	if options.QueueOptions == nil {
		defaultQueueOpts := DefaultQueueOptions()
		options.QueueOptions = &defaultQueueOpts
	}

	// Assert queue
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

	// Marshal payload to JSON
	message, err := json.Marshal(payload)
	if err != nil {
		logger.Error("Failed to marshal payload", map[string]interface{}{
			"error": err.Error(),
			"queue": queue,
		})
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	// Prepare publishing options
	publishing := amqp.Publishing{
		ContentType:  "application/json",
		Body:         message,
		DeliveryMode: amqp.Transient,
		Priority:     options.Priority,
		Headers:      options.Headers,
	}

	if options.Persistent {
		publishing.DeliveryMode = amqp.Persistent
	}

	if options.Expiration != "" {
		publishing.Expiration = options.Expiration
	}

	// Publish message
	err = channel.PublishWithContext(
		context.Background(),
		"",    // exchange
		queue, // routing key
		false, // mandatory
		false, // immediate
		publishing,
	)

	if err != nil {
		logger.Error("Failed to publish message to queue", map[string]interface{}{
			"error": err.Error(),
			"queue": queue,
		})
		return fmt.Errorf("failed to publish message to queue %s: %w", queue, err)
	}

	logger.Debug("Message published to queue", map[string]interface{}{
		"queue":       queue,
		"payloadSize": len(message),
	})

	return nil
}

// PublishToQueueRaw publishes raw bytes to a queue without JSON marshaling
func PublishToQueueRaw(
	conn *Connection,
	queue string,
	message []byte,
	options *PublishOptions,
) error {
	channel, err := conn.GetChannel()
	if err != nil {
		return err
	}

	logger := conn.GetLogger()

	// Use default options if not provided
	if options == nil {
		defaultOpts := DefaultPublishOptions()
		options = &defaultOpts
	}

	// Use default queue options if not provided
	if options.QueueOptions == nil {
		defaultQueueOpts := DefaultQueueOptions()
		options.QueueOptions = &defaultQueueOpts
	}

	// Assert queue
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

	// Prepare publishing options
	publishing := amqp.Publishing{
		ContentType:  "application/octet-stream",
		Body:         message,
		DeliveryMode: amqp.Transient,
		Priority:     options.Priority,
		Headers:      options.Headers,
	}

	if options.Persistent {
		publishing.DeliveryMode = amqp.Persistent
	}

	if options.Expiration != "" {
		publishing.Expiration = options.Expiration
	}

	// Publish message
	err = channel.PublishWithContext(
		context.Background(),
		"",    // exchange
		queue, // routing key
		false, // mandatory
		false, // immediate
		publishing,
	)

	if err != nil {
		logger.Error("Failed to publish raw message to queue", map[string]interface{}{
			"error": err.Error(),
			"queue": queue,
		})
		return fmt.Errorf("failed to publish raw message to queue %s: %w", queue, err)
	}

	logger.Debug("Raw message published to queue", map[string]interface{}{
		"queue":       queue,
		"payloadSize": len(message),
	})

	return nil
}
