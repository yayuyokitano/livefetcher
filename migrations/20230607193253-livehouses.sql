-- +migrate Up
CREATE TABLE livehouses (
	id TEXT PRIMARY KEY,
	url TEXT,
	description TEXT,
	areas_id INTEGER,
	FOREIGN KEY (areas_id) REFERENCES areas(id) ON DELETE SET NULL
);
CREATE INDEX idx_livehouses_areas_id ON livehouses(areas_id);
-- +migrate Down
DROP INDEX idx_livehouses_areas_id;
DROP TABLE livehouses;