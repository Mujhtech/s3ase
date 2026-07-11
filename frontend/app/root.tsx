import type { LinksFunction,LoaderFunctionArgs } from "@remix-run/node";
import {
Links,
Meta,
Outlet,
Scripts,
ScrollRestoration,
} from "@remix-run/react";

import { typedjson } from "remix-typedjson";
import { publicBackendUrl } from "./env.server";
import { getAppIdFromSession } from "./services/app.server";
import { getAuthTokenFromSession } from "./services/auth.server";
import { getFeatures } from "./services/feature.server";
import { getUser } from "./services/user.server";
import "./tailwind.css";

export const links: LinksFunction = () => [];

export const loader = async ({ request }: LoaderFunctionArgs) => {
  const backendUrl = publicBackendUrl;
  const [feature, accessToken, appId, user] = await Promise.all([
    getFeatures(request),
    getAuthTokenFromSession(request),
    getAppIdFromSession(request),
    getUser(request),
  ]);

  return typedjson({
    user: user,
    feature: feature,
    accessToken: accessToken,
    appId: appId,
    backendUrl,
  });
};

export type RootLoaderType = typeof loader;

export function Layout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="en" className="h-full" suppressHydrationWarning>
      <head>
        <meta charSet="utf-8" />
        <meta name="viewport" content="width=device-width, initial-scale=1" />
        <Meta />
        <Links />
      </head>
      <body
        className="h-full overflow-hidden bg-background text-foreground antialiased !m-0"
        suppressHydrationWarning
      >
        {children}
        <ScrollRestoration />
        <Scripts />
      </body>
    </html>
  );
}

export default function App() {
  return <Outlet />;
}
