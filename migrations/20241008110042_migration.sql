-- +migrate Up
CREATE TABLE IF NOT EXISTS posts (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  title TEXT,
  content TEXT,
  created_at TEXT,
  updated_at TEXT
);

-- +migrate Down
DROP TABLE IF EXISTS posts;