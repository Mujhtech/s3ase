import { createCookie, type LoaderFunction, redirect } from "@remix-run/node";
import {
  clearAuthSession,
  commitSession,
  setAuthSession,
} from "~/services/auth.server";

export const loader: LoaderFunction = async ({ request }) => {
  const session = await clearAuthSession(request);
  return redirect("/login", {
    headers: {
      "Set-Cookie": await commitSession(session),
    },
  });
};
