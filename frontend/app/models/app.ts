import { z } from "zod";
import { ServerResponse, ServerResponseSchema } from "./default";

export const AppSchema = z.object({
  id: z.string(),
  name: z.string(),
  owner_id: z.string(),
  bucket: z.string(),
  region: z.string(),
  domain: z.string(),
  slug: z.string(),
  description: z.string(),
  created_at: z.string(),
  updated_at: z.string(),
});

export type App = z.infer<typeof AppSchema>;

export const AppsSchema = z.array(AppSchema);

export type Apps = z.infer<typeof AppsSchema>;

export const GetAppsSchema = ServerResponseSchema.extend({
  data: AppsSchema,
});

export type GetApps = z.infer<typeof GetAppsSchema>;

export const CreateAppFormSchema = z.object({
  name: z.string().min(3),
  slug: z.string().optional(),
  description: z.string().optional(),
  domain: z.string().optional(),
  region: z.string(),
  intent: z.enum(["create", "update"]),
});

export type CreateAppForm = z.infer<typeof CreateAppFormSchema>;

export const CreateAppResponseSchema = ServerResponseSchema.extend({
  data: AppSchema,
});

export type CreateAppResponse = z.infer<typeof CreateAppResponseSchema>;
