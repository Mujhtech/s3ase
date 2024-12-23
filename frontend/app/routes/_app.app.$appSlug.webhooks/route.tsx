import React from "react";
import { Card, CardContent, CardHeader } from "~/components/ui/card";
import {
  Table,
  TableBody,
  TableCell,
  TableRow,
  TableHeader,
  TableHead,
} from "~/components/ui/table";
import { ActionFunction, json, LoaderFunctionArgs } from "@remix-run/node";
import { redirect, typedjson, useTypedLoaderData } from "remix-typedjson";
import { useApp } from "~/hooks/use-apps";
import { AppSlugParamSchema } from "../_app.app.$appSlug/route";
import {
  createWebhook,
  deleteWebhook,
  getWebhooks,
  updateWebhook,
} from "~/services/webhook.server";
import CreateWebhookDialog from "~/components/webhook/create-webhook-dialog";
import { parseWithZod } from "@conform-to/zod";
import { CreateWebhookFormSchema } from "~/models/webhook";
import { webhooksPath } from "~/lib/path";

export const action: ActionFunction = async ({ request, params }) => {
  const formData = await request.formData();
  const submission = parseWithZod(formData, {
    schema: CreateWebhookFormSchema,
  });

  if (submission.status !== "success") {
    return json(submission.reply());
  }

  try {
    const { appSlug } = AppSlugParamSchema.parse(params);

    switch (submission.value.intent) {
      case "create":
        await createWebhook(request, submission.value);
        break;
      case "delete":
        await deleteWebhook(request, submission.value.id!);
        break;
      case "update":
        await updateWebhook(request, submission.value.id!, submission.value);
        break;
    }

    return redirect(webhooksPath(appSlug), {});
  } catch (e) {
    return json(submission.reply());
  }
};

export const loader = async ({ request, params }: LoaderFunctionArgs) => {
  const { appSlug } = AppSlugParamSchema.parse(params);

  const webhooks = await getWebhooks(request);

  return typedjson({
    webhooks,
  });
};

export default function Page() {
  const { webhooks } = useTypedLoaderData<typeof loader>();
  const app = useApp();

  return (
    <div className="flex flex-col">
      <div className="mb-4 flex items-center border-border border-b pt-4 pb-2 justify-between">
        <h1 className="text-2xl font-bold">Webhooks</h1>
        <CreateWebhookDialog />
      </div>
      <Card className="mb-8 w-full max-w-5xl mx-auto">
        <CardContent>
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead className="">Name</TableHead>
                <TableHead>Url</TableHead>
                <TableHead className="w-[50px]">Action</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {webhooks.length > 0 ? (
                webhooks.map((webhook) => {
                  return (
                    <TableRow key={webhook.id}>
                      <TableCell className="font-medium py-2">
                        {webhook.name}
                      </TableCell>
                      <TableCell className="py-2">{webhook.url}</TableCell>

                      <TableCell className="py-2">
                        <CreateWebhookDialog webhook={webhook} />
                      </TableCell>
                    </TableRow>
                  );
                })
              ) : (
                <TableRow className="">
                  <TableCell
                    className="font-medium text-center py-2"
                    colSpan={3}
                  >
                    No webhooks found, add a webhook to get started
                  </TableCell>
                </TableRow>
              )}
            </TableBody>
          </Table>
        </CardContent>
      </Card>
    </div>
  );
}
