// import * as Sentry from '@sentry/node'
// import { commonConfig } from '../configs/common'
// import { BaseError, ErrorCode } from '../errors'

// const { SENTRY_DSN, NODE_ENV } = commonConfig

// Sentry.init({
//   dsn: SENTRY_DSN,
//   environment: NODE_ENV,
//   tracesSampleRate: 1.0, // adjust sampling for performance monitoring
// })

// export function sentryOnError(error: Error | unknown, context?: any) {
//   if (error instanceof BaseError) {
//     Sentry.captureException(error, {
//       extra: {
//         code: error.code,
//         meta: error.metadata,
//         ...context,
//       },
//     })
//   }
//   else {
//     Sentry.captureException(error, {
//       extra: {
//         code: ErrorCode.UNKNOWN,
//         ...context,
//       },
//     })
//   }
// }
