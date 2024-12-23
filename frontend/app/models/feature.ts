import { z } from "zod";
import { ServerResponseSchema } from "./default";

export const FeatureSchema = z.object({
  name: z.string(),
  description: z.string(),
  is_github_auth_enabled: z.boolean(),
  is_google_auth_enabled: z.boolean(),
  is_aws_configured: z.boolean(),
  version: z.string(),
});

export type Feature = z.infer<typeof FeatureSchema>;

export const GetFeatureSchema = ServerResponseSchema.extend({
  data: FeatureSchema,
});
