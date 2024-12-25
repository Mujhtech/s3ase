import { describe, it, expect, beforeEach } from "vitest";
import {
  getApiKeys,
  createApiKey,
  updateApiKey,
  deleteApiKey,
} from "~/services/api_key.server";
import { mockFetch, mockRequest, setupFetchMock } from "../utils/test-setup";
import { ApiKeys, CreateApiKeyForm } from "~/models/api_key";

describe("API Key Service", () => {
  const mockApiKeys: ApiKeys = [
    {
      id: "1",
      name: "Test API Key",
      description: "Test description",
      access: "create",
      created_at: "",
      updated_at: "",
    },
  ];

  beforeEach(() => {
    setupFetchMock(mockFetch(200, mockApiKeys));
  });

  it("should get API keys", async () => {
    const result = await getApiKeys(mockRequest);
    expect(result).toEqual(mockApiKeys);
  });

  it("should create API key", async () => {
    const newKey: CreateApiKeyForm = {
      name: "New API Key",
      access: "read",
      intent: "create",
      description: "New description",
    };
    const result = await createApiKey(mockRequest, newKey);
    expect(result).toEqual([mockApiKeys[0]]);
  });

  it("should update API key", async () => {
    setupFetchMock(mockFetch(200, { message: "API key updated" }));
    const updatedKey: CreateApiKeyForm = {
      name: "Updated API Key",
      access: "read",
      intent: "update",
      description: "Updated description",
    };
    const result = await updateApiKey(mockRequest, "1", updatedKey);
    expect(result).toBeDefined();
  });

  it("should delete API key", async () => {
    setupFetchMock(mockFetch(200, { message: "API key deleted" }));
    const result = await deleteApiKey(mockRequest, "1");
    expect(result).toBeDefined();
  });
});
