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
import { Label } from "~/components/ui/label";
import Paragraph from "~/components/ui/paragraph";
import { SendInviteFormSchema } from "~/models/member";

export default function InviteMemberDialog() {
  const [isOpen, setIsOpen] = useState(false);

  const fetcher = useFetcher<{ error?: string }>();

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
        <Button className="!h-8">Add Member</Button>
      </DialogTrigger>
      <DialogContent
        leading={<DialogTitle>Add Member</DialogTitle>}
        aria-describedby="Add Member"
      >
        <fetcher.Form
          method="post"
          className="grid grid-cols-1 gap-3"
          onSubmit={form.onSubmit}
        >
          <input type="hidden" name="role" value="member" />
          <input type="hidden" name="intent" value="create" />
          <FormField>
            <Label>Email</Label>
            <Input
              key={fields.email.key}
              name={fields.email.name}
              defaultValue={fields.email.initialValue || ""}
            />
            <FormError>{fields.email.errors}</FormError>
          </FormField>
          <Paragraph>The person must sign in to S3ase once before they can be added.</Paragraph>
      {fetcher.data?.error && <FormError>{fetcher.data.error}</FormError>}

          <div className="flex gap-2">
            <Button
              variant="outline"
              type="button"
              onClick={() => setIsOpen(false)}
            >
              Close
            </Button>

            <Button type="submit">
              Add Member
            </Button>
          </div>
        </fetcher.Form>
      </DialogContent>
    </Dialog>
  );
}
