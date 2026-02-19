import React from "react";
import "../assets/styles/dashboard.css";

export default function Dashboard() {
  return (
    <div className="dashboard-container">
      <h1>Dashboard</h1>
      <div className="dashboard-cards">
        <div className="card">Total Conversations: 12</div>
        <div className="card">Active Users: 5</div>
        <div className="card">Tickets Pending: 3</div>
      </div>
    </div>
  );
}
