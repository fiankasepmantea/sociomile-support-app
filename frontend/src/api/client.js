import axios from "axios";

const client = axios.create({
  baseURL: "http://localhost:8080/api",
});

client.interceptors.request.use((config) => {
  const token = localStorage.getItem("token");
  const tenant = localStorage.getItem("tenant") || "default";

  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  config.headers["X-Tenant-ID"] = tenant;
  return config;
});

export default client;
