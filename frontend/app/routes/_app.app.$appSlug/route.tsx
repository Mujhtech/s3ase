import { Outlet, useLocation } from "@remix-run/react";
import { parseWithZod } from "@conform-to/zod";
import { ActionFunction, LoaderFunctionArgs, redirect } from "@remix-run/node";
import React, { useCallback, useMemo } from "react";
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
import { useSSE } from "~/hooks/use-sse";
import { env } from "~/env.server";
import { getAuthTokenFromSession } from "~/services/auth.server";
import FilesUploadState, {
  FileUploadProgress,
} from "~/components/file/files-upload-state";

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

  const backendUrl = env.BACKEND_URL;

  const accessToken = await getAuthTokenFromSession(request);

  const user = await requireUser(request);

  const { appSlug } = AppSlugParamSchema.parse(params);

  const apps = await getApps(request);

  return typedjson({
    user,
    appSlug,
    backendUrl,
    apps,
    accessToken,
    app: apps.find((app) => app.slug === appSlug),
  });
};

export default function App() {
  const data = useTypedLoaderData<typeof loader>();
  const [open, setOpen] = React.useState(false);
  const [files, setFiles] = React.useState<Array<FileUploadProgress>>([]);
  const DEFAULT_EVENTS = [
    "upload_started",
    "upload_progress",
    "upload_completed",
  ];

  const location = useLocation();

  useSSE({
    baseUrl: data.backendUrl,
    shouldRun: true,
    appId: data.app?.id,
    accessToken: data.accessToken,
    path: "/api/ui/sse",
    events: useMemo(() => DEFAULT_EVENTS, []),
    onEvent: useCallback(
      (type: string, data: any) => {
        // check if file already exists using id, if exist update data else add new file

        if (data && DEFAULT_EVENTS.includes(type)) {
          setFiles((prev) => {
            const file = prev.find((file) => file.id === data.id);
            if (file) {
              return prev.map((file) =>
                file.id === data.id ? { ...file, ...data } : file
              );
            } else {
              return [...prev, data];
            }
          });

          if (open == false) {
            setOpen(true);
          }
        }
      },
      [open, files]
    ),
  });

  return (
    <div className="h-full w-full grid grid-rows-1 overflow-hidden">
      <div className="grid grid-cols-[5rem_1fr] overflow-hidden">
        <SideMenu user={data.user} appSlug={data.appSlug} />
        <div className="grid grid-rows-1 overflow-hidden">
          <div
            className={cn(
              "w-full relative",
              location.pathname != filesPath(data.appSlug) &&
                "p-3 overflow-y-auto scrollbar-thin scrollbar-track-transparent scrollbar-thumb-black/60"
            )}
          >
            <Outlet />

            <FilesUploadState files={files} open={open} setOpen={setOpen} />
          </div>
        </div>
      </div>
    </div>
  );
}
