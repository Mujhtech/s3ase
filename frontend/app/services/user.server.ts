import { GetUserSchema } from "~/models/user";
import { api } from "./api.server";
import { getAuthSession,logout } from "./auth.server";

export async function getUser(request: Request) {
  const session = await getAuthSession(request);

  const token = session.get("token");

  if (token == undefined || token == null) {
    return null;
  }

  const response = await api.get({
    request,
    path: "/ui/user",
    schema: GetUserSchema,
  });

  if (response.data) {
    return response.data;
  }

  throw await logout();
}

export async function requireUser(request: Request) {
  const user = await getUser(request);

  if (user) {
    return user;
  }

  throw await logout();
}
