-- +goose Up
CREATE TABLE PARKING_RECORDS (
  id TEXT PRIMARY KEY NOT NULL, -- ULID
  user_id TEXT NOT NULL,
  vehicle_type_id TEXT NOT NULL,
  license_plate TEXT NOT NULL,
  entry_time DATETIME NOT NULL,
  exit_time DATETIME,
  total_charge REAL,
  calculated_hours INTEGER,

  FOREIGN KEY (user_id) REFERENCES USERS(id) ON DELETE RESTRICT,
  FOREIGN KEY (vehicle_type_id) REFERENCES VEHICLE_TYPES(id) ON DELETE RESTRICT
);

CREATE UNIQUE INDEX idx_parking_records_one_active ON PARKING_RECORDS(license_plate COLLATE NOCASE) WHERE exit_time IS NULL;
CREATE INDEX idx_parking_records_plate_active ON PARKING_RECORDS(license_plate, exit_time);
CREATE INDEX idx_parking_records_exit_time ON PARKING_RECORDS(exit_time);


-- +goose Down
DROP INDEX IF EXISTS idx_parking_records_one_active;
DROP INDEX IF EXISTS idx_parking_records_plate_active;
DROP INDEX IF EXISTS idx_parking_records_exit_time;

DROP TABLE PARKING_RECORDS;
