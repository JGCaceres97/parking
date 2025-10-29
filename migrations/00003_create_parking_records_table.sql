-- +goose Up
CREATE TABLE PARKING_RECORDS (
  id VARCHAR(26) PRIMARY KEY NOT NULL, -- ULID
  user_id VARCHAR(26) NOT NULL,
  vehicle_type_id VARCHAR(26) NOT NULL,
  license_plate VARCHAR(10) NOT NULL,
  entry_time DATETIME NOT NULL,
  exit_time DATETIME NULL,
  total_charge DECIMAL(10, 2) NULL,
  calculated_hours INT NULL,

  FOREIGN KEY (user_id) REFERENCES USERS(id),
  FOREIGN KEY (vehicle_type_id) REFERENCES VEHICLE_TYPES(id)
);

-- +goose Down
DROP TABLE PARKING_RECORDS;
