import client from "./client";

export const ConversationsAPI = {
  list: () => client.get("/conversations"),
  get: (id) => client.get(`/conversations/${id}`),
  assign: (id, agentId) => client.patch(`/conversations/${id}/assign`, { agent_id: agentId }),
  updateStatus: (id, status) => client.patch(`/conversations/${id}/status`, { status }),
  sendMessage: (id, message) => client.post(`/conversations/${id}/messages`, { message }),
  listMessages: (id) => client.get(`/conversations/${id}/messages`),
};
