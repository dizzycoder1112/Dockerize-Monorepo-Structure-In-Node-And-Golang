export interface LoggerStrategy {
  log(level: string, message?: string, context?: Record<string, unknown>): void;
}
