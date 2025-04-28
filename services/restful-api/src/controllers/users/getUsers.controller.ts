import type { Request, Response } from "express";



export async function getUsers(req: Request, res: Response) {
  res.status(200).json({ message: "GET /users" });
}