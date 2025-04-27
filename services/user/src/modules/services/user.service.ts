import { UserRepository } from './user.repository'

export class UserService {
  constructor(private repo: UserRepository) {}

  async getUser(id: number) {
    return this.repo.findById(id)
  }

  async getUsers() {
    return this.repo.findAll()
  }

  async createUser(nickname: string) {
    return this.repo.create(nickname)
  }
}