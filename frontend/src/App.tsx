<<<<<<< Updated upstream
// import React from "react";
// import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
// import Login from "./pages/Login";
// import Register from "./pages/Register";
// import Home from "./pages/Home";

// const App: React.FC = () => {
//   return (
//     <Router>
//       <Routes>
//         <Route path="/login" element={<Login />} />
//         <Route path="/register" element={<Register />} />
//         <Route path="/home" element={<Home />} />
//       </Routes>
//     </Router>
//   );
// };

// export default App;

import React, { useEffect } from "react";
import { BrowserRouter as Router, Route, Routes } from "react-router-dom";
import { useAuth } from "./context/AuthContext";
import { refreshAuthToken } from "./service/authService";
import Home from "./pages/Home";
import Login from "./pages/Login";

const App: React.FC = () => {
  const { accessToken, logout } = useAuth();

  useEffect(() => {
    if (!accessToken) {
      refreshAuthToken()
        .then((newToken) => {
          // If token refresh is successful, update context
        })
        .catch(() => {
          logout(); // Log out if token refresh fails
        });
    }
  }, [accessToken, logout]);

=======
import React from "react";
import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import Login from "./pages/Login";
import Register from "./pages/Register";
import Home from "./pages/Home";
import VerifyEmail from "./pages/VerifyEmail";

const App: React.FC = () => {
>>>>>>> Stashed changes
  return (
    <Router>
      <Routes>
        <Route path="/login" element={<Login />} />
<<<<<<< Updated upstream
        <Route path="/home" element={<Home />} />
=======
        <Route path="/register" element={<Register />} />
        <Route path="/home" element={<Home />} />
        <Route path="/verify-email" element={<VerifyEmail />} />
>>>>>>> Stashed changes
      </Routes>
    </Router>
  );
};

export default App;
