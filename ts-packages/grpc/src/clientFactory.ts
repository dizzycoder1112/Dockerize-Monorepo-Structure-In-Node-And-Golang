import { createConnectTransport } from '@connectrpc/connect-node';
import { createClient } from '@connectrpc/connect';
import { Greeter } from './proto/hello_pb';
import { ElizaService } from './proto/eliza_pb';

export interface GrpcClientOptions {
  baseUrl?: string;
  httpVersion?: '2' | '1.1';
}

export function createGreeterClient(options: GrpcClientOptions = {}) {
  const transport = createConnectTransport({
    baseUrl: options.baseUrl ?? 'http://localhost:8080',
    httpVersion: options.httpVersion ?? '2',
  });
  return createClient(Greeter, transport);
}

export function createElizaClient(options: GrpcClientOptions = {}) {
  const transport = createConnectTransport({
    baseUrl: options.baseUrl ?? 'http://localhost:8080',
    httpVersion: options.httpVersion ?? '2',
  });
  return createClient(ElizaService, transport);
}
