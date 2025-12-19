import { ConsumeMessage } from 'amqplib';
import { RabbitMQConnection } from './connection';
import { ConsumeOptions, MessageHandler, QueueOptions } from './types';

export async function consumeQueue<T = unknown>(
  connection: RabbitMQConnection,
  queue: string,
  handler: MessageHandler<T>,
  options?: ConsumeOptions & { queueOptions?: QueueOptions },
): Promise<string> {
  const channel = await connection.getChannel(options?.channelId);
  const logger = connection.getLogger();

  try {
    await channel.assertQueue(queue, options?.queueOptions || { durable: true });

    const { consumerTag } = await channel.consume(
      queue,
      async (msg: ConsumeMessage | null) => {
        if (!msg) return;

        try {
          const payload = JSON.parse(msg.content.toString()) as T;

          logger.debug(
            {
              queue,
              messageId: msg.properties.messageId,
              channelId: options?.channelId || 'default',
            },
            'Processing message',
          );

          await handler(payload, msg);

          if (!options?.noAck) {
            channel.ack(msg);
          }
        } catch (error) {
          logger.error(
            {
              error,
              queue,
              message: msg.content.toString(),
              channelId: options?.channelId || 'default',
            },
            'Error processing message',
          );

          // NACK and requeue on error
          if (!options?.noAck) {
            channel.nack(msg, false, true);
          }
        }
      },
      {
        noAck: options?.noAck || false,
        exclusive: options?.exclusive,
        consumerTag: options?.consumerTag,
      },
    );

    connection.registerConsumerTag(queue, consumerTag);

    logger.info({ queue, consumerTag, channelId: options?.channelId || 'default' }, 'Consumer registered');

    return consumerTag;
  } catch (error) {
    logger.error({ error, queue, channelId: options?.channelId || 'default' }, 'Failed to register consumer');
    throw error;
  }
}

export async function cancelConsumer(
  connection: RabbitMQConnection,
  queue: string,
): Promise<void> {
  const channel = await connection.getChannel();
  const logger = connection.getLogger();
  const consumerTag = connection.getConsumerTag(queue);

  if (!consumerTag) {
    logger.warn({ queue }, 'No consumer tag found for queue');
    return;
  }

  try {
    await channel.cancel(consumerTag);
    connection.removeConsumerTag(queue);
    logger.info({ queue, consumerTag }, 'Consumer cancelled');
  } catch (error) {
    logger.error({ error, queue }, 'Failed to cancel consumer');
    throw error;
  }
}