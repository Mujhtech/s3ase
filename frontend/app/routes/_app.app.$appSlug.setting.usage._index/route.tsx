import { LoaderFunctionArgs } from "@remix-run/node";
import { typedjson,useTypedLoaderData } from "remix-typedjson";
import { UsageChart } from "~/components/usage/usage-chart";
import { getUsage } from "~/services/usage.server";

export const loader = async ({ request }: LoaderFunctionArgs) => {
  const usage = await getUsage(request);
  return typedjson({ usage });
};

export default function Page() {
  const { usage } = useTypedLoaderData<typeof loader>();
  return (
    <div className="flex flex-col">
    <UsageChart usage={usage} />
    </div>
  );
}
