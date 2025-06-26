import {
  LoggerStrategy,
  PinoStrategy,
  createLogger
} from "@ts-packages/logger";
import { SERVICE_NAME } from "@ts-packages/shared/constants";

import config from "../config";

const {
  NODE_ENV,
  LOG_LEVEL
} = config;

const strategies: LoggerStrategy[] = [];
const pinoStrategy = new PinoStrategy(SERVICE_NAME.CRON, {
  isPretty: NODE_ENV === "local",
  level: LOG_LEVEL
});
strategies.push(pinoStrategy);


export const logger = createLogger(strategies);
