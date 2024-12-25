import { describe, it, expect, beforeEach } from "vitest";
import { mockFetch, mockRequest, setupFetchMock } from "../utils/test-setup";
import {
  createFolder,
  deleteFolder,
  getFolders,
  updateFolder,
} from "~/services/folder.server";
import { CreateFolderForm, Folders } from "~/models/folder";

describe("Folder Service", () => {
  const mockFolders: Folders = [
    {
      id: "1",
      name: "Test Folder",
      description: "",
      url: "",
      created_at: "",
      updated_at: "",
    },
  ];

  beforeEach(() => {
    setupFetchMock(mockFetch(200, mockFolders));
  });

  it("should get folders", async () => {
    const result = await getFolders(mockRequest);
    expect(result).toEqual(mockFolders);
  });

  it("should create folder", async () => {
    const newFolder: CreateFolderForm = {
      name: "New Folder",
      type: "folder",
      intent: "create",
    };
    const result = await createFolder(mockRequest, newFolder);
    expect(result).toEqual([mockFolders[0]]);
  });

  it("should update folder", async () => {
    setupFetchMock(mockFetch(200, { message: "Folder updated" }));
    const result = await updateFolder(mockRequest, "1", {
      name: "New Folder",
      type: "folder",
      intent: "create",
    });
    expect(result).toBeDefined();
  });

  it("should delete folder", async () => {
    setupFetchMock(mockFetch(200, { message: "Folder deleted" }));
    const result = await deleteFolder(mockRequest, "1");
    expect(result).toBeDefined();
  });
});
