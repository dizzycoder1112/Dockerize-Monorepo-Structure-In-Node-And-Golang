package rabbitmq

import (
	"errors"
	"fmt"
	"net/url"
	"sync"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Connection manages RabbitMQ connection and channel
type Connection struct {
	config       Config
	logger       Logger
	conn         *amqp.Connection
	channel      *amqp.Channel
	consumerTags map[string]string
	mu           sync.RWMutex
	closed       bool
}

// NewConnection creates a new RabbitMQ connection instance
func NewConnection(config Config, logger Logger) *Connection {
	return &Connection{
		config:       config,
		logger:       logger,
		consumerTags: make(map[string]string),
		closed:       false,
	}
}

// Connect establishes connection to RabbitMQ and creates a channel
func (c *Connection) Connect() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.conn != nil && c.channel != nil {
		return nil
	}

	c.logger.Info("Connecting to RabbitMQ", map[string]interface{}{
		"url": c.maskURL(c.config.URL),
	})

	conn, err := amqp.Dial(c.config.URL)
	if err != nil {
		c.logger.Error("Failed to connect to RabbitMQ", map[string]interface{}{
			"error": err.Error(),
		})
		return fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		c.logger.Error("Failed to create channel", map[string]interface{}{
			"error": err.Error(),
		})
		return fmt.Errorf("failed to create channel: %w", err)
	}

	if c.config.Prefetch > 0 {
		if err := channel.Qos(c.config.Prefetch, 0, false); err != nil {
			channel.Close()
			conn.Close()
			c.logger.Error("Failed to set QoS", map[string]interface{}{
				"error":    err.Error(),
				"prefetch": c.config.Prefetch,
			})
			return fmt.Errorf("failed to set QoS: %w", err)
		}
	}

	c.conn = conn
	c.channel = channel

	c.setupConnectionHandlers()

	c.logger.Info("RabbitMQ connected successfully", nil)
	return nil
}

// setupConnectionHandlers sets up error and close handlers
func (c *Connection) setupConnectionHandlers() {
	go func() {
		if c.conn == nil {
			return
		}
		closeErr := <-c.conn.NotifyClose(make(chan *amqp.Error))
		if closeErr != nil {
			c.logger.Error("RabbitMQ connection error", map[string]interface{}{
				"error": closeErr.Error(),
			})
		} else {
			c.logger.Warn("RabbitMQ connection closed", nil)
		}
	}()

	go func() {
		if c.channel == nil {
			return
		}
		closeErr := <-c.channel.NotifyClose(make(chan *amqp.Error))
		if closeErr != nil {
			c.logger.Error("RabbitMQ channel error", map[string]interface{}{
				"error": closeErr.Error(),
			})
		} else {
			c.logger.Warn("RabbitMQ channel closed", nil)
		}
	}()
}

// GetChannel returns the active channel
func (c *Connection) GetChannel() (*amqp.Channel, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.channel == nil {
		return nil, errors.New("channel not initialized. Call Connect() first")
	}
	return c.channel, nil
}

// GetLogger returns the logger instance
func (c *Connection) GetLogger() Logger {
	return c.logger
}

// RegisterConsumerTag registers a consumer tag for a queue
func (c *Connection) RegisterConsumerTag(queue, consumerTag string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.consumerTags[queue] = consumerTag
}

// GetConsumerTag retrieves a consumer tag for a queue
func (c *Connection) GetConsumerTag(queue string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	tag, exists := c.consumerTags[queue]
	return tag, exists
}

// RemoveConsumerTag removes a consumer tag for a queue
func (c *Connection) RemoveConsumerTag(queue string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.consumerTags, queue)
}

// IsConnected checks if the connection and channel are active
func (c *Connection) IsConnected() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.conn != nil && c.channel != nil && !c.closed
}

// Close closes the channel and connection
func (c *Connection) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.closed {
		return nil
	}

	var errs []error

	if c.channel != nil {
		if err := c.channel.Close(); err != nil {
			c.logger.Error("Error closing channel", map[string]interface{}{
				"error": err.Error(),
			})
			errs = append(errs, err)
		}
		c.channel = nil
	}

	if c.conn != nil {
		if err := c.conn.Close(); err != nil {
			c.logger.Error("Error closing connection", map[string]interface{}{
				"error": err.Error(),
			})
			errs = append(errs, err)
		}
		c.conn = nil
	}

	c.consumerTags = make(map[string]string)
	c.closed = true

	c.logger.Info("RabbitMQ connection closed", nil)

	if len(errs) > 0 {
		return fmt.Errorf("errors during close: %v", errs)
	}

	return nil
}

// maskURL masks the password in the URL for logging
func (c *Connection) maskURL(rawURL string) string {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return "invalid-url"
	}

	if parsed.User != nil {
		if _, hasPassword := parsed.User.Password(); hasPassword {
			parsed.User = url.UserPassword(parsed.User.Username(), "***")
		}
	}

	return parsed.String()
}
