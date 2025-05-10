import React from "react";
import { useAuth } from "../context/AuthContext"; // Adjust path if needed
import { useNavigate, Link } from "react-router-dom";

export default function Home() {
  const { logout } = useAuth();
  const navigate = useNavigate();

  const handleLogout = () => {
    logout(); // Clear context and localStorage
    navigate("/login"); // Redirect to login page
  };

  return (
    <div style={{ padding: "20px" }}>
      <h1
        style={{ fontSize: "24px", fontWeight: "bold", marginBottom: "16px" }}
      >
        Home
      </h1>

      <button
        onClick={handleLogout}
        style={{
          padding: "10px 16px",
          backgroundColor: "#e74c3c",
          color: "#fff",
          border: "none",
          borderRadius: "4px",
          cursor: "pointer",
          marginBottom: "10px",
        }}
      >
        Logout
      </button>
    </div>
  );
}
