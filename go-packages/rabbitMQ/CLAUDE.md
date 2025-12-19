# RabbitMQ Package - Known Issues & Solutions

## Issue: Silent Message Loss in Producer

### Problem Description

When using `PublishToQueue()` with default settings, messages may be **silently discarded** if the target queue does not exist.

**Root Cause**:
- Default `mandatory = false` in RabbitMQ publish
- If queue doesn't exist or message cannot be routed, RabbitMQ silently discards the message
- No error is returned to the publisher

**When This Happens**:
1. Queue is deleted after producer starts
2. Queue name is misconfigured
3. Consumer hasn't started yet (queue not created)
4. Network partition causes queue to be unavailable

---

## Solutions

### Solution 1: Enable Mandatory Flag (Recommended for Most Cases)

**What**: Set `mandatory = true` to make RabbitMQ return unroutable messages

**Implementation**:
```go
// In app layer (e.g., apps/go-promotion-worker)
publishOptions := rabbitmqlib.DefaultPublishOptions()
publishOptions.Mandatory = true  // Enable mandatory flag

err := rabbitmqlib.PublishToQueue(conn, queue, payload, &publishOptions)
if err != nil {
    // Handle error: queue doesn't exist or connection broken
}
```

**Pros**:
- ✅ Detects queue deletion immediately
- ✅ Simple to implement (just set a flag)
- ✅ Works for connection failures too

**Cons**:
- ⚠️ For full handling, need to implement return handler (complex)
- ⚠️ Basic error from `PublishWithContext` may not give detailed reason

**Use When**:
- You need to know if messages are unroutable
- Queue should always exist during runtime

---

### Solution 2: Queue Verification at Initialization (Fail-Fast)

**What**: Check if queue exists when producer starts

**Implementation**:
```go
// In app layer producer initialization
func (p *NotificationProducer) VerifyQueue() error {
    channel, _ := p.rabbit.GetConnection().GetChannel()

    // Passive declare: check existence without creating
    _, err := channel.QueueDeclarePassive(
        queueName,
        true,  // durable
        false, // auto-delete
        false, // exclusive
        false, // no-wait
        nil,   // args
    )
    return err
}

// In Factory
func NewFactory(args FactoryArgs) *Factory {
    producer := NewNotificationProducer(...)

    if err := producer.VerifyQueue(); err != nil {
        log.Fatalf("Queue verification failed: %v", err)
    }

    return &Factory{Notification: producer}
}
```

**Pros**:
- ✅ Fail-fast: know immediately on startup
- ✅ Prevents silent failures at startup
- ✅ Simple implementation

**Cons**:
- ❌ Only checks at initialization
- ❌ Won't detect if queue is deleted during runtime

**Use When**:
- You want to catch configuration errors early
- Combined with Solution 1 for comprehensive coverage

---

### Solution 3: Publisher Confirms (For High Reliability)

**What**: Wait for broker acknowledgment that message was processed

**Implementation**:
```go
// Enable confirm mode
channel.Confirm(false)

// Publish
channel.PublishWithContext(...)

// Wait for confirmation
select {
case confirm := <-channel.NotifyPublish(...):
    if !confirm.Ack {
        // Message was not confirmed
    }
case <-time.After(5 * time.Second):
    // Timeout
}
```

**Pros**:
- ✅ Highest reliability
- ✅ Detects disk full, queue full, etc.

**Cons**:
- ❌ Performance overhead (synchronous wait)
- ❌ Complex implementation
- ❌ **Does NOT detect queue deletion** (still returns Ack=true)

**Use When**:
- You need guaranteed message persistence
- You care about disk/queue full scenarios
- **NOT recommended for detecting queue existence**

---

### Solution 4: Alternate Exchange (Advanced)

**What**: Route unroutable messages to a fallback exchange

**Implementation**:
```go
// When declaring exchange
args := amqp.Table{
    "alternate-exchange": "fallback-exchange",
}
channel.ExchangeDeclare(name, type, durable, autoDelete, internal, noWait, args)
```

**Pros**:
- ✅ Automatic handling of unroutable messages
- ✅ Can implement dead letter patterns

**Cons**:
- ❌ Requires additional infrastructure
- ❌ Complex setup
- ❌ Overkill for simple use cases

**Use When**:
- You have complex routing requirements
- You want centralized handling of unroutable messages

---

## Comparison Table

| Scenario | Mandatory Flag | Queue Verify | Publisher Confirms | Alternate Exchange |
|----------|---------------|--------------|-------------------|-------------------|
| **Queue Deleted** | ✅ Detect | ❌ No | ❌ No | ✅ Reroute |
| **Queue Not Created** | ✅ Detect | ✅ Detect | ❌ No | ✅ Reroute |
| **Connection Broken** | ✅ Detect | ✅ Detect | ✅ Detect | ✅ Detect |
| **Queue Full** | ❌ No | ❌ No | ✅ Detect (Nack) | ✅ Reroute |
| **Disk Full** | ❌ No | ❌ No | ✅ Detect (Nack) | ❌ No |
| **Performance** | ✅ Fast | ✅ Fast | ⚠️ Slow | ✅ Fast |
| **Complexity** | ✅ Simple | ✅ Simple | ❌ Complex | ❌ Complex |

---

## Recommended Strategy

### For Most Use Cases: **Mandatory + Initialization Verify**

```go
// 1. Add Mandatory support to shared layer (go-packages/rabbitMQ)
type PublishOptions struct {
    // ... existing fields
    Mandatory bool  // Add this field
}

// 2. In app layer: Verify at startup
func NewFactory(args FactoryArgs) *Factory {
    producer := NewNotificationProducer(...)
    if err := producer.VerifyQueue(); err != nil {
        log.Fatalf("Queue not found: %v", err)
    }
    return &Factory{Notification: producer}
}

// 3. In app layer: Enable Mandatory when publishing
publishOptions := rabbitmqlib.DefaultPublishOptions()
publishOptions.Mandatory = true
err := rabbitmqlib.PublishToQueue(conn, queue, payload, &publishOptions)
```

**This covers**:
- ✅ Startup: Queue not created → Fail-fast
- ✅ Runtime: Queue deleted → Error returned
- ✅ Runtime: Connection broken → Error returned

---

## Implementation Checklist

- [ ] Add `Mandatory bool` field to `PublishOptions` in `go-packages/rabbitMQ/types.go`
- [ ] Update `producer.go` to use `options.Mandatory` in `channel.PublishWithContext()`
- [ ] In app layer: Implement `VerifyQueue()` method in producer
- [ ] In app layer: Call `VerifyQueue()` in factory initialization
- [ ] In app layer: Set `Mandatory = true` when calling `PublishToQueue()`
- [ ] Test: Start producer without consumer (should fail at startup)
- [ ] Test: Delete queue during runtime (publish should return error)

---

## Related Issues

- [Issue #XXX]: PRECONDITION_FAILED when producer tries to declare queue with different parameters
  - **Solution**: Set `EnableQueueDeclare = false` in producer (let consumer create queue)

---

## References

- [RabbitMQ Publishers Guide](https://www.rabbitmq.com/publishers.html)
- [RabbitMQ Reliability Guide](https://www.rabbitmq.com/reliability.html)
- [Publisher Confirms](https://www.rabbitmq.com/confirms.html)
