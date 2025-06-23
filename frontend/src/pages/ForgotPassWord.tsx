import React, { useState, type FormEvent, type ChangeEvent } from "react";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import Layout from "@/components/layouts/Layout";

const ForgotPassword: React.FC = () => {
  const [email, setEmail] = useState<string>("");
  const [error, setError] = useState<string>("");
  const [successMessage, setSuccessMessage] = useState<string>("");

  const handleSubmit = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    try {
      const response = await fetch(
        "http://localhost:3000/api/auth/forgot-password",
        {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({ email }),
        }
      );

      if (response.ok) {
        setSuccessMessage("Password reset link sent to your email.");
        setEmail("");
        setError("");
      } else {
        const data = await response.json();
        setError(data?.error || "Failed to send reset link.");
        setSuccessMessage("");
      }
    } catch (err) {
      setError("Something went wrong. Please try again later.");
      setSuccessMessage("");
    }
  };

  return (
    <Layout title="Forgot Password" subtitle="We'll send you a reset link.">
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
        {error && <p className="text-red-400 text-sm">{error}</p>}
        {successMessage && (
          <p className="text-green-400 text-sm">{successMessage}</p>
        )}
        <Button type="submit" className="w-full">
          Send Reset Link
        </Button>
      </form>
    </Layout>
  );
};

export default ForgotPassword;
