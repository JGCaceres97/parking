-- +goose Up
CREATE TABLE VEHICLE_TYPES (
  id TEXT PRIMARY KEY NOT NULL, -- ULID
  name TEXT NOT NULL,
  hourly_rate REAL NOT NULL,
  description TEXT
);

CREATE UNIQUE INDEX idx_vehicle_types_name ON VEHICLE_TYPES(name COLLATE NOCASE);

INSERT INTO VEHICLE_TYPES (id, name, hourly_rate, description) VALUES
('01K8M9658PVKMJEBR2GD218M87', 'Normal', 15.00, 'Tarifa regular de $15 USD/hr'),
('01K8M9ADC6PZJPN7KFKNV2B3B1', 'Motocicleta', 0.00, 'Exento de pago'),
('01K8M9B2XWF1KHKVF6CZAB6KBH', 'Especial', 5.00, 'Vehículo especial con tarifa reducida de $5 USD/hr');

-- +goose Down
DROP INDEX IF EXISTS idx_vehicle_types_name;

DROP TABLE VEHICLE_TYPES;
