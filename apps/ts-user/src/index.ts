import { createServer } from 'http2';
import { createGrpcServer } from '@ts-packages/grpc';



async function main() {
  const handler = createGrpcServer({
    greeterImpl: {
      sayHello: (req) => ({ message: `You said ${req.name}` }),
    },
  });
  

  const server = createServer(handler);

  server.listen(8080, () => {
    console.log('✅ gRPC server is listening on http://localhost:8080');
  });
}

main();
