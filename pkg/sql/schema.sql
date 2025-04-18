CREATE TABLE IF NOT EXISTS contacts (
  blob BLOB NOT NULL,
  created_at TEXT DEFAULT CURRENT_TIMESTAMP NOT NULL,
  id INTEGER PRIMARY KEY,
  info BLOB NOT NULL,  -- override to models.Info
  name TEXT NOT NULL
) STRICT
