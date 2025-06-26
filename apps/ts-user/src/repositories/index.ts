import { UserRepository } from './users.repository';
import { RepositoryDeps } from './types';

export class RepositoryFactory {
  private db: any;
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
