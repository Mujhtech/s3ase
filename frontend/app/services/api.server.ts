import { z } from "zod";
import { env } from "~/env.server";
import { getAppIdFromSession } from "./app.server";
import { getAuthTokenFromSession } from "./auth.server";

type QueryValue = string | number | boolean | null | undefined;

type ClientRequestOptions<T> = {
  request: Request;
  path: string;
  body?: unknown;
  query?: Record<string, QueryValue>;
  schema?: z.ZodType<T, z.ZodTypeDef, unknown>;
  headers?: HeadersInit;
};

export class ApiError extends Error {
  constructor(
  message: string,
  readonly status: number,
  readonly payload: unknown,
  ) {
  super(message);
  this.name = "ApiError";
  }
}

async function parseResponse<T>(
  response: Response,
  schema?: z.ZodType<T, z.ZodTypeDef, unknown>,
): Promise<T> {
  const contentType = response.headers.get("content-type") ?? "";
  const payload: unknown = contentType.includes("application/json")
  ? await response.json()
  : await response.text();

  if (!response.ok) {
  const message =
    payload && typeof payload === "object" && "error" in payload
    ? String(payload.error)
    : `API request failed with status ${response.status}`;
  throw new ApiError(message, response.status, payload);
  }

  return schema ? schema.parse(payload) : (payload as T);
}

function buildQuery(query?: Record<string, QueryValue>) {
  const search = new URLSearchParams();
  for (const [key, value] of Object.entries(query ?? {})) {
  if (value !== undefined && value !== null) search.set(key, String(value));
  }
  return search.toString();
}

async function request<T>(
  method: "POST" | "GET" | "PUT" | "DELETE" | "HEAD",
  options: ClientRequestOptions<T>,
): Promise<T> {
  const query = buildQuery(options.query);
  const path = query ? `${options.path}?${query}` : options.path;
  const response = await clientRequest({
  request: options.request,
  method,
  path,
  body: options.body,
  headers: options.headers,
  });
  return parseResponse(response, options.schema);
}

async function clientRequest(options: {
  request: Request;
  method: "POST" | "GET" | "PUT" | "DELETE" | "HEAD";
  path: string;
  body?: unknown;
  headers?: HeadersInit;
}) {
  const [accessToken, appId] = await Promise.all([
  getAuthTokenFromSession(options.request),
  getAppIdFromSession(options.request),
  ]);
  const headers = new Headers(options.headers);
  if (!headers.has("Content-Type")) headers.set("Content-Type", "application/json");
  if (accessToken) headers.set("Authorization", `Bearer ${accessToken}`);
  if (appId) headers.set("x-app-id", appId);

  return fetch(`${env.BACKEND_URL}${options.path}`, {
  method: options.method,
  headers,
  body:
    options.method === "GET" || options.method === "HEAD"
    ? undefined
    : JSON.stringify(options.body ?? {}),
  });
}

export const api = {
  get: <T>(options: ClientRequestOptions<T>) => request("GET", options),
  post: <T>(options: ClientRequestOptions<T>) => request("POST", options),
  put: <T>(options: ClientRequestOptions<T>) => request("PUT", options),
  delete: <T>(options: ClientRequestOptions<T>) => request("DELETE", options),
};
