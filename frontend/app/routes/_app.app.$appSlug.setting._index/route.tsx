import { parseWithZod } from "@conform-to/zod";
import { ActionFunction, json, LoaderFunctionArgs } from "@remix-run/node";
import React from "react";
import { redirect, typedjson, useTypedLoaderData } from "remix-typedjson";
import AppDomainForm from "~/components/app/app-domain-form";
import { CreateAppForm } from "~/components/app/create-app-dialog";
import DeleteAppDialog from "~/components/app/delete-app-dialog";
import { Button } from "~/components/ui/button";
import { Card, CardContent, CardHeader } from "~/components/ui/card";
import Paragraph from "~/components/ui/paragraph";
import { useApp } from "~/hooks/use-apps";
import { AppSlugParamSchema } from "../_app.app.$appSlug/route";
import { getAppIdFromSession, updateApp } from "~/services/app.server";
import { CreateAppFormSchema } from "~/models/app";
import { settingsPath } from "~/lib/path";
import { requestUrl } from "~/services/request-url.server";
import { createOrUpdateDomain, getDomain } from "~/services/domain.server";
import { CreateOrUpdateDomainFormSchema, Domain } from "~/models/domain";
import { Badge } from "~/components/ui/badge";

export const action: ActionFunction = async ({ request, params }) => {
  const { appSlug } = AppSlugParamSchema.parse(params);

  const formData = await request.formData();

  const type = formData.get("type");

  switch (type) {
    case "app":
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
    case "domain":
      const domainSubmission = parseWithZod(formData, {
        schema: CreateOrUpdateDomainFormSchema,
      });

      if (domainSubmission.status !== "success") {
        return json(domainSubmission.reply());
      }

      try {
        await createOrUpdateDomain(request, domainSubmission.value);

        return redirect(settingsPath(appSlug), {});
      } catch (e) {
        return json(domainSubmission.reply());
      }

    default:
      return redirect(settingsPath(appSlug));
  }
};

export const loader = async ({ request, params }: LoaderFunctionArgs) => {
  const { appSlug } = AppSlugParamSchema.parse(params);

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
