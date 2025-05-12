import { connectNodeAdapter } from '@connectrpc/connect-node';
import type { ConnectRouter } from '@connectrpc/connect';
import { Greeter } from './proto/hello_pb';
import type { HelloRequest, HelloReply } from './proto/hello_pb';
import { ElizaService } from './proto/eliza_pb';
import { GrpcServer } from './types';

export interface GrpcServerOptions {
  greeterImpl?: {
    sayHello: (req: HelloRequest) => Promise<HelloReply>;
  };
  elizaImpl?: Record<string, any>;
}

export function createGrpcRoutes(options: GrpcServerOptions) {
  return (router: ConnectRouter) => {
    if (options.greeterImpl) {
      router.service(Greeter, options.greeterImpl);
    }
    if (options.elizaImpl) {
      router.service(ElizaService, options.elizaImpl);
    }
  };
}

export function createGrpcServer(options: GrpcServerOptions): GrpcServer {
  return connectNodeAdapter({
    routes: createGrpcRoutes(options),
  });
}
