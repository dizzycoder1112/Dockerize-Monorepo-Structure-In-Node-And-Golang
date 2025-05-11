import { Generated, Insertable, Selectable, Updateable } from 'kysely'

export interface Database {
  users: UsersTable
}

export interface UsersTable {
  id: Generated<number>
  nickname: string
  created_at: Date
  updated_at: Date
}

export type User = Selectable<UsersTable>
export type NewUser = Insertable<UsersTable>
export type UpdateUser = Updateable<UsersTable>
