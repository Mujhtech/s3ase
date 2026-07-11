import { LayoutGrid,LayoutList } from "lucide-react";
import { cn } from "~/lib/utils";

export default function FileLayout({
  layout,
  setLayout,
}: {
  layout: "grid" | "list";
  setLayout: (val: "grid" | "list") => void;
}) {
  return (
    <div className="flex items-center divide-x border">
      <button
        type="button"
        className={cn("p-2", layout === "list" && "bg-accent")}
        onClick={() => setLayout("list")}
      >
        <LayoutList className="h-4 w-4" />
      </button>
      <button
        type="button"
        className={cn("p-2", layout === "grid" && "bg-accent")}
        onClick={() => setLayout("grid")}
      >
        <LayoutGrid className="h-4 w-4" />
      </button>
    </div>
  );
}
