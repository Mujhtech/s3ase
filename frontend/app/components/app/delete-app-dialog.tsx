import { useForm } from "@conform-to/react";
import { parseWithZod } from "@conform-to/zod";
import { useFetcher } from "@remix-run/react";
import { useState } from "react";
import { Button } from "~/components/ui/button";
import {
Dialog,
DialogContent,
DialogTitle,
DialogTrigger,
} from "~/components/ui/dialog";
import FormError from "~/components/ui/form-error";
import FormField from "~/components/ui/form-field";
import { Input } from "~/components/ui/input";
import Paragraph from "~/components/ui/paragraph";
import { App,DeleteAppFormSchema } from "~/models/app";

export default function DeleteAppDialog({ app }: { app: App }) {
  const [isOpen, setIsOpen] = useState(false);

  const fetcher = useFetcher<{ error?: string }>();

  const [form, fields] = useForm({
    id: "delete-app-form",
    shouldValidate: "onBlur",
    shouldRevalidate: "onSubmit",

    onValidate({ formData }) {
      return parseWithZod(formData, {
        schema: DeleteAppFormSchema,
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
      <input type="hidden" name="type" value="delete-app" />
          <FormField>
            <Input
              key={fields.name.key}
              name={fields.name.name}
              placeholder="App Name"
              defaultValue={fields.name.initialValue || ""}
            />
            <FormError>{fields.name.errors}</FormError>
          </FormField>

          <div>
            <Paragraph>
              Type in the name of the app{" "}
              <code className="py-0.5 px-1 border-border border">
        {app.name}
              </code>
              . Note, this action is irreversible.
            </Paragraph>
          </div>
      {fetcher.data?.error && <FormError>{fetcher.data.error}</FormError>}

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
              value="delete"
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
