# RabbitMQ Package

A Go library for RabbitMQ messaging with built-in retry strategies, dead letter queues, and channel management.

---

## Table of Contents

- [Installation](#installation)
- [Quick Start](#quick-start)
- [Core Components](#core-components)
  - [Connection](#connection)
  - [Producer](#producer)
  - [Consumer](#consumer)
- [Retry Strategies](#retry-strategies)
  - [Immediate Retry](#1-immediate-retry)
  - [Fixed Delay Retry](#2-fixed-delay-retry)
  - [Exponential Backoff](#3-exponential-backoff)
- [Dead Letter Queue (DLQ)](#dead-letter-queue-dlq)
- [Configuration Options](#configuration-options)
- [Known Issues & Solutions](#known-issues--solutions)

---

## Installation

```go
import rabbitmqlib "go-packages/rabbitMQ"
```

---

## Quick Start

```go
package main

import (
    "log"
    rabbitmqlib "go-packages/rabbitMQ"
)

func main() {
    // 1. Create connection
    config := rabbitmqlib.Config{
        URL:      "amqp://guest:guest@localhost:5672/",
        Prefetch: 10,
    }
    conn := rabbitmqlib.NewConnection(config, nil)

    if err := conn.Connect(); err != nil {
        log.Fatal(err)
    }
    defer conn.Close()

    // 2. Publish a message
    payload := map[string]string{"message": "hello"}
    opts := rabbitmqlib.DefaultPublishOptions()

    err := rabbitmqlib.PublishToQueue(conn, "my-queue", payload, &opts)
    if err != nil {
        log.Fatal(err)
    }

    // 3. Consume messages
    handler := func(payload []byte, delivery amqp.Delivery) error {
        log.Printf("Received: %s", string(payload))
        return nil
    }

    consumeOpts := &rabbitmqlib.ConsumeOptions{
        RetryStrategy: rabbitmqlib.NewExponentialBackoff(3, 1000, 2.0),
    }

    err = rabbitmqlib.ConsumeQueue(conn, "my-queue", handler, consumeOpts)
    if err != nil {
        log.Fatal(err)
    }

    select {} // Block forever
}
```

---

## Core Components

### Connection

Manages RabbitMQ connection and channels with automatic error handling.

```go
// Create connection with custom logger
config := rabbitmqlib.Config{
    URL:      "amqp://user:pass@host:5672/vhost",
    Prefetch: 10,  // QoS prefetch count
}

conn := rabbitmqlib.NewConnection(config, customLogger)

// Connect
if err := conn.Connect(); err != nil {
    log.Fatal(err)
}

// Check connection status
if conn.IsConnected() {
    log.Println("Connected!")
}

// Get named channel (for isolation)
channel, err := conn.GetChannel("producer-channel")

// Close when done
defer conn.Close()
```

**Channel Isolation**: Use named channels to isolate different operations (e.g., separate channels for producers and consumers).

---

### Producer

Publish messages to queues or exchanges.

#### PublishToQueue

```go
payload := map[string]interface{}{
    "userId": "123",
    "action": "signup",
}

opts := rabbitmqlib.DefaultPublishOptions()
opts.Persistent = true           // Survive broker restart
opts.ChannelID = "producer"      // Use named channel

err := rabbitmqlib.PublishToQueue(conn, "user-events", payload, &opts)
```

#### PublishToExchange

```go
exchangeOpts := rabbitmqlib.DefaultExchangeOptions()
exchangeOpts.Type = "topic"

publishOpts := rabbitmqlib.DefaultPublishOptions()

err := rabbitmqlib.PublishToExchange(
    conn,
    "events",           // exchange name
    "user.created",     // routing key
    payload,
    &exchangeOpts,
    &publishOpts,
)
```

---

### Consumer

Consume messages with optional retry strategies and dead letter queues.

```go
handler := func(payload []byte, delivery amqp.Delivery) error {
    var data MyEvent
    if err := json.Unmarshal(payload, &data); err != nil {
        return err  // Will trigger retry if configured
    }

    // Process message...
    return nil  // Success - message will be ACK'd
}

opts := &rabbitmqlib.ConsumeOptions{
    ConsumerTag:   "my-consumer",
    RetryStrategy: rabbitmqlib.NewExponentialBackoff(5, 1000, 2.0),
    EnableDLQ:     true,   // Send failed messages to DLQ
    ChannelID:     "consumer",
}

err := rabbitmqlib.ConsumeQueue(conn, "my-queue", handler, opts)
```

---

## Retry Strategies

Three built-in strategies for handling failed messages.

### 1. Immediate Retry

Retries immediately by nacking with requeue. Simple but can cause tight retry loops.

```go
strategy := rabbitmqlib.NewImmediateRetry(3)  // Max 3 attempts
```

**Flow**:
```
Message fails → Nack with requeue → Immediately redelivered → Retry
```

**Use when**: Failures are transient and likely to succeed on immediate retry.

---

### 2. Fixed Delay Retry

Retries after a fixed delay using Dead Letter Exchange (DLX).

```go
strategy := rabbitmqlib.NewFixedDelayRetry(
    5,      // Max 5 attempts
    5000,   // 5 second delay between retries
)
```

**Infrastructure created**:
```
original-queue.wait     (TTL: 5000ms, DLX → original-queue.dlx)
original-queue.dlx      (Exchange, routes back to original-queue)
```

**Flow**:
```
Message fails → Publish to wait queue → Wait 5s → DLX routes back → Retry
```

---

### 3. Exponential Backoff

Retries with exponentially increasing delays. **Recommended for most use cases**.

```go
// Basic usage
strategy := rabbitmqlib.NewExponentialBackoff(
    5,      // Max 5 attempts
    1000,   // Initial delay: 1 second
    2.0,    // Multiplier: delay doubles each attempt
)

// With custom max delay
strategy := rabbitmqlib.NewExponentialBackoffWithMaxDelay(
    5,       // Max 5 attempts
    1000,    // Initial delay: 1 second
    2.0,     // Multiplier
    60000,   // Max delay: 60 seconds (cap)
)
```

#### Delay Calculation

Formula: `delay = initialDelay × multiplier^attemptCount`

| Attempt | Calculation | Delay |
|---------|-------------|-------|
| 0 | 1000 × 2^0 | 1s |
| 1 | 1000 × 2^1 | 2s |
| 2 | 1000 × 2^2 | 4s |
| 3 | 1000 × 2^3 | 8s |
| 4 | 1000 × 2^4 | 16s |

#### Infrastructure

Creates **N wait queues** (one per retry level) with different TTLs:

```
┌─────────────────────────────────────────────────────────────────┐
│                                                                 │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐          │
│  │ queue.wait.0 │  │ queue.wait.1 │  │ queue.wait.2 │   ...    │
│  │ TTL: 1000ms  │  │ TTL: 2000ms  │  │ TTL: 4000ms  │          │
│  └──────┬───────┘  └──────┬───────┘  └──────┬───────┘          │
│         │                 │                 │                   │
│         └────────────┬────┴────────────┬────┘                   │
│                      ▼                 ▼                        │
│              ┌───────────────────────────────┐                  │
│              │       queue.dlx               │                  │
│              │   (Dead Letter Exchange)      │                  │
│              └───────────────┬───────────────┘                  │
│                              │                                  │
│                              ▼                                  │
│              ┌───────────────────────────────┐                  │
│              │      original-queue           │                  │
│              │   (bound to DLX exchange)     │                  │
│              └───────────────────────────────┘                  │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

#### Complete Flow

```
1. Consumer receives message
              │
              ▼
2. Handler returns error (attempt 0)
              │
              ▼
3. ShouldRetry() → attemptCount(0) < maxAttempts(5)? → Yes
              │
              ▼
4. HandleFailure():
   - Update headers: x-retry-count = 1
   - Publish to queue.wait.0 (TTL: 1000ms)
   - ACK original message
              │
              ▼
5. Message waits in queue.wait.0 for 1 second
              │
              ▼
6. TTL expires → DLX routes message back to original-queue
              │
              ▼
7. Consumer receives message again (attempt 1)
              │
              ▼
8. Handler fails again → Publish to queue.wait.1 (TTL: 2000ms)
              │
              ▼
9. Wait 2 seconds → DLX routes back
              │
              ▼
10. Repeat until success or max attempts reached
              │
              ▼
11. If max attempts reached → Nack without requeue (or send to DLQ)
```

#### Retry Metadata

Track retry information via message headers:

```go
metadata := rabbitmqlib.GetRetryMetadata(delivery)

metadata.AttemptCount   // int: Current attempt number (0-indexed)
metadata.OriginalQueue  // string: Original queue name
metadata.FirstFailedAt  // int64: Unix timestamp of first failure
```

Headers stored in message:
- `x-retry-count`: Current attempt count
- `x-original-queue`: Original queue name
- `x-first-failed-at`: Timestamp of first failure

---

## Dead Letter Queue (DLQ)

Capture messages that fail all retry attempts for manual inspection.

```go
opts := &rabbitmqlib.ConsumeOptions{
    RetryStrategy: rabbitmqlib.NewExponentialBackoff(3, 1000, 2.0),
    EnableDLQ:     true,  // Enable DLQ
}

err := rabbitmqlib.ConsumeQueue(conn, "my-queue", handler, opts)
```

**Infrastructure created**:
```
my-queue.failed.dlx    (Dead Letter Exchange)
my-queue.failed        (DLQ - messages stay here permanently)
```

**Flow when all retries exhausted**:
```
Message fails all retries → Nack without requeue → DLX routes to my-queue.failed
```

---

## Configuration Options

### Config

```go
type Config struct {
    URL      string  // AMQP connection URL
    Prefetch int     // QoS prefetch count (0 = unlimited)
}
```

### QueueOptions

```go
type QueueOptions struct {
    Durable    bool        // Survive broker restart (default: true)
    AutoDelete bool        // Delete when unused (default: false)
    Exclusive  bool        // Exclusive to connection (default: false)
    NoWait     bool        // Don't wait for server confirmation
    Args       amqp.Table  // Additional arguments
}

// Get defaults
opts := rabbitmqlib.DefaultQueueOptions()
```

### ExchangeOptions

```go
type ExchangeOptions struct {
    Type       string      // direct, topic, fanout, headers
    Durable    bool        // Survive broker restart (default: true)
    AutoDelete bool        // Delete when unused
    Internal   bool        // Internal exchange
    NoWait     bool        // Don't wait for confirmation
    Args       amqp.Table  // Additional arguments
}

// Get defaults (topic exchange, durable)
opts := rabbitmqlib.DefaultExchangeOptions()
```

### PublishOptions

```go
type PublishOptions struct {
    Persistent         bool          // Message survives restart
    Priority           uint8         // Message priority (0-9)
    Expiration         string        // Message TTL
    Headers            amqp.Table    // Custom headers
    QueueOptions       *QueueOptions // Queue declaration options
    EnableQueueDeclare bool          // Declare queue before publish
    ChannelID          string        // Named channel for isolation
}

// Get defaults
opts := rabbitmqlib.DefaultPublishOptions()
```

### ConsumeOptions

```go
type ConsumeOptions struct {
    NoAck         bool           // Auto-ack messages
    Exclusive     bool           // Exclusive consumer
    ConsumerTag   string         // Consumer identifier
    NoWait        bool           // Don't wait for confirmation
    Args          amqp.Table     // Additional arguments
    QueueOptions  *QueueOptions  // Queue declaration options
    RetryStrategy RetryStrategy  // Retry strategy to use
    EnableDLQ     bool           // Enable Dead Letter Queue
    ChannelID     string         // Named channel for isolation
}
```

---

## Known Issues & Solutions

See [CLAUDE.md](./CLAUDE.md) for detailed documentation on:

- Silent message loss when queue doesn't exist
- Solutions: Mandatory flag, Queue verification, Publisher confirms
- Recommended patterns for reliable messaging
