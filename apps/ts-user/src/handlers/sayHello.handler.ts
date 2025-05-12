import { HelloReply, HelloReplySchema, HelloRequest } from 'ts-packages/grpc/src/proto/hello_pb';
import { create } from '@ts-packages/grpc';
import { UserService } from '../services';

export function sayHelloHandler(userService: UserService) {
  return async (req: HelloRequest): Promise<HelloReply> => {
    const helloRes = await userService.sayHello({ name: req.name });
    console.log(helloRes);
    return create(HelloReplySchema, { message: helloRes });
  };
}
