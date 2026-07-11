import {
CaretDown,
CaretUp,
CheckCircle,
Image,
Pause,
Play,
Spinner,
Trash,
WarningCircle,
X,
} from "@phosphor-icons/react";
import React,{ useCallback,useMemo } from "react";
import Paragraph from "../ui/paragraph";
import { Progress } from "../ui/progress";
import { useUploadManager } from "./upload-manager";

export type FileUploadProgress = {
  id: string;
  name: string;
  folder_id?: string;
  progress: number;
  status:
    | "started"
    | "uploading"
    | "completed"
    | "failed"
    | "pending"
    | "paused"
    | "cancelled";
};

export default function FilesUploadState({
  open,
  setOpen,
  files,
}: {
  files: FileUploadProgress[];
  open: boolean;
  setOpen: React.Dispatch<React.SetStateAction<boolean>>;
}) {
  const [expanded, setExpanded] = React.useState(true);
  const handleClose = useCallback(() => {
    setOpen(false);
  }, [setOpen]);

  const title = useMemo(() => {
    const totalUploadingFile = files.filter(
      (file) => file.status === "uploading"
    ).length;

    const totalCompletedFile = files.filter(
      (file) => file.status === "completed"
    ).length;

    if (totalUploadingFile > 0) {
      return `${totalUploadingFile} Upload files`;
    }

    return `${totalCompletedFile} Completed files`;
  }, [files]);

  if (!open) return null;

  return (
    <div className="absolute bottom-0 right-3">
      <div className="z-50 flex flex-col min-w-80 max-w-lg border bg-background shadow-lg duration-200">
        <div className="flex justify-between items-center border-b py-2 px-4">
          <Paragraph>{title}</Paragraph>
          <div className="flex gap-3 items-center">
            <button
              onClick={() => setExpanded(!expanded)}
              className="flex items-center gap-2"
            >
              {expanded ? <CaretUp /> : <CaretDown />}
            </button>
            <button onClick={handleClose} className="flex items-center gap-2">
              <X />
            </button>
          </div>
        </div>
        {expanded && (
          <div className="flex flex-col w-full">
            {files.map((file, i) => (
              <FileUploadStateCard key={i} file={file} />
            ))}
          </div>
        )}
      </div>
    </div>
  );
}

const FileUploadStateCard = ({ file }: { file: FileUploadProgress }) => {
  const { pauseUpload, resumeUpload, cancelUpload } = useUploadManager();
  const canPause = file.status === "uploading" || file.status === "pending";
  const canResume = file.status === "paused" || file.status === "failed";
  const canCancel = canPause || canResume;

  return (
    <div className="flex flex-col gap-2 px-4 py-2 w-full">
      <div className="flex items-center justify-between w-full">
        <div className="flex items-center gap-2">
          <Image className="w-5 h-5 text-muted-foreground" />
          <Paragraph>{file.name}</Paragraph>
        </div>
        <div className="flex items-center gap-2">
          {canPause && (
            <button
              type="button"
              aria-label={`Pause ${file.name}`}
              onClick={() => void pauseUpload(file.id)}
            >
              <Pause className="w-4 h-4 text-muted-foreground" />
            </button>
          )}
          {canResume && (
            <button
              type="button"
              aria-label={`Resume ${file.name}`}
              onClick={() => resumeUpload(file.id)}
            >
              <Play className="w-4 h-4 text-muted-foreground" />
            </button>
          )}
          {canCancel && (
            <button
              type="button"
              aria-label={`Cancel ${file.name}`}
              onClick={() => void cancelUpload(file.id)}
            >
              <Trash className="w-4 h-4 text-muted-foreground" />
            </button>
          )}
          {file.status === "completed" ? (
            <CheckCircle className="w-5 h-5 text-primary" />
          ) : file.status === "failed" || file.status === "cancelled" ? (
            <WarningCircle className="w-5 h-5 text-primary" />
          ) : file.status === "paused" ? (
            <Pause className="w-5 h-5 text-muted-foreground" />
          ) : (
            <Spinner className="w-5 h-5 text-primary animate-spin" />
          )}
        </div>
      </div>

      {(file.status === "uploading" || file.status === "paused") && (
        <Progress value={file.progress} />
      )}
    </div>
  );
};
