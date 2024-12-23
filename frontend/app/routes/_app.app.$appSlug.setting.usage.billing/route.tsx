import React from "react";
import {
  EnterpisePlanCard,
  BasicPlanCard,
  ProPlanCard,
} from "~/components/usage/plan-card";

export default function Page() {
  return (
    <div className="flex flex-col">
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        <BasicPlanCard />
        <ProPlanCard />
        <EnterpisePlanCard />
      </div>
    </div>
  );
}
