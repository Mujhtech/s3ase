import { z } from "zod";
import { ServerResponseSchema } from "./default";

export const DailyUsageSchema = z.object({
  date: z.string(),
  uploads: z.number(),
  storage_bytes: z.number(),
});

export const UsageSchema = z.object({
  file_count: z.number(),
  storage_bytes: z.number(),
  daily: z.array(DailyUsageSchema),
});

export type Usage = z.infer<typeof UsageSchema>;

export const GetUsageSchema = ServerResponseSchema.extend({ data: UsageSchema });
