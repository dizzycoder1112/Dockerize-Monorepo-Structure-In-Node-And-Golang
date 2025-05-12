import { createServer } from 'http2';
import { handler } from './handlers';

async function main() {
  const server = createServer(handler);

  server.listen(8080, () => {
    console.log('✅ gRPC server is listening on http://localhost:8080');
  });
}

main();
