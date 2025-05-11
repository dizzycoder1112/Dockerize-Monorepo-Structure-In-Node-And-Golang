import { RepositoryFactory } from "../repositories"


export class UserService {
  constructor(private repoFactory: RepositoryFactory) {}

  async getUser(id: number) {
    const userRepo = this.repoFactory.usersRepository()
    return userRepo.findById(id)
  }

  async getUsers() {
    const userRepo = this.repoFactory.usersRepository()
    return userRepo.findAll()
  }

  async createUser(nickname: string) {
    const userRepo = this.repoFactory.usersRepository()
    return userRepo.create(nickname)
  }
}