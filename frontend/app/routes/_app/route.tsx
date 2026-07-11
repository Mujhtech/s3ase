import { LoaderFunctionArgs } from "@remix-run/node";
import { Outlet } from "@remix-run/react";
import { typedjson } from "remix-typedjson";
import { clearRedirectTo,commitSession } from "~/services/redirect-to.server";
import { requireUser } from "~/services/user.server";

export const loader = async ({ request }: LoaderFunctionArgs) => {
  await requireUser(request);

  //you have to confirm basic details before you can do anything
  //   if (!user.confirmedBasicDetails) {
  //     return redirect(confirmBasicDetailsPath());
  //   }

  return typedjson(
    {},
    {
      headers: {
        "Set-Cookie": await commitSession(await clearRedirectTo(request)),
      },
    }
  );
};

export default function App() {
  return <Outlet />;
}
