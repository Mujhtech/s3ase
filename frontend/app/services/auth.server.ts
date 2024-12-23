import { redirect } from "@remix-run/node";

import { createCookieSessionStorage } from "@remix-run/node";
import { env } from "~/env.server";

// export the whole sessionStorage object
const { commitSession, getSession, destroySession } =
  createCookieSessionStorage({
    cookie: {
      name: "__auth_session", // use any name you want here
      sameSite: "lax", // this helps with CSRF
      path: "/", // remember to add this so the cookie will work in all routes
      httpOnly: true, // for security reasons, make this cookie http only
      secrets: [env.SESSION_SECRET], // replace this with an actual secret
      secure: process.env.NODE_ENV === "production", // enable this in prod only
    },
  });

export { commitSession };

export function getAuthSession(request: Request) {
  return getSession(request.headers.get("Cookie"));
}

export async function getAuthTokenFromSession(request: Request) {
  const session = await getAuthSession(request);

  return session.get("token");
}

export async function setAuthSession(request: Request, token: string) {
  const session = await getAuthSession(request);

  if (session) {
    session.set("token", token);
  }

  return session;
}

export async function clearAuthSession(request: Request) {
  const session = await getAuthSession(request);

  if (session) {
    session.unset("token");
  }

  return session;
}

export async function logout(request: Request) {
  return redirect("/logout");
}
