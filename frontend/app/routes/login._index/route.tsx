import { SiGithub,SiGoogle } from "@icons-pack/react-simple-icons";
import { LoaderFunctionArgs } from "@remix-run/node";
import { Form } from "@remix-run/react";
import { typedjson,useTypedLoaderData } from "remix-typedjson";
import { Boxes } from "~/components/login-bg-boxes";
import { Button } from "~/components/ui/button";
import { Input } from "~/components/ui/input";
import { publicBackendUrl } from "~/env.server";
import { useFeature } from "~/hooks/use-feature";
import { commitSession,setRedirectTo } from "~/services/redirect-to.server";
import { requestUrl } from "~/services/request-url.server";

export async function loader({ request }: LoaderFunctionArgs) {
  // redirect user to home is already logged in

  const url = requestUrl(request);
  const redirectTo = url.searchParams.get("redirectTo");
  const backendUrl = publicBackendUrl;

  if (redirectTo) {
    const session = await setRedirectTo(request, redirectTo);

    return typedjson(
      {
        backendUrl,
      },
      {
        headers: {
          "Set-Cookie": await commitSession(session),
        },
      }
    );
  } else {
    return typedjson({
      backendUrl,
    });
  }
}

export default function Page() {
  const data = useTypedLoaderData<typeof loader>();
  const { is_github_auth_enabled, is_google_auth_enabled } = useFeature();

  const isSocialAuthEnabled = is_github_auth_enabled || is_google_auth_enabled;

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
                  method="post"
                  action="/login"
                  className="grid grid-cols-1 gap-4"
                >
                  <Input placeholder="Email Address" />
                  <Button
                    className="w-full"
                    type="button"
                    onClick={() => {
                      window.location.href = `${data.backendUrl}/ui/auth/github`;
                      window.close();
                    }}
                  >
                    Continue
                  </Button>
                </Form>

                {isSocialAuthEnabled && (
                  <>
                    <p className="text-sm text-center">Or</p>

                    {is_github_auth_enabled && (
                      <Button
                        className="w-full"
                        type="button"
                        onClick={() => {
                          window.location.href = `${data.backendUrl}/ui/auth/github`;
                          window.close();
                        }}
                      >
                        <SiGithub className="mr-3 h-5 w-5" /> Continue with
                        Github
                      </Button>
                    )}

                    {is_google_auth_enabled && (
                      <Button
                        className="w-full"
                        type="button"
                        onClick={() => {
                          window.location.href = `${data.backendUrl}/ui/auth/google`;
                          window.close();
                        }}
                      >
                        <SiGoogle className="mr-3 h-5 w-5" />
                        Continue with Google
                      </Button>
                    )}
                  </>
                )}
              </div>
            </div>
          </div>
        </div>
      </div>
    </main>
  );
}
