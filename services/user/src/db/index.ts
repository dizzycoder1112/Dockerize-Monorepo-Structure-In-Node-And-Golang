import {  createKyselyLogger } from "@monorepo-packages/db";
import { CamelCasePlugin, Kysely, PostgresDialect } from 'kysely'
import { Pool } from 'pg'
import { Database } from './schema'
import { config } from "../config";
import { SERVICE_NAME } from "@monorepo-packages/shared/constants";


const { DB_HOST, DB_NAME, DB_PASSWORD, DB_PORT, DB_USER } = config

const dbInstance = new Kysely<Database>({
  dialect: new PostgresDialect({
    pool: new Pool({
      user: DB_USER,
      password: DB_PASSWORD,
      host: DB_HOST,
      port: Number(DB_PORT),
      database: DB_NAME,
    }),
  }),
  plugins: [new CamelCasePlugin()],
  log: createKyselyLogger(SERVICE_NAME.USER),
})

export const db = dbInstance.withSchema('public')
