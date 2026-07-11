import { ServerResponseSchema } from "~/models/default";
import {
CreateWebhookForm,
GetWebhookSchema,
GetWebhooksSchema,
} from "~/models/webhook";
import { api } from "./api.server";

export async function getWebhooks(request: Request) {
  const res = await api.get({
    request,
    path: "/ui/webhooks",
    schema: GetWebhooksSchema,
  });

  return res.data;
}

export async function createWebhook(request: Request, body: CreateWebhookForm) {
  const res = await api.post({
    request,
    path: "/ui/webhooks",
    body: body,
    schema: GetWebhookSchema,
  });

  return res.data;
}

export async function updateWebhook(
  request: Request,
  id: string,
  body: CreateWebhookForm
) {
  const res = await api.put({
    request,
    path: `/ui/webhooks/${id}`,
    body: body,
    schema: ServerResponseSchema,
  });

  return res.message;
}

export async function deleteWebhook(request: Request, id: string) {
  const res = await api.delete({
    request,
    path: `/ui/webhooks/${id}`,
    body: {},
    schema: ServerResponseSchema,
  });

  return res.message;
}
