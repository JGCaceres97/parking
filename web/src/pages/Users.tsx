import { useCallback, useEffect, useState } from "react";
import {
  FaCalendarAlt,
  FaCheckCircle,
  FaShieldAlt,
  FaTimesCircle,
  FaUserEdit,
} from "react-icons/fa";
import { FaPlus, FaTrash, FaUser } from "react-icons/fa6";
import Toolbar from "../components/Toolbar";

type Role = "admin" | "common";

interface User {
  id: string;
  username: string;
  role: "admin" | "common";
  is_active: boolean;
  created_at: string;
}

function User() {
  const token = localStorage.getItem("token");
  const userRole = localStorage.getItem("role") || "";

  const [users, setUsers] = useState<User[]>([]);

  const [showCreate, setShowCreate] = useState(false);
  const [showUpdate, setShowUpdate] = useState<null | User>(null);

  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [role, setRole] = useState<Role>("common");
  const [isActive, setIsActive] = useState(true);

  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  const fetchUsers = useCallback(async () => {
    try {
      const res = await fetch("/api/v1/admin/users", {
        headers: { Authorization: `Bearer ${token}` },
      });

      const data: User[] = await res.json();
      setUsers(data);
    } catch (err) {
      console.error(err);
      setError("error al cargar usuarios");
    }
  }, [token]);

  useEffect(() => {
    const fetchData = async () => await fetchUsers();
    fetchData();
  }, [fetchUsers]);

  const handleOpenCreate = () => {
    setUsername("");
    setPassword("");
    setRole("common");
    setIsActive(true);

    setShowCreate(true);
    setError("");
  };

  const handleCloseCreate = () => {
    setShowCreate(false);
    setError("");
  };

  const handleOpenUpdate = (user: User) => {
    setUsername(user.username);
    setRole(user.role);
    setIsActive(user.is_active);

    setShowUpdate(user);
    setError("");
  };

  const handleCloseUpdate = () => {
    setShowUpdate(null);
    setError("");
  };

  const handleCreateUser = async () => {
    if (!username.trim() || !password.trim()) {
      setError("se debe ingresar usuario y contraseña");
      return;
    }

    setLoading(true);
    setError("");

    try {
      const res = await fetch("/api/v1/admin/users", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({
          username,
          password,
          role,
          is_active: isActive,
        }),
      });

      const data = await res.json();
      if (!res.ok) throw new Error(data.error);

      fetchUsers();
      setShowCreate(false);
    } catch (err) {
      console.error(err);

      if (err instanceof Error) {
        setError(err.message);
      } else {
        setError("error al crear usuario");
      }
    } finally {
      setLoading(false);
    }
  };

  const handleUpdateUser = async () => {
    if (!showUpdate) return;

    if (!username.trim()) {
      setError("se debe ingresar usuario");
      return;
    }

    setLoading(true);
    setError("");

    try {
      const res = await fetch(`/api/v1/admin/users/${showUpdate.id}`, {
        method: "PUT",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({
          username,
          role,
          is_active: isActive,
        }),
      });

      const data = await res.json();
      if (!res.ok) throw new Error(data.error);

      fetchUsers();
      setShowUpdate(null);
    } catch (err) {
      console.error(err);

      if (err instanceof Error) {
        setError(err.message);
      } else {
        setError("error al actualizar usuario");
      }
    } finally {
      setLoading(false);
    }
  };

  const handleDeleteUser = async (id: string) => {
    setLoading(true);
    setError("");

    try {
      const res = await fetch(`/api/v1/admin/users/${id}`, {
        method: "DELETE",
        headers: { Authorization: `Bearer ${token}` },
      });

      const data = await res.json();
      if (!res.ok) throw new Error(data.error);

      fetchUsers();
    } catch (err) {
      console.error(err);

      if (err instanceof Error) {
        setError(err.message);
      } else {
        setError("error al eliminar/desactivar usuario");
      }
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen flex flex-col items-center bg-linear-to-br from-gray-200 via-gray-300 to-gray-400 text-gray-800">
      <Toolbar title="Gestión de Estacionamiento" role={userRole} />

      <div className="w-full max-w-5xl p-6 space-y-6">
        {error && !showCreate && !showUpdate && (
          <div
            className="bg-red-100 border border-red-400 text-red-700 px-4 py-3 shadow"
            role="alert"
          >
            <strong className="font-semibold">Error: </strong>
            <span className="block sm:inline">{error}</span>
          </div>
        )}

        <div className="flex justify-end mb-4">
          <button
            onClick={handleOpenCreate}
            className="flex items-center gap-2 bg-green-700 text-white px-4 py-2 rounded-md hover:bg-green-800 transition-colors hover:cursor-pointer"
          >
            <FaPlus /> Crear Usuario
          </button>
        </div>

        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
          {users.map((user) => (
            <div
              key={user.id}
              className="bg-white rounded-xl shadow-md p-5 flex flex-col justify-between hover:shadow-lg transition-shadow border border-gray-400"
            >
              <div className="mb-4">
                <p className="text-xl font-bold text-gray-800 flex items-center gap-2">
                  <FaUser className="text-gray-500" /> {user.username}
                </p>
                <p className="text-gray-600 font-medium flex items-center gap-2 capitalize">
                  <FaShieldAlt className="text-gray-500" /> {user.role}
                </p>
              </div>

              <div className="space-y-2 text-gray-700 text-sm">
                <div className="flex items-center justify-between">
                  <span className="flex items-center gap-1 font-semibold">
                    {user.is_active ? (
                      <FaCheckCircle className="text-green-500" />
                    ) : (
                      <FaTimesCircle className="text-red-500" />
                    )}{" "}
                    Estado:
                  </span>
                  <span>{user.is_active ? "Activo" : "Inactivo"}</span>
                </div>

                <div className="flex items-center justify-between">
                  <span className="flex items-center gap-1 font-semibold">
                    <FaCalendarAlt className="text-gray-500" /> Creado:
                  </span>
                  <span>{new Date(user.created_at).toLocaleString()}</span>
                </div>
              </div>

              <div className="flex justify-between mt-4">
                <button
                  onClick={() => handleOpenUpdate(user)}
                  className="flex-1 mr-2 flex items-center justify-center gap-2 bg-blue-600 text-white py-2 rounded-md hover:bg-blue-700 transition-colors hover:cursor-pointer"
                >
                  <FaUserEdit /> Editar
                </button>
                <button
                  onClick={() => handleDeleteUser(user.id)}
                  className="flex-1 ml-2 flex items-center justify-center gap-2 bg-red-600 text-white py-2 rounded-md hover:bg-red-500 transition-colors hover:cursor-pointer"
                >
                  <FaTrash /> Eliminar
                </button>
              </div>
            </div>
          ))}
        </div>

        {showCreate && (
          <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
            <div className="bg-white rounded-lg shadow-lg p-6 w-full max-w-md">
              <h3 className="text-lg font-semibold mb-2">Crear Usuario</h3>
              {error && (
                <div
                  className="bg-red-100 border border-red-400 text-red-700 px-4 py-3 shadow"
                  role="alert"
                >
                  <strong className="font-semibold">Error: </strong>
                  <span className="block sm:inline">{error}</span>
                </div>
              )}

              <input
                disabled={loading}
                type="text"
                placeholder="Usuario"
                value={username}
                onChange={(e) => setUsername(e.target.value)}
                className="w-full my-3 border border-gray-300 hover:border-gray-400 hover:shadow-sm rounded px-3 py-2"
              />
              <input
                disabled={loading}
                type="password"
                placeholder="Contraseña"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                className="w-full mb-3 border border-gray-300 hover:border-gray-400 hover:shadow-sm rounded px-3 py-2"
              />
              <select
                disabled={loading}
                value={role}
                onChange={(e) => setRole(e.target.value as Role)}
                className="w-full mb-3 border border-gray-300 hover:border-gray-400 hover:shadow-sm rounded px-3 py-2 hover:cursor-pointer"
              >
                <option value="common">Común</option>
                <option value="admin">Admin</option>
              </select>
              <div className="flex items-center mb-3">
                <input
                  disabled={loading}
                  type="checkbox"
                  checked={isActive}
                  onChange={(e) => setIsActive(e.target.checked)}
                  className="mr-2"
                />
                <span>Activo</span>
              </div>
              <div className="flex justify-between gap-2 mt-5">
                <button
                  disabled={loading}
                  onClick={handleCloseCreate}
                  className="px-4 py-2 rounded-md bg-gray-200 hover:bg-gray-300 hover:cursor-pointer"
                >
                  Cancelar
                </button>
                <button
                  disabled={loading}
                  onClick={handleCreateUser}
                  className="px-4 py-2 rounded-md bg-green-600 text-white hover:bg-green-700 hover:cursor-pointer"
                >
                  {loading ? "Creando..." : "Crear"}
                </button>
              </div>
            </div>
          </div>
        )}

        {showUpdate && (
          <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
            <div className="bg-white rounded-lg shadow-lg p-6 w-full max-w-md">
              <h3 className="text-lg font-semibold mb-2">Editar Usuario</h3>
              {error && (
                <div
                  className="bg-red-100 border border-red-400 text-red-700 px-4 py-3 shadow"
                  role="alert"
                >
                  <strong className="font-semibold">Error: </strong>
                  <span className="block sm:inline">{error}</span>
                </div>
              )}

              <input
                type="text"
                placeholder="Username"
                value={username}
                onChange={(e) => setUsername(e.target.value)}
                className="w-full my-3 border border-gray-300 hover:border-gray-400 hover:shadow-sm rounded px-3 py-2"
              />
              <select
                value={role}
                onChange={(e) => setRole(e.target.value as Role)}
                className="w-full mb-3 border border-gray-300 hover:border-gray-400 hover:shadow-sm rounded px-3 py-2 hover:cursor-pointer"
              >
                <option value="common">Común</option>
                <option value="admin">Admin</option>
              </select>
              <div className="flex items-center mb-3">
                <input
                  type="checkbox"
                  checked={isActive}
                  onChange={(e) => setIsActive(e.target.checked)}
                  className="mr-2"
                />
                <span>Activo</span>
              </div>
              <div className="flex justify-between gap-2 mt-5">
                <button
                  disabled={loading}
                  onClick={handleCloseUpdate}
                  className="px-4 py-2 rounded-md bg-gray-200 hover:bg-gray-300 hover:cursor-pointer"
                >
                  Cancelar
                </button>
                <button
                  disabled={loading}
                  onClick={handleUpdateUser}
                  className="px-4 py-2 rounded-md bg-blue-600 text-white hover:bg-blue-700 hover:cursor-pointer"
                >
                  {loading ? "Actualizando..." : "Actualizar"}
                </button>
              </div>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}

export default User;
