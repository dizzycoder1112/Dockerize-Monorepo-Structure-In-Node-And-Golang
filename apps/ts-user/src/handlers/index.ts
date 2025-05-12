import { createGrpcServer } from '@ts-packages/grpc';
import { sayHelloHandler } from './sayHello.handler';
import { userService } from '../services';

export const handler = createGrpcServer({
  greeterImpl: {
    sayHello: sayHelloHandler(userService),
  },
});
