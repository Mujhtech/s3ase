import React from "react";
import { cn } from "~/lib/utils";

export default function FormError({
  children,
  className,
}: {
  className?: string;
  children: React.ReactNode;
}) {
  return (
    <p className={cn("text-sm font-medium text-red-500", className)}>
      {children}
    </p>
  );
}
