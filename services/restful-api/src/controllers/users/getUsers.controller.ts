import type { Request, Response } from "express";
import {findUsers} from "@monorepo-services/user"



export async function getUsers(req: Request, res: Response) {
  const users = await findUsers();
  res.status(200).json({ message: "GET /users" });
}