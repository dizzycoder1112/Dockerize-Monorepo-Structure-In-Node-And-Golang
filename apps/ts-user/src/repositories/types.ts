import { Kysely } from 'kysely';
import { Database } from '../db/schema';

export interface RepositoryDeps {
  db: Kysely<Database>;
}
