CREATE TYPE user_status AS ENUM ('active', 'inactive', 'blocked', 'deleted');
CREATE TYPE user_role AS ENUM ('employee', 'employer', 'admin');

CREATE TABLE IF NOT EXISTS users (
    -- main fields
    id UUID PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,

    -- personal info
    first_name VARCHAR(255),
    last_name VARCHAR(255),
    company_name VARCHAR(255),

    -- status and role
    status user_status NOT NULL DEFAULT 'active',
    role user_role NOT NULL DEFAULT 'employee',

    BOOLEAN verificated NOT NULL DEFAULT false,

    -- time points 
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),

    CONSTRAINT not_null_personal_info CHECK (first_name IS NOT NULL AND last_name IS NOT NULL OR company_name IS NOT NULL)
);

-- indexes for fast search
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email) WHERE status != 'deleted';

-- trigger for auto update 'updated_at'
CREATE OR REPLACE FUNCTION update_updated_at_column() RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER IF NOT EXISTS update_users_updated_at
    AFTER UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();


CREATE OR REPLACE FUNCTION set_created_at_column() RETURNS TRIGGER AS $$
BEGIN
    NEW.created_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER IF NOT EXISTS set_users_created_at
    AFTER INSERT ON users
    FOR EACH ROW
    EXECUTE FUNCTION set_created_at_column();