import { LoaderFunctionArgs } from "@remix-run/node";
import { typedjson,useTypedLoaderData } from "remix-typedjson";
import {
BasicPlanCard,
EnterpisePlanCard,
ProPlanCard,
} from "~/components/usage/plan-card";
import { getSubscription } from "~/services/subscription.server";

export const loader = async ({ request }: LoaderFunctionArgs) => {
  const subscription = await getSubscription(request);
  return typedjson({ subscription });
};

export default function Page() {
  const { subscription } = useTypedLoaderData<typeof loader>();
  return (
    <div className="flex flex-col">
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
    <BasicPlanCard currentPlan={subscription.plan} />
    <ProPlanCard currentPlan={subscription.plan} />
    <EnterpisePlanCard currentPlan={subscription.plan} />
      </div>
    </div>
  );
}
