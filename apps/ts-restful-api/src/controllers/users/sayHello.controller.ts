import { Request, Response } from 'express';
import type { UserService } from '../../services/user.service';

export function sayHelloController(userService: UserService) {
  return async (req: Request, res: Response) => {
    let { name } = req.query;
    if (Array.isArray(name)) {
      name = name[0];
    }
    if (typeof name !== 'string') {
      name = '';
    }
    const helloRes = await userService.sayHello({ name });
    res.status(200).json({ message: helloRes.message });
  };
}
