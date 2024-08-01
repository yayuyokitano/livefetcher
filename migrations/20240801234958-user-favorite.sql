-- +migrate Up
CREATE TABLE userfavorites (
	id BIGSERIAL PRIMARY KEY,
	users_id BIGINT NOT NULL,
	lives_id BIGINT NOT NULL,
	favorited_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	FOREIGN KEY (users_id) REFERENCES users(id) ON DELETE CASCADE,
	FOREIGN KEY (lives_id) REFERENCES lives(id) ON DELETE CASCADE
);
CREATE UNIQUE INDEX idx_userfavorites_lives_users ON userfavorites(lives_id, users_id);
CREATE INDEX idx_userfavorites_users ON userfavorites(users_id);
-- +migrate Down
DROP INDEX idx_userfavorites_lives_users;
DROP INDEX idx_userfavorites_users;
DROP TABLE userfavorites;