import { RepositoryFactory } from '../repositories';
import { UserRepository } from '../repositories/user.repository';

export class UserService {
  private userRepository: UserRepository;

  constructor(private repositoryFactory: RepositoryFactory) {
    this.userRepository = this.repositoryFactory.getUserRepository();
  }

  async sayHello(args: { name: string }) {
    const res = await this.userRepository.sayHello(args);
    return res;
  }
}
