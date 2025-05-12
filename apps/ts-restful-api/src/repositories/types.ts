import { GreeterClient } from '@ts-packages/grpc';

export interface RepositoryDeps {
  userClient: GreeterClient;
  // ... could scale up
}
