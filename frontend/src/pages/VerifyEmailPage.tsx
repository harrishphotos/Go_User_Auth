import React, { useEffect, useState } from "react";
import { useSearchParams, useNavigate } from "react-router-dom";
import axios from "axios";

const VerifyEmailPage: React.FC = () => {
  const [status, setStatus] = useState<"pending" | "success" | "error">(
    "pending"
  );
  const [message, setMessage] = useState("");
  const [searchParams] = useSearchParams();
  const navigate = useNavigate();

  useEffect(() => {
    const token = searchParams.get("token");

    // If no token is provided in the URL, show an error
    if (!token) {
      setStatus("error");
      setMessage("❌ Invalid verification link.");
      return;
    }

    // Send a request to the backend to verify the token
    axios
      .get(`/api/auth/verify-email?token=${token}`)
      .then(() => {
        // On success, show a success message
        setStatus("success");
        setMessage(
          "✅ Your email has been successfully verified! You can now log in."
        );

        // Redirect to login page after 3 seconds
        setTimeout(() => {
          navigate("/login");
        }, 3000);
      })
      .catch(() => {
        // On error, show an error message
        setStatus("error");
        setMessage("❌ Invalid or expired verification link.");
      });
  }, [searchParams, navigate]);

  return (
    <div className="verification-container">
      <h1>Email Verification</h1>
      <p>{message}</p>
      {status === "pending" && <p>Verifying your email...</p>}
    </div>
  );
};

export default VerifyEmailPage;
