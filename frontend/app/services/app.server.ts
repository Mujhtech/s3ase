import { z } from "zod";
import { api } from "./api.server";
import {
  Apps,
  AppsSchema,
  CreateAppForm,
  CreateAppResponseSchema,
  GetApps,
  GetAppsSchema,
} from "~/models/app";
import { createCookieSessionStorage } from "@remix-run/node";
import { env } from "~/env.server";
import { ServerResponseSchema } from "~/models/default";

const { commitSession, getSession, destroySession } =
  createCookieSessionStorage({
    cookie: {
      name: "__app_session", // use any name you want here
      sameSite: "lax", // this helps with CSRF
      path: "/", // remember to add this so the cookie will work in all routes
      httpOnly: true, // for security reasons, make this cookie http only
      secrets: [env.SESSION_SECRET], // replace this with an actual secret
      secure: process.env.NODE_ENV === "production", // enable this in prod only
    },
  });

export { commitSession };

export function getAppSession(request: Request) {
  return getSession(request.headers.get("Cookie"));
}

export async function getAppIdFromSession(request: Request) {
  const session = await getAppSession(request);

  const appId = session.get("app-id");

  if (appId) {
    return appId as string;
  }

  return null;
}

export async function setAppSession(request: Request, token: string) {
  const session = await getAppSession(request);

  if (session) {
    session.set("app-id", token);
  }

  return session;
}

export async function clearAppSession(request: Request) {
  const session = await getAppSession(request);

  if (session) {
    session.unset("app-id");
  }

  return session;
}

export async function getApps(request: Request) {
  const res = await api.get<GetApps>({
    request,
    path: "/ui/apps",
    schema: GetAppsSchema,
  });

  return res.data;
}

export async function createApp(request: Request, data: CreateAppForm) {
  return await api.post({
    request,
    path: "/ui/apps",
    body: data,
    schema: CreateAppResponseSchema,
  });
}

export async function updateApp(
  request: Request,
  id: string,
  data: CreateAppForm
) {
  return await api.put({
    request,
    path: `/ui/apps/${id}`,
    body: data,
    schema: ServerResponseSchema,
  });
}
