import { RepositoryFactory } from '../repositories';
import { UserRepository } from '../repositories/users.repository';

export class UserService {
  private userRepo: UserRepository;
  constructor(private repoFactory: RepositoryFactory) {
    this.userRepo = this.repoFactory.usersRepository();
  }

  async sayHello(args: { name: string }) {
    return this.userRepo.sayHello(args);
  }

}
