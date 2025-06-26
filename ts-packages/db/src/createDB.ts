import { CamelCasePlugin, Kysely, PostgresDialect, sql } from "kysely";
import { Pool } from "pg";

import { Logger } from "@ts-packages/shared/types";

import { createKyselyLogger } from "./kyselyLog";

export function createDB<Database>(
  config: { user: string; password: string; host: string; port: number; database: string },
  logger: Logger | Console = console
) {
  const pool = new Pool(config);
  const db = new Kysely<Database>({
    dialect: new PostgresDialect({
      pool
    }),
    plugins: [new CamelCasePlugin()],
    log: createKyselyLogger(logger)
  });

  const { user, host, port, database } = config;

  db.connection()
    .execute(async db => {
      await sql`SELECT 1`.execute(db);
      logger.info({ database }, "Database connected successfully");
      logger.info(
        { totalCount: pool.totalCount, idleCount: pool.idleCount, waitingCount: pool.waitingCount },
        "PG pool status"
      );
    })
    .catch(error => {
      logger.error({ error, database, user, host, port }, "Failed to connect to database");
      throw error;
    });

  return db;
}
