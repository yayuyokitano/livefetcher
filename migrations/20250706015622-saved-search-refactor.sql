-- +migrate Up

DROP INDEX idx_saved_search_live_ids;
DROP TABLE saved_search_live;

DROP INDEX idx_saved_search_areas_saved_search;
DROP INDEX idx_saved_search_areas_area;
DROP TABLE saved_search_areas;

DROP INDEX idx_saved_searches_text_search;
DROP INDEX idx_saved_searches_users;
DROP TABLE saved_searches;

CREATE TABLE saved_searches (
	users_id BIGINT NOT NULL,
	keyword text NOT NULL,
	allow_all_locations BOOLEAN NOT NULL DEFAULT FALSE,
	FOREIGN KEY (users_id) REFERENCES users(id) ON DELETE CASCADE
);
CREATE UNIQUE INDEX idx_saved_searches_user_keyword ON saved_searches(users_id, keyword);
CREATE INDEX idx_saved_searches_keyword ON saved_searches(keyword);

CREATE TABLE user_saved_search_areas (
	users_id BIGINT NOT NULL,
	areas_id BIGINT NOT NULL,
	FOREIGN KEY (users_id) REFERENCES users(id) ON DELETE CASCADE,
	FOREIGN KEY (areas_id) REFERENCES areas(id) ON DELETE CASCADE
);
CREATE UNIQUE INDEX idx_user_saved_search_areas_user_area ON user_saved_search_areas(users_id, areas_id);
CREATE INDEX idx_user_saved_search_areas_area ON user_saved_search_areas(areas_id);

-- +migrate Down

DROP INDEX idx_user_saved_search_areas_area;
DROP INDEX idx_user_saved_search_areas_user_area;
DROP TABLE user_saved_search_areas;

DROP INDEX idx_saved_searches_keyword;
DROP INDEX idx_saved_searches_user_keyword;
DROP TABLE saved_searches;

CREATE TABLE saved_searches (
	id BIGSERIAL PRIMARY KEY,
	users_id BIGINT NOT NULL,
	text_search TEXT,
	FOREIGN KEY (users_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_saved_searches_users ON saved_searches(users_id);
CREATE INDEX idx_saved_searches_text_search ON saved_searches(text_search);

CREATE TABLE saved_search_areas (
	saved_searches_id BIGINT NOT NULL,
	areas_id BIGINT NOT NULL,
	FOREIGN KEY (saved_searches_id) REFERENCES saved_searches(id) ON DELETE CASCADE,
	FOREIGN KEY (areas_id) REFERENCES areas(id) ON DELETE CASCADE
);
CREATE UNIQUE INDEX idx_saved_search_areas_saved_search ON saved_search_areas(saved_searches_id);
CREATE INDEX idx_saved_search_areas_area ON saved_search_areas(areas_id);

CREATE TABLE saved_search_live (
	lives_id BIGINT,
	users_id BIGINT,
	FOREIGN KEY (lives_id) REFERENCES lives(id) ON DELETE CASCADE,
	FOREIGN KEY (users_id) REFERENCES users(id) ON DELETE CASCADE
);
CREATE UNIQUE INDEX idx_saved_search_live_ids ON saved_search_live(lives_id, users_id);
