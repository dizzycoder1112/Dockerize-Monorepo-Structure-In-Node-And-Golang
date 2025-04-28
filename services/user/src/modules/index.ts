import { db } from '../db'
import { RepositoryFactory } from './repositories'
import { UserService } from './services'


const factory = new RepositoryFactory(db)

export const userService = new UserService(factory)