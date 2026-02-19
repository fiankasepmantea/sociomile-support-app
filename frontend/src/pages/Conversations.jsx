import React, { useEffect, useState } from "react";
import { ConversationsAPI } from "../api/ConversationsAPI";
import "../assets/styles/conversations.css";

export default function Conversations() {
  const [conversations, setConversations] = useState([]);
  const [selectedConv, setSelectedConv] = useState(null);
  const [modalOpen, setModalOpen] = useState(false);
  const [newMessage, setNewMessage] = useState("");
  const [toast, setToast] = useState("");

  const fetchConversations = async () => {
    try {
      const res = await ConversationsAPI.list();
      setConversations(res.data);
    } catch (err) {
      console.error(err);
    }
  };

  useEffect(() => {
    fetchConversations();
  }, []);

  const openModal = (conv) => {
    setSelectedConv(conv);
    setNewMessage("");
    setModalOpen(true);
  };

  const sendMessage = async () => {
    try {
      await ConversationsAPI.sendMessage(selectedConv.id, newMessage);
      setToast("Message sent!");
      setModalOpen(false);
      fetchConversations();
      setTimeout(() => setToast(""), 2000);
    } catch (err) {
      console.error(err);
      setToast("Failed to send message");
      setTimeout(() => setToast(""), 2000);
    }
  };

  return (
    <div className="conversations-container">
      <h1>Conversations</h1>
      <ul className="conversation-list">
        {conversations.map((conv) => (
          <li key={conv.id} className="conversation-item">
            <span>{conv.customer_name || conv.customer}</span>
            <span>{conv.lastMessage}</span>
            <button onClick={() => openModal(conv)}>Reply</button>
          </li>
        ))}
      </ul>

      {modalOpen && (
        <div className="modal-overlay">
          <div className="modal">
            <h3>Send Message to {selectedConv.customer_name || selectedConv.customer}</h3>
            <textarea
              value={newMessage}
              onChange={(e) => setNewMessage(e.target.value)}
              rows={4}
            />
            <button onClick={sendMessage}>Send</button>
            <button onClick={() => setModalOpen(false)}>Cancel</button>
          </div>
        </div>
      )}

      {toast && <div className="toast">{toast}</div>}
    </div>
  );
}
