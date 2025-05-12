import { UserRepository } from './user.repository';
import { RepositoryDeps } from './types';

export class RepositoryFactory {
  constructor(private deps: RepositoryDeps) {}

  getUserRepository() {
    return new UserRepository(this.deps.userClient);
  }
}
