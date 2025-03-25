-- +migrate Up
CREATE TABLE artistaliases (
	alias TEXT NOT NULL,
	artists_name TEXT NOT NULL,
	FOREIGN KEY (artists_name) REFERENCES artists(name) ON DELETE CASCADE
);
CREATE UNIQUE INDEX idx_artistaliases_aliases_artists ON artistaliases(alias, artists_name);
CREATE INDEX idx_artistaliases_aliases ON artistaliases(alias);
-- +migrate Down
DROP INDEX idx_artistaliases_aliases;
DROP INDEX idx_artistaliases_aliases_artists;
DROP TABLE artistaliases;