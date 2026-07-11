import { LoaderFunctionArgs,redirect } from "@remix-run/node";
import { appPath,appsPath } from "~/lib/path";
import { getApps } from "~/services/app.server";
import { requireUser } from "~/services/user.server";

export const loader = async ({ request }: LoaderFunctionArgs) => {
  const user = await requireUser(request);

  if (!user) {
    return redirect("/logout");
  }

  const apps = await getApps(request);

  if (apps.length === 0) {
    return redirect(appsPath());
  }

  return redirect(appPath(apps[0].slug));
};
