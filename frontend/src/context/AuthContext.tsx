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

  // 🛑 Silent refresh temporarily disabled
  useEffect(() => {
    const storedToken = localStorage.getItem("access_token");
    if (storedToken) {
      setAccessToken(storedToken);
    }
    // Refresh logic is paused for now
  }, [location.pathname]);

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
