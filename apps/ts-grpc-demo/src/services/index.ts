export * from './user.service';

import { RepositoryFactory } from '../repositories';
import { UserService } from './user.service';

const factory = new RepositoryFactory({ db: null }); // Replace null with actual db instance or mock

export const userService = new UserService(factory);
