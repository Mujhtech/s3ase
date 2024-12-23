import { createCookie, type LoaderFunction, redirect } from "@remix-run/node";
import { commitSession, setAuthSession } from "~/services/auth.server";

export let loader: LoaderFunction = async ({ request }) => {
  const url = new URL(request.url);
  const token = url.searchParams.get("token");

  const cookie = request.headers.get("Cookie");
  const redirectValue = await redirectCookie.parse(cookie);
  const redirectTo = redirectValue ?? "/";

  if (!token) {
    return redirect("/login");
  }

  const session = await setAuthSession(request, token);

  let headers = new Headers({ "Set-Cookie": await commitSession(session) });

  return redirect(redirectTo, {
    headers,
  });
};

export const redirectCookie = createCookie("redirect-to", {
  maxAge: 60 * 60, // 1 hour
  httpOnly: true,
});
