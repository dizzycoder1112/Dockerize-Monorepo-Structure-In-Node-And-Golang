import express from "express";
import config  from "./config";
import { Server } from "node:http";
import { logger } from "@monorepo-packages/logger";


const { PORT } = config;

function bootstrap() {
  const app = express();

  app.get("/health-check", async (_req, res) => {
    res.status(200).send("OK");
  });

  app.get('/', (req, res) => {
    res.send('Hello from Express API in pnpm workspace!');
  });

  const server = app.listen(PORT, () => {
    logger.info(`API server running at http://localhost:${PORT}`);
  });
  setupGracefulShutdown(server);
}

function setupGracefulShutdown(server: Server) {
  process.on("SIGINT", handleShutdown);
  process.on("SIGTERM", handleShutdown);

  async function handleShutdown(signal: string) {
    logger.info(`Received ${signal}, shutting down gracefully...`);

    server.close(async (err: unknown) => {
      if (err) {
        logger.error("Error during shutdown", err);
        process.exit(1);
      }
      logger.info("Closed out remaining connections");
      process.exit(0);
    });
  }
}

bootstrap();