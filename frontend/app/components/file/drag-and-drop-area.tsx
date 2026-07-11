import { AnimatePresence,motion } from "framer-motion";
import { CloudUpload } from "lucide-react";
import React,{ useCallback,useState } from "react";
import { useDropzone } from "react-dropzone-esm";
import { useUploadManager } from "./upload-manager";

export default function DragAndDropArea({
  children,
  folderId,
}: {
  children: React.ReactNode;
  folderId?: string;
}) {
  const [isDragActive, setIsDragActive] = useState(false);
  const { startUpload } = useUploadManager();

  const onDrop = useCallback(
  (acceptedFiles: File[]) => {
    for (const file of acceptedFiles) {
    void startUpload(file, folderId).catch((error: unknown) => {
      console.error("File upload failed", error);
    });
    }
  },
  [folderId, startUpload]
  );

  const {
    getRootProps,
    getInputProps,
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
