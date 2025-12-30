-- +goose Up
CREATE TABLE USERS (
  id TEXT PRIMARY KEY NOT NULL, -- ULID
  username TEXT NOT NULL,
  password_hash TEXT NOT NULL,
  role TEXT NOT NULL CHECK(role IN ('admin', 'common')),
  is_active INTEGER NOT NULL DEFAULT 1,
  created_at DATETIME NOT NULL DEFAULT (CURRENT_TIMESTAMP)
);

CREATE UNIQUE INDEX idx_users_username ON USERS(username COLLATE NOCASE);

-- +goose Down
DROP INDEX IF EXISTS idx_users_username;

DROP TABLE USERS;
