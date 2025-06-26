import fs from "node:fs";
import path from "node:path";
import process from "node:process";

const migrationsDir = path.resolve(__dirname, "./migrations");

const name = process.argv[2];

if (!name) {
  console.error("❌ please enter migration name, ex: pnpm run db:migrate:create add_users_table");
  process.exit(1);
}

const timestamp = new Date().toISOString().replace(/[-:T]/g, "").slice(0, 14); // e.g., 20250408153200

const fileName = `${timestamp}_${name}.ts`;
const filePath = path.join(migrationsDir, fileName);

const template = `
import { type Kysely, sql } from 'kysely'

export async function up(db: Kysely<any>): Promise<void> {
  // Migration code
  await sql\`\`.execute(db)
}

export async function down(db: Kysely<any>): Promise<void> {
  // Migration code
  await sql\`\`.execute(db)
}
`;

fs.mkdirSync(migrationsDir, { recursive: true });
fs.writeFileSync(filePath, template);

console.info(`✅ Migration created: ${filePath}`);
