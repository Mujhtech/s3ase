import { Card,CardContent } from "~/components/ui/card";
import {
Table,
TableBody,
TableCell,
TableHead,
TableHeader,
TableRow,
} from "~/components/ui/table";

import { parseWithZod } from "@conform-to/zod";
import { ActionFunction,json,LoaderFunctionArgs } from "@remix-run/node";
import { redirect,typedjson,useTypedLoaderData } from "remix-typedjson";
import CreateApiKeyDialog from "~/components/api-key/create-api-key-dialog";
import { apiKeysPath } from "~/lib/path";
import { CreateApiKeyFormSchema } from "~/models/api_key";
import {
createApiKey,
deleteApiKey,
getApiKeys,
updateApiKey,
} from "~/services/api_key.server";
import { AppSlugParamSchema } from "../_app.app.$appSlug/route";

export const action: ActionFunction = async ({ request, params }) => {
  const formData = await request.formData();
  const submission = parseWithZod(formData, { schema: CreateApiKeyFormSchema });

  if (submission.status !== "success") {
    return json(submission.reply());
  }

  try {
    const { appSlug } = AppSlugParamSchema.parse(params);

    switch (submission.value.intent) {
      case "create": {
        const created = await createApiKey(request, submission.value);
        return json({ secret: created.secret });
      }
      case "delete":
        await deleteApiKey(request, submission.value.id!);
        break;
      case "update":
        await updateApiKey(request, submission.value.id!, submission.value);
        break;
    }

    return redirect(apiKeysPath(appSlug), {});
  } catch (error) {
    return json({ error: error instanceof Error ? error.message : "Unable to save API key" });
  }
};

export const loader = async ({ request, params }: LoaderFunctionArgs) => {
  AppSlugParamSchema.parse(params);

  const apiKeys = await getApiKeys(request);

  return typedjson({
    apiKeys,
  });
};

export default function Page() {
  const { apiKeys } = useTypedLoaderData<typeof loader>();

  return (
    <div className="flex flex-col w-full">
      <div className="mb-4 flex items-center border-border border-b pt-4 pb-2 justify-between">
        <h1 className="text-2xl font-bold">API Keys</h1>
        <CreateApiKeyDialog />
      </div>
      <Card className="mb-8 w-full max-w-5xl mx-auto">
        <CardContent>
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead className="">Name</TableHead>
                <TableHead className="w-[80px]">Access</TableHead>
                <TableHead className="w-[100px]">Last Used</TableHead>
                <TableHead className="w-[50px]">Action</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {apiKeys.length > 0 ? (
                apiKeys.map((apiKey) => {
                  return (
                    <TableRow key={apiKey.id}>
                      <TableCell className="font-medium py-2">
            <div className="flex flex-col">
              <span>{apiKey.name}</span>
              {apiKey.key_prefix && (
              <span className="font-mono text-xs text-muted-foreground">
                {apiKey.key_prefix}…{apiKey.last_four}
              </span>
              )}
            </div>
                      </TableCell>
                      <TableCell className="capitalize py-2">
                        {apiKey.access}
                      </TableCell>
                      <TableCell className="capitalize py-2">
                        {apiKey.last_used || "Never"}
                      </TableCell>

                      <TableCell className="py-2">
                        <CreateApiKeyDialog apiKey={apiKey} />
                      </TableCell>
                    </TableRow>
                  );
                })
              ) : (
                <TableRow className="">
                  <TableCell
                    className="font-medium text-center py-2"
                    colSpan={4}
                  >
                    No api keys found, create api key to get started
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
