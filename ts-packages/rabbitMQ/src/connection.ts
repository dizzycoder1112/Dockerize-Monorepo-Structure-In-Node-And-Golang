import { connect } from 'amqplib';
import type { Channel, ChannelModel } from 'amqplib';

import { RabbitMQConfig } from './types';

export class RabbitMQConnection {
  private connection: ChannelModel | null = null;
  private defaultChannel: Channel | null = null;
  private channels: Map<string, Channel> = new Map(); // Named channels for isolation
  private consumerTags: Map<string, string> = new Map();
  private checkedQueues: Set<string> = new Set(); // Cache for queue existence checks

  constructor(
    private config: RabbitMQConfig,
    private logger: any | Console = console,
  ) {}

  async connect(): Promise<void> {
    if (this.connection && this.defaultChannel) {
      return;
    }

    try {
      this.logger.info({ url: this.maskUrl(this.config.url) }, 'Connecting to RabbitMQ');

      this.connection = await connect(this.config.url);
      this.defaultChannel = await this.connection.createChannel();

      if (this.config.prefetch) {
        await this.defaultChannel.prefetch(this.config.prefetch);
      }

      this.setupConnectionHandlers();

      this.logger.info('RabbitMQ connected successfully');
    } catch (error) {
      this.logger.error({ error }, 'Failed to connect to RabbitMQ');
      throw error;
    }
  }

  private setupConnectionHandlers(): void {
    if (!this.connection || !this.defaultChannel) return;

    this.connection.on('error', (error) => {
      this.logger.error({ error }, 'RabbitMQ connection error');
    });

    this.connection.on('close', () => {
      this.logger.warn('RabbitMQ connection closed');
    });

    this.setupChannelHandlers(this.defaultChannel, 'default');
  }

  private setupChannelHandlers(channel: Channel, channelId: string): void {
    channel.on('error', (error) => {
      this.logger.error({ error, channelId }, 'RabbitMQ channel error');
    });

    channel.on('close', () => {
      this.logger.warn({ channelId }, 'RabbitMQ channel closed');
      // Remove from map if it's a named channel
      if (channelId !== 'default') {
        this.channels.delete(channelId);
      }
    });
  }

  /**
   * Get a channel by ID
   * @param channelId - Optional channel ID. If not provided, returns default channel
   * @returns Channel instance
   */
  async getChannel(channelId?: string): Promise<Channel> {
    // Return default channel if no channelId specified
    if (!channelId) {
      if (!this.defaultChannel) {
        throw new Error('Default channel not initialized. Call connect() first.');
      }
      return this.defaultChannel;
    }

    // Check if named channel already exists
    if (this.channels.has(channelId)) {
      return this.channels.get(channelId)!;
    }

    // Create new named channel
    if (!this.connection) {
      throw new Error('Connection not initialized. Call connect() first.');
    }

    this.logger.info({ channelId }, 'Creating new named channel');
    const channel = await this.connection.createChannel();

    if (this.config.prefetch) {
      await channel.prefetch(this.config.prefetch);
    }

    this.setupChannelHandlers(channel, channelId);
    this.channels.set(channelId, channel);

    this.logger.info({ channelId }, 'Named channel created successfully');
    return channel;
  }

  getLogger(): any {
    return this.logger;
  }

  /**
   * Ensure queue exists (with caching to avoid repeated checks)
   * Only checks once per queue, then caches the result
   * @param queue - Queue name to check
   * @param channelId - Optional channel ID to use for checking
   */
  async ensureQueueExists(queue: string, channelId?: string): Promise<void> {
    // Check cache first
    if (this.checkedQueues.has(queue)) {
      return;
    }

    const channel = await this.getChannel(channelId);

    try {
      // Use passive check (doesn't create queue)
      await channel.checkQueue(queue);
      this.checkedQueues.add(queue);
      this.logger.debug(
        { queue, channelId: channelId || 'default' },
        'Queue exists - cached for future use',
      );
    } catch (error) {
      this.logger.error(
        { queue, channelId: channelId || 'default', error },
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
      // Close all named channels
      for (const [channelId, channel] of this.channels.entries()) {
        try {
          await channel.close();
          this.logger.debug({ channelId }, 'Named channel closed');
        } catch (error) {
          this.logger.error({ error, channelId }, 'Error closing named channel');
        }
      }
      this.channels.clear();

      // Close default channel
      if (this.defaultChannel) {
        await this.defaultChannel.close();
        this.defaultChannel = null;
      }

      // Close connection
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
    return this.connection !== null && this.defaultChannel !== null;
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
