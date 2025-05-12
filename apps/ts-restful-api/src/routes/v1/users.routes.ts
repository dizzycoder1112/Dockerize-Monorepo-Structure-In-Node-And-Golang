import { Router } from 'express';
import { getUsers, sayHelloController } from '../../controllers/users';
import { userService } from '../../services';

const router: Router = Router();

router.get('/', getUsers);
router.get('/sayHello', sayHelloController(userService));

export default router;
