import { describe, it, expect, beforeEach } from "vitest";
import { mockFetch, mockRequest, setupFetchMock } from "../utils/test-setup";
import { getFeatures } from "~/services/feature.server";

describe("Feature Service", () => {
  const mockFeatures = {
    name: "S3ase",
    description: "Simplifying S3 usage through open source.",
    is_github_auth_enabled: true,
    is_google_auth_enabled: true,
    is_aws_configured: true,
    is_object_storage_configured: true,
    object_storage_provider: "r2",
    object_storage_region: "auto",
    version: "1.0.0",
  };

  beforeEach(() => {
    setupFetchMock(mockFetch(200, mockFeatures));
  });

  it("should get features", async () => {
    const result = await getFeatures(mockRequest);
    expect(result).toEqual(mockFeatures);
  });
});
