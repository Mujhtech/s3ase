import { vi } from "vitest";

export function mockFetch(
	status: number,
	responseData: unknown,
	message = "request successful"
) {
	const body = JSON.stringify({ data: responseData, message });
	return vi.fn().mockImplementation(() =>
		Promise.resolve(
			new Response(body, {
				status,
				headers: { "Content-Type": "application/json" },
			})
		)
	);
}

export const mockRequest = new Request("http://localhost", {
	headers: { Authorization: "Bearer test-token" },
});

export function setupFetchMock(impl: typeof fetch) {
	global.fetch = impl;
}
