import { UIMatch } from "@remix-run/react";
import { useTypedMatchesData } from "./use-typed-match";
import type { loader as appLoader } from "~/routes/_app.app.$appSlug/route";
import invariant from "tiny-invariant";

export function useOptionalApps(matches?: UIMatch[]) {
  const data = useTypedMatchesData<typeof appLoader>({
    id: "routes/_app.app.$appSlug",
    matches,
  });
  return data?.apps;
}

export function useApps(matches?: UIMatch[]) {
  const apps = useOptionalApps(matches);
  invariant(apps, "No apps found in loader.");
  return apps;
}

export function useOptionalApp(matches?: UIMatch[]) {
  const data = useTypedMatchesData<typeof appLoader>({
    id: "routes/_app.app.$appSlug",
    matches,
  });
  return data?.app;
}

export function useApp(matches?: UIMatch[]) {
  const app = useOptionalApp(matches);
  invariant(app, "No app found in loader.");
  return app;
}
