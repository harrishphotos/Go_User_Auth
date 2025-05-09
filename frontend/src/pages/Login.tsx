import React, { useState } from "react";
import { useNavigate } from "react-router-dom";
import { useAuth } from "../context/AuthContext";
import useAxiosWrapper from "../api/axiosWrapper";
<<<<<<< Updated upstream

const Login = () => {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
=======
import type { FormEvent } from "react";

type User = {
  id: string;
  email: string;
};

type LoginProps = {};

const Login: React.FC<LoginProps> = () => {
  const [email, setEmail] = useState<string>("");
  const [password, setPassword] = useState<string>("");
  const [error, setError] = useState<string>("");
>>>>>>> Stashed changes
  const { setAuth } = useAuth();
  const axios = useAxiosWrapper();
  const navigate = useNavigate();

<<<<<<< Updated upstream
  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      // Call backend API to authenticate user
      const response = await axios.post("/auth/login", {
=======
  const handleSubmit = async (e: FormEvent<HTMLFormElement>): Promise<void> => {
    e.preventDefault();
    try {
      const response = await axios.post<{
        access_token: string;
        user: User;
      }>("/auth/login", {
>>>>>>> Stashed changes
        email,
        password,
      });

      const { access_token, user } = response.data;

<<<<<<< Updated upstream
      //  Store access token in localStorage
      localStorage.setItem("access_token", access_token);

      //  Set user info and token in context
      setAuth(user, access_token); // assuming setAuth(user, token)

      //  Navigate to home
      navigate("/home");
    } catch (err: any) {
=======
      localStorage.setItem("access_token", access_token);
      setAuth(user, access_token);

      navigate("/home");
    } catch (err) {
>>>>>>> Stashed changes
      setError("Login failed. Please check your credentials.");
      console.error(err);
    }
  };

  return (
    <div>
      <h2>Login</h2>
      {error && <p>{error}</p>}
      <form onSubmit={handleSubmit}>
        <div>
          <label htmlFor="email">Email</label>
          <input
            type="email"
            id="email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            required
          />
        </div>
<<<<<<< Updated upstream
        <div className="mb-4">
=======
        <div>
>>>>>>> Stashed changes
          <label htmlFor="password">Password</label>
          <input
            type="password"
            id="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            required
          />
        </div>
        <button type="submit">Login</button>
      </form>
    </div>
  );
};

export default Login;
