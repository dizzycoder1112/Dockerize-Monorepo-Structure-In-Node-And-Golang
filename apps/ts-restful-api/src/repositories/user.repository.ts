import { createGreeterClient } from '@ts-packages/grpc';

export class UserRepository {
  constructor(private client: ReturnType<typeof createGreeterClient>) {}

  async sayHello(args: { name: string }) {
    return this.client.sayHello(args);
  }
}
