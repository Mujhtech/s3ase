import type { LoaderFunctionArgs } from "@remix-run/node";
import { z } from "zod";
import { env } from "~/env.server";
import { getAppIdFromSession } from "~/services/app.server";
import { getAuthTokenFromSession } from "~/services/auth.server";

const Params = z.object({ appSlug: z.string(), fileId: z.string().uuid() });

export async function loader({ request, params }: LoaderFunctionArgs) {
  const { fileId } = Params.parse(params);
  const [accessToken, appId] = await Promise.all([
  getAuthTokenFromSession(request),
  getAppIdFromSession(request),
  ]);
  if (!accessToken || !appId) return new Response("Unauthorized", { status: 401 });

  const upstream = await fetch(`${env.BACKEND_URL}/ui/files/${fileId}/content`, {
  headers: { Authorization: `Bearer ${accessToken}`, "x-app-id": appId },
  });
  if (!upstream.ok || !upstream.body) {
  return new Response("File is unavailable", { status: upstream.status });
  }

  const headers = new Headers();
  for (const name of ["Content-Type", "Content-Length", "Content-Disposition"]) {
  const value = upstream.headers.get(name);
  if (value) headers.set(name, value);
  }
  return new Response(upstream.body, { status: upstream.status, headers });
}
