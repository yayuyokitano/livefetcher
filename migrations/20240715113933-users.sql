-- +migrate Up
CREATE TABLE users (
	id BIGSERIAL PRIMARY KEY,
	email TEXT NOT NULL,
	username TEXT NOT NULL,
	nickname TEXT,
	password_hash TEXT NOT NULL,
	bio TEXT,
	location TEXT,
	is_verified BOOLEAN DEFAULT FALSE NOT NULL
);
CREATE UNIQUE INDEX idx_users_email ON users(email);
CREATE UNIQUE INDEX idx_users_username ON users(username);
-- +migrate Down
DROP INDEX idx_users_username;
DROP INDEX idx_users_email;
DROP TABLE users;