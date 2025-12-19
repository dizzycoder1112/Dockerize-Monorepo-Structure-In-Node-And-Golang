import { RabbitMQConnection } from './connection';
import { PublishOptions, QueueOptions, ExchangeOptions } from './types';

/**
 * Publish message to an exchange with routing key
 * For topic/fanout/direct exchanges
 */
export async function publishToExchange<T = unknown>(
  connection: RabbitMQConnection,
  exchange: string,
  routingKey: string,
  payload: T,
  options?: PublishOptions & { exchangeOptions?: ExchangeOptions },
): Promise<boolean> {
  const channel = await connection.getChannel(options?.channelId);
  const logger = connection.getLogger();

  try {
    // Ensure exchange exists
    const exchangeType = options?.exchangeOptions?.type || 'topic';
    await channel.assertExchange(exchange, exchangeType, {
      durable: options?.exchangeOptions?.durable !== false,
      autoDelete: options?.exchangeOptions?.autoDelete || false,
    });

    const message = Buffer.from(JSON.stringify(payload));
    const publishOptions = {
      persistent: options?.persistent !== false,
      priority: options?.priority,
      expiration: options?.expiration,
      headers: options?.headers,
    };

    const result = channel.publish(exchange, routingKey, message, publishOptions);

    logger.debug(
      {
        exchange,
        routingKey,
        payloadSize: message.length,
        channelId: options?.channelId || 'default',
      },
      'Message published to exchange',
    );

    return result;
  } catch (error) {
    logger.error(
      {
        error,
        exchange,
        routingKey,
        channelId: options?.channelId || 'default',
      },
      'Failed to publish message to exchange',
    );
    throw error;
  }
}

export async function publishToQueue<T = unknown>(
  connection: RabbitMQConnection,
  queue: string,
  payload: T,
  options?: PublishOptions & { queueOptions?: QueueOptions },
): Promise<boolean> {
  const channel = await connection.getChannel(options?.channelId);
  const logger = connection.getLogger();

  try {
    // Check queue exists (with caching - only checks once per queue)
    // This prevents PRECONDITION_FAILED when consumer has DLX config
    await connection.ensureQueueExists(queue, options?.channelId);

    const message = Buffer.from(JSON.stringify(payload));
    const publishOptions = {
      persistent: options?.persistent !== false,
      priority: options?.priority,
      expiration: options?.expiration,
      headers: options?.headers,
    };

    const result = channel.sendToQueue(queue, message, publishOptions);

    logger.debug(
      {
        queue,
        payloadSize: message.length,
        channelId: options?.channelId || 'default',
      },
      'Message published to queue',
    );

    return result;
  } catch (error) {
    logger.error(
      {
        error,
        queue,
        channelId: options?.channelId || 'default',
      },
      'Failed to publish message to queue',
    );
    throw error;
  }
}