import { describe, it, expect, beforeEach } from "vitest";
import {
  getWebhooks,
  createWebhook,
  updateWebhook,
  deleteWebhook,
} from "~/services/webhook.server";
import { mockFetch, mockRequest, setupFetchMock } from "../utils/test-setup";

describe("Webhook Service", () => {
  const mockWebhooks = [
    {
      id: "1",
      name: "Test webhook",
      description: "Test webhook description",
      url: "https://test.com",
      created_at: "",
      updated_at: "",
      metadata: { events: [] },
    },
  ];

  beforeEach(() => {
    setupFetchMock(mockFetch(200, mockWebhooks));
  });

  it("should get webhooks", async () => {
    const result = await getWebhooks(mockRequest);
    expect(result).toEqual(mockWebhooks);
  });

  it("should create webhook", async () => {
    setupFetchMock(mockFetch(201, mockWebhooks[0]));
    const result = await createWebhook(mockRequest, {
      url: "https://test.com",
      name: "Text",
      intent: "create",
      events: [],
    });
    expect(result).toEqual(mockWebhooks[0]);
  });

  it("should update webhook", async () => {
    setupFetchMock(mockFetch(200, null, "Webhook updated"));
    const result = await updateWebhook(mockRequest, "1", {
      url: "https://test.com",
      name: "Text",
      intent: "create",
      events: [],
    });
    expect(result).toBe("Webhook updated");
  });

  it("should delete webhook", async () => {
    setupFetchMock(mockFetch(200, null, "Webhook deleted"));
    const result = await deleteWebhook(mockRequest, "1");
    expect(result).toBe("Webhook deleted");
  });
});
