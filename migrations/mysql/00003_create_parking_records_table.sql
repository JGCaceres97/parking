-- +goose Up
CREATE TABLE PARKING_RECORDS (
  id VARCHAR(26) PRIMARY KEY NOT NULL, -- ULID
  user_id VARCHAR(26) NOT NULL,
  vehicle_type_id VARCHAR(26) NOT NULL,
  license_plate VARCHAR(10) NOT NULL COLLATE utf8mb4_general_ci,
  entry_time DATETIME NOT NULL,
  exit_time DATETIME NULL,
  total_charge DECIMAL(10, 2) NULL,
  calculated_hours INT NULL,

  FOREIGN KEY (user_id) REFERENCES USERS(id),
  FOREIGN KEY (vehicle_type_id) REFERENCES VEHICLE_TYPES(id)
);

CREATE UNIQUE INDEX idx_parking_records_one_active ON PARKING_RECORDS(license_plate, (CASE WHEN exit_time IS NULL THEN 1 ELSE NULL END));
CREATE INDEX idx_parking_records_plate_active ON PARKING_RECORDS(license_plate, exit_time);
CREATE INDEX idx_parking_records_exit_time ON PARKING_RECORDS(exit_time);

-- +goose Down
DROP INDEX IF EXISTS idx_parking_records_one_active;
DROP INDEX IF EXISTS idx_parking_records_plate_active;
DROP INDEX IF EXISTS idx_parking_records_exit_time;

DROP TABLE PARKING_RECORDS;
