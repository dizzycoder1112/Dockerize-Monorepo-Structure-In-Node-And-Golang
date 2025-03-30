import { PrismaEvent, prismaPrefix } from "@monorepo-packages/shared/constants";
import { asyncContext } from ".";
import { prismaEventEmitter } from "./prismaEventEmiiter";

interface PrismaClientEvent {
  timestamp: Date;
  message?: string;
  query?: string;
  params?: string;
  duration?: number;
  target?: string;
}

interface PrismaLike {
  $on(event: PrismaEvent, callback: (e: any) => void): void;
  $extends(extension: any): any;
}


export function prismaEventRegister(prisma: PrismaLike, serviceName = 'unknown') {
  const getTracePrefix = () => {
    const traceId = asyncContext.getStore()?.traceId ?? 'no-trace-id';
    return `[${traceId}] [${serviceName}]`;
  };

  prisma.$on(PrismaEvent.QUERY, (e) => {
    prismaEventEmitter.emit(`${prismaPrefix}:${PrismaEvent.QUERY}`, JSON.stringify(e));
  });

  prisma.$on(PrismaEvent.INFO, (e: PrismaClientEvent) => {
    const prefix = getTracePrefix();
    console.info(`${prefix} INFO: ${e.message}`);
  });

  prisma.$on(PrismaEvent.WARN, (e: PrismaClientEvent) => {
    const prefix = getTracePrefix();
    console.warn(`${prefix} WARN: ${e.message}`);
  });

  prisma.$on(PrismaEvent.ERROR, (e: PrismaClientEvent) => {
    const prefix = getTracePrefix();
    console.error(`${prefix} ERROR: ${e.message}`);
  });
}
