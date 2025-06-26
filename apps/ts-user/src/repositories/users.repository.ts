
export class UserRepository {
  constructor(private db: any) {}

  async sayHello(args: { name: string }) {
    return `You said ${args.name}`;
  }
}
