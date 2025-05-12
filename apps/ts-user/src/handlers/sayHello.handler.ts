import { helloProto } from '@ts-packages/grpc';
import { create } from '@ts-packages/grpc';
import { UserService } from '../services';

export function sayHelloHandler(userService: UserService) {
  return async (req: helloProto.HelloRequest): Promise<helloProto.HelloReply> => {
    const helloRes = await userService.sayHello({ name: req.name });
    console.log(helloRes);
    return create(helloProto.HelloReplySchema, { message: helloRes });
  };
}
