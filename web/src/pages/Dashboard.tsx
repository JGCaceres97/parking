import { useCallback, useEffect, useState } from "react";
import {
  FaArrowCircleDown,
  FaArrowCircleUp,
  FaClock,
  FaDollarSign,
} from "react-icons/fa";
import { useNavigate } from "react-router-dom";
import Toolbar from "../components/Toolbar";

type Tabs = "current" | "history";

interface VehicleType {
  id: string;
  name: string;
  hourly_rate: number;
  description: string;
}

interface Record {
  vehicle_type_id: string;
  license_plate: string;
  entry_time: string;
  exit_time?: string | null;
  total_charge?: number | null;
  calculated_hours?: number | null;
}

function Dashboard() {
  const navigate = useNavigate();

  const token = localStorage.getItem("token");
  const role = localStorage.getItem("role") || "";

  const [now, setNow] = useState(new Date());
  const [activeTab, setActiveTab] = useState<Tabs>("current");

  const [types, setTypes] = useState<VehicleType[]>([]);
  const [records, setRecords] = useState<Record[]>([]);

  const [licensePlate, setLicensePlate] = useState("");
  const [selectedTypeId, setSelectedTypeId] = useState("");

  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  const isCurrent = activeTab === "current";

  const fetchTypes = useCallback(async () => {
    try {
      const res = await fetch("/api/v1/vehicle-types", {
        headers: { Authorization: `Bearer ${token}` },
      });

      const data: VehicleType[] = await res.json();
      setTypes(data);

      if (data.length > 0) setSelectedTypeId(data[0].id);
    } catch (err) {
      console.error(err);
      setError("error al cargar tipos de vehículo");
    }
  }, [token]);

  const fetchRecords = useCallback(async () => {
    try {
      const res = await fetch(`/api/v1/parking/${activeTab}`, {
        headers: { Authorization: `Bearer ${token}` },
      });

      const data: Record[] = await res.json();
      setRecords(data);
    } catch (err) {
      console.error(err);
      setError("error al cargar información de vehículos");
    }
  }, [token, activeTab]);

  useEffect(() => {
    const interval = setInterval(() => {
      setNow(new Date());
    }, 60 * 1000);

    const fetchData = async () => await fetchTypes();
    fetchData();

    return () => clearInterval(interval);
  }, [fetchTypes]);

  useEffect(() => {
    const fetchData = async () => await fetchRecords();
    fetchData();
  }, [navigate, token, fetchRecords]);

  const handleVehicleEntry = async () => {
    if (!licensePlate.trim() || !selectedTypeId) {
      setError("ingresa una placa válida y selecciona un tipo de vehículo");
      return;
    }

    setLoading(true);

    try {
      const res = await fetch("/api/v1/parking/entry", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({
          license_plate: licensePlate,
          vehicle_type_id: selectedTypeId,
        }),
      });

      const data = await res.json();

      if (!res.ok) throw new Error(data.error);
      setLicensePlate("");

      if (isCurrent) await fetchRecords();
    } catch (err) {
      console.error(err);

      if (err instanceof Error) {
        setError(err.message);
      } else {
        setError("error registrando entrada de vehículo");
      }
    } finally {
      setLoading(false);
    }
  };

  const handleVehicleExit = async (plate: string) => {
    setLoading(true);

    try {
      const res = await fetch("/api/v1/parking/exit", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({ license_plate: plate }),
      });

      const data = await res.json();

      if (!res.ok) throw new Error(data.error);
      await fetchRecords();
    } catch (err) {
      console.error(err);

      if (err instanceof Error) {
        setError(err.message);
      } else {
        setError("error registrando salida de vehículo");
      }
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen flex flex-col items-center bg-linear-to-br from-gray-200 via-gray-300 to-gray-400 text-gray-800">
      <Toolbar title="Gestión de Estacionamiento" role={role} />

      <div className="w-full max-w-5xl p-6 space-y-6">
        {error && (
          <div
            className="bg-red-100 border border-red-400 text-red-700 px-4 py-3 shadow"
            role="alert"
          >
            <strong className="font-semibold">Error: </strong>
            <span className="block sm:inline">{error}</span>
          </div>
        )}

        <section className="bg-white rounded-lg shadow p-4 max-w-md mx-auto border border-gray-400">
          <h2 className="text-lg font-semibold mb-3 text-center">
            Tarifa de Vehículos
          </h2>
          <ul className="space-y-2">
            {types.map((type) => (
              <li
                key={type.id}
                className="flex items-center justify-between border-b border-gray-200 py-2 px-3 hover:bg-gray-50 rounded transition-colors"
              >
                <div>
                  <span className="font-medium">{type.name}</span>
                  <p className="text-sm text-gray-500">{type.description}</p>
                </div>
                <span className="font-semibold">${type.hourly_rate}/hr</span>
              </li>
            ))}
          </ul>
        </section>

        <section className="space-y-4">
          <h2 className="text-xl font-semibold mb-2 text-center">
            Vehículos Parqueados
          </h2>
          <div className="flex justify-center mb-4">
            {["current", "history"].map((tab) => (
              <button
                disabled={loading}
                key={tab}
                onClick={() => setActiveTab(tab as Tabs)}
                className={`px-4 py-2 -mb-px font-medium text-gray-800 border-b-2 transition-colors hover:cursor-pointer ${
                  activeTab === tab
                    ? "bg-gray-800 text-white border-gray-800"
                    : "bg-white text-gray-800 border-gray-300 hover:bg-gray-100"
                }`}
              >
                {tab === "current" ? "Vigentes" : "Historial"}
              </button>
            ))}
          </div>

          <div className="flex flex-col sm:flex-row mb-4 space-y-2 sm:space-y-0 sm:space-x-2">
            <input
              type="text"
              value={licensePlate}
              onChange={(e) => setLicensePlate(e.target.value.toUpperCase())}
              placeholder="Placa"
              className="flex-1 border border-gray-300 rounded-lg px-3 py-2 focus:outline-none focus:ring-1 focus:ring-gray-400 hover:border-gray-400 hover:shadow-sm transition-all bg-white"
            />
            <select
              value={selectedTypeId}
              onChange={(e) => setSelectedTypeId(e.target.value)}
              className="border border-gray-300 rounded-lg px-3 py-2 focus:outline-none focus:ring-1 focus:ring-gray-400 hover:border-gray-400 hover:shadow-sm transition-all bg-white hover:cursor-pointer w-full sm:w-1/4"
            >
              {types.map((type) => (
                <option key={type.id} value={type.id}>
                  {type.name}
                </option>
              ))}
            </select>
            <button
              disabled={loading}
              onClick={handleVehicleEntry}
              className="bg-green-700 text-white px-4 py-2 rounded-md hover:cursor-pointer hover:bg-green-800 transition-colors"
            >
              Registrar Entrada
            </button>
          </div>

          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
            {records.map((record) => {
              const type = types.find(
                ({ id }) => id === record.vehicle_type_id,
              );

              const entryDate = new Date(record.entry_time);
              const diffMs = now.getTime() - entryDate.getTime();
              const diffHours = Math.floor(diffMs / (1000 * 60 * 60));
              const diffMinutes = Math.floor(
                (diffMs % (1000 * 60 * 60)) / (1000 * 60),
              );
              const elapsedTime =
                diffHours > 0
                  ? `${diffHours}h ${diffMinutes}m`
                  : `${diffMinutes}m`;

              return (
                <div
                  key={record.license_plate + record.entry_time}
                  className={`bg-white rounded-xl shadow-md p-5 flex flex-col justify-between hover:shadow-lg transition-shadow border border-gray-400`}
                >
                  <div className="mb-4">
                    <p className="text-xl font-bold text-gray-800">
                      {record.license_plate}
                    </p>
                    <p className="text-gray-600 font-medium">
                      {type ? type.name : "Desconocido"}
                    </p>
                  </div>

                  <div className="space-y-2 text-gray-700 text-sm">
                    <div className="flex items-center justify-between">
                      <span className="flex items-center gap-1 font-semibold">
                        <FaArrowCircleDown className="text-gray-500" /> Entrada:
                      </span>
                      <span>{entryDate.toLocaleString()}</span>
                    </div>

                    {isCurrent && (
                      <div className="flex items-center justify-between">
                        <span className="flex items-center gap-1 font-semibold">
                          <FaClock className="text-gray-500" />
                          Estadía:
                        </span>
                        <span>~{elapsedTime}</span>
                      </div>
                    )}

                    {!isCurrent && (
                      <>
                        <div className="flex items-center justify-between">
                          <span className="flex items-center gap-1 font-semibold">
                            <FaArrowCircleUp className="text-gray-500" />{" "}
                            Salida:
                          </span>
                          <span>
                            {record.exit_time
                              ? new Date(record.exit_time).toLocaleString()
                              : "-"}
                          </span>
                        </div>
                        <div className="flex items-center justify-between">
                          <span className="flex items-center gap-1 font-semibold">
                            <FaClock className="text-gray-500" /> Horas:
                          </span>
                          <span>{record.calculated_hours ?? "-"}</span>
                        </div>
                        <div className="flex items-center justify-between">
                          <span className="flex items-center gap-1 font-semibold">
                            <FaDollarSign className="text-gray-500" /> Cobro
                            total:
                          </span>
                          <span>
                            {record.total_charge != null
                              ? `$${record.total_charge}`
                              : "-"}
                          </span>
                        </div>
                      </>
                    )}
                  </div>

                  {isCurrent && (
                    <button
                      disabled={loading}
                      onClick={() => handleVehicleExit(record.license_plate)}
                      className="mt-4 bg-red-600 text-white py-2 rounded-md hover:bg-red-800 transition-colors font-medium hover:cursor-pointer"
                    >
                      Salida
                    </button>
                  )}
                </div>
              );
            })}
          </div>
        </section>
      </div>
    </div>
  );
}

export default Dashboard;
