import {
  Dialog,
  DialogContent,
  DialogTitle,
  DialogTrigger,
} from "~/components/ui/dialog";
import { Button } from "~/components/ui/button";
import { useFetcher, useLocation, useNavigation } from "@remix-run/react";
import { Label } from "~/components/ui/label";
import { Input, InputGroup } from "~/components/ui/input";
import FormField from "~/components/ui/form-field";
import { useState } from "react";
import { useForm, useInputControl } from "@conform-to/react";
import { parseWithZod } from "@conform-to/zod";
import FormError from "~/components/ui/form-error";
import { CreateWebhookFormSchema, Webhook } from "~/models/webhook";
import { Textarea } from "../ui/textarea";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "~/components/ui/dropdown-menu";
import { CircleEllipsis } from "lucide-react";
import { Checkbox } from "../ui/checkbox";

const events = [
  {
    name: "upload.started",
    description: "When an upload is started",
  },
  {
    name: "upload.completed",
    description: "When an upload is completed",
  },
  {
    name: "upload.failed",
    description: "When an upload is failed",
  },
  {
    name: "upload.deleted",
    description: "When an upload is deleted",
  },
];

export default function CreateWebhookDialog({
  webhook,
}: {
  webhook?: Webhook;
}) {
  const [isOpen, setIsOpen] = useState(false);
  const [isDropdownOpen, setIsDropdownOpen] = useState(false);
  const [selectedEvents, setSelectedEvents] = useState<string[]>([]);
  const [url, setUrl] = useState<string | undefined>(
    webhook?.url.replace("https://", "") || ""
  );

  const fetcher = useFetcher();

  const [form, fields] = useForm({
    id: "create-webhook-form",
    shouldValidate: "onBlur",
    shouldRevalidate: "onSubmit",

    onValidate({ formData }) {
      return parseWithZod(formData, {
        schema: CreateWebhookFormSchema,
      });
    },
  });

  const title = webhook ? "Update" : "Create" + " Webhook";

  return (
    <Dialog open={isOpen} onOpenChange={(value) => setIsOpen(value)}>
      {webhook ? (
        <DropdownMenu open={isDropdownOpen} onOpenChange={setIsDropdownOpen}>
          <DropdownMenuTrigger asChild>
            <button className="!h-8">
              <CircleEllipsis className="h-4 w-4" />
            </button>
          </DropdownMenuTrigger>
          <DropdownMenuContent>
            <DropdownMenuItem>
              <DialogTrigger asChild>
                <button className="relative flex items-center rounded-sm px-2 py-1.5 text-sm outline-none transition-colors focus:bg-accent focus:text-accent-foreground">
                  Edit
                </button>
              </DialogTrigger>
            </DropdownMenuItem>
            <DropdownMenuItem>
              <DropdownMenuLabel>Remove</DropdownMenuLabel>
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      ) : (
        <DialogTrigger asChild>
          <Button className="!h-8">Add Webhook</Button>
        </DialogTrigger>
      )}

      <DialogContent
        leading={<DialogTitle>{title}</DialogTitle>}
        aria-describedby={title}
      >
        <fetcher.Form
          method="post"
          className="grid grid-cols-1 gap-3"
          onSubmit={form.onSubmit}
        >
          <input
            type="hidden"
            key={fields.id.key}
            name={fields.id.name}
            defaultValue={webhook?.id || ""}
          />
          <input
            type="hidden"
            key={fields.url.key}
            name={fields.url.name}
            defaultValue={`https://${url?.trim()}`}
          />
          <input
            type="hidden"
            key={fields.events.key}
            name={fields.events.name}
            defaultValue={selectedEvents}
          />
          <FormField>
            <Label>Name</Label>
            <Input
              key={fields.name.key}
              name={fields.name.name}
              defaultValue={fields.name.initialValue || webhook?.name || ""}
            />
            <FormError>{fields.name.errors}</FormError>
          </FormField>

          <FormField>
            <Label>Description (Optional)</Label>
            <Textarea
              key={fields.description.key}
              name={fields.description.name}
              defaultValue={
                fields.description.initialValue || webhook?.description || ""
              }
            />
            <FormError>{fields.description.errors}</FormError>
          </FormField>

          <FormField>
            <Label>URL</Label>
            <InputGroup
              leading={<Label>https://</Label>}
              // key={fields.url.key}
              // name={fields.url.name}
              defaultValue={url || ""}
              onChange={(e) => setUrl(e.target.value)}
            />
            <FormError>{fields.url.errors}</FormError>
          </FormField>

          <FormField>
            <Label className="mb-2">Event to listens to</Label>
            {events.map((item, index) => (
              <div
                key={index}
                className="flex flex-row items-start space-x-3 space-y-0"
              >
                <Checkbox
                  checked={selectedEvents.includes(item.name)}
                  onCheckedChange={(checked) => {
                    if (checked) {
                      setSelectedEvents([...selectedEvents, item.name]);
                    } else {
                      setSelectedEvents(
                        selectedEvents.filter((event) => event !== item.name)
                      );
                    }
                    return checked;
                  }}
                />
                <Label className="text-sm font-normal">{item.name}</Label>
              </div>
            ))}
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
              value={webhook ? "update" : "create"}
            >
              {title}
            </Button>
          </div>
        </fetcher.Form>
      </DialogContent>
    </Dialog>
  );
}
