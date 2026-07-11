import type { ActionFunctionArgs,LoaderFunctionArgs } from "@remix-run/node";
import { Form } from "@remix-run/react";
import { redirect,typedjson,useTypedLoaderData } from "remix-typedjson";
import { z } from "zod";
import { Button } from "~/components/ui/button";
import { Card,CardContent,CardHeader } from "~/components/ui/card";
import Paragraph from "~/components/ui/paragraph";
import { fileDownloadPath,filesPath,folderPath } from "~/lib/path";
import { deleteFile,getFile } from "~/services/file.server";
import { AppSlugParamSchema } from "../_app.app.$appSlug/route";

const FileRouteParams = AppSlugParamSchema.extend({ fileId: z.string().uuid() });

export async function loader({ request, params }: LoaderFunctionArgs) {
  const { appSlug, fileId } = FileRouteParams.parse(params);
  return typedjson({ appSlug, file: await getFile(request, fileId) });
}

export async function action({ request, params }: ActionFunctionArgs) {
  const { appSlug, fileId } = FileRouteParams.parse(params);
  const formData = await request.formData();
  if (formData.get("intent") !== "delete") return redirect(filesPath(appSlug));

  const file = await getFile(request, fileId);
  await deleteFile(request, fileId);
  return redirect(
  file.folder_id ? folderPath(appSlug, file.folder_id) : filesPath(appSlug),
  );
}

export default function Page() {
  const { appSlug, file } = useTypedLoaderData<typeof loader>();
  const size = new Intl.NumberFormat(undefined, {
  style: "unit",
  unit: "byte",
  unitDisplay: "narrow",
  }).format(file.size);

  return (
  <div className="flex flex-col w-full max-w-4xl">
    <Card className="mb-8">
    <CardHeader className="flex flex-col">
      <h1 className="text-2xl font-bold">{file.name}</h1>
      <Paragraph>{file.mime_type || "Unknown file type"}</Paragraph>
    </CardHeader>
    <CardContent className="flex flex-col gap-4">
      <div className="grid grid-cols-2 gap-3 text-sm">
      <Paragraph>Size</Paragraph>
      <Paragraph>{size}</Paragraph>
      <Paragraph>Status</Paragraph>
      <Paragraph className="capitalize">{file.status}</Paragraph>
      </div>
      <div className="flex gap-2">
      <Button asChild>
        <a href={fileDownloadPath(appSlug, file.id)}>Download</a>
      </Button>
      <Form method="post">
        <Button type="submit" name="intent" value="delete" variant="destructive">
        Delete
        </Button>
      </Form>
      </div>
    </CardContent>
    </Card>
  </div>
  );
}
