import {
  CreateApiKeyForm,
  GetApiKeySchema,
  GetApiKeysSchema,
} from "~/models/api_key";
import { api } from "./api.server";
import { ServerResponseSchema } from "~/models/default";

export async function getApiKeys(request: Request) {
  const res = await api.get({
    request,
    path: "/ui/api_keys",
    schema: GetApiKeysSchema,
  });

  return res.data;
}

export async function createApiKey(request: Request, body: CreateApiKeyForm) {
  const res = await api.post({
    request,
    path: "/ui/api_keys",
    body: body,
    schema: GetApiKeySchema,
  });

  return res.data;
}

export async function updateApiKey(
  request: Request,
  id: string,
  body: CreateApiKeyForm
) {
  const res = await api.put({
    request,
    path: `/ui/api_keys/${id}`,
    body: body,
    schema: ServerResponseSchema,
  });

  return res.message;
}

export async function deleteApiKey(request: Request, id: string) {
  const res = await api.delete({
    request,
    path: `/ui/api_keys/${id}`,
    body: {},
    schema: ServerResponseSchema,
  });

  return res.message;
}
