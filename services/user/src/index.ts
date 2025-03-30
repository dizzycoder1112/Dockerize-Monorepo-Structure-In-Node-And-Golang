import { PrismaClient } from '@prisma/client';
import { prismaEventRegister } from '@monorepo-packages/db';



const prisma = new PrismaClient({
  log: [
    { level: 'query', emit: 'event' },
    { level: 'info', emit: 'event' },
    { level: 'warn', emit: 'event' },
    { level: 'error', emit: 'event' },
  ],
});


prismaEventRegister(prisma, 'user-service');


export async function findUsers() {


  const users = await prisma.users.findMany({
  })

  return users

}