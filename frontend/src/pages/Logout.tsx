import React, { useEffect } from "react";
import { useAuth } from "../context/AuthContext"; // Adjust path if necessary
import { useNavigate } from "react-router-dom";

const Logout: React.FC = () => {
  const { logout } = useAuth();
  const navigate = useNavigate();

  useEffect(() => {
    // Clear all auth-related data from localStorage
    localStorage.removeItem("access_token");
    localStorage.removeItem("user"); // if stored
    localStorage.clear(); // optional full clear

    // Call context logout to reset state
    logout();

    // Redirect to login
    navigate("/login");
  }, [logout, navigate]);

  return null; // No UI needed
};

export default Logout;
