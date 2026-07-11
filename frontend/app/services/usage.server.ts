import { GetUsageSchema } from "~/models/usage";
import { api } from "./api.server";

export async function getUsage(request: Request) {
  const response = await api.get({ request, path: "/ui/usage", schema: GetUsageSchema });
  return response.data;
}
