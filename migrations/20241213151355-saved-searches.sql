-- +migrate Up
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

-- +migrate Down
DROP INDEX idx_saved_search_areas_saved_search;
DROP INDEX idx_saved_search_areas_area;
DROP TABLE saved_search_areas;

DROP INDEX idx_saved_searches_text_search;
DROP INDEX idx_saved_searches_users;
DROP TABLE saved_searches;
