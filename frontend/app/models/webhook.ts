import { z } from "zod";
import { ServerResponse, ServerResponseSchema } from "./default";

export const WebhookSchema = z.object({
  id: z.string(),
  name: z.string(),
  description: z.string(),
  url: z.string(),
  created_at: z.string(),
  updated_at: z.string(),
});

export type Webhook = z.infer<typeof WebhookSchema>;

export const WebhooksSchema = z.array(WebhookSchema);

export type Webhooks = z.infer<typeof WebhooksSchema>;

export const GetWebhookSchema = ServerResponseSchema.extend({
  data: WebhookSchema,
});

export const GetWebhooksSchema = ServerResponseSchema.extend({
  data: WebhooksSchema,
});

export const CreateWebhookFormSchema = z.object({
  name: z.string().min(3),
  description: z.string().optional(),
  url: z.string().url(),
  intent: z.enum(["create", "update", "delete"]).default("create"),
  events: z.array(z.string()).default([]),
  id: z.string().optional(),
});

export type CreateWebhookForm = z.infer<typeof CreateWebhookFormSchema>;
