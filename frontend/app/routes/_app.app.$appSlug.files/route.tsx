import { parseWithZod } from "@conform-to/zod";
import {
  ActionFunction,
  json,
  LoaderFunctionArgs,
  redirect,
} from "@remix-run/node";

import React, { useState } from "react";
import CreateFolderDropdown from "~/components/file/create-folder-dropdown";
import { Button } from "~/components/ui/button";
import { AppSlugParamSchema } from "../_app.app.$appSlug/route";
import { MediaUploadSchema } from "~/models/file";
import { filesPath } from "~/lib/path";
import { typedjson, useTypedLoaderData } from "remix-typedjson";
import {
  createFolder,
  deleteFolder,
  getFolders,
  updateFolder,
} from "~/services/folder.server";
import { useApp } from "~/hooks/use-apps";
import { CreateFolderFormSchema } from "~/models/folder";
import FolderCard from "~/components/file/folder-card";
import { getFiles, uploadFile } from "~/services/file.server";
import CreateFolderDialog from "~/components/file/create-folder-dialog";
import FileCard from "~/components/file/file-card";
import DragAndDropArea from "~/components/file/drag-and-drop-area";
import FilesPageContext from "~/components/file/files-page-context";
import Paragraph from "~/components/ui/paragraph";
import FileLayout from "~/components/file/file-layout";

export const action: ActionFunction = async ({ request, params }) => {
  const formData = await request.formData();
  const submission = parseWithZod(formData, { schema: MediaUploadSchema });

  if (submission.status !== "success") {
    return json(submission.reply());
  }

  try {
    const { appSlug } = AppSlugParamSchema.parse(params);

    switch (submission.value.type) {
      case "folder":
        const submission = parseWithZod(formData, {
          schema: CreateFolderFormSchema,
        });

        if (submission.status !== "success") {
          return json(submission.reply());
        }

        switch (submission.value.intent) {
          case "create":
            await createFolder(request, submission.value);
            break;
          case "delete":
            await deleteFolder(request, submission.value.id!);
            break;
          case "update":
            await updateFolder(request, submission.value.id!, submission.value);
            break;
        }

        break;
      case "file":
        break;
    }

    return redirect(filesPath(appSlug), {});
  } catch (e) {
    return json(submission.reply());
  }
};

export const loader = async ({ request, params }: LoaderFunctionArgs) => {
  const { appSlug } = AppSlugParamSchema.parse(params);

  const folders = await getFolders(request);

  const files = await getFiles(request);

  return typedjson({
    folders,
    files,
  });
};

export default function Page() {
  const { folders, files } = useTypedLoaderData<typeof loader>();
  const app = useApp();
  const [layout, setLayout] = useState<"list" | "grid">("grid");

  return (
    <div className="flex flex-col">
      <div className="flex m-3 mb-3 items-center border-border border-b pt-4 pb-2 justify-between">
        <h1 className="text-2xl font-bold">Files</h1>
        {/* <CreateFolderDropdown /> */}
        <div className="flex items-center gap-2">
          <FileLayout layout={layout} setLayout={setLayout} />
          <CreateFolderDialog />
        </div>
      </div>
      <DragAndDropArea>
        <div className="relative min-h-[80vh] h-full ">
          <FilesPageContext>
            <div className="m-3 h-full overflow-y-auto scrollbar-thin scrollbar-track-transparent scrollbar-thumb-black/60">
              <div className="">
                <h4 className="text-sm font-medium mb-3">Folders</h4>
                {folders.length > 0 ? (
                  <div className="grid grid-cols-3 md:grid-cols-4 lg:grid-cols-6 xl:grid-cols-7 2xl:grid-cols-8 gap-3">
                    {folders.map((folder, i) => (
                      <FolderCard key={i} folder={folder} appSlug={app.slug} />
                    ))}
                  </div>
                ) : (
                  <div className="flex">
                    <Paragraph>
                      No folders found, start by creating one.
                    </Paragraph>
                  </div>
                )}
              </div>
              <div className="mt-6">
                <h4 className="text-sm font-medium mb-3">Files</h4>
                {files.length > 0 ? (
                  <div className="grid grid-cols-3 md:grid-cols-4 lg:grid-cols-6 xl:grid-cols-7 2xl:grid-cols-8 gap-3">
                    {files.map((file, i) => (
                      <FileCard key={i} file={file} />
                    ))}
                  </div>
                ) : (
                  <div className="flex">
                    <Paragraph>
                      No files found, start by uploading a file.
                    </Paragraph>
                  </div>
                )}
              </div>
            </div>
          </FilesPageContext>
        </div>
      </DragAndDropArea>
    </div>
  );
}
