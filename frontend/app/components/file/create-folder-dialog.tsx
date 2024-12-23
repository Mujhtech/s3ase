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
import { Folder, CreateFolderFormSchema } from "~/models/folder";
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

export default function CreateFolderDialog({
  folder,
  onDialogOpen,
}: {
  folder?: Folder;
  onDialogOpen?: () => void;
}) {
  const [isOpen, setIsOpen] = useState(false);
  const [isDropdownOpen, setIsDropdownOpen] = useState(false);

  const title = folder ? "Update" : "Create" + " Folder";

  return (
    <Dialog open={isOpen} onOpenChange={(value) => setIsOpen(value)}>
      <DialogTrigger asChild>
        <Button className="!h-8">New Folder</Button>
      </DialogTrigger>
      <DialogContent
        leading={<DialogTitle>{title}</DialogTitle>}
        aria-describedby={title}
        onClick={(e) => e.stopPropagation()}
      >
        <CreateFolderForm
          folder={folder}
          onCloseDialog={() => setIsOpen(false)}
          title={title}
        />
      </DialogContent>
    </Dialog>
  );
}

export function CreateFolderForm({
  folder,
  onCloseDialog,
  title,
}: {
  folder?: Folder;
  onCloseDialog?: () => void;
  title: string;
}) {
  const fetcher = useFetcher();

  const [form, fields] = useForm({
    id: folder ? "update-folder-form" : "create-folder-form",
    shouldValidate: "onBlur",
    shouldRevalidate: "onSubmit",

    onValidate({ formData }) {
      return parseWithZod(formData, {
        schema: CreateFolderFormSchema,
      });
    },
  });

  return (
    <fetcher.Form
      method="post"
      className="grid grid-cols-1 gap-3"
      onSubmit={form.onSubmit}
    >
      <input
        type="hidden"
        key={fields.type.key}
        name={fields.type.name}
        defaultValue={"folder"}
      />

      <FormField>
        <Label>Name</Label>
        <Input
          key={fields.name.key}
          name={fields.name.name}
          defaultValue={fields.name.initialValue || folder?.name || ""}
        />
        <FormError>{fields.name.errors}</FormError>
      </FormField>

      <FormField>
        <Label>Description (Optional)</Label>
        <Textarea
          key={fields.description.key}
          name={fields.description.name}
          defaultValue={
            fields.description.initialValue || folder?.description || ""
          }
        />
        <FormError>{fields.description.errors}</FormError>
      </FormField>

      <div className="flex gap-2">
        {onCloseDialog && (
          <Button
            variant="outline"
            type="button"
            onClick={() => onCloseDialog()}
          >
            Close
          </Button>
        )}

        <Button
          type="submit"
          name="intent"
          value={folder ? "update" : "create"}
        >
          {title}
        </Button>
      </div>
    </fetcher.Form>
  );
}
