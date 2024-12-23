import React from "react";
import { File } from "~/models/file";
import { Folder as FolderIcon } from "lucide-react";
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from "~/components/ui/tooltip";

export default function FileCard({ file }: { file: File }) {
  return (
    <div className="border border-b p-8">
      <div className="flex flex-col items-center justify-center">
        <FolderIcon className="w-16 h-16" />
        <span className="text-xs mt-2">{file.name}</span>
      </div>
    </div>
  );
}
