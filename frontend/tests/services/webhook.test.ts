import { describe, it, expect, beforeEach } from "vitest";
import {
  getWebhooks,
  createWebhook,
  updateWebhook,
  deleteWebhook,
} from "~/services/webhook.server";
import { mockFetch, mockRequest, setupFetchMock } from "../utils/test-setup";

describe("Webhook Service", () => {
  const mockWebhooks = [{ id: "1", url: "https://test.com" }];

  beforeEach(() => {
    setupFetchMock(mockFetch(200, mockWebhooks));
  });

  it("should get webhooks", async () => {
    const result = await getWebhooks(mockRequest);
    expect(result).toEqual(mockWebhooks);
  });

  it("should create webhook", async () => {
    const result = await createWebhook(mockRequest, {
      url: "https://test.com",
      name: "Text",
      intent: "create",
      events: [],
    });
    expect(result).toEqual([mockWebhooks[0]]);
  });

  it("should update webhook", async () => {
    setupFetchMock(mockFetch(200, { message: "Webhook updated" }));
    const result = await updateWebhook(mockRequest, "1", {
      url: "https://test.com",
      name: "Text",
      intent: "create",
      events: [],
    });
    expect(result).toBeDefined();
  });

  it("should delete webhook", async () => {
    setupFetchMock(mockFetch(200, { message: "Webhook deleted" }));
    const result = await deleteWebhook(mockRequest, "1");
    expect(result).toBeDefined();
  });
});
