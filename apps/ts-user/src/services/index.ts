export * from './user.service';

import { db } from '../db';
import { RepositoryFactory } from '../repositories';
import { UserService } from './user.service';

const factory = new RepositoryFactory({ db });

export const userService = new UserService(factory);
