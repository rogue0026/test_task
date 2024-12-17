CREATE TABLE IF NOT EXISTS unregistered_users (
    id UUID PRIMARY KEY,
    login VARCHAR(40) NOT NULL,
    name VARCHAR(40) NOT NULL,
    email VARCHAR(40) NOT NULL,
    password VARCHAR(128) NOT NULL,
    is_verified BOOLEAN NOT NULL DEFAULT FALSE,
    verification_code VARCHAR(10) NOT NULL DEFAULT '',
    verification_code_expires VARCHAR(50) NOT NULL DEFAULT ''
);