import { describe, it, expect, beforeEach } from "vitest";
import { mockFetch, mockRequest, setupFetchMock } from "../utils/test-setup";
import { getFiles } from "~/services/file.server";
import { Files } from "~/models/file";

describe("File Service", () => {
  const mockFiles: Files = [
    {
      id: "1",
      name: "test.txt",
      app_id: "app-1",
      folder_id: null,
      mime_type: "text/plain",
      extension: "txt",
      size: 9,
      is_public: true,
      public_id: "public-file-1",
      status: "completed",
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
