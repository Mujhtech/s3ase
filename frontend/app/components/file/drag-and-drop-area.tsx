import React, { useState, useCallback } from "react";
import { motion, AnimatePresence } from "framer-motion";
import { useDropzone } from "react-dropzone-esm";
import { CloudUpload } from "lucide-react";
import { useFetcher } from "@remix-run/react";
import { useApp } from "~/hooks/use-apps";

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
  const { slug } = useApp();

  const onDrop = useCallback((acceptedFiles: File[]) => {
    if (acceptedFiles.length > 0) {
      const formData = new FormData();
      formData.append("file", acceptedFiles[0]);
      formData.append("intent", "create");
      formData.append("type", "file");
      formData.append("name", acceptedFiles[0].name);

      if (folderId) {
        formData.append("folder_id", folderId);
      }

      fetcher.submit(formData, {
        method: "POST",
        action: `/resources/${slug}/files`,
      });
    }
    // setIsUploading(true);
    // // Simulate upload progress
    // let progress = 0;
    // const interval = setInterval(() => {
    //   progress += 10;
    //   setUploadProgress(progress);
    //   if (progress >= 100) {
    //     clearInterval(interval);
    //     setTimeout(() => {
    //       setIsUploading(false);
    //       setUploadProgress(0);
    //     }, 500);
    //   }
    // }, 200);
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
