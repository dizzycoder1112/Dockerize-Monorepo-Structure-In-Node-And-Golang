import config from './config';
import express from 'express';
import { Server } from 'node:http';
import { logger } from '@ts-packages/logger';
import routerV1 from './routes/v1';
import { traceMiddleware } from './middleware';

import { createConnectTransport } from '@connectrpc/connect-node';
import { createClient } from '@connectrpc/connect';
import { ElizaService } from '@ts-packages/grpc/src/proto/eliza_pb';
import { Greeter } from '@ts-packages/grpc/src/proto/hello_pb';

const { PORT } = config;

async function bootstrap() {
  // 建立 Connect transport
  const transport = createConnectTransport({
    baseUrl: 'http://localhost:8080',
    httpVersion: '2',
  });

  // 建立 Promise-based gRPC client
  const client = createClient(Greeter, transport);
  try {
    const res = await client.sayHello({ name: 'Champer' });
    console.log(res.message);
  } catch (err) {
    console.error('gRPC call failed:', err);
  }

  //   const app = express();

  //   app.use(traceMiddleware);
  //   app.use("/api/v1", routerV1);

  //   app.get("/health-check", async (_req, res) => {
  //     res.status(200).send("OK");
  //   });

  //   app.get('/', (req, res) => {
  //     res.send('Hello from Express API in pnpm workspace!');
  //   });

  //   const server = app.listen(PORT, () => {
  //     logger.info(`API server running at http://localhost:${PORT}`);
  //   });
  //   setupGracefulShutdown(server);
  // }

  // function setupGracefulShutdown(server: Server) {
  //   process.on("SIGINT", handleShutdown);
  //   process.on("SIGTERM", handleShutdown);

  //   async function handleShutdown(signal: string) {
  //     logger.info(`Received ${signal}, shutting down gracefully...`);

  //     server.close(async (err: unknown) => {
  //       if (err) {
  //         logger.error("Error during shutdown", err);
  //         process.exit(1);
  //       }
  //       logger.info("Closed out remaining connections");
  //       process.exit(0);
  //     });
  //   }
}

bootstrap();
