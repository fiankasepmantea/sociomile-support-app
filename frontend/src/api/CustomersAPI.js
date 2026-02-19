import client from "./client";

export const CustomersAPI = {
  list: () => client.get("/customers"),
  get: (id) => client.get(`/customers/${id}`),
  create: (data) => client.post("/customers", data),
  update: (id, data) => client.put(`/customers/${id}`, data),
  delete: (id) => client.delete(`/customers/${id}`),
};
