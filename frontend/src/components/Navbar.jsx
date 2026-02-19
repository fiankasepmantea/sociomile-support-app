import React from "react";
import { Link, useNavigate } from "react-router-dom";
import { useAuth } from "../context/AuthContext";
import "../assets/styles/navbar.css";

export default function Navbar() {
  const { user, logout, tenant } = useAuth();
  const navigate = useNavigate();

  const handleLogout = async () => {
    await logout();
    navigate("/login");
  };

  return (
    <nav className="navbar">
      <div className="navbar-brand">Sociomile Support</div>
      <div className="navbar-links">
        <Link to="/dashboard">Dashboard</Link>
        <Link to="/conversations">Conversations</Link>
        <Link to="/customers">Customers</Link>
        <Link to="/tickets">Tickets</Link>
        {user && (
          <>
            <span>Tenant: {tenant}</span>
            <button onClick={handleLogout}>Logout</button>
          </>
        )}
      </div>
    </nav>
  );
}
