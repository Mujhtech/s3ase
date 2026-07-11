import { z } from "zod";
import { ServerResponseSchema } from "./default";

export const ApiKeySchema = z.object({
  id: z.string(),
  name: z.string(),
  description: z.string().nullable().transform((value) => value ?? ""),
  access: z.string(),
  expired_at: z.number().nullable().optional(),
  created_at: z.string(),
  updated_at: z.string(),
  last_used: z.string().nullable().optional(),
  key_prefix: z.string().optional(),
  last_four: z.string().optional(),
  secret: z.string().optional(),
});

export type ApiKey = z.infer<typeof ApiKeySchema>;

export const ApiKeysSchema = z.array(ApiKeySchema);

export type ApiKeys = z.infer<typeof ApiKeysSchema>;

export const GetApiKeySchema = ServerResponseSchema.extend({
  data: ApiKeySchema,
});

export const GetApiKeysSchema = ServerResponseSchema.extend({
  data: ApiKeysSchema,
});

export const CreateApiKeyFormSchema = z.object({
  name: z.string().min(3),
  description: z.string().optional(),
  access: z.enum(["read", "write", "full"]).default("read"),
  intent: z.enum(["create", "update", "delete"]).default("create"),
  id: z.string().optional(),
});

export type CreateApiKeyForm = z.infer<typeof CreateApiKeyFormSchema>;
