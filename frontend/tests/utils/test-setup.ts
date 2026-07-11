import { vi } from "vitest";

export function mockFetch(status: number, responseData: any) {
  return vi.fn().mockImplementation(() =>
    Promise.resolve({
      status,
      json: () => Promise.resolve({ data: responseData }),
    })
  );
}

export const mockRequest = {
  headers: new Headers({
    Authorization: "Bearer test-token",
  }),
} as Request;

export function setupFetchMock(impl: any) {
  global.fetch = impl as any;
}
