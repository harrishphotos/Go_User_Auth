import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import Login from "./pages/Login";
import Home from "./pages/Home";
import Register from "./pages/Register";
import VerifyEmailPage from "./pages/VerifyEmailPage";
import { AuthProvider } from "./context/AuthContext";
import ForgotPassword from "./pages/ForgotPassWord";

const App = () => {
  return (
    <AuthProvider>
      {/* Tailwind CSS v4 Test - Remove this once you verify it's working */}
      <div className="bg-blue-500 text-white p-4 text-center">
        🎉 Tailwind CSS v4 is working! You can remove this test div.
      </div>

      <Routes>
        <Route path="/login" element={<Login />} />
        <Route path="/home" element={<Home />} />
        <Route path="/register" element={<Register />} />
        <Route path="/verify-email" element={<VerifyEmailPage />} />
        <Route path="/forgot-password" element={<ForgotPassword />} />
      </Routes>
    </AuthProvider>
  );
};

export default App;
