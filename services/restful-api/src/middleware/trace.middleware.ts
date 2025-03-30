import { v4 as uuidv4 } from 'uuid'
import { asyncContext, prismaEventEmitter } from "@monorepo-packages/db";
import { PrismaEvent, prismaPrefix } from '@monorepo-packages/shared/constants';
import { logger } from '@monorepo-packages/logger';


export function traceMiddleware(req: any, res: any, next: any) {
  const traceId = uuidv4()
    res.setHeader('X-Trace-Id', traceId)
    prismaEventEmitter.once(`${prismaPrefix}:${PrismaEvent.QUERY}`, (event: string) => {
      const parseEvent = JSON.parse(event);
      logger.debug({...parseEvent, traceId}, "Prisma event listen")
    });
    asyncContext.run({ traceId }, () => next());
}