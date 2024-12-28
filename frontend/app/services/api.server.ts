import { env } from "~/env.server";
import { getAuthSession, getAuthTokenFromSession } from "./auth.server";
import { z } from "zod";
import { getAppIdFromSession } from "./app.server";

type ClientRequestOptions<T> = {
  request: Request;
  path: string;
  body?: any;
  query?: any;
  schema?: z.ZodType<T>;
  headers?: any;
};

async function postRequest<T = any>(
  options: ClientRequestOptions<T>
): Promise<T> {
  return handleClientUnauthorized<T>(
    () => {
      return clientRequest({
        method: "POST",
        path: options.path,
        body: options.body,
        request: options.request,
        headers: options.headers,
      });
    },
    {
      schema: options.schema,
    }
  );
}

async function getRequest<T = any>(
  options: ClientRequestOptions<T>
): Promise<T> {
  return handleClientUnauthorized<T>(
    () => {
      const query = new URLSearchParams(options.query).toString();

      const path = query ? `${options.path}?${query}` : options.path;

      return clientRequest({
        request: options.request,
        headers: options.headers,
        method: "GET",
        path: path,
        body: {},
      });
    },
    {
      schema: options.schema,
    }
  );
}

async function putRequest<T = any>(
  options: ClientRequestOptions<T>
): Promise<T> {
  return handleClientUnauthorized<T>(
    () => {
      return clientRequest({
        request: options.request,
        headers: options.headers,
        method: "PUT",
        path: options.path,
        body: options.body,
      });
    },
    {
      schema: options.schema,
    }
  );
}

async function deleteRequest<T = any>(
  options: ClientRequestOptions<T>
): Promise<T> {
  return handleClientUnauthorized<T>(
    () => {
      return clientRequest({
        method: "DELETE",
        path: options.path,
        headers: options.headers,
        body: options.body,
        request: options.request,
      });
    },
    {
      schema: options.schema,
    }
  );
}

async function handleClientUnauthorized<T = any>(
  sendRequest: () => Promise<Response>,
  options: {
    schema?: z.ZodType<T>;
  }
): Promise<T> {
  let response = await sendRequest();

  if (response.status === 401) {
    response = await sendRequest();
  }

  if ([200, 201].includes(response.status) == false) {
    const body = await response.json();

    console.log(body);

    throw body;
  }

  const data = await response.json();

  console.log(data);

  if (options.schema) {
    try {
      return options.schema.parse(data.data);
    } catch (error) {
      return data as T;
    }
  }

  return data as T;
}

const clientRequest = async (options: {
  request: Request;
  method: "POST" | "GET" | "PUT" | "DELETE" | "HEAD";
  path: string;
  body: any;
  headers?: any;
}): Promise<Response> => {
  const accessToken = await getAuthTokenFromSession(options.request);

  const appId = await getAppIdFromSession(options.request);

  let headers = {
    "Content-Type": "application/json",
    Authorization: `Bearer ${accessToken}`,
    ...options.headers,
  };

  if (appId) {
    headers = {
      ...headers,
      "x-app-id": appId,
    };
  }

  const url = `${env.BACKEND_URL}${options.path}`;

  const response = fetch(url, {
    method: options.method,
    headers: headers,
    body:
      options.method == "GET" || options.method == "HEAD"
        ? undefined
        : JSON.stringify({
            ...options.body,
          }),
  });

  return response;
};

export const api = {
  post: postRequest,
  get: getRequest,
  put: putRequest,
  delete: deleteRequest,
};
