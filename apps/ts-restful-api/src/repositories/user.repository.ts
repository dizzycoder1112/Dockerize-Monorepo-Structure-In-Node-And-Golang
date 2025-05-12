import { GreeterClient } from '@ts-packages/grpc';

export class UserRepository {
  constructor(private client: GreeterClient) {}

  async sayHello(args: { name: string }) {
    return this.client.sayHello(args);
  }
}
