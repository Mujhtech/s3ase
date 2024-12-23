import React from "react";
import { cn } from "~/lib/utils";

export default function Paragraph({
  className,
  children,
}: {
  className?: string;
  children: React.ReactNode;
}) {
  return (
    <p className={cn("text-xs text-muted-foreground", className)}>{children}</p>
  );
}
