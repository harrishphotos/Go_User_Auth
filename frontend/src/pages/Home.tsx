import React from "react";
import { useAuth } from "../context/AuthContext"; // Adjust path if needed
import { useNavigate } from "react-router-dom";

export default function Home() {
  const { logout } = useAuth();
  const navigate = useNavigate();

  const handleLogout = () => {
    logout(); // Clear context and localStorage
    navigate("/login"); // Redirect to login page
  };

  return (
    <div className="p-4">
      <h1 className="text-xl font-bold mb-4">Home</h1>
      <button
        onClick={handleLogout}
        className="px-4 py-2 bg-red-500 text-white rounded hover:bg-red-600"
      >
        Logout
      </button>
    </div>
  );
}
