import { parseWithZod } from "@conform-to/zod";
import { ActionFunction,json,LoaderFunctionArgs } from "@remix-run/node";
import { redirect,typedjson,useTypedLoaderData } from "remix-typedjson";
import AppDomainForm from "~/components/app/app-domain-form";
import { CreateAppForm } from "~/components/app/create-app-dialog";
import DeleteAppDialog from "~/components/app/delete-app-dialog";
import { Badge } from "~/components/ui/badge";
import { Card,CardContent,CardHeader } from "~/components/ui/card";
import Paragraph from "~/components/ui/paragraph";
import { useApp } from "~/hooks/use-apps";
import { appsPath,settingsPath } from "~/lib/path";
import { CreateAppFormSchema,DeleteAppFormSchema } from "~/models/app";
import { CreateOrUpdateDomainFormSchema,Domain } from "~/models/domain";
import {
clearAppSession,
commitSession,
deleteApp,
getAppIdFromSession,
updateApp,
} from "~/services/app.server";
import { createOrUpdateDomain,getDomain,verifyDomain } from "~/services/domain.server";
import { AppSlugParamSchema } from "../_app.app.$appSlug/route";

export const action: ActionFunction = async ({ request, params }) => {
  const { appSlug } = AppSlugParamSchema.parse(params);

  const formData = await request.formData();

  const type = formData.get("type");

  switch (type) {
  case "delete-app": {
    const submission = parseWithZod(formData, { schema: DeleteAppFormSchema });
    if (submission.status !== "success") {
    return json(submission.reply());
    }
    try {
    const appId = await getAppIdFromSession(request);
    if (!appId) return redirect(appsPath());
    await deleteApp(request, appId, submission.value.name);
    const session = await clearAppSession(request);
    return redirect(appsPath(), {
      headers: { "Set-Cookie": await commitSession(session) },
    });
    } catch (error) {
    return json({ error: error instanceof Error ? error.message : "Unable to delete app" });
    }
  }
    case "app": {
      const appSubmission = parseWithZod(formData, {
        schema: CreateAppFormSchema,
      });

      if (appSubmission.status !== "success") {
        return json(appSubmission.reply());
      }

      try {
        const appId = await getAppIdFromSession(request);

        await updateApp(request, appId!, appSubmission.value);

        return redirect(settingsPath(appSlug), {});
      } catch (e) {
        return json(appSubmission.reply());
      }
    }
    case "domain": {
      const domainSubmission = parseWithZod(formData, {
        schema: CreateOrUpdateDomainFormSchema,
      });

      if (domainSubmission.status !== "success") {
        return json(domainSubmission.reply());
      }

      try {
    if (domainSubmission.value.intent === "refresh") {
      await verifyDomain(request);
    } else {
      await createOrUpdateDomain(request, domainSubmission.value);
    }

        return redirect(settingsPath(appSlug), {});
      } catch (e) {
        return json(domainSubmission.reply());
      }
    }

    default:
      return redirect(settingsPath(appSlug));
  }
};

export const loader = async ({ request, params }: LoaderFunctionArgs) => {
  AppSlugParamSchema.parse(params);

  let domain: Domain | null = null;

  try {
    domain = await getDomain(request);
  } catch (e) {
    //
  }

  return typedjson({
    domain,
  });
};

export default function Page() {
  const app = useApp();

  const { domain } = useTypedLoaderData<typeof loader>();

  return (
    <div className="flex flex-col w-full max-w-4xl">
      <Card className="mb-8">
        <CardHeader className="flex flex-col">
          <h1 className="font-semibold">App Information</h1>
        </CardHeader>
        <CardContent>
          <CreateAppForm app={app} />
        </CardContent>
      </Card>
      <Card className="mb-8">
        <CardHeader className="flex !flex-row justify-between items-center">
          <h1 className="font-semibold">Domain</h1>
          {domain && <Badge className="capitalize">{domain.status}</Badge>}
        </CardHeader>
        <CardContent>
          <AppDomainForm domain={domain} />
        </CardContent>
      </Card>
      <Card className="border-red-500 mb-8">
        <CardHeader className="flex flex-col">
          <h1 className="font-semibold">Delete App</h1>
          <Paragraph>Are you sure you want to delete this app?</Paragraph>
        </CardHeader>
        <CardContent>
          <div className="flex justify-between">
            <DeleteAppDialog app={app} />
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
