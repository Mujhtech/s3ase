import { UIMatch } from "@remix-run/react";
import { loader } from "~/root";
import { useTypedMatchesData } from "./use-typed-match";
import { Feature } from "~/models/feature";

export function useOptionalFeature(matches?: UIMatch[]): Feature | undefined {
  const routeMatch = useTypedMatchesData<typeof loader>({
    id: "root",
    matches,
  });

  return routeMatch?.feature ?? undefined;
}

export function useFeature(matches?: UIMatch[]): Feature {
  const maybeFeature = useOptionalFeature(matches);
  if (!maybeFeature) {
    throw new Error(
      "No feature found in root loader, but feature is required by useFeature. If feature is optional, try useOptionalFeature instead."
    );
  }
  return maybeFeature;
}
