import React from "react";
import { cn } from "~/lib/utils";

export default function FormField({
  className,
  children,
}: {
  className?: string;
  children: React.ReactNode;
}) {
  return (
    <div
      className={cn(
        "grid w-full items-center gap-1.5",

        className
      )}
    >
      {children}
    </div>
  );
}
