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

  async getUser(id: number) {
    return this.userRepo.findById(id);
  }

  async getUsers() {
    return this.userRepo.findAll();
  }

  async createUser(nickname: string) {
    return this.userRepo.create(nickname);
  }
}
