import { LoaderFunctionArgs } from "@remix-run/node";
import { AppSlugParamSchema } from "../_app.app.$appSlug/route";
import { typedjson } from "remix-typedjson";

export const loader = async ({ request, params }: LoaderFunctionArgs) => {
  const { appSlug } = AppSlugParamSchema.parse(params);

  return typedjson({});
};

export default function Page() {
  return <div className="text-white">route</div>;
}
