import React, { useMemo } from "react";
import { Button } from "../ui/button";
import { useFetcher } from "@remix-run/react";
import { useForm } from "@conform-to/react";
import { parseWithZod } from "@conform-to/zod";
import FormField from "../ui/form-field";
import { Label } from "../ui/label";
import FormError from "../ui/form-error";
import { Input, InputGroup } from "../ui/input";
import Paragraph from "../ui/paragraph";
import { CreateOrUpdateDomainFormSchema, Domain } from "~/models/domain";

export default function AppDomainForm({ domain }: { domain: Domain | null }) {
  const fetcher = useFetcher();
  const [url, setUrl] = React.useState(
    domain?.domain.replace("https://", "") ?? ""
  );

  const domainInfo = useMemo(() => {
    // return subdomain and domain with tld
    if (domain) {
      const [subDomain, ...rest] = domain.domain.split(".");
      return { subDomain, domain: rest.join(".") };
    }
    return { subDomain: "", domain: "" };
  }, [url]);

  const [form, fields] = useForm({
    id: "create-or-update-domain-form",
    shouldValidate: "onBlur",
    shouldRevalidate: "onSubmit",

    onValidate({ formData }) {
      return parseWithZod(formData, {
        schema: CreateOrUpdateDomainFormSchema,
      });
    },
  });

  return (
    <fetcher.Form
      method="post"
      className="grid grid-cols-1 gap-3"
      onSubmit={form.onSubmit}
    >
      <input type="hidden" name={"type"} value={"domain"} />
      <input
        type="hidden"
        key={fields.domain.key}
        name={fields.domain.name}
        defaultValue={`https://${url?.trim()}`}
      />
      <FormField>
        <InputGroup
          leading={<Label>https://</Label>}
          defaultValue={url || ""}
          onChange={(e) => setUrl(e.target.value)}
        />
        <FormError>{fields.domain.errors}</FormError>
      </FormField>
      <div>
        <Paragraph>
          Note: This domain will be attach to the app engine that serve the
          files publicly.
        </Paragraph>
      </div>
      {domain && (
        <div className="flex flex-col gap-3">
          <Paragraph>
            Please set the following CNAME record on
            <code className="mx-1 py-0.5 px-1 border-border border">
              {domainInfo.domain}
            </code>
            to prove ownership of
            <code className="ml-1 py-0.5 px-1 border-border border">{url}</code>
          </Paragraph>
          <div className="grid grid-cols-[repeat(3,min-content)] items-end gap-x-10 gap-y-1 border p-2">
            <Paragraph>Type</Paragraph>
            <Paragraph>Name</Paragraph>
            <Paragraph>Value</Paragraph>
            <Paragraph>CNAME</Paragraph>
            <Paragraph>{domainInfo.subDomain}</Paragraph>
            <Paragraph>{domain.cname_record}</Paragraph>
          </div>
        </div>
      )}
      <div className="flex gap-2">
        <Button
          type="submit"
          name="intent"
          value={domain ? "update" : "create"}
        >
          {domain ? "Update" : "Create"}
        </Button>
        <Button
          type="submit"
          disabled={domain == null}
          name="intent"
          value="refresh"
        >
          Refresh
        </Button>
      </div>
    </fetcher.Form>
  );
}
