import { Kysely } from 'kysely'
import { Database } from '../../db/schema'
import { UserRepository } from './users.repository'



export class RepositoryFactory {
  constructor(private db: Kysely<Database>) {}

  getDb() {
    return this.db
  }

  usersRepository() {
    return new UserRepository(this.db)
  }

  
}
