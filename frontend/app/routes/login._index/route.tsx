import { SiGithub, SiGoogle } from "@icons-pack/react-simple-icons";
import React from "react";
import { Boxes } from "~/components/login-bg-boxes";
import { Button } from "~/components/ui/button";
import { Form } from "@remix-run/react";
import { LoaderFunctionArgs } from "@remix-run/node";
import { redirect, typedjson, useTypedLoaderData } from "remix-typedjson";
import { commitSession, setRedirectTo } from "~/services/redirect-to.server";
import { requestUrl } from "~/services/request-url.server";

export async function loader({ request }: LoaderFunctionArgs) {
  // redirect user to home is already logged in

  const url = requestUrl(request);
  const redirectTo = url.searchParams.get("redirectTo");

  if (redirectTo) {
    const session = await setRedirectTo(request, redirectTo);

    return typedjson(
      {
        redirectTo,
      },
      {
        headers: {
          "Set-Cookie": await commitSession(session),
        },
      }
    );
  } else {
    return typedjson({
      redirectTo: null,
    });
  }
}

export default function Page() {
  const data = useTypedLoaderData<typeof loader>();

  return (
    <main className="h-full text-white">
      <div className="relative h-full ">
        <div className="absolute inset-0 w-full h-full z-20 [mask-image:radial-gradient(transparent,white)] pointer-events-none" />
        <Boxes />

        <div className="flex justify-center w-full h-full ">
          <div className="m-auto min-w-[400px] flex-col flex items-center bg-black z-30 ">
            <div className="flex flex-col px-12 py-20 w-full">
              <h1 className="text-4xl font-sans font-bold">Welcome</h1>
              <p className="text-sm">Create an account or login to continue</p>
              <div className="flex flex-col gap-y-3 mt-6">
                <Form
                  action={`/auth/github${
                    data.redirectTo ? `?redirectTo=${data.redirectTo}` : ""
                  }`}
                  method="post"
                  className="w-full"
                >
                  <Button className="w-full">
                    <SiGithub className="mr-3 h-5 w-5" /> Continue with Github
                  </Button>
                </Form>
                <Form
                  action={`/auth/google${
                    data.redirectTo ? `?redirectTo=${data.redirectTo}` : ""
                  }`}
                  method="post"
                  className="w-full"
                >
                  <Button className="w-full">
                    <SiGoogle className="mr-3 h-5 w-5" />
                    Continue with Google
                  </Button>
                </Form>
              </div>
            </div>
          </div>
        </div>
      </div>
    </main>
  );
}
