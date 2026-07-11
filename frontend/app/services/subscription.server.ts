import { GetSubscriptionSchema } from "~/models/subscription";
import { api } from "./api.server";

export async function getSubscription(request: Request) {
  const response = await api.get({ request, path: "/ui/subscription", schema: GetSubscriptionSchema });
  return response.data;
}
