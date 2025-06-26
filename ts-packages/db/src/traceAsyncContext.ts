import { AsyncLocalStorage } from "node:async_hooks";

interface TraceContext {
  traceId: string;
}

export const asyncContext = new AsyncLocalStorage<TraceContext>();
