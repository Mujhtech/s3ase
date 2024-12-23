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
import { SendInviteFormSchema } from "~/models/member";
import { App } from "~/models/app";

export default function DeleteAppDialog({ app }: { app: App }) {
  const [isOpen, setIsOpen] = useState(false);

  const fetcher = useFetcher();

  const [form, fields] = useForm({
    id: "invite-member-form",
    shouldValidate: "onBlur",
    shouldRevalidate: "onSubmit",

    onValidate({ formData }) {
      return parseWithZod(formData, {
        schema: SendInviteFormSchema,
      });
    },
  });

  return (
    <Dialog open={isOpen} onOpenChange={(value) => setIsOpen(value)}>
      <DialogTrigger asChild>
        <Button className="!bg-red-500">Delete</Button>
      </DialogTrigger>
      <DialogContent
        leading={<DialogTitle>Delete App</DialogTitle>}
        aria-describedby="Delete App"
      >
        <fetcher.Form
          method="post"
          className="grid grid-cols-1 gap-3"
          onSubmit={form.onSubmit}
        >
          <FormField>
            {/* <Label>Email</Label> */}
            <Input
              key={fields.email.key}
              name={fields.email.name}
              placeholder="App Name"
              defaultValue={fields.email.initialValue || ""}
            />
            <FormError>{fields.email.errors}</FormError>
          </FormField>

          <div>
            <Paragraph>
              Type in the name of the app{" "}
              <code className="py-0.5 px-1 border-border border">
                {app.name.toLowerCase()}
              </code>
              . Note, this action is irreversible.
            </Paragraph>
          </div>

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
              value="create"
              className="!bg-red-500"
            >
              Delete
            </Button>
          </div>
        </fetcher.Form>
      </DialogContent>
    </Dialog>
  );
}
