import React from "react";
import { useAuth } from "../context/AuthContext";
import { useNavigate } from "react-router-dom";
import { useEffect } from "react";

const Home: React.FC = () => {
  const { logout } = useAuth();
  const navigate = useNavigate();

  const { accessToken } = useAuth();
  useEffect(() => {
    if (!accessToken) return; // ✅ prevent API hit after logout

    // ✅ existing API logic here
  }, [accessToken]);

  const handleLogout = (): void => {
    logout();
    localStorage.clear();
    navigate("/login");
  };

  return (
    <div
      style={{
        display: "flex",
        flexDirection: "column",
        alignItems: "center",
        marginTop: "50px",
      }}
    >
      <h1>Home</h1>
      <button onClick={handleLogout}>Logout</button>
    </div>
  );
};

export default Home;
