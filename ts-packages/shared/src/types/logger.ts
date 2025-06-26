export interface Logger {
  info(contextOrMessage: object | string, message?: string): void;
  warn(contextOrMessage: object | string, message?: string): void;
  error(contextOrMessage: object | string, message?: string): void;
  debug(contextOrMessage: object | string, message?: string): void;
  log(level: string, contextOrMessage: object | string, message?: string): void;
}
