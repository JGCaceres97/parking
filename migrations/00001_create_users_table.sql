-- +goose Up
CREATE TABLE USERS (
  id VARCHAR(26) PRIMARY KEY NOT NULL, -- ULID
  username VARCHAR(255) UNIQUE NOT NULL,
  password_hash VARCHAR(255) NOT NULL,
  role ENUM('admin', 'common') NOT NULL,
  is_active BOOLEAN DEFAULT TRUE,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- USR:admin - PASS:admin
INSERT INTO USERS (id, username, password_hash, role) VALUES
('01K8M9EX953V6Z68Q96X645BPC', 'admin', '$2a$10$M3wzhyw1mDcspx75s7IwlOXtiH/wVHsAEDSQT8iZd8XEFMKVGv4Ai', 'admin');

-- +goose Down
DROP TABLE USERS;
