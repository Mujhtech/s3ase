import { parseWithZod } from "@conform-to/zod";
import { ActionFunction,LoaderFunctionArgs,redirect } from "@remix-run/node";
import { Outlet,useLocation,useRevalidator } from "@remix-run/react";
import React,{ useCallback } from "react";
import { typedjson,useTypedLoaderData } from "remix-typedjson";
import invariant from "tiny-invariant";
import { z } from "zod";
import FilesUploadState,{
FileUploadProgress,
} from "~/components/file/files-upload-state";
import { UploadManagerProvider } from "~/components/file/upload-manager";
import SideMenu from "~/components/layout/side-menu";
import { publicBackendUrl } from "~/env.server";
import { useSSE } from "~/hooks/use-sse";
import { appPath,appsPath,filesPath } from "~/lib/path";
import { cn } from "~/lib/utils";
import {
commitSession,
getAppIdFromSession,
getApps,
setAppSession,
} from "~/services/app.server";
import { getAuthTokenFromSession } from "~/services/auth.server";
import { requireUser } from "~/services/user.server";

export const AppSlugParamSchema = z.object({
  appSlug: z.string(),
});

const SwitchAppSchema = z.object({
  id: z.string(),
});

const UPLOAD_EVENTS = new Set([
  "upload_started",
  "upload_progress",
  "upload_completed",
  "upload_failed",
  "upload_cancelled",
  "upload_deleted",
]);

const UPLOAD_EVENT_NAMES = Array.from(UPLOAD_EVENTS);

function isUploadProgress(value: unknown): value is FileUploadProgress {
  if (!value || typeof value !== "object") return false;
  const upload = value as Partial<FileUploadProgress>;
  return typeof upload.id === "string" && typeof upload.name === "string";
}

export const action: ActionFunction = async ({ request, params }) => {
  const formData = await request.formData();

  const submission = parseWithZod(formData, { schema: SwitchAppSchema });

  if (submission.status !== "success") {
    return redirect(appsPath());
  }

  try {
    const { appSlug } = AppSlugParamSchema.parse(params);

    const session = await setAppSession(request, submission.value.id);

    const headers = new Headers({ "Set-Cookie": await commitSession(session) });

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

  const backendUrl = publicBackendUrl;

  const { appSlug } = AppSlugParamSchema.parse(params);
  const [accessToken, user, apps] = await Promise.all([
    getAuthTokenFromSession(request),
    requireUser(request),
    getApps(request),
  ]);

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

  const location = useLocation();
  const revalidator = useRevalidator();

  const handleUploadChange = useCallback(
    (eventData: FileUploadProgress) => {
      setFiles((prev) => {
        const existing = prev.some((file) => file.id === eventData.id);
        return existing
          ? prev.map((file) =>
              file.id === eventData.id ? { ...file, ...eventData } : file
            )
          : [...prev, eventData];
      });
      setOpen(true);
    },
    []
  );

  useSSE({
    baseUrl: data.backendUrl,
    shouldRun: true,
    appId: data.app?.id,
    accessToken: data.accessToken,
    path: "/api/ui/sse",
    events: UPLOAD_EVENT_NAMES,
    onEvent: useCallback(
    (type: string, eventData: unknown) => {
      if (UPLOAD_EVENTS.has(type) && isUploadProgress(eventData)) {
      handleUploadChange(eventData);
      if (type === "upload_completed" || type === "upload_deleted") {
        revalidator.revalidate();
      }
      }
    },
    [handleUploadChange, revalidator]
    ),
  });

  return (
    <UploadManagerProvider onUploadChange={handleUploadChange}>
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
    </UploadManagerProvider>
  );
}
