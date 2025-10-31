import { useState } from "react";
import { FaCarSide } from "react-icons/fa";
import { useNavigate } from "react-router-dom";
import { useAuth } from "../context/AuthContext";

function Login() {
  const navigate = useNavigate();
  const { setIsLoggedIn } = useAuth();

  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");

  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");

    if (!username.trim() || !password.trim()) {
      setError("se debe ingresar usuario y contraseña.");
      return;
    }

    setLoading(true);

    try {
      const res = await fetch("/api/v1/login", {
        body: JSON.stringify({ username, password }),
        headers: { "Content-Type": "application/json" },
        method: "POST",
      });

      const data = await res.json();

      if (!res.ok) {
        setError(data.error || "error al iniciar sesión");
        return;
      }

      localStorage.setItem("token", data.token);
      localStorage.setItem("role", data.role);

      setIsLoggedIn(true);
      navigate("/dashboard");
    } catch (err) {
      console.error(err);
      setError("error de conexión con el servidor");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen flex flex-col items-center justify-center bg-linear-to-br from-gray-200 via-gray-300 to-gray-400 text-gray-800">
      <h1 className="text-3xl font-bold mb-8 text-gray-800">
        Gestión de Estacionamiento
      </h1>

      <div className="w-full max-w-md bg-white shadow-lg rounded-2xl p-8">
        <div className="flex justify-center text-gray-700">
          <FaCarSide className="w-16 h-16" />
        </div>

        <h2 className="text-2xl font-semibold text-center mb-6 text-gray-700">
          Inicio de Sesión
        </h2>

        {error && (
          <div
            className="bg-red-100 border border-red-400 text-red-700 px-4 py-2 shadow mb-4"
            role="alert"
          >
            <strong className="font-semibold">Error: </strong>
            <span className="block sm:inline">{error}</span>
          </div>
        )}

        <form onSubmit={handleSubmit} className="space-y-5">
          <div>
            <label
              htmlFor="username"
              className="block text-sm font-medium text-gray-600 mb-1"
            >
              Usuario
            </label>
            <input
              disabled={loading}
              id="username"
              type="text"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              className="w-full border border-gray-300 rounded-lg px-3 py-2 focus:outline-none focus:ring-1 focus:ring-gray-400 hover:border-gray-400 hover:shadow-sm transition-all"
              placeholder="Ingresa tu usuario"
            />
          </div>

          <div>
            <label
              htmlFor="password"
              className="block text-sm font-medium text-gray-600 mb-1"
            >
              Contraseña
            </label>
            <input
              disabled={loading}
              id="password"
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              className="w-full border border-gray-300 rounded-lg px-3 py-2 focus:outline-none focus:ring-1 focus:ring-gray-400 hover:border-gray-400 hover:shadow-sm transition-all"
              placeholder="••••••••"
            />
          </div>

          <button
            type="submit"
            disabled={loading}
            className="w-full bg-gray-800 text-white py-2 rounded-lg hover:bg-gray-700 transition-colors cursor-pointer disabled:opacity-50"
          >
            {loading ? "Iniciando..." : "Iniciar sesión"}
          </button>
        </form>
      </div>
    </div>
  );
}

export default Login;
