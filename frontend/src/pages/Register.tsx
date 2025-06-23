import React, { useState, type FormEvent, type ChangeEvent } from "react";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import Layout from "@/components/layouts/Layout";
import { useNavigate } from "react-router-dom";
import useAxiosWrapper from "../api/axiosWrapper";

const Register: React.FC = () => {
  const [email, setEmail] = useState<string>("");
  const [username, setUsername] = useState<string>("");
  const [password, setPassword] = useState<string>("");
  const [error, setError] = useState<string>("");
  const [message, setMessage] = useState<string>("");
  const navigate = useNavigate();
  const axios = useAxiosWrapper();

  const handleSubmit = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setError("");
    setMessage("");
    try {
      await axios.post("/auth/register", {
        email,
        username,
        password,
      });
      setMessage("Registration successful! Please check your email to verify.");
    } catch (err) {
      setError("Registration failed. Try a different email or username.");
      console.error(err);
    }
  };

  return (
    <Layout title="Register" subtitle="Create your new account.">
      {error && <p className="text-red-400 text-sm mb-2">{error}</p>}
      {message && <p className="text-green-400 text-sm mb-2">{message}</p>}

      <form onSubmit={handleSubmit} className="space-y-5">
        <div>
          <label
            htmlFor="email"
            className="block text-sm font-medium text-white mb-1"
          >
            Email
          </label>
          <Input
            type="email"
            id="email"
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
            htmlFor="username"
            className="block text-sm font-medium text-white mb-1"
          >
            Username
          </label>
          <Input
            type="text"
            id="username"
            value={username}
            onChange={(e: ChangeEvent<HTMLInputElement>) =>
              setUsername(e.target.value)
            }
            required
            placeholder="Choose a username"
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
            value={password}
            onChange={(e: ChangeEvent<HTMLInputElement>) =>
              setPassword(e.target.value)
            }
            required
            placeholder="Create a password"
            className="bg-white text-black"
          />
        </div>

        <Button type="submit" className="w-full">
          Register
        </Button>
      </form>

      <p className="text-sm text-white text-center mt-4">
        Already have an account?{" "}
        <a href="/login" className="text-blue-400 hover:underline">
          Login here
        </a>
      </p>
    </Layout>
  );
};

export default Register;
