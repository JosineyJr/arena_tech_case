import { start } from "./server";


start(parseInt(process.env.PORT ?? '3000', 10));
