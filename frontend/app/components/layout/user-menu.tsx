import {
DropdownMenu,
DropdownMenuContent,
DropdownMenuItem,
DropdownMenuLabel,
DropdownMenuSeparator,
DropdownMenuTrigger,
} from "~/components/ui/dropdown-menu";

import { Link } from "@remix-run/react";
import { Avatar,AvatarFallback,AvatarImage } from "~/components/ui/avatar";
import { appsPath,logoutPath } from "~/lib/path";
import { cn } from "~/lib/utils";
import { User } from "~/models/user";

export default function UserMenu({
  user,
  isSidebar,
}: {
  user: User;
  isSidebar: boolean;
}) {
  return (
    <DropdownMenu>
      <DropdownMenuTrigger>
        <UserAvatar user={user} />
      </DropdownMenuTrigger>
      <DropdownMenuContent className={cn(isSidebar ? "ml-2" : "mr-2")}>
        <DropdownMenuLabel>My Account</DropdownMenuLabel>
        <DropdownMenuSeparator />
        {isSidebar && (
          <>
            <DropdownMenuItem>
              <Link to={appsPath()}>Apps</Link>
            </DropdownMenuItem>
            <DropdownMenuItem>Profile</DropdownMenuItem>
            <DropdownMenuItem>Billing</DropdownMenuItem>
          </>
        )}
        <DropdownMenuItem>
          <Link to={logoutPath()}>Logout</Link>
        </DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  );
}

export const UserAvatar = ({ user }: { user: User }) => {
  return (
    <Avatar className="!rounded-none p-1 flex items-center md:justify-center border border-border bg-background hover:bg-background hover:border-border hover:border">
      {user.avatar_url && <AvatarImage src={user?.avatar_url} />}
      <AvatarFallback>{user.name[0]}</AvatarFallback>
    </Avatar>
  );
};
