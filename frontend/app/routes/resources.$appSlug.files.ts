import {
  ActionFunction,
  json,
  LoaderFunctionArgs,
  redirect,
  unstable_composeUploadHandlers as composeUploadHandlers,
  unstable_createMemoryUploadHandler as createMemoryUploadHandler,
  unstable_parseMultipartFormData as parseMultipartFormData,
} from "@remix-run/node";
import { AppSlugParamSchema } from "./_app.app.$appSlug/route";
import { CreateFileFormSchema } from "~/models/file";
import { parseWithZod } from "@conform-to/zod";
import { filesPath, folderPath } from "~/lib/path";
import { uploadFile } from "~/services/file.server";

export const action: ActionFunction = async ({ request, params }) => {
  const formData = await request.formData();
  const submission = parseWithZod(formData, {
    schema: CreateFileFormSchema,
  });

  if (submission.status !== "success") {
    return json(submission.reply());
  }

  try {
    const { appSlug } = AppSlugParamSchema.parse(params);

    await uploadFile(request, submission.value);

    if (submission.value.folder_id) {
      return redirect(folderPath(appSlug, submission.value.folder_id), {});
    }

    return true;
  } catch (e) {
    return json(submission.reply());
  }
};
