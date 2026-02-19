import React, { useEffect, useState } from "react";
import { CustomersAPI } from "../api/CustomersAPI";

export default function Customers() {
  const [customers, setCustomers] = useState([]);
  const [modalOpen, setModalOpen] = useState(false);
  const [editCustomer, setEditCustomer] = useState(null);
  const [form, setForm] = useState({ name: "", email: "" });
  const [toast, setToast] = useState("");

  const fetchCustomers = async () => {
    try {
      const res = await CustomersAPI.list();
      setCustomers(res.data);
    } catch (err) {
      console.error(err);
    }
  };

  useEffect(() => {
    fetchCustomers();
  }, []);

  const openModal = (customer = null) => {
    setEditCustomer(customer);
    setForm(customer ? { name: customer.name, email: customer.email } : { name: "", email: "" });
    setModalOpen(true);
  };

  const saveCustomer = async () => {
    try {
      if (editCustomer) {
        await CustomersAPI.update(editCustomer.id, form);
        setToast("Customer updated!");
      } else {
        await CustomersAPI.create(form);
        setToast("Customer created!");
      }
      fetchCustomers();
      setModalOpen(false);
      setTimeout(() => setToast(""), 2000);
    } catch (err) {
      console.error(err);
      setToast("Failed to save customer");
      setTimeout(() => setToast(""), 2000);
    }
  };

  const deleteCustomer = async (id) => {
    if (!window.confirm("Delete this customer?")) return;
    try {
      await CustomersAPI.delete(id);
      setToast("Customer deleted!");
      fetchCustomers();
      setTimeout(() => setToast(""), 2000);
    } catch (err) {
      console.error(err);
      setToast("Failed to delete");
      setTimeout(() => setToast(""), 2000);
    }
  };

  return (
    <div className="container">
      <h1>Customers</h1>
      <button onClick={() => openModal()}>New Customer</button>
      <ul>
        {customers.map((c) => (
          <li key={c.id}>
            {c.name} - {c.email}
            <button onClick={() => openModal(c)}>Edit</button>
            <button onClick={() => deleteCustomer(c.id)}>Delete</button>
          </li>
        ))}
      </ul>

      {modalOpen && (
        <div className="modal-overlay">
          <div className="modal">
            <h3>{editCustomer ? "Edit Customer" : "New Customer"}</h3>
            <input
              type="text"
              placeholder="Name"
              value={form.name}
              onChange={(e) => setForm({ ...form, name: e.target.value })}
            />
            <input
              type="email"
              placeholder="Email"
              value={form.email}
              onChange={(e) => setForm({ ...form, email: e.target.value })}
            />
            <button onClick={saveCustomer}>{editCustomer ? "Update" : "Create"}</button>
            <button onClick={() => setModalOpen(false)}>Cancel</button>
          </div>
        </div>
      )}

      {toast && <div className="toast">{toast}</div>}
    </div>
  );
}
