import { z } from "zod";
import { ServerResponseSchema } from "./default";
import { UserSchema } from "./user";

export const MemberSchema = z.object({
  id: z.string(),
  user_id: z.string(),
  app_id: z.string(),
  user: UserSchema,
  role: z.string(),
  created_at: z.string(),
  updated_at: z.string(),
});

export type Member = z.infer<typeof MemberSchema>;

export const MembersSchema = z.array(MemberSchema);

export type Members = z.infer<typeof MembersSchema>;

export const GetMembersSchema = ServerResponseSchema.extend({
  data: MembersSchema,
});

export const SendInviteFormSchema = z.object({
  email: z.string().email(),
  role: z.literal("member").default("member"),
});

export const RemoveMemberFormSchema = z.object({
  memberId: z.string().min(1),
  intent: z.literal("remove"),
});
