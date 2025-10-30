import {
  Navigate,
  Route,
  BrowserRouter as Router,
  Routes,
} from "react-router-dom";
import { useAuth } from "./context/AuthContext";
import Dashboard from "./pages/Dashboard";
import Login from "./pages/Login";
import User from "./pages/Users";

function App() {
  const { isLoggedIn } = useAuth();

  return (
    <Router>
      <Routes>
        <Route
          path="/"
          element={isLoggedIn ? <Navigate to="/dashboard" /> : <Login />}
        />

        <Route
          path="/dashboard"
          element={isLoggedIn ? <Dashboard /> : <Navigate to="/" />}
        />

        <Route
          path="/users"
          element={isLoggedIn ? <User /> : <Navigate to="/" />}
        />

        <Route
          path="*"
          element={
            isLoggedIn ? <Navigate to="/dashboard" /> : <Navigate to="/" />
          }
        />
      </Routes>
    </Router>
  );
}

export default App;
