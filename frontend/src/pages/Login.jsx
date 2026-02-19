import React, { useState } from "react";
import { useAuth } from "../context/AuthContext";
import "../assets/styles/login.css";
import { useNavigate } from "react-router-dom";
export default function Login() {
  const { login } = useAuth();
  const [email, setEmail] = useState(""); 
  const [password, setPassword] = useState(""); 
  const [tenantId, setTenantId] = useState("test-tenant");
  const [error, setError] = useState("");
  const navigate = useNavigate();
  const handleSubmit = async (e) => {
    e.preventDefault();
    setError("");
    try {
      await login(email, password, tenantId);
      navigate("/dashboard"); 
    } catch (err) {
      console.error(err);
      setError(err.response?.data?.error || "Login failed");
    }
  };

  return (
    <div className="login-container">
      <h1>Login</h1>
      {error && <div style={{ color: "red", marginBottom: "10px" }}>{error}</div>}
      <form onSubmit={handleSubmit}>
        <label>
          Tenant:
          <select value={tenantId} onChange={(e) => setTenantId(e.target.value)}>
            <option value="test-tenant">Test Tenant</option>
            <option value="demo-company">Demo Company</option>
          </select>
        </label>

        <input
          type="email"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          placeholder="admin@gmail.com"
          required
        />

        <input
          type="password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          placeholder="123456"
          required
        />

        <button type="submit">Login</button>
      </form>
    </div>
  );
}
