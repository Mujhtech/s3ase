import { z } from "zod";
import { ServerResponse, ServerResponseSchema } from "./default";

export const MediaUploadSchema = z.object({
  type: z.enum(["file", "folder"]),
});

export const FileSchema = z.object({
  id: z.string(),
  name: z.string(),
  created_at: z.string(),
  updated_at: z.string(),
});

export type File = z.infer<typeof FileSchema>;

export const FilesSchema = z.array(FileSchema);

export type Files = z.infer<typeof FilesSchema>;

export const GetFileSchema = ServerResponseSchema.extend({
  data: FileSchema,
});

export const GetFilesQuerySchema = z.object({
  folder_id: z.string().optional(),
  search: z.string().optional(),
  page: z.number().optional(),
  per_page: z.number().optional(),
});

export type GetFilesQuery = z.infer<typeof GetFilesQuerySchema>;

export const GetFilesSchema = ServerResponseSchema.extend({
  data: FilesSchema,
});

export const CreateFileFormSchema = MediaUploadSchema.extend({
  name: z.string().min(3),
  description: z.string().optional(),
  intent: z.enum(["create", "update", "delete"]).default("create"),
  id: z.string().optional(),
});

export type CreateFileForm = z.infer<typeof CreateFileFormSchema>;
