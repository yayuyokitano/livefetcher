-- +migrate Up
CREATE TABLE artists (
	name TEXT PRIMARY KEY,
	url TEXT,
	description TEXT,
	socials TEXT
);
-- +migrate Down
DROP TABLE artists;