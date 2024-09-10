CREATE TABLE metadata (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  ready BOOLEAN
);

CREATE TABLE sessions (
  token TEXT PRIMARY KEY,
  data BLOB NOT NULL,
  expiry REAL NOT NULL
);

CREATE INDEX sessions_expiry_idx ON sessions(expiry);

CREATE TABLE users (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  firstname TEXT,
  lastname TEXT,
  email TEXT,
  password TEXT,
  first_login DATETIME,
  last_login DATETIME
);

INSERT INTO metadata (ready) 
VALUES (false);
