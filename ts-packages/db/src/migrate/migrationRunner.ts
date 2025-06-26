import type { Kysely } from "kysely";
import { FileMigrationProvider, Migrator } from "kysely";
import { promises as fs } from "node:fs";
import * as path from "node:path";

import type { Database } from "../schema";

export async function runMigrations(db: Kysely<Database>, direction: "up" | "down" = "up") {
  const migrator = new Migrator({
    db,
    provider: new FileMigrationProvider({
      fs,
      path,
      migrationFolder: path.join(__dirname, "migrations")
    })
  });

  const result = direction === "up" ? await migrator.migrateToLatest() : await migrator.migrateDown();

  const { error, results } = result;

  if (!results || results.length === 0) {
    console.info(`📦 Database is already up to date. No migrations applied.`);
  } else {
    results.forEach(it => {
      if (it.status === "Success") {
        console.info(`✅ Migration "${it.migrationName}" ${direction === "up" ? "executed" : "reverted"} successfully`);
      } else if (it.status === "Error") {
        console.error(`❌ Failed to ${direction === "up" ? "execute" : "revert"} migration "${it.migrationName}"`);
      }
    });
  }

  if (error) {
    console.error({ error }, `❌ Migration ${direction} failed`);
    throw error;
  }
}
