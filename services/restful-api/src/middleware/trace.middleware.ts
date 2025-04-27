import { v4 as uuidv4 } from 'uuid'
import { asyncContext } from "@monorepo-packages/db";


export function traceMiddleware(req: any, res: any, next: any) {
  const traceId = uuidv4()
    res.setHeader('X-Trace-Id', traceId)
    asyncContext.run({ traceId }, () => next());
}