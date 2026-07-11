import { parseWithZod } from "@conform-to/zod";
import { ActionFunction,json,LoaderFunctionArgs } from "@remix-run/node";
import { Form } from "@remix-run/react";
import { redirect,typedjson,useTypedLoaderData } from "remix-typedjson";
import AnimatedLogo from "~/components/animated-logo";
import CreateAppDialog from "~/components/app/create-app-dialog";
import UserMenu from "~/components/layout/user-menu";
import { Card,CardContent,CardHeader } from "~/components/ui/card";
import Paragraph from "~/components/ui/paragraph";
import { useUser } from "~/hooks/use-user";
import { appPath } from "~/lib/path";
import { CreateAppFormSchema } from "~/models/app";
import {
commitSession,
createApp,
getApps,
setAppSession,
} from "~/services/app.server";

export const action: ActionFunction = async ({ request }) => {
  const formData = await request.formData();
  const submission = parseWithZod(formData, { schema: CreateAppFormSchema });

  if (submission.status !== "success") {
    return json(submission.reply());
  }

  try {
    const app = await createApp(request, submission.value);

    const session = await setAppSession(request, app.data.id);

    const headers = new Headers({ "Set-Cookie": await commitSession(session) });

    return redirect(appPath(app.data.slug), {
      headers,
    });
  } catch (e) {
    return json(submission.reply());
  }
};

export const loader = async ({ request }: LoaderFunctionArgs) => {
  const apps = await getApps(request);

  return typedjson({
    apps,
  });
};

export default function Page() {
  const user = useUser();

  const { apps } = useTypedLoaderData<typeof loader>();

  return (
    <div className="mx-6 md:mx-10 flex flex-col">
      <div className="flex flex-col">
        <div className="flex items-center border-border border-b pt-4 pb-2 justify-between mb-3">
          <h1 className="text-2xl font-bold">Apps</h1>
          <div className="flex items-center gap-2">
            <CreateAppDialog />
            <UserMenu user={user} isSidebar={false} />
          </div>
        </div>

        {apps.length > 0 ? (
          <div className="grid grid-cols-1 md:grid-cols-3 lg:grid-cols-4 gap-8 md:gap-10">
            {apps.map((app) => (
              <Form key={app.id} action={appPath(app.slug)} method="POST">
                <input type="hidden" name="id" value={app.id} />
                <button type="submit" className="w-full">
                  <Card key={app.slug}>
                    <CardHeader className="bg-background">
                      <AnimatedLogo className="text-5xl" />
                    </CardHeader>
                    <CardContent className="!pt-2">
                      <div className="flex flex-col items-start">
                        <h1 className="text-xl font-semibold">{app.name}</h1>
                        <div className="flex divide-x divide-border">
                          <Paragraph className="pr-2">{app.bucket}</Paragraph>
                          <Paragraph className="pl-2">{app.region}</Paragraph>
                        </div>
                      </div>
                    </CardContent>
                  </Card>
                </button>
              </Form>
            ))}
          </div>
        ) : (
          <div className="flex">
            <Paragraph>No app found, create one to get started.</Paragraph>
          </div>
        )}
      </div>
    </div>
  );
}
