import { z } from "zod";

export const UserSchema = z.object({
  id: z.string(),
  email: z.string(),
  email_verified: z.boolean(),
  avatar_url: z.string().optional(),
  name: z.string(),
  authentication_method: z.string().optional(),
  created_at: z.string().optional(),
  updated_at: z.string().optional(),
});

export type User = z.infer<typeof UserSchema>;
