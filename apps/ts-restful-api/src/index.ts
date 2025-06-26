import config from './config';
import express from 'express';
import { Server } from 'node:http';
import { logger } from './libs/logger';
import routerV1 from './routes/v1';
import { traceMiddleware } from './middleware';

const { PORT } = config;

async function bootstrap() {
  const app = express();

  app.use(traceMiddleware);
  app.use('/api/v1', routerV1);
  app.get('/health-check', async (_req, res) => {
    res.status(200).send('OK');
  });

  const server = app.listen(PORT, () => {
    logger.info(`API server running at http://localhost:${PORT}`);
  });
  setupGracefulShutdown(server);
}

function setupGracefulShutdown(server: Server) {
  process.on('SIGINT', handleShutdown);
  process.on('SIGTERM', handleShutdown);

  async function handleShutdown(signal: string) {
    logger.info(`Received ${signal}, shutting down gracefully...`);

    server.close(async (err: unknown) => {
      if (err) {
        logger.error('Error during shutdown', err);
        process.exit(1);
      }
      logger.info('Closed out remaining connections');
      process.exit(0);
    });
  }
}

bootstrap();
