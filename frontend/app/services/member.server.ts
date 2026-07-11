import { ServerResponseSchema } from "~/models/default";
import { GetMembersSchema } from "~/models/member";
import { api } from "./api.server";

export async function getMembers(request: Request) {
  const res = await api.get({
    request,
    path: "/ui/members",
    schema: GetMembersSchema,
  });

  return res.data;
}

export async function createMember(
  request: Request,
  data: { email: string; role: "member" },
) {
  return api.post({
    request,
    path: "/ui/members",
    body: data,
    schema: ServerResponseSchema,
  });
}

export async function removeMember(request: Request, memberId: string) {
  return api.delete({
    request,
    path: `/ui/members/${memberId}`,
    schema: ServerResponseSchema,
  });
}
