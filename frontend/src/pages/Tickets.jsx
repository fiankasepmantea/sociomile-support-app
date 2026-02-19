import React, { useEffect, useState } from "react";
import { TicketsAPI } from "../api/TicketsAPI";

export default function Tickets() {
  const [tickets, setTickets] = useState([]);
  const [modalOpen, setModalOpen] = useState(false);
  const [editTicket, setEditTicket] = useState(null);
  const [form, setForm] = useState({ title: "", description: "" });
  const [toast, setToast] = useState("");

  const fetchTickets = async () => {
    try {
      const res = await TicketsAPI.list();
      setTickets(res.data);
    } catch (err) {
      console.error(err);
    }
  };

  useEffect(() => {
    fetchTickets();
  }, []);

  const openModal = (ticket = null) => {
    setEditTicket(ticket);
    setForm(ticket ? { title: ticket.title, description: ticket.description } : { title: "", description: "" });
    setModalOpen(true);
  };

  const saveTicket = async () => {
    try {
      if (editTicket) {
        // Update status if needed
        await TicketsAPI.updateStatus(editTicket.id, form.status || "open");
        setToast("Ticket updated!");
      } else {
        await TicketsAPI.create(form);
        setToast("Ticket created!");
      }
      fetchTickets();
      setModalOpen(false);
      setTimeout(() => setToast(""), 2000);
    } catch (err) {
      console.error(err);
      setToast("Failed to save ticket");
      setTimeout(() => setToast(""), 2000);
    }
  };

  return (
    <div className="container">
      <h1>Tickets</h1>
      <button onClick={() => openModal()}>New Ticket</button>
      <ul>
        {tickets.map((t) => (
          <li key={t.id}>
            {t.title} - {t.status || "open"}
            <button onClick={() => openModal(t)}>Edit</button>
          </li>
        ))}
      </ul>

      {modalOpen && (
        <div className="modal-overlay">
          <div className="modal">
            <h3>{editTicket ? "Edit Ticket" : "New Ticket"}</h3>
            <input
              type="text"
              placeholder="Title"
              value={form.title}
              onChange={(e) => setForm({ ...form, title: e.target.value })}
            />
            <textarea
              placeholder="Description"
              value={form.description}
              onChange={(e) => setForm({ ...form, description: e.target.value })}
            />
            <button onClick={saveTicket}>{editTicket ? "Update" : "Create"}</button>
            <button onClick={() => setModalOpen(false)}>Cancel</button>
          </div>
        </div>
      )}

      {toast && <div className="toast">{toast}</div>}
    </div>
  );
}
