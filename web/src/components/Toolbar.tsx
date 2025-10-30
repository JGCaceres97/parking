import { useNavigate } from "react-router-dom";
import { useAuth } from "../context/AuthContext";
import { useState } from "react";
import { FaHouse, FaUsers } from "react-icons/fa6";
import { FaUserCog } from "react-icons/fa";

interface ToolbarProps {
  title: string;
  role: string;
}

function Toolbar({ title, role }: ToolbarProps) {
  const navigate = useNavigate();
  const { setIsLoggedIn } = useAuth();

  const [showModal, setShowModal] = useState(false);
  const [username, setUsername] = useState("");

  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  const handleUpdateUsername = async () => {
    if (!username.trim()) {
      setError("debe ingresar un nuevo nombre de usuario");
      return;
    }

    setLoading(true);
    setError("");

    try {
      const token = localStorage.getItem("token");
      const res = await fetch("/api/v1/users/me", {
        method: "PUT",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({ username }),
      });

      const data = await res.json();
      if (!res.ok) throw new Error(data.error);

      setShowModal(false);
      setUsername("");
    } catch (err) {
      console.error(err);

      if (err instanceof Error) {
        setError(err.message);
      } else {
        setError("error al actualizar el usuario");
      }
    } finally {
      setLoading(false);
    }
  };

  const handleLogout = () => {
    localStorage.removeItem("token");
    localStorage.removeItem("role");

    setIsLoggedIn(false);
    navigate("/");
  };

  const goToDashboard = () => navigate("/dashboard");
  const goToUserManagement = () => navigate("/users");

  return (
    <div className="w-full bg-gray-800 text-white flex items-center justify-between px-6 py-3 shadow-md">
      <div className="flex items-center space-x-4">
        <h1 className="text-xl font-bold">{title}</h1>
        <button
          onClick={goToDashboard}
          className="flex items-center gap-2 bg-gray-700 hover:bg-gray-600 px-3 py-1 rounded transition-colors hover:cursor-pointer"
        >
          <FaHouse /> Dashboard
        </button>

        {role === "admin" && (
          <button
            onClick={goToUserManagement}
            className="flex items-center gap-2 bg-gray-700 hover:bg-gray-600 px-3 py-1 rounded transition-colors hover:cursor-pointer"
          >
            <FaUsers />
            Usuarios
          </button>
        )}

        <button
          onClick={() => setShowModal(true)}
          className="flex items-center gap-2 bg-gray-700 hover:bg-gray-600 px-3 py-1 rounded transition-colors hover:cursor-pointer"
        >
          <FaUserCog />
          Editar usuario
        </button>
      </div>

      <button
        onClick={handleLogout}
        className="bg-red-600 hover:bg-red-800 hover:cursor-pointer px-3 py-1 rounded transition-colors"
      >
        Cerrar sesión
      </button>

      {showModal && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center">
          <div className="bg-white rounded-lg shadow-lg p-6 w-full max-w-md">
            <h2 className="text-lg font-semibold text-gray-700 mb-4">
              Actualizar nombre de usuario
            </h2>

            {error && (
              <div
                className="bg-red-100 border border-red-400 text-red-700 px-4 py-3 shadow mb-4"
                role="alert"
              >
                <strong className="font-semibold">Error: </strong>
                <span className="block sm:inline">{error}</span>
              </div>
            )}

            <input
              type="text"
              placeholder="Nuevo nombre de usuario"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              className="text-gray-700 w-full border border-gray-300 hover:border-gray-400 hover:shadow-sm rounded px-3 py-2 mb-4 focus:outline-none focus:ring-2 focus:ring-gray-400"
            />

            <div className="flex justify-between gap-2 mt-2">
              <button
                disabled={loading}
                onClick={() => setShowModal(false)}
                className="px-4 py-2 text-gray-700 rounded-md bg-gray-200 hover:bg-gray-300 hover:cursor-pointer"
              >
                Cancelar
              </button>
              <button
                disabled={loading}
                onClick={handleUpdateUsername}
                className="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 transition disabled:opacity-50 hover:cursor-pointer"
              >
                {loading ? "Actualizando..." : "Actualizar"}
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}

export default Toolbar;
