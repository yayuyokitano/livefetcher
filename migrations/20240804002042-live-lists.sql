-- +migrate Up
CREATE TABLE livelists (
	id BIGSERIAL PRIMARY KEY,
	users_id BIGINT NOT NULL,
	title TEXT NOT NULL,
	list_description TEXT NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	FOREIGN KEY (users_id) REFERENCES users(id) ON DELETE CASCADE
);
CREATE INDEX idx_livelists_users ON livelists(users_id);

CREATE TABLE livelistlives (
	id BIGSERIAL PRIMARY KEY,
	livelists_id BIGINT NOT NULL,
	lives_id BIGINT NOT NULL,
	live_description TEXT NOT NULL,
	added_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	FOREIGN KEY (livelists_id) REFERENCES livelists(id) ON DELETE CASCADE,
	FOREIGN KEY (lives_id) REFERENCES lives(id) ON DELETE CASCADE
);
CREATE UNIQUE INDEX idx_livelistlives_lives_livelists ON livelistlives(lives_id, livelists_id);
CREATE INDEX idx_livelistlives_livelists ON livelistlives(livelists_id);

CREATE TABLE livelistfavorites (
	id BIGSERIAL PRIMARY KEY,
	users_id BIGINT NOT NULL,
	livelists_id BIGINT NOT NULL,
	favorited_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	FOREIGN KEY (users_id) REFERENCES users(id) ON DELETE CASCADE,
	FOREIGN KEY (livelists_id) REFERENCES livelists(id) ON DELETE CASCADE
);
CREATE UNIQUE INDEX idx_livelistfavorites_livelists_users ON livelistfavorites(livelists_id, users_id);
CREATE INDEX idx_livelistfavorites_users ON livelistfavorites(users_id);
-- +migrate Down
DROP INDEX idx_livelistfavorites_livelists_users;
DROP INDEX idx_livelistfavorites_users;
DROP TABLE livelistfavorites;

DROP INDEX idx_livelistlives_lives_livelists;
DROP INDEX idx_livelistlives_livelists;
DROP TABLE livelistlives;

DROP INDEX idx_livelists_users;
DROP TABLE livelists;