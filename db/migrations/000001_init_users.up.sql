CREATE TABLE IF NOT EXISTS
    users (
        id SERIAL PRIMARY KEY,
        email VARCHAR(255) NOT NULL,
        name VARCHAR(50) NOT NULL,
        hashed_password BYTEA NOT NULL,
        created_at TIMESTAMP
    );

CREATE INDEX IF NOT EXISTS users_email
	ON users USING HASH (email);

ALTER TABLE users ADD CONSTRAINT users_email_unique UNIQUE (email);
ALTER TABLE users ALTER COLUMN created_at SET DEFAULT current_timestamp;
