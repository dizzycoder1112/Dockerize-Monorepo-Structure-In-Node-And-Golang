import { createConnectTransport } from '@connectrpc/connect-node';
import { createClient } from '@connectrpc/connect';
import { Greeter } from './proto/hello_connect';
import { ElizaService } from './proto/eliza_connect';

// HTTP/2 is hard-coded across this template (server uses node:http2, clients
// use httpVersion: '2'). Native gRPC requires HTTP/2 by spec, so locking to
// H2 keeps Go ↔ Node interop trivial. Connect protocol also runs fine over H2,
// so we lose nothing by removing the H1.1 escape hatch.
export interface GrpcClientOptions {
  baseUrl?: string;
}

export function createGreeterClient(options: GrpcClientOptions = {}) {
  const transport = createConnectTransport({
    baseUrl: options.baseUrl ?? 'http://localhost:8080',
    httpVersion: '2',
  });
  return createClient(Greeter, transport);
}

export function createElizaClient(options: GrpcClientOptions = {}) {
  const transport = createConnectTransport({
    baseUrl: options.baseUrl ?? 'http://localhost:8080',
    httpVersion: '2',
  });
  return createClient(ElizaService, transport);
}
