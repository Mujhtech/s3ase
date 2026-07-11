import { describe, expect, test } from "vitest";

import { fileDownloadPath, filePath, folderPath } from "../app/lib/path";
import { ApiKeySchema } from "../app/models/api_key";
import { FileSchema } from "../app/models/file";
import { UsageSchema } from "../app/models/usage";

describe("frontend API contracts", () => {
  test("accepts a completed root file", () => {
    const file = FileSchema.parse({
      id: "file-id",
      app_id: "app-id",
      uploaded_by: "user-id",
      folder_id: null,
      name: "report.pdf",
      mime_type: "application/pdf",
      extension: "pdf",
      size: 42,
      is_public: true,
      public_id: "public-id",
      status: "completed",
      created_at: "2026-07-11T00:00:00Z",
      updated_at: "2026-07-11T00:00:00Z",
    });

    expect(file.folder_id).toBeNull();
    expect(file.status).toBe("completed");
  });

  test("accepts one-time API key credentials", () => {
    const key = ApiKeySchema.parse({
      id: "key-id",
      name: "Production",
      description: null,
      access: "full",
      created_at: "2026-07-11T00:00:00Z",
      updated_at: "2026-07-11T00:00:00Z",
      key_prefix: "s3ase_abc123",
      last_four: "wxyz",
      secret: "s3ase_abc123wxyz",
    });

    expect(key.description).toBe("");
    expect(key.secret).toContain("s3ase_");
  });

  test("parses usage returned by the backend", () => {
    const usage = UsageSchema.parse({
      file_count: 1,
      storage_bytes: 42,
      daily: [{ date: "2026-07-11", uploads: 1, storage_bytes: 42 }],
    });

    expect(usage.daily).toHaveLength(1);
  });
});

describe("established UI paths", () => {
  test("keeps file and folder routes under the selected app", () => {
    expect(folderPath("demo", "folder-id")).toBe("/app/demo/folder/folder-id");
    expect(filePath("demo", "file-id")).toBe("/app/demo/file/file-id");
    expect(fileDownloadPath("demo", "file-id")).toBe(
      "/resources/demo/files/file-id",
    );
  });
});
