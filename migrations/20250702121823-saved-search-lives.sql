-- +migrate Up
CREATE TABLE saved_search_live (
	lives_id BIGINT,
	users_id BIGINT,
	FOREIGN KEY (lives_id) REFERENCES lives(id) ON DELETE CASCADE,
	FOREIGN KEY (users_id) REFERENCES users(id) ON DELETE CASCADE
);
CREATE UNIQUE INDEX idx_saved_search_live_ids ON saved_search_live(lives_id, users_id);

-- +migrate Down
DROP INDEX idx_saved_search_live_ids;
DROP TABLE saved_search_live;
