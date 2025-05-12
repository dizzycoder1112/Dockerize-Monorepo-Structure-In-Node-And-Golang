import { createServer } from 'http2';
import { handler } from './handlers';
import { config } from './config';

const { PORT } = config;

async function main() {
  const server = createServer(handler);

  server.listen(PORT, () => {
    console.log(`gRPC server is listening on http://localhost:${PORT}`);
  });
}

main();
