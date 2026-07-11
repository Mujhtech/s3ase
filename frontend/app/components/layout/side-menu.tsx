import {
Folder,
Gauge,
LockKeyhole,
Settings,
Webhook
} from "lucide-react";
import React from "react";

import { Link,useLocation } from "@remix-run/react";
import {
apiKeysPath,
appPath,
filesPath,
settingsPath,
webhooksPath
} from "~/lib/path";
import { cn } from "~/lib/utils";
import { User } from "~/models/user";
import AnimatedLogo from "../animated-logo";
import UserMenu from "./user-menu";

export default function SideMenu({
  user,
  appSlug,
}: {
  user: User;
  appSlug: string;
}) {
  return (
    <div className="h-full overflow-hidden hidden md:flex flex-col bg-background">
      <div className="flex flex-col justify-between h-full">
        <div className="flex flex-col px-4 mt-6">
          <Link to={appPath(appSlug)} className="relative flex items-center md:justify-center">
            <AnimatedLogo />
          </Link>
          <nav className="mt-8">
            <ul className="flex flex-col gap-3">
              <MenuItem
                to={appPath(appSlug)}
                icon={<Gauge className="w-6 h-6" />}
              />
              <MenuItem
                to={filesPath(appSlug)}
                icon={<Folder className="w-6 h-6" />}
              />
              <MenuItem
                to={webhooksPath(appSlug)}
                icon={<Webhook className="w-6 h-6" />}
              />
              <MenuItem
                to={apiKeysPath(appSlug)}
                icon={<LockKeyhole className="w-6 h-6" />}
              />
              <MenuItem
                to={settingsPath(appSlug)}
                icon={<Settings className="w-6 h-6" />}
              />
            </ul>
          </nav>
        </div>
        <div className="px-4 mb-4">
          <UserMenu user={user} isSidebar={true} />
        </div>
      </div>
    </div>
  );
}

const MenuItem = ({ to, icon }: { to: string; icon: React.ReactNode }) => {
  const location = useLocation();

  let isActive = location.pathname == to;

  if (to.includes("setting") && location.pathname.includes("setting")) {
    isActive = true;
  }

  if (
    to.includes("files") &&
    (location.pathname.includes("folder") || location.pathname.includes("file"))
  ) {
    isActive = true;
  }

  return (
    <li>
      <Link
        to={to}
        className={cn(
          "relative p-2 flex items-center md:justify-center hover:bg-background hover:border-border hover:border",
          isActive && "bg-background border-border border"
        )}
      >
        {icon}
      </Link>
    </li>
  );
};
