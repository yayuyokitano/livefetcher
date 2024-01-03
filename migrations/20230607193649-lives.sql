-- +migrate Up
CREATE TABLE lives (
	id BIGSERIAL PRIMARY KEY,
	opentime TIMESTAMP(0),
	starttime TIMESTAMP(0),
	url TEXT,
	price TEXT,
	livehouses_id TEXT,
	FOREIGN KEY (livehouses_id) REFERENCES livehouses(id) ON DELETE SET NULL
);
CREATE INDEX idx_lives_starttime ON lives(starttime ASC);
CREATE INDEX idx_lives_livehouses_id ON lives(livehouses_id);
-- +migrate Down
DROP INDEX idx_lives_starttime;
DROP TABLE lives;