import pino from "pino";

import { LoggerStrategy } from "./interface";

export class PinoStrategy implements LoggerStrategy {
  private logger: pino.Logger;
  private serviceName: string;
  constructor(serviceName = "default", customOptions?: { isPretty?: boolean; level?: string }) {
    const { isPretty = false, level = "debug" } = customOptions ?? {};
    this.serviceName = serviceName;
    const transport = isPretty
      ? {
          targets: [
            {
              target: "pino-pretty",
              options: {
                ignore: "pid,hostname",
                colorize: true,
                translateTime: "SYS:standard"
              }
            }
          ]
        }
      : undefined;
    const options = {
      level,
      timestamp: pino.stdTimeFunctions.isoTime,
      serializers: {
        error: pino.stdSerializers.errWithCause
      },
      transport
    };

    this.logger = pino(options);
  }

  log(level: string, message?: string, context?: Record<string, unknown>) {
    if (context && message) {
      this.logger[level]({ ...context, serviceName: this.serviceName }, message);
    } else if (context) {
      this.logger[level]({ ...context, serviceName: this.serviceName });
    } else if (message) {
      this.logger[level]({ serviceName: this.serviceName }, message);
    }
  }
}
