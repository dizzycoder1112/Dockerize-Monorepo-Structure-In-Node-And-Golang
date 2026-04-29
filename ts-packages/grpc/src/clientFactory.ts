import { createGrpcTransport } from '@connectrpc/connect-node';
import { createClient } from '@connectrpc/connect';
import { Greeter } from './proto/hello_connect';
import { ElizaService } from './proto/eliza_connect';

// Client always speaks native gRPC over HTTP/2. ts-grpc-demo's connectNodeAdapter
// is multi-protocol and routes application/grpc to its native gRPC handler;
// go-grpc-demo speaks native gRPC directly. So a single transport works for both.
export interface GrpcClientOptions {
  baseUrl?: string;
}

function buildTransport(options: GrpcClientOptions) {
  return createGrpcTransport({
    baseUrl: options.baseUrl ?? 'http://localhost:8080',
    httpVersion: '2',
  });
}

export function createGreeterClient(options: GrpcClientOptions = {}) {
  return createClient(Greeter, buildTransport(options));
}

export function createElizaClient(options: GrpcClientOptions = {}) {
  return createClient(ElizaService, buildTransport(options));
}
