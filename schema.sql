CREATE TABLE IF NOT EXISTS contacts (
  blob BLOB NOT NULL,
  created_at REAL DEFAULT(JULIANDAY(CURRENT_TIMESTAMP)) NOT NULL,
  id INTEGER PRIMARY KEY,
  name TEXT NOT NULL
)
