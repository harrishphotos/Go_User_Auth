import React, { createContext, useContext, useState, useEffect } from "react";
import axios from "../api/axios";
import { useLocation } from "react-router-dom";

type User = {
  id: string;
  email: string;
};

type AuthContextType = {
  user: User | null;
  accessToken: string | null;
  setAuth: (user: User, accessToken: string) => void;
  logout: () => void;
};

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const AuthProvider: React.FC<{ children: React.ReactNode }> = ({
  children,
}) => {
  const [user, setUser] = useState<User | null>(null);
  const [accessToken, setAccessToken] = useState<string | null>(null);
  const location = useLocation();

  const setAuth = (user: User, token: string) => {
    setUser(user);
    setAccessToken(token);
    localStorage.setItem("access_token", token);
  };

  const logout = () => {
    setUser(null);
    setAccessToken(null);
    localStorage.removeItem("access_token");
    axios
      .post("/api/auth/logout", {}, { withCredentials: true })
      .catch(() => {});
  };

  // ⏱️ Refresh token periodically since we can't read expiry (PASETO)
  useEffect(() => {
    if (!accessToken) return;

    const interval = setInterval(async () => {
      try {
        const res = await axios.post<{ accessToken: string }>(
          "/api/auth/refresh",
          {},
          { withCredentials: true }
        );
        const newToken = res.data.accessToken;
        setAccessToken(newToken);
        localStorage.setItem("access_token", newToken);
      } catch {
        logout(); // fallback on refresh failure
      }
    }, 54000); // ⏱️ every 45 seconds for 1-min tokens

    return () => clearInterval(interval); // cleanup on unmount/token change
  }, [accessToken]);

  return (
    <AuthContext.Provider value={{ user, accessToken, setAuth, logout }}>
      {children}
    </AuthContext.Provider>
  );
};

export const useAuth = (): AuthContextType => {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error("useAuth must be used within an AuthProvider");
  }
  return context;
};
