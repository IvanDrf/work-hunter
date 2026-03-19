BEGIN;

CREATE TABLE IF NOT EXISTS tokens (
    email TEXT PRIMARY KEY,
    token TEXT NOT NULL,
    exp DATE NOT NULL
);

CREATE INDEX idx_token ON tokens(token);

COMMIT;
