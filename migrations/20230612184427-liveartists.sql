-- +migrate Up
CREATE TABLE liveartists (
	lives_id BIGINT NOT NULL,
	artists_name TEXT NOT NULL,
	FOREIGN KEY (lives_id) REFERENCES lives(id) ON DELETE CASCADE,
	FOREIGN KEY (artists_name) REFERENCES artists(name) ON DELETE CASCADE
);
CREATE UNIQUE INDEX idx_liveartists_lives_artists ON liveartists(lives_id, artists_name);
CREATE INDEX idx_liveartists_artists ON liveartists(artists_name);
-- +migrate Down
DROP INDEX idx_liveartists_artists;
DROP INDEX idx_liveartists_lives_artists;
DROP TABLE liveartists;