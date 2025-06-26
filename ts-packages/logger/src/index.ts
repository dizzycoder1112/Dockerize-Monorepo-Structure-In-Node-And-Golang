import type { Logger } from "@ts-packages/shared/types";

import { LoggerStrategy } from "./strategies";

export * from "./strategies";

type LogContext = Record<string, unknown>;
type LogArg = string | LogContext;

export function createLogger(strategies: LoggerStrategy[]): Logger {
  function callStrategies(level: string, contextOrMessage: LogArg, message?: string) {
    if (typeof contextOrMessage === "string") {
      // Only message
      strategies.forEach(s => s.log(level, contextOrMessage));
    } else if (typeof message === "string") {
      // context + message
      strategies.forEach(s => s.log(level, message, contextOrMessage));
    } else {
      // Only context
      strategies.forEach(s => s.log(level, undefined, contextOrMessage));
    }
  }

  return {
    info: (contextOrMessage: LogArg, message?: string) => callStrategies("info", contextOrMessage, message),
    warn: (contextOrMessage: LogArg, message?: string) => callStrategies("warn", contextOrMessage, message),
    error: (contextOrMessage: LogArg, message?: string) => callStrategies("error", contextOrMessage, message),
    debug: (contextOrMessage: LogArg, message?: string) => callStrategies("debug", contextOrMessage, message),
    // Optionally, expose the generic log API:
    log: (level: string, contextOrMessage: LogArg, message?: string) => callStrategies(level, contextOrMessage, message)
  };
}
