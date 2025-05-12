import { connectNodeAdapter } from '@connectrpc/connect-node';
import type { ConnectRouter } from '@connectrpc/connect';
import { Greeter } from './proto/hello_pb';
import { ElizaService } from './proto/eliza_pb';

export interface GrpcServerOptions {
  greeterImpl?: {
    sayHello: (req: any) => any;
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

export function createGrpcServer(options: GrpcServerOptions) {
  const handler = connectNodeAdapter({
    routes: createGrpcRoutes(options),
  });
  return handler;
}
