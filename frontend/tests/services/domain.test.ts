import { describe, it, expect, beforeEach } from "vitest";
import { mockFetch, mockRequest, setupFetchMock } from "../utils/test-setup";
import { createOrUpdateDomain, getDomain } from "~/services/domain.server";
import { CreateOrUpdateDomainForm, Domain } from "~/models/domain";

describe("Domain Service", () => {
  const mockDomain: Domain = {
    domain: "test.com",
    created_at: "",
    updated_at: "",
    id: "1",
    txt_record: "",
    cname_record: "",
    status: "pending",
  };

  beforeEach(() => {
    setupFetchMock(mockFetch(200, mockDomain));
  });

  it("should get domains", async () => {
    const result = await getDomain(mockRequest);
    expect(result).toEqual(mockDomain);
  });

  it("should create domain", async () => {
    const newDomain: CreateOrUpdateDomainForm = {
      domain: "new.domain.com",
      intent: "create",
    };
    const result = await createOrUpdateDomain(mockRequest, newDomain);
    expect(result).toEqual(mockDomain);
  });

  it("should update domain", async () => {
    const newDomain: CreateOrUpdateDomainForm = {
      domain: "new.domain.com",
      intent: "update",
    };
    const result = await createOrUpdateDomain(mockRequest, newDomain);
    expect(result).toEqual(mockDomain);
  });

  it("should refresh domain", async () => {
    const newDomain: CreateOrUpdateDomainForm = {
      domain: "new.domain.com",
      intent: "refresh",
    };
    const result = await createOrUpdateDomain(mockRequest, newDomain);
    expect(result).toEqual(mockDomain);
  });
});
