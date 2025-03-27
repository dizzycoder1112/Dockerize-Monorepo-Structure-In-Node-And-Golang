import { Router } from "express";
import { getUsers } from "../../controllers/users";


const router: Router = Router();

router.get("/", getUsers)

export default router;