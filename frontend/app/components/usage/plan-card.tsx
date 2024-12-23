import React from "react";
import { Card, CardContent, CardHeader } from "../ui/card";
import { Form } from "@remix-run/react";
import { Button } from "../ui/button";
import { CheckIcon, X } from "lucide-react";
import { cn } from "~/lib/utils";
import Paragraph from "../ui/paragraph";
import { Separator } from "../ui/separator";

export function BasicPlanCard() {
  return (
    <PlanCard
      title="Basic"
      price={1}
      storage={2}
      feature={
        <ul className="flex flex-col gap-3">
          <Item checked>Shared Storage</Item>
          <Item checked>Engine</Item>
          <Item checked>Unlimited upload</Item>
          <Item checked>Unlimited download</Item>
          <Item checked>Unlimited team member</Item>
          <Item checked>Unlimited regions</Item>
          <Item checked>Community support</Item>
          <Item checked={false}>Bucket Migration</Item>
          <Item checked={false}>Webhook</Item>
          <Item checked={false}>SSO</Item>
        </ul>
      }
    />
  );
}

export function ProPlanCard() {
  return (
    <PlanCard
      recommended
      title="Pro"
      price={30}
      storage={100}
      feature={
        <ul className="flex flex-col gap-3">
          <Item checked>Dedicated Storage</Item>
          <Item checked>Engine</Item>
          <Item checked>Unlimited upload</Item>
          <Item checked>Unlimited download</Item>
          <Item checked>Unlimited team member</Item>
          <Item checked>Unlimited regions</Item>
          <Item checked>Dedicated support</Item>
          <Item checked>Bucket Migration</Item>
          <Item checked>Webhook</Item>
          <Item checked={false}>SSO</Item>
        </ul>
      }
    />
  );
}

export function EnterpisePlanCard() {
  return (
    <PlanCard
      title="Enterprise"
      feature={
        <ul className="flex flex-col gap-3">
          <Item checked>Dedicated Storage</Item>
          <Item checked>Engine</Item>
          <Item checked>Unlimited upload</Item>
          <Item checked>Unlimited download</Item>
          <Item checked>Unlimited team member</Item>
          <Item checked>Unlimited regions</Item>
          <Item checked>Priority support</Item>
          <Item checked>Bucket Migration</Item>
          <Item checked>Webhook</Item>
          <Item checked>SSO</Item>
        </ul>
      }
    />
  );
}

export default function PlanCard({
  recommended,
  title,
  feature,
  price,
  storage,
}: {
  recommended?: boolean;
  title: string;
  feature: React.ReactNode;
  price?: number;
  storage?: number;
}) {
  return (
    <Card className={cn("w-full", recommended && "border-white border-2")}>
      <Form method="post">
        <CardHeader className="!gap-3">
          <h3 className="font-medium text-muted-foreground">{title}</h3>
          <h1 className="text-4xl font-semibold">
            {price != undefined ? `$${price}` : "Custom"}
            <span className="text-sm font-medium text-muted-foreground">
              /month
            </span>
          </h1>
        </CardHeader>
        <CardContent className="!pb-16">
          <Paragraph className="mb-4 !text-sm">
            {storage ? `Up to ${storage}gb storage` : "Custom storage"}
          </Paragraph>
          <Separator />
          <Button className="my-5 !h-8 w-full" disabled={true}>
            Upgrade
          </Button>
          {feature}
        </CardContent>
      </Form>
    </Card>
  );
}

const Item = ({
  checked,
  children,
}: {
  checked: boolean;
  children: React.ReactNode;
}) => {
  return (
    <li className="flex items-center gap-2">
      {checked ? (
        <CheckIcon className="h-4 w-4 text-primary" />
      ) : (
        <X className="h-4 w-4 text-muted-foreground" />
      )}
      <div className={cn("text-sm", checked ? "" : "text-muted-foreground")}>
        {children}
      </div>
    </li>
  );
};
