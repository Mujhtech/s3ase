import { z } from "zod";
import { ServerResponseSchema } from "./default";

export const DomainSchema = z.object({
  id: z.string(),
  domain: z.string(),
  txt_record: z.string(),
  cname_record: z.string(),
  status: z.enum(["pending", "review", "failed", "verified"]),
  created_at: z.string(),
  updated_at: z.string(),
});

export type Domain = z.infer<typeof DomainSchema>;

export const GetDomainSchema = ServerResponseSchema.extend({
  data: DomainSchema,
});

export const CreateOrUpdateDomainFormSchema = z.object({
  domain: z.string().url(),
  intent: z.enum(["create", "update", "refresh"]).default("create"),
});

export type CreateOrUpdateDomainForm = z.infer<
  typeof CreateOrUpdateDomainFormSchema
>;
