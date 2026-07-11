import { parseWithZod } from "@conform-to/zod";
import { ActionFunctionArgs,json,LoaderFunctionArgs,redirect } from "@remix-run/node";
import { CircleEllipsis } from "lucide-react";
import { typedjson,useTypedLoaderData } from "remix-typedjson";
import { UserAvatar } from "~/components/layout/user-menu";
import InviteMemberDialog from "~/components/member/invite-member-dialog";
import { Card,CardContent,CardHeader } from "~/components/ui/card";
import {
DropdownMenu,
DropdownMenuContent,
DropdownMenuItem,
DropdownMenuTrigger,
} from "~/components/ui/dropdown-menu";
import { Table,TableBody,TableCell,TableRow } from "~/components/ui/table";
import { useApp } from "~/hooks/use-apps";
import { settingsMenuPath } from "~/lib/path";
import { RemoveMemberFormSchema,SendInviteFormSchema } from "~/models/member";
import { createMember,getMembers,removeMember } from "~/services/member.server";
import { AppSlugParamSchema } from "../_app.app.$appSlug/route";

export const action = async ({ request, params }: ActionFunctionArgs) => {
  const { appSlug } = AppSlugParamSchema.parse(params);
  const formData = await request.formData();
  const intent = formData.get("intent");
  if (intent === "remove") {
    const submission = parseWithZod(formData, { schema: RemoveMemberFormSchema });
    if (submission.status !== "success") return json(submission.reply());
    await removeMember(request, submission.value.memberId);
  } else {
    const submission = parseWithZod(formData, { schema: SendInviteFormSchema });
    if (submission.status !== "success") return json(submission.reply());
    try {
      await createMember(request, submission.value);
    } catch (error) {
    return json({ error: error instanceof Error ? error.message : "Unable to add member" });
    }
  }
  return redirect(settingsMenuPath(appSlug, "members"));
};

export const loader = async ({ request, params }: LoaderFunctionArgs) => {
  AppSlugParamSchema.parse(params);

  const members = await getMembers(request);

  return typedjson({
    members,
  });
};

export default function Page() {
  const { members } = useTypedLoaderData<typeof loader>();
  const app = useApp();

  return (
    <div className="flex flex-col w-full max-w-4xl">
      <Card className="mb-8">
        <CardHeader className="flex flex-row items-center justify-between">
          <h1 className="font-semibold">Members</h1>
          <InviteMemberDialog />
        </CardHeader>
        <CardContent>
          <Table>
            <TableBody>
              {members.map((member) => {
                return (
                  <TableRow key={member.id}>
                    <TableCell className="font-medium">
                      <div className="flex gap-2">
                        <UserAvatar user={member.user} />
                        <div className="flex flex-col">
                          <p className="font-semibold">{member.user.name}</p>
                          <p className="text-xs text-muted-foreground">
                            {member.user.email}
                          </p>
                        </div>
                      </div>
                    </TableCell>
                    <TableCell className="capitalize">{member.role}</TableCell>

                    {member.user_id !== app.owner_id && (
                      <TableCell className="w-[100px]">
                        <DropdownMenu>
                          <DropdownMenuTrigger asChild>
                            <button className="!h-8">
                              <CircleEllipsis className="h-4 w-4" />
                            </button>
                          </DropdownMenuTrigger>
                          <DropdownMenuContent>
                            <DropdownMenuItem asChild>
                              <form method="post">
                                <input type="hidden" name="intent" value="remove" />
                                <input type="hidden" name="memberId" value={member.id} />
                                <button type="submit" className="w-full text-left">Remove</button>
                              </form>
                            </DropdownMenuItem>
                          </DropdownMenuContent>
                        </DropdownMenu>
                      </TableCell>
                    )}
                  </TableRow>
                );
              })}
            </TableBody>
          </Table>
        </CardContent>
      </Card>
    </div>
  );
}
