import { Outlet, useLocation } from "@remix-run/react";
import { parseWithZod } from "@conform-to/zod";
import { ActionFunction, LoaderFunctionArgs, redirect } from "@remix-run/node";
import React from "react";
import { appPath, appsPath, filesPath } from "~/lib/path";
import {
  commitSession,
  getAppIdFromSession,
  setAppSession,
} from "~/services/app.server";
import { requireUser } from "~/services/user.server";
import { typedjson, useTypedLoaderData } from "remix-typedjson";
import { z } from "zod";
import SideMenu from "~/components/layout/side-menu";
import { getApps } from "~/services/app.server";
import invariant from "tiny-invariant";
import { cn } from "~/lib/utils";

export const AppSlugParamSchema = z.object({
  appSlug: z.string(),
});

const SwitchAppSchema = z.object({
  id: z.string(),
});

export const action: ActionFunction = async ({ request, params }) => {
  const formData = await request.formData();
  const submission = parseWithZod(formData, { schema: SwitchAppSchema });

  if (submission.status !== "success") {
    return redirect(appsPath());
  }

  try {
    const { appSlug } = AppSlugParamSchema.parse(params);

    const session = await setAppSession(request, submission.value.id);

    let headers = new Headers({ "Set-Cookie": await commitSession(session) });

    return redirect(appPath(appSlug), {
      headers,
    });
  } catch (e) {
    return redirect(appsPath());
  }
};

export const loader = async ({ request, params }: LoaderFunctionArgs) => {
  const appId = await getAppIdFromSession(request);

  invariant(appId, "No app found in session.");

  // if (!appId) {
  //   return redirect(appsPath());
  // }

  const user = await requireUser(request);

  const { appSlug } = AppSlugParamSchema.parse(params);

  const apps = await getApps(request);

  return typedjson({
    user,
    appSlug,
    apps,
    app: apps.find((app) => app.slug === appSlug),
  });
};

export default function App() {
  const data = useTypedLoaderData<typeof loader>();

  const location = useLocation();

  return (
    <div className="h-full w-full grid grid-rows-1 overflow-hidden">
      <div className="grid grid-cols-[5rem_1fr] overflow-hidden">
        <SideMenu user={data.user} appSlug={data.appSlug} />
        <div className="grid grid-rows-1 overflow-hidden">
          <div
            className={cn(
              "w-full",
              location.pathname != filesPath(data.appSlug) &&
                "p-3 overflow-y-auto scrollbar-thin scrollbar-track-transparent scrollbar-thumb-black/60"
            )}
          >
            <Outlet />
          </div>
        </div>
      </div>
    </div>
  );
}
