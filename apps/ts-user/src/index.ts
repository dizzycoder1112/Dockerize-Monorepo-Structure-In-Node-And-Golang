import { connectNodeAdapter } from '@connectrpc/connect-node';
import { createServer } from 'http2';
import type { ConnectRouter } from '@connectrpc/connect';

import { Greeter, HelloRequest } from '@ts-packages/grpc/src/proto/hello_pb';

function routes(router: ConnectRouter) {
  router.service(Greeter, {
    sayHello: (req: HelloRequest) => {
      console.log(req.name);
      return {
        message: `You said ${req.name}`,
      };
    },
  });
}

async function main() {
  const handler = connectNodeAdapter({
    routes,
  });

  const server = createServer(handler);

  server.listen(8080, () => {
    console.log('✅ gRPC server is listening on http://localhost:8080');
  });
}

main();
