import { Kysely } from 'kysely';
import { Database } from '../db/schema';

export class UserRepository {
  constructor(private db: Kysely<Database>) {}

  async findById(id: number) {
    return this.db.selectFrom('users').selectAll().where('id', '=', id).executeTakeFirst();
  }

  async findAll() {
    return this.db.selectFrom('users').selectAll().execute();
  }

  async create(nickname: string) {
    return this.db.insertInto('users').values({ nickname }).returningAll().executeTakeFirst();
  }

  async sayHello(args: { name: string }) {
    return `You said ${args.name}`;
  }
}
