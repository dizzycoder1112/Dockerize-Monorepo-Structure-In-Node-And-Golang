import { UserRepository } from './users.repository';

export class RepositoryFactory {
  usersRepository() {
    return new UserRepository();
  }
}
