export * from './user.service';

import { RepositoryFactory } from '../repositories';
import { UserService } from './user.service';

const factory = new RepositoryFactory();

export const userService = new UserService(factory);
