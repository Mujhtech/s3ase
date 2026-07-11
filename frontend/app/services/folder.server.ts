import { ServerResponseSchema } from "~/models/default";
import {
CreateFolderForm,
GetFolderSchema,
GetFoldersSchema,
} from "~/models/folder";
import { api } from "./api.server";

export async function getFolders(request: Request) {
  const res = await api.get({
    request,
    path: "/ui/folders",
    schema: GetFoldersSchema,
  });

  return res.data;
}

export async function createFolder(request: Request, body: CreateFolderForm) {
  const res = await api.post({
    request,
    path: "/ui/folders",
    body: body,
    schema: GetFolderSchema,
  });

  return res.data;
}

export async function getFolder(request: Request, id: string) {
  const res = await api.get({
    request,
    path: `/ui/folders/${id}`,
    schema: GetFolderSchema,
  });

  return res.data;
}

export async function updateFolder(
  request: Request,
  id: string,
  body: CreateFolderForm
) {
  const res = await api.put({
    request,
    path: `/ui/folders/${id}`,
    body: body,
    schema: ServerResponseSchema,
  });

  return res.message;
}

export async function deleteFolder(request: Request, id: string) {
  const res = await api.delete({
    request,
    path: `/ui/folders/${id}`,
    body: {},
    schema: ServerResponseSchema,
  });

  return res.message;
}
