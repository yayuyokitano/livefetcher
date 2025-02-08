-- +migrate Up
ALTER TABLE users ADD COLUMN calendar_type SMALLINT;
ALTER TABLE users ADD COLUMN calendar_token TEXT;
ALTER TABLE users ADD COLUMN calendar_id TEXT;

CREATE TABLE calendarevents (
	lives_id BIGINT NOT NULL,
	users_id BIGINT NOT NULL,
	open_id TEXT NOT NULL,
	start_id TEXT NOT NULL,
	FOREIGN KEY (lives_id) REFERENCES lives(id) ON DELETE CASCADE,
	FOREIGN KEY (users_id) REFERENCES users(id) ON DELETE CASCADE
);
CREATE UNIQUE INDEX idx_calendarevents_users_lives ON calendarevents(users_id, lives_id);
CREATE INDEX idx_calendarevents_lives ON calendarevents(lives_id);

-- +migrate Down
ALTER TABLE users DROP COLUMN calendar_id;
ALTER TABLE users DROP COLUMN calendar_token;
ALTER TABLE users DROP COLUMN calendar_type;
DROP INDEX idx_calendarevents_users_lives;
DROP INDEX idx_calendarevents_lives;