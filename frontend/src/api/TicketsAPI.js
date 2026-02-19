import client from "./client";

export const TicketsAPI = {
  list: () => client.get("/tickets"),
  get: (id) => client.get(`/tickets/${id}`),
  create: (data) => client.post("/tickets", data),
  updateStatus: (id, status) => client.patch(`/tickets/${id}/status`, { status }),
};
