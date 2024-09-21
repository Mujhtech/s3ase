import { z } from "zod";

const EnvironmentSchema = z.object({
  NODE_ENV: z.union([
    z.literal("development"),
    z.literal("production"),
    z.literal("test"),
  ]),
  SESSION_SECRET: z.string().default("s3cr3t"),
  BACKEND_URL: z.string().default("http://localhost:5555"),
});

export type Environment = z.infer<typeof EnvironmentSchema>;
export const env = EnvironmentSchema.parse(process.env);
