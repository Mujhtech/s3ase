import { CreateOrUpdateDomainForm, GetDomainSchema } from "~/models/domain";
import { api } from "./api.server";

export async function getDomain(request: Request) {
  const res = await api.get({
    request,
    path: `/ui/domain`,
    schema: GetDomainSchema,
  });

  return res.data;
}

export async function createOrUpdateDomain(
  request: Request,
  body: CreateOrUpdateDomainForm
) {
  const res = await api.post({
    request,
    path: "/ui/domain",
    body: body,
    schema: GetDomainSchema,
  });

  return res.data;
}
