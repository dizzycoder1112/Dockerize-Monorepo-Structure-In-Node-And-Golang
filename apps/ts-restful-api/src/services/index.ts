import { createGreeterClient } from '@ts-packages/grpc';
import { RepositoryFactory } from '../repositories';
import { UserService } from './user.service';

const greeterClient = createGreeterClient({
  baseUrl: 'http://localhost:8080',
});

const repositoryFactory = new RepositoryFactory({ userClient: greeterClient });

export const userService = new UserService(repositoryFactory);
