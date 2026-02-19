import { createContext, useContext, useState } from "react";
import client from "../api/client";

const AuthContext = createContext();

export const AuthProvider = ({ children }) => {
  // Ambil user & tenant dari localStorage saat inisialisasi
  const [user, setUser] = useState(() => {
    const saved = localStorage.getItem("user");
    return saved ? JSON.parse(saved) : null;
  });
  const [tenant, setTenant] = useState(() => localStorage.getItem("tenant") || "test-tenant");

  const login = async (email, password, tenantId) => {
    const res = await client.post(
      "/auth/login",
      { email, password },
      { headers: { "X-Tenant-ID": tenantId } }
    );

    // Simpan data ke localStorage & state
    localStorage.setItem("token", res.data.token);
    localStorage.setItem("user", JSON.stringify(res.data.user));
    localStorage.setItem("tenant", tenantId);

    setUser(res.data.user);
    setTenant(tenantId);

    return res.data.user;
  };

  // Fungsi logout
  const logout = async () => {
    try {
      await client.post("/logout", {}, { headers: { "X-Tenant-ID": tenant } });
    } catch {}
    localStorage.removeItem("token");
    localStorage.removeItem("user");
    localStorage.removeItem("tenant");
    setUser(null);
  };

  return (
    <AuthContext.Provider value={{ user, login, logout, tenant, setTenant }}>
      {children}
    </AuthContext.Provider>
  );
};

// Hook helper
export const useAuth = () => useContext(AuthContext);
