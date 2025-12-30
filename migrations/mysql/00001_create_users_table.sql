-- +goose Up
CREATE TABLE USERS (
  id VARCHAR(26) PRIMARY KEY NOT NULL, -- ULID
  username VARCHAR(255) UNIQUE NOT NULL COLLATE utf8mb4_general_ci,
  password_hash VARCHAR(255) NOT NULL,
  role ENUM('admin', 'common') NOT NULL,
  is_active BOOLEAN DEFAULT TRUE,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX idx_users_username ON USERS(username);

-- +goose Down
DROP INDEX IF EXISTS idx_users_username;

DROP TABLE USERS;
