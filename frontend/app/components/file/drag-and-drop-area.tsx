import React, { useState, useCallback } from "react";
import { motion, AnimatePresence } from "framer-motion";
import { useDropzone } from "react-dropzone-esm";
import { CloudUpload } from "lucide-react";
import { useFetcher } from "@remix-run/react";
import { useApp } from "~/hooks/use-apps";
import { useAuthToken, useBackendUrl } from "~/hooks/use-user";
import * as tus from "tus-js-client";

export default function DragAndDropArea({
  children,
  folderId,
}: {
  children: React.ReactNode;
  folderId?: string;
}) {
  const [isUploading, setIsUploading] = useState(false);
  const [uploadProgress, setUploadProgress] = useState(0);
  const [isDragActive, setIsDragActive] = useState(false);
  const fetcher = useFetcher();
  const { slug, id } = useApp();
  const backendUrl = useBackendUrl();
  const accessToken = useAuthToken();

  const uploadeFile = useCallback(
    async (file: File) => {
      let headers = {
        "Content-Type": "application/offset+octet-stream",
        Authorization: `Bearer ${accessToken}`,
        "x-app-id": id,
      };

      const url = `${backendUrl}/ui/files?name=${file.name}&folder_id=${
        folderId ?? ""
      }&intent=create&type=file`;

      const response = await fetch(url, {
        method: "POST",
        headers: headers,
        body: file,
      });

      if (!response.ok) {
        throw new Error(`Upload failed with status: ${response.status}`);
      }

      const data = await response.json();

      return data;

      // const upload = new tus.Upload(file, {
      //   // Replace this with tusd's upload creation URL
      //   endpoint: url,
      //   headers: headers,

      //   onError: function (error) {
      //     console.log("Failed because: " + error);
      //   },
      //   onSuccess: function () {
      //     //console.log("Download %s from %s", upload.file.name, upload.url);
      //   },
      // });

      // upload.start();
    },
    [accessToken, backendUrl]
  );

  const onDrop = useCallback((acceptedFiles: File[]) => {
    if (acceptedFiles.length > 0) {
      uploadeFile(acceptedFiles[0])
        .then((data) => {
          console.log(data);
        })
        .catch((error) => {
          console.log(error);
        });
    }
  }, []);

  const {
    getRootProps,
    getInputProps,
    isDragActive: dropzoneIsDragActive,
  } = useDropzone({
    onDrop,
    onDragEnter: () => setIsDragActive(true),
    onDragLeave: () => setIsDragActive(false),
    onDropAccepted: () => setIsDragActive(false),
    onDropRejected: () => setIsDragActive(false),
  });

  return (
    <div className="relative min-h-[90vh] h-full w-full">
      <div className="absolute w-full h-full" {...getRootProps()}>
        <input {...getInputProps()} />
        <AnimatePresence>
          {isDragActive == true && (
            <motion.div
              initial={{ opacity: 0, y: 50 }}
              animate={{ opacity: 1, y: 0 }}
              exit={{ opacity: 0, y: 50 }}
              transition={{ duration: 0.3 }}
              className="absolute bottom-6 left-0 right-0 top-0 bg-white bg-opacity-5 border-2 flex flex-col items-center justify-end rounded-lg z-10 mr-3 pb-10"
            >
              <motion.div
                initial={{ y: 50, opacity: 0 }}
                animate={{ y: 0, opacity: 1 }}
                exit={{ y: 50, opacity: 0 }}
                transition={{ delay: 0.1, duration: 0.3 }}
                className="flex flex-col items-center"
              >
                <motion.div
                  animate={{
                    y: [0, -10, 0],
                    transition: {
                      duration: 1,
                      repeat: Infinity,
                      ease: "easeInOut",
                    },
                  }}
                >
                  <CloudUpload className="w-16 h-16 text-white" />
                </motion.div>
                <motion.p className="text-white text-lg font-medium mt-4">
                  Drop files to upload here
                </motion.p>
              </motion.div>
            </motion.div>
          )}
        </AnimatePresence>
      </div>
      {children}
    </div>
  );
}
