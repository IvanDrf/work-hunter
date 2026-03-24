CREATE TYPE user_status AS ENUM ('active', 'inactive', 'blocked', 'deleted');
CREATE TYPE user_role AS ENUM ('user', 'moderator', 'admin');

CREATE TABLE IF NOT EXISTS users (
    -- main fields
    id UUID PRIMARY KEY,
    username VARCHAR(255) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,

    -- personal info
    first_name VARCHAR(255) NOT NULL DEFAULT '',
    last_name VARCHAR(255) NOT NULL DEFAULT '',
    phone_number VARCHAR(50) NOT NULL DEFAULT '',
    avatar_url TEXT NOT NULL DEFAULT '',

    -- status and role
    status user_status NOT NULL DEFAULT 'active',
    role user_role NOT NULL DEFAULT 'user',

    metadata JSONB NOT NULL DEFAULT '{}',

    -- time points 
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),

    deleted_at TIMESTAMP
);

-- indexes for fast search
CREATE INDEX IF NOT EXISTS idx_users_username ON users(username) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_users_status ON users(status) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_users_role ON users(role) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_users_created_at ON users(created_at) WHERE deleted_at IS NULL;

-- index for fill-text search
CREATE INDEX IF NOT EXISTS idx_users_search ON users
USING GIN(to_tsvector('russian', username || ' ' || email || ' ' || first_name || ' ' || last_name));

-- trigger for auto update 'updated_at'
CREATE OR REPLACE FUNCTION update_updated_at_column() RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER IF NOT EXISTS update_users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();