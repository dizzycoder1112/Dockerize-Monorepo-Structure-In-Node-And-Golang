import process from 'node:process';
import z from 'zod';
import { parseEnv, port } from 'znv';

try {
  process.loadEnvFile();
} catch (e) {
  console.warn('No .env file found');
}

function createConfigFromEnvironment(environment: NodeJS.ProcessEnv) {
  const config = parseEnv(environment, {
    NODE_ENV: z.enum(['development', 'production', 'local']),
    LOG_LEVEL: z
      .enum(['trace', 'debug', 'info', 'warn', 'error', 'fatal', 'silent'])
      .default('info'),
    PORT: port().default(3000),
  });

  return {
    ...config,
    isDev: process.env.NODE_ENV === 'development',
    isProd: process.env.NODE_ENV === 'production',
  };
}

export type Config = ReturnType<typeof createConfigFromEnvironment>;

export default createConfigFromEnvironment(process.env);
