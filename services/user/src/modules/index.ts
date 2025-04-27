import { db } from '../db'
import { UserRepository } from './repositories'
import { UserService } from './services'

export const userRepository = new UserRepository(db)
export const userService = new UserService(userRepository)