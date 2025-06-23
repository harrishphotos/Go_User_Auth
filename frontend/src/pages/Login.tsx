import React, { useState, type FormEvent, type ChangeEvent } from "react";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import Layout from "@/components/layouts/Layout";
import { useNavigate, Link } from "react-router-dom";
import { useAuth } from "../context/AuthContext";
import useAxiosWrapper from "../api/axiosWrapper";

type User = {
  id: string;
  email: string;
};

const Login: React.FC = () => {
  const [email, setEmail] = useState<string>("");
  const [password, setPassword] = useState<string>("");
  const [error, setError] = useState<string>("");
  const { setAuth } = useAuth();
  const axios = useAxiosWrapper();
  const navigate = useNavigate();

  const handleSubmit = async (e: FormEvent<HTMLFormElement>): Promise<void> => {
    e.preventDefault();
    try {
      const response = await axios.post<{
        access_token: string;
        refresh_token: string;
        user: User;
      }>("/auth/login", {
        email,
        password,
      });

      const { access_token, refresh_token, user } = response.data;

      localStorage.setItem("access_token", access_token);
      localStorage.setItem("refresh_token", refresh_token);
      setAuth(user, access_token, refresh_token);

      navigate("/home");
    } catch (err) {
      setError("Login failed. Please check your credentials.");
      console.error(err);
    }
  };

  return (
    <Layout title="Login" subtitle="Welcome back! Please login to continue.">
      <form onSubmit={handleSubmit} className="space-y-5">
        <div>
          <label
            htmlFor="email"
            className="block text-sm font-medium text-white mb-1"
          >
            Email Address
          </label>
          <Input
            type="email"
            id="email"
            name="email"
            value={email}
            onChange={(e: ChangeEvent<HTMLInputElement>) =>
              setEmail(e.target.value)
            }
            required
            placeholder="you@example.com"
            className="bg-white text-black"
          />
        </div>

        <div>
          <label
            htmlFor="password"
            className="block text-sm font-medium text-white mb-1"
          >
            Password
          </label>
          <Input
            type="password"
            id="password"
            name="password"
            value={password}
            onChange={(e: ChangeEvent<HTMLInputElement>) =>
              setPassword(e.target.value)
            }
            required
            placeholder="password"
            className="bg-white text-black"
          />
        </div>

        {error && <p className="text-red-400 text-sm">{error}</p>}

        <Button type="submit" className="w-full">
          Login
        </Button>

        <div className="text-center mt-4">
          <Link
            to="/forgot-password"
            className="text-blue-400 hover:underline text-sm"
          >
            Forgot Password?
          </Link>
        </div>
      </form>
    </Layout>
  );
};

export default Login;
