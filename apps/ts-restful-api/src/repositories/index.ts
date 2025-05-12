import { GreeterClient } from '@ts-packages/grpc';
import { UserRepository } from './user.repository';

export class RepositoryFactory {
  constructor(private userClient: GreeterClient) {}

  getUserRepository() {
    return new UserRepository(this.userClient);
  }
}
