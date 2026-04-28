export class UserRepository {
  async sayHello(args: { name: string }) {
    return `You said ${args.name}`;
  }
}
