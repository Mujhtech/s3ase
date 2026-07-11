import { parseWithZod } from "@conform-to/zod";
import {
ActionFunction,
json,
redirect
} from "@remix-run/node";
import { folderPath } from "~/lib/path";
import { CreateFileFormSchema } from "~/models/file";
import { uploadFile } from "~/services/file.server";
import { AppSlugParamSchema } from "./_app.app.$appSlug/route";

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
