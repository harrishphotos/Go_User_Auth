import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import Login from "./pages/Login";
import Home from "./pages/Home";
import Register from "./pages/Register";
import VerifyEmailPage from "./pages/VerifyEmailPage";
import { AuthProvider } from "./context/AuthContext";
import ForgotPassword from "./pages/ForgotPassWord";
import ResetPassword from "./pages/ResetPassword";
import Logout from "./pages/Logout";

const App = () => {
  return (
    <AuthProvider>
      <Routes>
        <Route path="/login" element={<Login />} />
        <Route path="/home" element={<Home />} />
        <Route path="/register" element={<Register />} />
        <Route path="/verify-email" element={<VerifyEmailPage />} />
        <Route path="/forgot-password" element={<ForgotPassword />} />
        <Route path="/reset-password" element={<ResetPassword />} />
        <Route path="/logout" element={<Logout />} />
      </Routes>
    </AuthProvider>
  );
};

export default App;
