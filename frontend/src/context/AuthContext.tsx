// import React, {
//   createContext,
//   useContext,
//   useState,
//   type ReactNode,
// } from "react";
// import { useNavigate } from "react-router-dom";
// import {
//   getAccessToken,
//   setAccessToken,
//   removeAccessToken,
// } from "../utils/tokenUtil"; // Import token utilities

// type User = {
//   id: string;
//   email: string;
// };

// type AuthContextType = {
//   user: User | null;
//   accessToken: string | null;
//   setAuth: (user: User, token: string) => void;
//   logout: () => void;
// };

// const AuthContext = createContext<AuthContextType | undefined>(undefined);

// export const AuthProvider = ({ children }: { children: ReactNode }) => {
//   const [user, setUser] = useState<User | null>(null);
//   const [accessToken, setAccessTokenState] = useState<string | null>(
//     getAccessToken()
//   );

//   const setAuth = (user: User, token: string) => {
//     setUser(user);
//     setAccessTokenState(token);
//     setAccessToken(token); // Store the new token in localStorage
//   };

//   const logout = () => {
//     setUser(null);
//     setAccessTokenState(null);
//     removeAccessToken(); // Remove token from localStorage
//   };

//   return (
//     <AuthContext.Provider value={{ user, accessToken, setAuth, logout }}>
//       {children}
//     </AuthContext.Provider>
//   );
// };

// // Hook for accessing auth context
// export const useAuth = () => {
//   const context = useContext(AuthContext);
//   if (!context) {
//     throw new Error("useAuth must be used within an AuthProvider");
//   }
//   return context;
// };
// AuthContext.tsx
// AuthContext.tsx
import React, { createContext, useContext, useState, useEffect } from "react";

type User = {
  id: string;
  email: string;
};

type AuthContextType = {
  user: User | null;
  accessToken: string | null;
  refreshToken: string | null;
  setAuth: (user: User, accessToken: string, refreshToken: string) => void;
  logout: () => void;
};

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const AuthProvider: React.FC<{ children: React.ReactNode }> = ({
  children,
}) => {
  const [user, setUser] = useState<User | null>(null);
  const [accessToken, setAccessToken] = useState<string | null>(null);
  const [refreshToken, setRefreshToken] = useState<string | null>(null);

  // ✅ Restore from localStorage on mount
  useEffect(() => {
    const storedAccessToken = localStorage.getItem("access_token");
    const storedRefreshToken = localStorage.getItem("refresh_token");
    
    if (storedAccessToken) {
      setAccessToken(storedAccessToken);
    }
    if (storedRefreshToken) {
      setRefreshToken(storedRefreshToken);
    }
  }, []);

  const setAuth = (user: User, accessToken: string, refreshToken: string) => {
    setUser(user);
    setAccessToken(accessToken);
    localStorage.setItem('refresh_token', refreshToken);
  };

  const logout = () => {
    setUser(null);
    setAccessToken(null);
    setRefreshToken(null);
    localStorage.removeItem('refresh_token');
    localStorage.removeItem("access_token");
  };

  return (
    <AuthContext.Provider value={{ user, accessToken, refreshToken, setAuth, logout }}>
      {children}
    </AuthContext.Provider>
  );
};

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error("useAuth must be used within an AuthProvider");
  }
  return context;
};
