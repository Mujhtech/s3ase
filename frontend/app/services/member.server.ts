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
