CREATE TABLE IF NOT EXISTS contacts (
  created_at REAL DEFAULT(JULIANDAY(CURRENT_TIMESTAMP)) NOT NULL,
  id INTEGER PRIMARY KEY,
  blob BLOB NOT NULL,
  name TEXT NOT NULL
)
