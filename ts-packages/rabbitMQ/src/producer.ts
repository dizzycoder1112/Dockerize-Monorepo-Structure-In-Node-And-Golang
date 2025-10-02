import { RabbitMQConnection } from './connection';
import { PublishOptions, QueueOptions } from './types';

export async function publishToQueue<T = unknown>(
  connection: RabbitMQConnection,
  queue: string,
  payload: T,
  options?: PublishOptions & { queueOptions?: QueueOptions },
): Promise<boolean> {
  const channel = connection.getChannel();
  const logger = connection.getLogger();

  try {
    // Check queue exists (with caching - only checks once per queue)
    // This prevents PRECONDITION_FAILED when consumer has DLX config
    await connection.ensureQueueExists(queue);

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
      },
      'Message published to queue',
    );

    return result;
  } catch (error) {
    logger.error(
      {
        error,
        queue,
      },
      'Failed to publish message to queue',
    );
    throw error;
  }
}