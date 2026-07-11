import { LoaderFunctionArgs } from "@remix-run/node";
import { Link } from "@remix-run/react";
import { FileText, Folder, HardDrive, Users } from "lucide-react";
import { typedjson, useTypedLoaderData } from "remix-typedjson";
import { Card, CardContent, CardHeader } from "~/components/ui/card";
import { useApp } from "~/hooks/use-apps";
import { filePath, filesPath, settingsMenuPath } from "~/lib/path";
import { getFiles } from "~/services/file.server";
import { getFolders } from "~/services/folder.server";
import { getMembers } from "~/services/member.server";
import { getUsage } from "~/services/usage.server";
import { AppSlugParamSchema } from "../_app.app.$appSlug/route";

export const loader = async ({ request, params }: LoaderFunctionArgs) => {
  AppSlugParamSchema.parse(params);
  const [usage, folders, members, files] = await Promise.all([
    getUsage(request),
    getFolders(request),
    getMembers(request),
    getFiles(request, { per_page: 4 }),
  ]);
  return typedjson({ usage, folders, members, files });
};

function formatBytes(bytes: number) {
  if (bytes === 0) return "0 B";
  const units = ["B", "KB", "MB", "GB", "TB"];
  const index = Math.min(Math.floor(Math.log(bytes) / Math.log(1024)), units.length - 1);
  return `${(bytes / 1024 ** index).toFixed(index === 0 ? 0 : 1)} ${units[index]}`;
}

export default function Page() {
  const { usage, folders, members, files } = useTypedLoaderData<typeof loader>();
  const app = useApp();

  const metrics = [
    { label: "Files", value: usage.file_count.toLocaleString(), icon: FileText },
    { label: "Storage", value: formatBytes(usage.storage_bytes), icon: HardDrive },
    { label: "Folders", value: folders.length.toLocaleString(), icon: Folder },
    { label: "Members", value: members.length.toLocaleString(), icon: Users },
  ];

  return (
    <div className="flex w-full max-w-5xl flex-col gap-6">
      <div className="border-b border-border pb-3 pt-4">
        <h1 className="text-2xl font-bold">{app.name}</h1>
        <p className="text-sm text-muted-foreground">App overview</p>
      </div>

      <div className="grid grid-cols-1 gap-3 sm:grid-cols-2 lg:grid-cols-4">
        {metrics.map(({ label, value, icon: Icon }) => (
          <Card key={label}>
            <CardContent className="flex items-center justify-between p-5">
              <div>
                <p className="text-sm text-muted-foreground">{label}</p>
                <p className="mt-1 text-2xl font-semibold">{value}</p>
              </div>
              <Icon className="h-5 w-5 text-muted-foreground" />
            </CardContent>
          </Card>
        ))}
      </div>

      <Card>
        <CardHeader className="flex flex-row items-center justify-between">
          <h2 className="font-semibold">Recent files</h2>
          <Link className="text-sm text-muted-foreground hover:text-foreground" to={filesPath(app.slug)}>
            View all
          </Link>
        </CardHeader>
        <CardContent>
          {files.length === 0 ? (
            <p className="text-sm text-muted-foreground">No files uploaded yet.</p>
          ) : (
            <div className="divide-y divide-border border border-border">
              {files.map((file) => (
                <Link
                  key={file.id}
                  to={filePath(app.slug, file.id)}
                  className="flex items-center justify-between p-3 hover:bg-muted/50"
                >
                  <span className="truncate text-sm font-medium">{file.name}</span>
                  <span className="text-xs text-muted-foreground">{formatBytes(file.size)}</span>
                </Link>
              ))}
            </div>
          )}
        </CardContent>
      </Card>

      <Link className="w-fit text-sm text-muted-foreground hover:text-foreground" to={settingsMenuPath(app.slug, "members")}>
        Manage members
      </Link>
    </div>
  );
}
