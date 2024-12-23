import { z } from "zod";
import { ServerResponse, ServerResponseSchema } from "./default";
import { MediaUploadSchema } from "./file";

export const FolderSchema = z.object({
  id: z.string(),
  name: z.string(),
  description: z.string(),
  url: z.string(),
  created_at: z.string(),
  updated_at: z.string(),
});

export type Folder = z.infer<typeof FolderSchema>;

export const FoldersSchema = z.array(FolderSchema);

export type Folders = z.infer<typeof FoldersSchema>;

export const GetFolderSchema = ServerResponseSchema.extend({
  data: FolderSchema,
});

export const GetFoldersSchema = ServerResponseSchema.extend({
  data: FoldersSchema,
});

export const CreateFolderFormSchema = MediaUploadSchema.extend({
  name: z.string().min(2),
  description: z.string().optional(),
  intent: z.enum(["create", "update", "delete"]).default("create"),
  id: z.string().optional(),
});

export type CreateFolderForm = z.infer<typeof CreateFolderFormSchema>;
