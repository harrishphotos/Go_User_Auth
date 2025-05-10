import { useEffect, useState } from "react";
import { useSearchParams, useNavigate } from "react-router-dom";
import axios from "../api/axios";
import { toast } from "react-hot-toast";

interface ResetPasswordForm {
  newPassword: string;
  confirmPassword: string;
}

const ResetPassword = () => {
  const [form, setForm] = useState<ResetPasswordForm>({
    newPassword: "",
    confirmPassword: "",
  });
  const [error, setError] = useState<string>("");
  const [loading, setLoading] = useState<boolean>(false);
  const [searchParams] = useSearchParams();
  const navigate = useNavigate();

  const token = searchParams.get("token");

  useEffect(() => {
    if (!token) {
      setError("Reset token is missing in the URL.");
    }
  }, [token]);

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setForm({ ...form, [e.target.name]: e.target.value });
    setError("");
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!token) {
      setError("Reset token is missing. Cannot submit the form.");
      return;
    }

    if (!form.newPassword || !form.confirmPassword) {
      setError("Please fill out both password fields.");
      return;
    }

    if (form.newPassword !== form.confirmPassword) {
      setError("Passwords do not match.");
      return;
    }

    try {
      setLoading(true);
      const res = await axios.post("/api/auth/reset-password", {
        token,
        newPassword: form.newPassword,
      });

      toast.success("✅ Password reset successful! Redirecting to login...");
      setTimeout(() => navigate("/login"), 2500);
    } catch (err: any) {
      if (err.response) {
        const msg = err.response.data?.error || "An error occurred.";
        setError(`❌ Reset failed: ${msg}`);
      } else {
        setError("❌ Server is not responding. Please try again later.");
      }
    } finally {
      setLoading(false);
    }
  };

  return (
    <div style={{ padding: "2rem", display: "flex", justifyContent: "center" }}>
      <form
        onSubmit={handleSubmit}
        style={{ width: "100%", maxWidth: "400px" }}
      >
        <h2
          style={{
            fontSize: "24px",
            fontWeight: "bold",
            marginBottom: "1rem",
            textAlign: "center",
          }}
        >
          Reset Password
        </h2>

        {error && (
          <div style={{ marginBottom: "1rem", fontSize: "14px", color: "red" }}>
            {error}
          </div>
        )}

        <div style={{ marginBottom: "1rem" }}>
          <label
            htmlFor="newPassword"
            style={{ display: "block", marginBottom: "0.5rem" }}
          >
            New Password
          </label>
          <input
            type="password"
            name="newPassword"
            id="newPassword"
            value={form.newPassword}
            onChange={handleChange}
            placeholder="Enter new password"
            required
            style={{ width: "100%", padding: "0.5rem" }}
          />
        </div>

        <div style={{ marginBottom: "1.5rem" }}>
          <label
            htmlFor="confirmPassword"
            style={{ display: "block", marginBottom: "0.5rem" }}
          >
            Confirm Password
          </label>
          <input
            type="password"
            name="confirmPassword"
            id="confirmPassword"
            value={form.confirmPassword}
            onChange={handleChange}
            placeholder="Confirm new password"
            required
            style={{ width: "100%", padding: "0.5rem" }}
          />
        </div>

        <button
          type="submit"
          disabled={loading}
          style={{
            width: "100%",
            padding: "0.75rem",
            fontWeight: "bold",
            cursor: "pointer",
            opacity: loading ? 0.6 : 1,
          }}
        >
          {loading ? "Resetting..." : "Reset Password"}
        </button>
      </form>
    </div>
  );
};

export default ResetPassword;
