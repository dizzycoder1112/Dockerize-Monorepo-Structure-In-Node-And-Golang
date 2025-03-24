import { pino } from 'pino'
// import { holdemConfig as config } from '../../configs'

// const { LOG_LEVEL } = config

export const logger = pino({
  level: 'debug',
  timestamp: pino.stdTimeFunctions.isoTime,
  serializers: {
    error: pino.stdSerializers.errWithCause,

  },
  transport: {
    targets: [
      {
        target: 'pino-pretty',
        options: {
          ignore: 'pid,hostname',
          colorize: true,
          translateTime: 'SYS:standard',
        },
      },
    ],
  },
})



// import { pino } from "pino";
// import { config } from "@src/config";

// // Define a transport for writing to the error log file and another for console
// const fileTransport = pino.transport({
//   targets: [
//     {
//       target: "pino/file",
//       level: "error",
//       options: {
//         destination: "./logs/error.log",
//         translateTime: "SYS:standard",
//         mkdir: true,
//       },
//     },
//     {
//       target: "pino/file",
//       options: {
//         destination: "./logs/combined.log",
//         mkdir: true,
//       },
//     },
//     {
//       target: "pino-pretty",
//       level: "info",
//       options: {
//         colorize: true, // Enable colorization for console output
//         translateTime: "SYS:standard",
//         ignore: "pid,hostname", // Ignore pid and hostname
//       },
//     },
//   ],
// });

// const logger = pino(
//   {
//     level: config.LOG_LEVEL,
//     timestamp: pino.stdTimeFunctions.isoTime,
//   },
//   fileTransport
// );

// Example usage of the logger
// logger.info({ hah: "hah" }, "This is an info message");
// logger.error("This is an error message");

// export default logger;

// export type Logger = typeof logger;

