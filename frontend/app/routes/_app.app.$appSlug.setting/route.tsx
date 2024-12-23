import { LoaderFunctionArgs } from "@remix-run/node";
import { Link, Outlet, useLocation } from "@remix-run/react";
import React from "react";
import { typedjson, useTypedLoaderData } from "remix-typedjson";
import { cn } from "~/lib/utils";
import { AppSlugParamSchema } from "../_app.app.$appSlug/route";
import { settingsMenuPath, settingsPath } from "~/lib/path";

export const loader = async ({ request, params }: LoaderFunctionArgs) => {
  const { appSlug } = AppSlugParamSchema.parse(params);

  return typedjson({
    appSlug,
  });
};

export default function Setting() {
  const { appSlug } = useTypedLoaderData<typeof loader>();

  return (
    <div className="flex flex-col">
      <div className="flex flex-col mb-3">
        <div className="flex items-center border-border border-b pt-4 pb-2 justify-between">
          <h1 className="text-2xl font-bold">Setting</h1>
        </div>
      </div>
      <SettingMenuNav appSlug={appSlug} />
      <Outlet />
    </div>
  );
}

const SettingMenuNav = ({ appSlug }: { appSlug: string }) => {
  return (
    <nav className="max-w-md mb-4">
      <ul className="flex gap-3">
        <SettingMenuNavItem title="General" to={settingsPath(appSlug)} />
        <SettingMenuNavItem
          title="Usage"
          to={settingsMenuPath(appSlug, "usage")}
        />
        <SettingMenuNavItem
          title="Members"
          to={settingsMenuPath(appSlug, "members")}
        />
      </ul>
    </nav>
  );
};

export const SettingMenuNavItem = ({
  title,
  to,
}: {
  title: string;
  to: string;
}) => {
  const location = useLocation();

  const isActive = location.pathname == to;

  return (
    <li>
      <Link
        to={to}
        className={cn(
          "text-sm ",
          isActive ? "text-white font-medium" : "text-muted-foreground"
        )}
      >
        {title}
      </Link>
    </li>
  );
};
