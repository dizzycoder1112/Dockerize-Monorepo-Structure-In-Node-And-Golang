import type { LogEvent } from "kysely";
import { Buffer } from "node:buffer";

import type { Logger } from "@ts-packages/shared/types";

import { asyncContext } from "./traceAsyncContext";

function maskPII(value: unknown): string | unknown {
  if (typeof value === "string" && looksSensitive(value)) {
    return "***";
  }
  return value;
}

function convertBufferToBase64(value: unknown): string | unknown {
  if (Buffer.isBuffer(value)) {
    return `${value.toString("base64")}`;
  }
  return value;
}

function looksSensitive(value: string): boolean {
  return value.includes("password"); // password
  // value.includes('@') || // email
  // value.length > 50 ||   // long string (possibly a token)
  // value.match(/^(\d{4}[- ]?){4}\d{4}$/) // simple credit-card pattern
}

export function createKyselyLogger(logger: Logger) {
  return (event: LogEvent) => {
    const traceId = asyncContext.getStore()?.traceId ?? "no-trace-id";
    let formattedParams = event.query.parameters.map(convertBufferToBase64);
    formattedParams = formattedParams.map(maskPII);
    if (event.level === "error") {
      logger.error(
        {
          durationMs: event.queryDurationMillis,
          error: event.error,
          sql: event.query.sql,
          params: formattedParams,
          traceId
        },
        `Query Failed`
      );
    } else {
      logger.debug(
        {
          durationMs: event.queryDurationMillis,
          sql: event.query.sql,
          params: formattedParams,
          traceId
        },
        `Query Executed`
      );
    }
  };
}
