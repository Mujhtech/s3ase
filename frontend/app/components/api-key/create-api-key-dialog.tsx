import {
  Dialog,
  DialogContent,
  DialogTitle,
  DialogTrigger,
} from "~/components/ui/dialog";
import { Button } from "~/components/ui/button";
import { useFetcher, useLocation, useNavigation } from "@remix-run/react";
import { Label } from "~/components/ui/label";
import { Input } from "~/components/ui/input";
import { Textarea } from "~/components/ui/textarea";
import FormField from "~/components/ui/form-field";
import Paragraph from "~/components/ui/paragraph";
import { useState } from "react";
import { useForm, useInputControl } from "@conform-to/react";
import { parseWithZod } from "@conform-to/zod";
import FormError from "~/components/ui/form-error";
import { ApiKey, CreateApiKeyFormSchema } from "~/models/api_key";
import {
  Command,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
  CommandList,
} from "~/components/ui/command";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "~/components/ui/popover";
import { Check, ChevronsUpDown, CircleEllipsis } from "lucide-react";
import { cn } from "~/lib/utils";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "~/components/ui/dropdown-menu";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "~/components/ui/select";

const accesses = ["full", "read", "write"];

export default function CreateApiKeyDialog({
  apiKey,
  onDialogOpen,
}: {
  apiKey?: ApiKey;
  onDialogOpen?: () => void;
}) {
  const [isOpen, setIsOpen] = useState(false);
  const [isDropdownOpen, setIsDropdownOpen] = useState(false);
  const [access, setAccess] = useState<string | undefined>(
    apiKey?.access ?? "read"
  );

  const fetcher = useFetcher();

  const [form, fields] = useForm({
    id: apiKey ? "update-api-key-form" : "create-api-key-form",
    shouldValidate: "onBlur",
    shouldRevalidate: "onSubmit",

    onValidate({ formData }) {
      return parseWithZod(formData, {
        schema: CreateApiKeyFormSchema,
      });
    },
  });

  const title = apiKey ? "Update" : "Create" + " API Key";

  return (
    <Dialog open={isOpen} onOpenChange={(value) => setIsOpen(value)}>
      {apiKey ? (
        <DropdownMenu open={isDropdownOpen} onOpenChange={setIsDropdownOpen}>
          <DropdownMenuTrigger asChild>
            <button className="!h-8">
              <CircleEllipsis className="h-4 w-4" />
            </button>
          </DropdownMenuTrigger>
          <DropdownMenuContent>
            <DropdownMenuItem onSelect={(e) => e.preventDefault()}>
              <DialogTrigger asChild>
                <div
                  // onClick={(e) => {
                  //   onDialogOpen?.call(null);
                  //   e.stopPropagation();
                  //   setIsOpen(true);
                  // }}
                  className="relative flex items-center cursor-pointer rounded-sm px-2 py-1.5 text-sm outline-none transition-colors focus:bg-accent focus:text-accent-foreground"
                >
                  Edit
                </div>
              </DialogTrigger>
            </DropdownMenuItem>
            <DropdownMenuItem>
              <DropdownMenuLabel>Remove</DropdownMenuLabel>
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      ) : (
        <DialogTrigger asChild>
          <Button className="!h-8">Create API Key</Button>
        </DialogTrigger>
      )}
      <DialogContent
        leading={<DialogTitle>{title}</DialogTitle>}
        aria-describedby={title}
        //onClick={(e) => e.stopPropagation()}
      >
        <fetcher.Form
          method="post"
          className="grid grid-cols-1 gap-3"
          onSubmit={form.onSubmit}
        >
          <input
            type="hidden"
            key={fields.access.key}
            name={fields.access.name}
            defaultValue={access || apiKey?.access || ""}
          />
          <input
            type="hidden"
            key={fields.id.key}
            name={fields.id.name}
            defaultValue={apiKey?.id || ""}
          />
          <FormField>
            <Label>Name</Label>
            <Input
              key={fields.name.key}
              name={fields.name.name}
              defaultValue={fields.name.initialValue || apiKey?.name || ""}
            />
            <FormError>{fields.name.errors}</FormError>
          </FormField>

          <FormField>
            <Label>Description (Optional)</Label>
            <Textarea
              key={fields.description.key}
              name={fields.description.name}
              defaultValue={
                fields.description.initialValue || apiKey?.description || ""
              }
            />
            <FormError>{fields.description.errors}</FormError>
          </FormField>

          <FormField>
            <Label>Access</Label>
            <Select defaultValue={access} onValueChange={setAccess}>
              <SelectTrigger className={cn("capitalize")}>
                <SelectValue placeholder="Select access" className="" />
              </SelectTrigger>
              <SelectContent>
                {accesses.map((access) => (
                  <SelectItem
                    key={access}
                    value={access}
                    className="capitalize"
                  >
                    {access}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>

            {/* <SelectApiAccess
              value={access}
              onChange={(val) => {
                setAccess(val);
              }}
            /> */}
            <FormError>{fields.access.errors}</FormError>
          </FormField>

          <div className="flex gap-2">
            <Button
              variant="outline"
              type="button"
              onClick={() => setIsOpen(false)}
            >
              Close
            </Button>

            <Button
              type="submit"
              name="intent"
              value={apiKey ? "update" : "create"}
            >
              {title}
            </Button>
          </div>
        </fetcher.Form>
      </DialogContent>
    </Dialog>
  );
}
