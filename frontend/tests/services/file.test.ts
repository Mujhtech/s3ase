import { describe, it, expect, beforeEach } from "vitest";
import { mockFetch, mockRequest, setupFetchMock } from "../utils/test-setup";
import { getFiles } from "~/services/file.server";
import { Files } from "~/models/file";

describe("File Service", () => {
  const mockFiles: Files = [
    {
      id: "1",
      name: "test.txt",
      created_at: "",
      updated_at: "",
    },
  ];

  beforeEach(() => {
    setupFetchMock(mockFetch(200, mockFiles));
  });

  it("should get files", async () => {
    const result = await getFiles(mockRequest);
    expect(result).toEqual(mockFiles);
  });
});
