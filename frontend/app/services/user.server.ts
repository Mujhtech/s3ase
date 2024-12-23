import { User, UserSchema } from "~/models/user";
import { api } from "./api.server";
import { getAuthSession, logout } from "./auth.server";
import { redirect } from "remix-typedjson";

export async function getUser(request: Request) {
  const session = await getAuthSession(request);

  const token = session.get("token");

  if (token == undefined || token == null) {
    return null;
  }

  const user = await api.get<User>({
    request,
    path: "/ui/user",
    schema: UserSchema,
  });

  if (user) {
    return user;
  }

  throw await logout(request);
}

export async function requireUser(request: Request) {
  const user = await getUser(request);

  if (user) {
    return user;
  }

  throw await logout(request);
}
