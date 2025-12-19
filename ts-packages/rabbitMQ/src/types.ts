import { ConsumeMessage } from 'amqplib';

export interface RabbitMQConfig {
  url: string;
  prefetch?: number;
}

export interface QueueOptions {
  durable?: boolean;
  exclusive?: boolean;
  autoDelete?: boolean;
  arguments?: Record<string, unknown>;
}

export interface ExchangeOptions {
  type?: 'direct' | 'topic' | 'fanout' | 'headers';
  durable?: boolean;
  autoDelete?: boolean;
  arguments?: Record<string, unknown>;
}

export interface PublishOptions {
  persistent?: boolean;
  priority?: number;
  expiration?: string | number;
  headers?: Record<string, unknown>;
  /**
   * Optional channel ID for channel isolation
   * If not provided, uses default channel
   */
  channelId?: string;
}

export interface ConsumeOptions {
  noAck?: boolean;
  exclusive?: boolean;
  consumerTag?: string;
  /**
   * Optional channel ID for channel isolation
   * If not provided, uses default channel
   */
  channelId?: string;
}

export interface MessageHandler<T = unknown> {
  (payload: T, message: ConsumeMessage): Promise<void>;
}