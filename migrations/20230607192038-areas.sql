-- +migrate Up
CREATE TABLE areas (
	id SERIAL PRIMARY KEY,
	prefecture TEXT NOT NULL,
	name TEXT NOT NULL,
	description TEXT
);
CREATE UNIQUE INDEX idx_areas_prefecture_name ON areas(prefecture ASC, name ASC);
-- +migrate Down
DROP INDEX idx_areas_prefecture_name;
DROP TABLE areas;
