import { Kysely, PostgresDialect } from "kysely";
import process from "node:process";
import { Pool } from "pg";

import type { Database } from "../schema";

try {
  process.loadEnvFile();
} catch (error) {
  console.error("No .env file found.");
}

export function createDbInstance(): Kysely<Database> {
  const { DB_HOST, DB_NAME, DB_PASSWORD, DB_PORT, DB_USER } = process.env;

  return new Kysely<Database>({
    dialect: new PostgresDialect({
      pool: new Pool({
        user: DB_USER,
        password: DB_PASSWORD,
        host: DB_HOST,
        port: Number(DB_PORT),
        database: DB_NAME
      })
    })
  });
}
