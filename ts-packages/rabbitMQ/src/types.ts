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

export interface PublishOptions {
  persistent?: boolean;
  priority?: number;
  expiration?: string | number;
  headers?: Record<string, unknown>;
}

export interface ConsumeOptions {
  noAck?: boolean;
  exclusive?: boolean;
  consumerTag?: string;
}

export interface MessageHandler<T = unknown> {
  (payload: T, message: ConsumeMessage): Promise<void>;
}