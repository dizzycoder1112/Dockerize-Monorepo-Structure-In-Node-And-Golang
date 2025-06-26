import { sql } from "kysely";
import process, { stdin as input, stdout as output } from "node:process";
import { createInterface } from "node:readline/promises";

import { createDbInstance } from "./dbInstance";
import { runMigrations } from "./migrationRunner";

async function askConfirmation(): Promise<boolean> {
  const rl = createInterface({ input, output });

  const answer = await rl.question(
    "⚠️  This will DROP ALL TABLES and DELETE ALL DATA. Are you sure you want to reset the database? (Y/N): "
  );

  rl.close();
  return answer.trim().toLowerCase() === "y";
}

async function resetDatabase() {
  const confirmed = await askConfirmation();

  if (!confirmed) {
    console.info("❌ Operation cancelled.");
    process.exit(0);
  }

  const db = createDbInstance();

  console.info("⚠️  Dropping all tables...");
  await sql`DROP SCHEMA public CASCADE;`.execute(db);
  await sql`CREATE SCHEMA public;`.execute(db);
  console.info("✅  Schema reset");

  console.info("🚀 Running migrations...");
  await runMigrations(db, "up");
  console.info("✅  Database ready");

  await db.destroy();
}

resetDatabase().catch(err => {
  console.error("❌ Reset failed:", err);
  process.exit(1);
});
