import { useForm } from "@conform-to/react";
import { parseWithZod } from "@conform-to/zod";
import { useFetcher } from "@remix-run/react";
import { Check,ChevronsUpDown } from "lucide-react";
import { useState } from "react";
import { Button } from "~/components/ui/button";
import {
Command,
CommandEmpty,
CommandGroup,
CommandInput,
CommandItem,
CommandList,
} from "~/components/ui/command";
import {
Dialog,
DialogContent,
DialogTitle,
DialogTrigger,
} from "~/components/ui/dialog";
import FormError from "~/components/ui/form-error";
import FormField from "~/components/ui/form-field";
import { Input } from "~/components/ui/input";
import { Label } from "~/components/ui/label";
import Paragraph from "~/components/ui/paragraph";
import {
Popover,
PopoverContent,
PopoverTrigger,
} from "~/components/ui/popover";
import { Textarea } from "~/components/ui/textarea";
import { cn } from "~/lib/utils";
import { App,CreateAppFormSchema } from "~/models/app";

const regions = [
  "us-east-1",
  "us-east-2",
  "us-west-1",
  "us-west-2",
  "eu-central-1",
  "eu-west-1",
  "eu-west-2",
  "eu-west-3",
];

export default function CreateAppDialog() {
  const [isOpen, setIsOpen] = useState(false);

  return (
    <Dialog open={isOpen} onOpenChange={setIsOpen}>
      <DialogTrigger asChild>
        <Button className="!h-8">Create App</Button>
      </DialogTrigger>
      <DialogContent
        leading={<DialogTitle>Create New App</DialogTitle>}
        aria-describedby="Create New App"
      >
        <CreateAppForm closeDialog={() => setIsOpen(false)} />
      </DialogContent>
    </Dialog>
  );
}

export function CreateAppForm({
  closeDialog,
  app,
}: {
  closeDialog?: () => void;
  app?: App;
}) {
  const [region, setRegion] = useState<string>(app?.region || "eu-west-2");

  const fetcher = useFetcher();

  const [form, fields] = useForm({
    id: "create-app-form",
    shouldValidate: "onBlur",
    shouldRevalidate: "onSubmit",

    onValidate({ formData }) {
      return parseWithZod(formData, {
        schema: CreateAppFormSchema,
      });
    },
  });

  return (
    <fetcher.Form
      method="post"
      className="grid grid-cols-1 gap-3"
      onSubmit={form.onSubmit}
    >
      {app && <input type="hidden" name={"type"} value={"app"} />}
      <input
        type="hidden"
        key={fields.region.key}
        name={fields.region.name}
        value={region}
      />
      <input
        type="hidden"
        key={fields.slug.key}
        name={fields.slug.name}
        defaultValue={fields.slug.initialValue || app?.slug || ""}
      />
      <FormField>
        <Label>Name</Label>
        <Input
          key={fields.name.key}
          name={fields.name.name}
          defaultValue={fields.name.initialValue || app?.name || ""}
        />
        <FormError>{fields.name.errors}</FormError>
      </FormField>
      <FormField>
        <Label>Description (Optional)</Label>
        <Textarea
          key={fields.description.key}
          name={fields.description.name}
          defaultValue={
            fields.description.initialValue || app?.description || ""
          }
        />
        <FormError>{fields.description.errors}</FormError>
      </FormField>

      {app && (
        <FormField>
          <Label>Bucket</Label>
          <Input defaultValue={app?.bucket || ""} readOnly={true} />
        </FormField>
      )}

      <FormField>
        <Label>Region</Label>
        {app ? (
          <Input defaultValue={app?.region || ""} readOnly={true} />
        ) : (
          <SelectRegionDialog
            value={region}
            onChange={(val) => {
              setRegion(val);
            }}
          />
        )}

        <FormError>{fields.region.errors}</FormError>
      </FormField>
      <div>
        <Paragraph>
          Note: This action will auto create a new bucket if it does not exist with
          the following slug{" "}
          {fields.name.value?.toLowerCase().replace(/\s+/g, "-")}
        </Paragraph>
      </div>
      <div className="flex gap-2">
        {closeDialog && (
          <Button variant="outline" type="button" onClick={() => closeDialog()}>
            Close
          </Button>
        )}

        <Button type="submit" name="intent" value="create">
          {app ? "Update" : "Create"}
        </Button>
      </div>
    </fetcher.Form>
  );
}

export function SelectRegionDialog({
  value,
  onChange,
}: {
  value?: string;
  onChange: (value: string) => void;
}) {
  const [open, setOpen] = useState(false);

  return (
    <Popover open={open} onOpenChange={setOpen}>
      <PopoverTrigger asChild>
        <Button
          variant="outline"
          role="combobox"
          aria-expanded={open}
          className="w-full justify-between"
        >
          {value
            ? regions.find((region) => region === value)
            : "Select region..."}
          <ChevronsUpDown className="ml-2 h-4 w-4 shrink-0 opacity-50" />
        </Button>
      </PopoverTrigger>
      <PopoverContent className="w-[200px] p-0">
        <Command>
          <CommandInput placeholder="Search region..." />
          <CommandList>
            <CommandEmpty>No region found.</CommandEmpty>
            <CommandGroup>
              {regions.map((region) => (
                <CommandItem
                  key={region}
                  value={region}
                  onSelect={(currentValue) => {
                    onChange(currentValue === value ? "" : currentValue);
                    setOpen(false);
                  }}
                >
                  <Check
                    className={cn(
                      "mr-2 h-4 w-4",
                      value === region ? "opacity-100" : "opacity-0"
                    )}
                  />
                  {region}
                </CommandItem>
              ))}
            </CommandGroup>
          </CommandList>
        </Command>
      </PopoverContent>
    </Popover>
  );
}
