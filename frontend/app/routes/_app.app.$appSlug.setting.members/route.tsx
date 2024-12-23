import { LoaderFunctionArgs } from "@remix-run/node";
import React from "react";
import { Button } from "~/components/ui/button";
import { Card, CardContent, CardHeader } from "~/components/ui/card";
import { Table, TableBody, TableCell, TableRow } from "~/components/ui/table";
import { AppSlugParamSchema } from "../_app.app.$appSlug/route";
import { typedjson, useTypedLoaderData } from "remix-typedjson";
import { getMembers } from "~/services/member.server";
import { UserAvatar } from "~/components/layout/user-menu";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "~/components/ui/dropdown-menu";
import { useApp } from "~/hooks/use-apps";
import { CircleEllipsis } from "lucide-react";
import InviteMemberDialog from "~/components/member/invite-member-dialog";

export const loader = async ({ request, params }: LoaderFunctionArgs) => {
  const { appSlug } = AppSlugParamSchema.parse(params);

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
                  <TableRow>
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
                            <DropdownMenuItem>
                              <DropdownMenuLabel>Remove</DropdownMenuLabel>
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
