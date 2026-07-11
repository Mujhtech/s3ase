import { z } from "zod";
import { ServerResponseSchema } from "./default";

export const SubscriptionSchema = z.object({
  id: z.string(),
  app_id: z.string(),
  plan: z.enum(["basic", "pro", "enterprise"]),
  status: z.string(),
});

export type Subscription = z.infer<typeof SubscriptionSchema>;
export const GetSubscriptionSchema = ServerResponseSchema.extend({ data: SubscriptionSchema });
