package rabbitmq

import (
	"errors"
	"fmt"
	"net/url"
	"sync"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Connection manages RabbitMQ connection and channel(s)
type Connection struct {
	config         Config
	logger         Logger
	conn           *amqp.Connection
	defaultChannel *amqp.Channel
	channels       map[string]*amqp.Channel // Named channels for isolation
	consumerTags   map[string]string
	mu             sync.RWMutex
	closed         bool
}

// NewConnection creates a new RabbitMQ connection instance
// If logger is nil, a default simple logger will be used
func NewConnection(config Config, logger Logger) *Connection {
	if logger == nil {
		logger = defaultLogger
	}
	return &Connection{
		config:       config,
		logger:       logger,
		channels:     make(map[string]*amqp.Channel),
		consumerTags: make(map[string]string),
		closed:       false,
	}
}

// Connect establishes connection to RabbitMQ and creates a default channel
func (c *Connection) Connect() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.conn != nil && c.defaultChannel != nil {
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
		c.logger.Error("Failed to create default channel", map[string]interface{}{
			"error": err.Error(),
		})
		return fmt.Errorf("failed to create default channel: %w", err)
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
	c.defaultChannel = channel

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

	c.setupChannelHandlers(c.defaultChannel, "default")
}

// setupChannelHandlers sets up error and close handlers for a specific channel
func (c *Connection) setupChannelHandlers(channel *amqp.Channel, channelID string) {
	go func() {
		if channel == nil {
			return
		}
		closeErr := <-channel.NotifyClose(make(chan *amqp.Error))
		if closeErr != nil {
			c.logger.Error("RabbitMQ channel error", map[string]interface{}{
				"error":     closeErr.Error(),
				"channelId": channelID,
			})
		} else {
			c.logger.Warn("RabbitMQ channel closed", map[string]interface{}{
				"channelId": channelID,
			})
		}

		// Remove from map if it's a named channel
		if channelID != "default" {
			c.mu.Lock()
			delete(c.channels, channelID)
			c.mu.Unlock()
		}
	}()
}

// GetChannel returns a channel by ID
// If channelID is empty, returns the default channel
// If channelID is specified and doesn't exist, creates a new named channel
func (c *Connection) GetChannel(channelID string) (*amqp.Channel, error) {
	// Return default channel if no channelID specified
	if channelID == "" {
		c.mu.RLock()
		defer c.mu.RUnlock()

		if c.defaultChannel == nil {
			return nil, errors.New("default channel not initialized. Call Connect() first")
		}
		return c.defaultChannel, nil
	}

	// Check if named channel already exists
	c.mu.RLock()
	if channel, exists := c.channels[channelID]; exists {
		c.mu.RUnlock()
		return channel, nil
	}
	c.mu.RUnlock()

	// Create new named channel
	c.mu.Lock()
	defer c.mu.Unlock()

	// Double-check after acquiring write lock
	if channel, exists := c.channels[channelID]; exists {
		return channel, nil
	}

	if c.conn == nil {
		return nil, errors.New("connection not initialized. Call Connect() first")
	}

	c.logger.Info("Creating new named channel", map[string]interface{}{
		"channelId": channelID,
	})

	channel, err := c.conn.Channel()
	if err != nil {
		c.logger.Error("Failed to create named channel", map[string]interface{}{
			"error":     err.Error(),
			"channelId": channelID,
		})
		return nil, fmt.Errorf("failed to create named channel %s: %w", channelID, err)
	}

	if c.config.Prefetch > 0 {
		if err := channel.Qos(c.config.Prefetch, 0, false); err != nil {
			channel.Close()
			c.logger.Error("Failed to set QoS on named channel", map[string]interface{}{
				"error":     err.Error(),
				"channelId": channelID,
				"prefetch":  c.config.Prefetch,
			})
			return nil, fmt.Errorf("failed to set QoS on channel %s: %w", channelID, err)
		}
	}

	c.setupChannelHandlers(channel, channelID)
	c.channels[channelID] = channel

	c.logger.Info("Named channel created successfully", map[string]interface{}{
		"channelId": channelID,
	})

	return channel, nil
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

// IsConnected checks if the connection and default channel are active
func (c *Connection) IsConnected() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.conn != nil && c.defaultChannel != nil && !c.closed
}

// Close closes all channels and connection
func (c *Connection) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.closed {
		return nil
	}

	var errs []error

	// Close all named channels
	for channelID, channel := range c.channels {
		if err := channel.Close(); err != nil {
			c.logger.Error("Error closing named channel", map[string]interface{}{
				"error":     err.Error(),
				"channelId": channelID,
			})
			errs = append(errs, err)
		} else {
			c.logger.Debug("Named channel closed", map[string]interface{}{
				"channelId": channelID,
			})
		}
	}
	c.channels = make(map[string]*amqp.Channel)

	// Close default channel
	if c.defaultChannel != nil {
		if err := c.defaultChannel.Close(); err != nil {
			c.logger.Error("Error closing default channel", map[string]interface{}{
				"error": err.Error(),
			})
			errs = append(errs, err)
		}
		c.defaultChannel = nil
	}

	// Close connection
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
