import { LoaderFunctionArgs } from "@remix-run/node";
import { Outlet } from "@remix-run/react";
import { typedjson,useTypedLoaderData } from "remix-typedjson";
import { usageMenuPath,usagesMenuPath } from "~/lib/path";
import { SettingMenuNavItem } from "../_app.app.$appSlug.setting/route";
import { AppSlugParamSchema } from "../_app.app.$appSlug/route";

export const loader = async ({ params }: LoaderFunctionArgs) => {
  const { appSlug } = AppSlugParamSchema.parse(params);

  return typedjson({
    appSlug,
  });
};

export default function Usage() {
  const { appSlug } = useTypedLoaderData<typeof loader>();

  return (
    <div className="flex flex-col">
      <UsageMenuNav appSlug={appSlug} />
      <Outlet />
    </div>
  );
}

const UsageMenuNav = ({ appSlug }: { appSlug: string }) => {
  return (
    <nav className="max-w-md mb-4">
      <ul className="flex gap-3">
        <SettingMenuNavItem title="Usage" to={usageMenuPath(appSlug)} />
        <SettingMenuNavItem
          title="Billing"
          to={usagesMenuPath(appSlug, "billing")}
        />
      </ul>
    </nav>
  );
};
