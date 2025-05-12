import { connectNodeAdapter } from '@connectrpc/connect-node';
import { createElizaClient, createGreeterClient } from './clientFactory';

export type GreeterClient = ReturnType<typeof createGreeterClient>;
export type ElizaClient = ReturnType<typeof createElizaClient>;
export type GrpcServer = ReturnType<typeof connectNodeAdapter>;
