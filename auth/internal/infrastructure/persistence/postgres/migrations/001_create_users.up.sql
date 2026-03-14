BEGIN;

CREATE TABLE IF NOT EXISTS users (
    user_id UUID PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    hashed_password TEXT NOT NULL,
    verificated BOOL DEFAULT FALSE
);

CREATE INDEX IF NOT EXISTS indx_email ON users(email);

COMMIT;