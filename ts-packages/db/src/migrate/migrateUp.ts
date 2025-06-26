import process from "node:process";

import { createDbInstance } from "./dbInstance";
import { runMigrations } from "./migrationRunner";

async function migrateUp() {
  const db = createDbInstance();
  await runMigrations(db, "up");
  await db.destroy();
}

migrateUp().catch(err => {
  console.error("❌ Migration failed:", err);
  process.exit(1);
});
