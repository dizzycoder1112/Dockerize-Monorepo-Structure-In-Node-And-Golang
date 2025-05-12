import { Kysely } from 'kysely';
import { Database } from '../db/schema';
import { UserRepository } from './users.repository';
import { RepositoryDeps } from './types';

export class RepositoryFactory {
  private db: Kysely<Database>;
  constructor(private repoDeps: RepositoryDeps) {
    this.db = repoDeps.db;
  }

  getDb() {
    return this.db;
  }

  usersRepository() {
    return new UserRepository(this.db);
  }
}
