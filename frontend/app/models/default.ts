import { z } from "zod";

export const ServerResponseSchema = z.object({
  message: z.string(),
});

export type ServerResponse = z.infer<typeof ServerResponseSchema>;
