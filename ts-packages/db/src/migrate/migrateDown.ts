import process from "node:process";

import { createDbInstance } from "./dbInstance";
import { runMigrations } from "./migrationRunner";

async function migrateDown() {
  const db = createDbInstance();
  await runMigrations(db, "down");
  await db.destroy();
}

migrateDown().catch(err => {
  console.error("❌ Migration failed:", err);
  process.exit(1);
});
