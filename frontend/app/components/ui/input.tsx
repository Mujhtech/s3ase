import * as React from "react";

import { cn } from "~/lib/utils";

export interface InputProps
  extends React.InputHTMLAttributes<HTMLInputElement> {}

const Input = React.forwardRef<HTMLInputElement, InputProps>(
  ({ className, type, ...props }, ref) => {
    return (
      <input
        type={type}
        className={cn(
          "flex h-9 w-full rounded-md border border-input bg-transparent px-3 py-1 text-sm shadow-sm transition-colors file:border-0 file:bg-transparent file:text-sm file:font-medium file:text-foreground placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-0 focus-visible:ring-ring disabled:cursor-not-allowed disabled:opacity-50",
          className
        )}
        ref={ref}
        {...props}
      />
    );
  }
);
Input.displayName = "Input";

export interface InputWithAddonsProps
  extends React.InputHTMLAttributes<HTMLInputElement> {
  leading?: React.ReactNode;
  trailing?: React.ReactNode;
}

const InputGroup = React.forwardRef<HTMLInputElement, InputWithAddonsProps>(
  ({ leading, trailing, className, ...props }, ref) => {
    return (
      <div className="border-input ring-offset-background focus-within:ring-ring group flex h-10 w-full rounded-md border bg-transparent text-sm focus-within:outline-none focus-within:ring-0 focus-within:ring-offset-0">
        {leading ? (
          <div className="border-input bg-muted border-r px-3 py-2">
            {leading}
          </div>
        ) : null}
        <input
          className={cn(
            "bg-transparent placeholder:text-muted-foreground w-full rounded-md px-3 py-2 focus:outline-none disabled:cursor-not-allowed disabled:opacity-50",
            className
          )}
          ref={ref}
          {...props}
        />
        {trailing ? (
          <div className="border-input bg-muted border-l px-3 py-2">
            {trailing}
          </div>
        ) : null}
      </div>
    );
  }
);
InputGroup.displayName = "InputGroup";

export { Input,InputGroup };
