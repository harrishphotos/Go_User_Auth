<<<<<<< Updated upstream
import React, { createContext, useContext, useState, ReactNode } from "react";
=======
import React, {
  createContext,
  useContext,
  useState,
  type ReactNode,
} from "react";
>>>>>>> Stashed changes
import { useNavigate } from "react-router-dom";
import {
  getAccessToken,
  setAccessToken,
  removeAccessToken,
} from "../utils/tokenUtil"; // Import token utilities

type User = {
  id: string;
  email: string;
};

type AuthContextType = {
  user: User | null;
  accessToken: string | null;
  setAuth: (user: User, token: string) => void;
  logout: () => void;
};

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const AuthProvider = ({ children }: { children: ReactNode }) => {
  const [user, setUser] = useState<User | null>(null);
  const [accessToken, setAccessTokenState] = useState<string | null>(
    getAccessToken()
  );

  const setAuth = (user: User, token: string) => {
    setUser(user);
    setAccessTokenState(token);
    setAccessToken(token); // Store the new token in localStorage
  };

  const logout = () => {
    setUser(null);
    setAccessTokenState(null);
    removeAccessToken(); // Remove token from localStorage
  };

  return (
    <AuthContext.Provider value={{ user, accessToken, setAuth, logout }}>
      {children}
    </AuthContext.Provider>
  );
};

// Hook for accessing auth context
export const useAuth = () => {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error("useAuth must be used within an AuthProvider");
  }
  return context;
};
