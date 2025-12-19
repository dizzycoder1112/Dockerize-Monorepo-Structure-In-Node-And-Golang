package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

// PublishToExchange publishes a message to an exchange with routing key
// For topic/fanout/direct exchanges
func PublishToExchange(
	conn *Connection,
	exchange string,
	routingKey string,
	payload interface{},
	exchangeOptions *ExchangeOptions,
	publishOptions *PublishOptions,
) error {
	// Use default options if not provided
	if publishOptions == nil {
		defaultPublishOpts := DefaultPublishOptions()
		publishOptions = &defaultPublishOpts
	}

	channel, err := conn.GetChannel(publishOptions.ChannelID)
	if err != nil {
		return err
	}

	logger := conn.GetLogger()

	// Use default exchange options if not provided
	if exchangeOptions == nil {
		defaultExchangeOpts := DefaultExchangeOptions()
		exchangeOptions = &defaultExchangeOpts
	}

	// Use default publish options if not provided
	if publishOptions == nil {
		defaultPublishOpts := DefaultPublishOptions()
		publishOptions = &defaultPublishOpts
	}

	// Ensure exchange exists
	err = channel.ExchangeDeclare(
		exchange,
		exchangeOptions.Type,
		exchangeOptions.Durable,
		exchangeOptions.AutoDelete,
		exchangeOptions.Internal,
		exchangeOptions.NoWait,
		exchangeOptions.Args,
	)
	if err != nil {
		logger.Error("Failed to declare exchange", map[string]interface{}{
			"error":    err.Error(),
			"exchange": exchange,
			"type":     exchangeOptions.Type,
		})
		return fmt.Errorf("failed to declare exchange %s: %w", exchange, err)
	}

	// Marshal payload to JSON
	message, err := json.Marshal(payload)
	if err != nil {
		logger.Error("Failed to marshal payload", map[string]interface{}{
			"error":      err.Error(),
			"exchange":   exchange,
			"routingKey": routingKey,
		})
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	// Prepare publishing options
	publishing := amqp.Publishing{
		ContentType:  "application/json",
		Body:         message,
		DeliveryMode: amqp.Transient,
		Priority:     publishOptions.Priority,
		Headers:      publishOptions.Headers,
	}

	if publishOptions.Persistent {
		publishing.DeliveryMode = amqp.Persistent
	}

	if publishOptions.Expiration != "" {
		publishing.Expiration = publishOptions.Expiration
	}

	// Publish message to exchange
	err = channel.PublishWithContext(
		context.Background(),
		exchange,   // exchange
		routingKey, // routing key
		false,      // mandatory
		false,      // immediate
		publishing,
	)

	if err != nil {
		logger.Error("Failed to publish message to exchange", map[string]interface{}{
			"error":      err.Error(),
			"exchange":   exchange,
			"routingKey": routingKey,
		})
		return fmt.Errorf("failed to publish message to exchange %s: %w", exchange, err)
	}

	channelID := "default"
	if publishOptions.ChannelID != "" {
		channelID = publishOptions.ChannelID
	}

	logger.Debug("Message published to exchange", map[string]interface{}{
		"exchange":    exchange,
		"routingKey":  routingKey,
		"payloadSize": len(message),
		"channelId":   channelID,
	})

	return nil
}

// PublishToQueue publishes a message to a queue
func PublishToQueue(
	conn *Connection,
	queue string,
	payload interface{},
	options *PublishOptions,
) error {
	// Use default options if not provided
	if options == nil {
		defaultOpts := DefaultPublishOptions()
		options = &defaultOpts
	}

	channel, err := conn.GetChannel(options.ChannelID)
	if err != nil {
		return err
	}

	logger := conn.GetLogger()

	// Use default options if not provided
	if options == nil {
		defaultOpts := DefaultPublishOptions()
		options = &defaultOpts
	}

	// Only declare queue if explicitly enabled
	if options.EnableQueueDeclare {
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
		channelID := "default"
		if options.ChannelID != "" {
			channelID = options.ChannelID
		}
		logger.Error("Failed to publish message to queue", map[string]interface{}{
			"error":     err.Error(),
			"queue":     queue,
			"channelId": channelID,
		})
		return fmt.Errorf("failed to publish message to queue %s: %w", queue, err)
	}

	channelID := "default"
	if options.ChannelID != "" {
		channelID = options.ChannelID
	}

	logger.Debug("Message published to queue", map[string]interface{}{
		"queue":       queue,
		"payloadSize": len(message),
		"channelId":   channelID,
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
	// Use default options if not provided
	if options == nil {
		defaultOpts := DefaultPublishOptions()
		options = &defaultOpts
	}

	channel, err := conn.GetChannel(options.ChannelID)
	if err != nil {
		return err
	}

	logger := conn.GetLogger()

	// Use default options if not provided
	if options == nil {
		defaultOpts := DefaultPublishOptions()
		options = &defaultOpts
	}

	// Only declare queue if explicitly enabled
	if options.EnableQueueDeclare {
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
		channelID := "default"
		if options.ChannelID != "" {
			channelID = options.ChannelID
		}
		logger.Error("Failed to publish raw message to queue", map[string]interface{}{
			"error":     err.Error(),
			"queue":     queue,
			"channelId": channelID,
		})
		return fmt.Errorf("failed to publish raw message to queue %s: %w", queue, err)
	}

	channelID := "default"
	if options.ChannelID != "" {
		channelID = options.ChannelID
	}

	logger.Debug("Raw message published to queue", map[string]interface{}{
		"queue":       queue,
		"payloadSize": len(message),
		"channelId":   channelID,
	})

	return nil
}
