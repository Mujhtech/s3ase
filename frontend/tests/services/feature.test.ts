import { describe, it, expect, beforeEach } from "vitest";
import { mockFetch, mockRequest, setupFetchMock } from "../utils/test-setup";
import { getFeatures } from "~/services/feature.server";

describe("Feature Service", () => {
  const mockFeatures = {
    name: "S3ase",
    description: "Simplifying S3 usage through open source.",
  };

  beforeEach(() => {
    setupFetchMock(mockFetch(200, mockFeatures));
  });

  it("should get features", async () => {
    const result = await getFeatures(mockRequest);
    expect(result).toEqual(mockFeatures);
  });
});
