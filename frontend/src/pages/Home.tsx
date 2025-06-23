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
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-slate-100 via-white to-slate-200 px-4">
      <div className="max-w-md w-full bg-white/90 rounded-xl shadow-xl p-8 backdrop-blur">
        <h1 className="text-2xl font-bold mb-6 text-center text-gray-800">
          Home
        </h1>

        <button
          onClick={handleLogout}
          className="w-full bg-red-500 hover:bg-red-600 text-white font-semibold py-2 px-4 rounded transition"
        >
          Logout
        </button>
      </div>
    </div>
  );
}
