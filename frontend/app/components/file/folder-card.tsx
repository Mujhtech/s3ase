import React from "react";
import { Folder } from "~/models/folder";
import { Folder as FolderIcon } from "lucide-react";
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from "~/components/ui/tooltip";
import { Link } from "@remix-run/react";
import { folderPath } from "~/lib/path";

export default function FolderCard({
  folder,
  appSlug,
}: {
  folder: Folder;
  appSlug: string;
}) {
  return (
    <TooltipProvider>
      <Tooltip>
        <TooltipTrigger>
          <Link to={folderPath(appSlug, folder.id)}>
            <div className="border border-b p-8">
              <div className="flex flex-col items-center justify-center">
                <FolderIcon className="w-16 h-16" />
                <span className="text-xs mt-2">{folder.name}</span>
              </div>
            </div>
          </Link>
        </TooltipTrigger>
        <TooltipContent className="w-[180px]">
          <p>{folder.description}</p>
        </TooltipContent>
      </Tooltip>
    </TooltipProvider>
  );
}
