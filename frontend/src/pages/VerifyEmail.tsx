import React from "react";
import { useNavigate, useSearchParams } from "react-router-dom";

const VerifyEmail: React.FC = () => {
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();
  const status = searchParams.get("status") || "failed"; // Default to failed if no status

  const isSuccess = status === "success";

  return (
    <div>
      <h1>{isSuccess ? "Email Verified!" : "Verification Failed"}</h1>
      <p>
        {isSuccess
          ? "Your email has been successfully verified. You can now log in to your account."
          : "We couldn't verify your email address. The verification link may be invalid or expired."}
      </p>
      <button onClick={() => navigate("/login")}>Go to Login</button>
    </div>
  );
};

export default VerifyEmail;
