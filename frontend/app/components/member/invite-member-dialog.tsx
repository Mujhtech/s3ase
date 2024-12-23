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

export default function InviteMemberDialog() {
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
        <Button className="!h-8">Send Invite</Button>
      </DialogTrigger>
      <DialogContent
        leading={<DialogTitle>Invite Member</DialogTitle>}
        aria-describedby="Invite Member"
      >
        <fetcher.Form
          method="post"
          className="grid grid-cols-1 gap-3"
          onSubmit={form.onSubmit}
        >
          <FormField>
            <Label>Email</Label>
            <Input
              key={fields.email.key}
              name={fields.email.name}
              defaultValue={fields.email.initialValue || ""}
            />
            <FormError>{fields.email.errors}</FormError>
          </FormField>

          <div className="flex gap-2">
            <Button
              variant="outline"
              type="button"
              onClick={() => setIsOpen(false)}
            >
              Close
            </Button>

            <Button type="submit" name="intent" value="create">
              Send Invite
            </Button>
          </div>
        </fetcher.Form>
      </DialogContent>
    </Dialog>
  );
}
