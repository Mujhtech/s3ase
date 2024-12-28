import { LoaderFunctionArgs } from "@remix-run/node";
import React, { useState } from "react";
import { typedjson, useTypedLoaderData } from "remix-typedjson";
import { getFiles, uploadFile } from "~/services/file.server";
import { getFolder } from "~/services/folder.server";
import { AppSlugParamSchema } from "../_app.app.$appSlug/route";
import { z } from "zod";
import { useApp } from "~/hooks/use-apps";
import FileLayout from "~/components/file/file-layout";
import { ChevronRight } from "lucide-react";
import FilesPageContext from "~/components/file/files-page-context";
import DragAndDropArea from "~/components/file/drag-and-drop-area";
import FileCard from "~/components/file/file-card";
import Paragraph from "~/components/ui/paragraph";

const AppWithFolderParamSchema = AppSlugParamSchema.extend({
  folderId: z.string(),
});

export const loader = async ({ request, params }: LoaderFunctionArgs) => {
  const { appSlug, folderId } = AppWithFolderParamSchema.parse(params);

  const folder = await getFolder(request, folderId);

  const files = await getFiles(request, {
    folder_id: folder.id,
  });

  return typedjson({
    folder,
    files,
  });
};

export default function Page() {
  const { folder, files } = useTypedLoaderData<typeof loader>();
  const app = useApp();
  const [layout, setLayout] = useState<"list" | "grid">("grid");

  return (
    <div className="flex flex-col">
      <div className="flex m-3 mb-3 items-center border-border border-b pt-4 pb-2 justify-between">
        <div className="flex items-center gap-1">
          <ChevronRight className="w-5 h-5 text-foreground text-opacity-20" />
          <h1 className="text-2xl font-bold">{folder.name}</h1>
        </div>
        {/* <CreateFolderDropdown /> */}
        <div className="flex items-center gap-2">
          <FileLayout layout={layout} setLayout={setLayout} />
        </div>
      </div>
      <DragAndDropArea folderId={folder.id}>
        <div className="relative min-h-[80vh] h-full ">
          <FilesPageContext>
            <div className="m-3 h-full overflow-y-auto scrollbar-thin scrollbar-track-transparent scrollbar-thumb-black/60">
              <div className="">
                <h4 className="text-sm font-medium mb-3">Files</h4>
                {files.length > 0 ? (
                  <div className="grid grid-cols-8 gap-3">
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
