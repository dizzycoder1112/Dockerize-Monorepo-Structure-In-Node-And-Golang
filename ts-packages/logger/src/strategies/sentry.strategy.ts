import * as Sentry from "@sentry/node";

import { LoggerStrategy } from "./interface";

export class SentryStrategy implements LoggerStrategy {
  private sentryClient: typeof Sentry;
  constructor(options: { dsn: string; environment: string }) {
    this.sentryClient = Sentry;
    this.sentryClient.init({
      ...options,
      tracesSampleRate: 1.0
    });
  }
  log(level: string, message?: string, context?: Record<string, unknown>): void {
    if (level !== "error" && level !== "warn") return;
    if (!context || typeof context !== "object") return;
    const { error, ...extra } = context;
    switch (level) {
      case "error":
        if (!(error instanceof Error)) {
          console.error("SentryStrategy: error is not an Error instance", error);
          return;
        }
        this.sentryClient.withScope(scope => {
          Object.entries(extra).forEach(([key, value]) => {
            scope.setExtra(key, value);
          });
          this.sentryClient.captureException(error);
        });
        break;
      case "warn":
        if (!message) {
          console.error("SentryStrategy: warn level requires a message");
          return;
        }
        this.sentryClient.withScope(scope => {
          Object.entries(extra).forEach(([key, value]) => {
            scope.setExtra(key, value);
          });
          this.sentryClient.captureMessage(message);
        });
        break;
    }
  }
}
