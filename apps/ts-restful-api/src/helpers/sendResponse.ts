import { Response } from "express";

export function sendResponse(res: Response, { statusCode = 200, success = true, payload = null, message = '', error = null }) {
  const response = {
      success,
      data: payload,
      message,
      error
  };
  res.status(statusCode).json(response);
}