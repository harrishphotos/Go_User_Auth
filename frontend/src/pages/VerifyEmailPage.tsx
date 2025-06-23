import React, { useEffect, useState } from "react";
import { useSearchParams, useNavigate } from "react-router-dom";
import axios from "axios";
import Layout from "@/components/layouts/Layout";

const VerifyEmailPage: React.FC = () => {
  const [status, setStatus] = useState<"pending" | "success" | "error">(
    "pending"
  );
  const [message, setMessage] = useState("");
  const [searchParams] = useSearchParams();
  const navigate = useNavigate();

  useEffect(() => {
    const token = searchParams.get("token");

    if (!token) {
      setStatus("error");
      setMessage("Invalid verification link.");
      return;
    }

    axios
      .get(`/api/auth/verify-email?token=${token}`)
      .then(() => {
        setStatus("success");
        setMessage(
          "Your email has been successfully verified! You can now log in."
        );
        setTimeout(() => {
          navigate("/login");
        }, 3000);
      })
      .catch(() => {
        setStatus("error");
        setMessage("Invalid or expired verification link.");
      });
  }, [searchParams, navigate]);

  return (
    <Layout
      title="Email Verification"
      subtitle="Just a moment while we confirm your email..."
    >
      <div className="text-white text-center space-y-4">
        <p className="text-lg font-medium">{message}</p>
        {status === "pending" && (
          <p className="text-sm">Verifying your email...</p>
        )}
      </div>
    </Layout>
  );
};

export default VerifyEmailPage;
