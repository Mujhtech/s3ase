import { useForm } from "@conform-to/react";
import { parseWithZod } from "@conform-to/zod";
import { useFetcher } from "@remix-run/react";
import React,{ useMemo } from "react";
import { CreateOrUpdateDomainFormSchema,Domain } from "~/models/domain";
import { Button } from "../ui/button";
import FormError from "../ui/form-error";
import FormField from "../ui/form-field";
import { InputGroup } from "../ui/input";
import { Label } from "../ui/label";
import Paragraph from "../ui/paragraph";

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
  }, [domain]);

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
    value={`https://${url?.trim()}`}
    readOnly
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
      Add both DNS records below on
            <code className="mx-1 py-0.5 px-1 border-border border">
              {domainInfo.domain}
            </code>
            to prove ownership of
            <code className="ml-1 py-0.5 px-1 border-border border">{url}</code>
          </Paragraph>
      <div className="grid grid-cols-[repeat(3,min-content)] items-end gap-x-10 gap-y-1 border p-2 overflow-x-auto">
            <Paragraph>Type</Paragraph>
            <Paragraph>Name</Paragraph>
            <Paragraph>Value</Paragraph>
            <Paragraph>CNAME</Paragraph>
            <Paragraph>{domainInfo.subDomain}</Paragraph>
            <Paragraph>{domain.cname_record}</Paragraph>
      <Paragraph>TXT</Paragraph>
      <Paragraph>_s3ase</Paragraph>
      <Paragraph>{domain.txt_record}</Paragraph>
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
