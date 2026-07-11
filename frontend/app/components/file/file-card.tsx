import { Link } from "@remix-run/react";
import { FileText } from "lucide-react";
import { useApp } from "~/hooks/use-apps";
import { filePath } from "~/lib/path";
import { File } from "~/models/file";

export default function FileCard({ file }: { file: File }) {
  const app = useApp();
  return (
    <Link to={filePath(app.slug, file.id)}>
      <div className="border border-b p-8">
        <div className="flex flex-col items-center justify-center">
          <FileText className="w-16 h-16" />
          <span className="text-xs mt-2 max-w-full truncate">{file.name}</span>
        </div>
      </div>
    </Link>
  );
}
