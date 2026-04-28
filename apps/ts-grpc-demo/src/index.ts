// HTTP/2 only. Native gRPC requires H2 by spec, so the server speaks H2 from
// node:http2 — clients (and any curl test) must opt into H2 explicitly.
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
