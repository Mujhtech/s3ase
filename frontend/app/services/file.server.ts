import { ServerResponseSchema } from "~/models/default";
import { api } from "./api.server";
import {
  CreateFileForm,
  GetFileSchema,
  GetFilesQuery,
  GetFilesSchema,
} from "~/models/file";

export async function getFiles(request: Request, query: GetFilesQuery = {}) {
  const res = await api.get({
    request,
    path: "/ui/files",
    schema: GetFilesSchema,
    query: query,
  });

  return res.data;
}

export async function getFile(request: Request, id: string) {
  const res = await api.get({
    request,
    path: `/ui/files/${id}`,
    schema: GetFileSchema,
  });

  return res.data;
}

export async function createFile(request: Request, body: CreateFileForm) {
  const res = await api.post({
    request,
    path: "/ui/files",
    body: body,
    schema: GetFileSchema,
  });

  return res.data;
}

export async function updateFile(
  request: Request,
  id: string,
  body: CreateFileForm
) {
  const res = await api.put({
    request,
    path: `/ui/files/${id}`,
    body: body,
    schema: ServerResponseSchema,
  });

  return res.message;
}

export async function deleteFile(request: Request, id: string) {
  const res = await api.delete({
    request,
    path: `/ui/files/${id}`,
    body: {},
    schema: ServerResponseSchema,
  });

  return res.message;
}
