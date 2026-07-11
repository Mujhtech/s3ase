import React, {
  createContext,
  useCallback,
  useContext,
  useMemo,
  useRef,
} from "react";
import * as tus from "tus-js-client";
import { useApp } from "~/hooks/use-apps";
import { useAuthToken, useBackendUrl } from "~/hooks/use-user";
import type { FileUploadProgress } from "./files-upload-state";

type UploadEntry = {
  file: File;
  upload: tus.Upload;
  folderId?: string;
  progress: number;
};

type UploadManagerValue = {
  startUpload: (file: File, folderId?: string) => Promise<void>;
  pauseUpload: (id: string) => Promise<void>;
  resumeUpload: (id: string) => void;
  cancelUpload: (id: string) => Promise<void>;
};

const UploadManagerContext = createContext<UploadManagerValue | null>(null);

export function UploadManagerProvider({
  children,
  onUploadChange,
}: {
  children: React.ReactNode;
  onUploadChange: (file: FileUploadProgress) => void;
}) {
  const app = useApp();
  const accessToken = useAuthToken();
  const backendUrl = useBackendUrl();
  const uploads = useRef(new Map<string, UploadEntry>());

  const startUpload = useCallback(
    async (file: File, folderId?: string) => {
      if (!accessToken || !app.id || !backendUrl) {
        throw new Error("Upload session is not ready");
      }

      const normalizedBase = backendUrl.endsWith("/")
        ? backendUrl
        : `${backendUrl}/`;
      const endpoint = new URL("ui/files/uploads", normalizedBase);
      endpoint.searchParams.set("name", file.name);
      if (folderId) endpoint.searchParams.set("folder_id", folderId);

      let fileID: string = crypto.randomUUID();
      const notify = (
        status: FileUploadProgress["status"],
        progress: number
      ) =>
        onUploadChange({
          id: fileID,
          name: file.name,
          folder_id: folderId,
          progress,
          status,
        });

      const upload = new tus.Upload(file, {
        endpoint: endpoint.toString(),
        headers: {
          Authorization: `Bearer ${accessToken}`,
          "x-app-id": app.id,
        },
        metadata: {
          filename: file.name,
          filetype: file.type || "application/octet-stream",
          file_id: fileID,
          ...(folderId ? { folder_id: folderId } : {}),
        },
        chunkSize: 8 * 1024 * 1024,
        retryDelays: [0, 1_000, 3_000, 5_000, 10_000],
        removeFingerprintOnSuccess: true,
        fingerprint: async (selectedFile) =>
          [
            "s3ase",
            app.id,
            folderId ?? "root",
            selectedFile.name,
            selectedFile.type,
            selectedFile.size,
            selectedFile.lastModified,
          ].join("-"),
        onProgress: (bytesSent, bytesTotal) => {
          const progress = bytesTotal > 0 ? Math.round((bytesSent / bytesTotal) * 100) : 0;
          const entry = uploads.current.get(fileID);
          if (entry) entry.progress = progress;
          notify("uploading", progress);
        },
        onSuccess: () => {
          notify("completed", 100);
          uploads.current.delete(fileID);
        },
        onError: () => notify("failed", uploads.current.get(fileID)?.progress ?? 0),
      });

      const previousUploads = await upload.findPreviousUploads();
      const previous = previousUploads[0];
      if (previous) {
        fileID = previous.metadata.file_id || fileID;
        upload.resumeFromPreviousUpload(previous);
      }

      uploads.current.set(fileID, { file, upload, folderId, progress: 0 });
      notify(previous ? "uploading" : "pending", 0);
      upload.start();
    },
    [accessToken, app.id, backendUrl, onUploadChange]
  );

  const pauseUpload = useCallback(
    async (id: string) => {
      const entry = uploads.current.get(id);
      if (!entry) return;
      await entry.upload.abort(false);
      onUploadChange({
        id,
        name: entry.file.name,
        folder_id: entry.folderId,
        progress: entry.progress,
        status: "paused",
      });
    },
    [onUploadChange]
  );

  const resumeUpload = useCallback(
    (id: string) => {
      const entry = uploads.current.get(id);
      if (!entry) return;
      onUploadChange({
        id,
        name: entry.file.name,
        folder_id: entry.folderId,
        progress: entry.progress,
        status: "uploading",
      });
      entry.upload.start();
    },
    [onUploadChange]
  );

  const cancelUpload = useCallback(
    async (id: string) => {
      const entry = uploads.current.get(id);
      if (!entry) return;
      await entry.upload.abort(true);
      uploads.current.delete(id);
      onUploadChange({
        id,
        name: entry.file.name,
        folder_id: entry.folderId,
        progress: 0,
        status: "cancelled",
      });
    },
    [onUploadChange]
  );

  const value = useMemo(
    () => ({ startUpload, pauseUpload, resumeUpload, cancelUpload }),
    [cancelUpload, pauseUpload, resumeUpload, startUpload]
  );

  return (
    <UploadManagerContext.Provider value={value}>
      {children}
    </UploadManagerContext.Provider>
  );
}

export function useUploadManager() {
  const value = useContext(UploadManagerContext);
  if (!value) {
    throw new Error("useUploadManager must be used inside UploadManagerProvider");
  }
  return value;
}
