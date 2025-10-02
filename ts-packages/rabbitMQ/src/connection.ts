import { connect } from 'amqplib';
import { RabbitMQConfig } from './types';
import type { ChannelModel, Channel } from 'amqplib';

export class RabbitMQConnection {
  private connection: ChannelModel | null = null;
  private channel: Channel | null = null;
  private consumerTags: Map<string, string> = new Map();
  private checkedQueues: Set<string> = new Set(); // Cache for queue existence checks

  constructor(
    private config: RabbitMQConfig,
    private logger: any,
  ) {}

  async connect(): Promise<void> {
    if (this.connection && this.channel) {
      return;
    }

    try {
      this.logger.info({ url: this.maskUrl(this.config.url) }, 'Connecting to RabbitMQ');

      this.connection = await connect(this.config.url);
      this.channel = await this.connection.createChannel();

      if (this.config.prefetch) {
        await this.channel.prefetch(this.config.prefetch);
      }

      this.setupConnectionHandlers();

      this.logger.info('RabbitMQ connected successfully');
    } catch (error) {
      this.logger.error({ error }, 'Failed to connect to RabbitMQ');
      throw error;
    }
  }

  private setupConnectionHandlers(): void {
    if (!this.connection || !this.channel) return;

    this.connection.on('error', (error) => {
      this.logger.error({ error }, 'RabbitMQ connection error');
    });

    this.connection.on('close', () => {
      this.logger.warn('RabbitMQ connection closed');
    });

    this.channel.on('error', (error) => {
      this.logger.error({ error }, 'RabbitMQ channel error');
    });

    this.channel.on('close', () => {
      this.logger.warn('RabbitMQ channel closed');
    });
  }

  getChannel(): Channel {
    if (!this.channel) {
      throw new Error('Channel not initialized. Call connect() first.');
    }
    return this.channel;
  }

  getLogger(): any {
    return this.logger;
  }

  /**
   * Ensure queue exists (with caching to avoid repeated checks)
   * Only checks once per queue, then caches the result
   */
  async ensureQueueExists(queue: string): Promise<void> {
    // Check cache first
    if (this.checkedQueues.has(queue)) {
      return;
    }

    if (!this.channel) {
      throw new Error('Channel not initialized. Call connect() first.');
    }

    try {
      // Use passive check (doesn't create queue)
      await this.channel.checkQueue(queue);
      this.checkedQueues.add(queue);
      this.logger.debug({ queue }, 'Queue exists - cached for future use');
    } catch (error) {
      this.logger.error(
        { queue, error },
        'Queue does not exist - consumer must be started first',
      );
      throw new Error(
        `Queue '${queue}' does not exist. Consumer must be started first to create queue infrastructure.`,
      );
    }
  }

  /**
   * Clear queue cache (useful for testing or when queue infrastructure changes)
   */
  clearQueueCache(): void {
    this.checkedQueues.clear();
  }

  // Consumer tag management
  registerConsumerTag(queue: string, consumerTag: string): void {
    this.consumerTags.set(queue, consumerTag);
  }

  getConsumerTag(queue: string): string | undefined {
    return this.consumerTags.get(queue);
  }

  removeConsumerTag(queue: string): void {
    this.consumerTags.delete(queue);
  }

  async close(): Promise<void> {
    try {
      if (this.channel) {
        await this.channel.close();
        this.channel = null;
      }
      if (this.connection) {
        await this.connection.close();
        this.connection = null;
      }
      this.consumerTags.clear();
      this.logger.info('RabbitMQ connection closed');
    } catch (error) {
      this.logger.error({ error }, 'Error closing RabbitMQ connection');
      throw error;
    }
  }

  isConnected(): boolean {
    return this.connection !== null && this.channel !== null;
  }

  private maskUrl(url: string): string {
    try {
      const parsed = new URL(url);
      if (parsed.password) {
        parsed.password = '***';
      }
      return parsed.toString();
    } catch {
      return 'invalid-url';
    }
  }
}
