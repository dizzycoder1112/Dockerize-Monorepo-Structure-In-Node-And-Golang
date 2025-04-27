import { userService } from "./modules"








async function main() {
  const users = await findUsers()

}

main()


export async function findUsers() {

  const users = userService.getUsers()

  return users

}