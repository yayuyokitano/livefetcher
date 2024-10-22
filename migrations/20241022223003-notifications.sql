-- +migrate Up
CREATE TABLE notifications (
	id BIGSERIAL PRIMARY KEY,
	users_id BIGINT NOT NULL,
	lives_id BIGINT,
	seen BOOLEAN NOT NULL DEFAULT FALSE,
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	FOREIGN KEY (users_id) REFERENCES users(id) ON DELETE CASCADE,
	FOREIGN KEY (lives_id) REFERENCES lives(id) ON DELETE CASCADE
);
CREATE INDEX idx_notifications_users_id ON notifications(users_id);
CREATE INDEX idx_notifications_id ON notifications(id);
CREATE INDEX idx_notifications_seen ON notifications(seen);

CREATE TABLE notification_contents (
	notifications_id BIGINT NOT NULL,
	notification_type SMALLINT NOT NULL,
	old_value TEXT NOT NULL,
	new_value TEXT NOT NULL
);
CREATE INDEX idx_notifications_contents_notifications_id ON notification_contents(notifications_id);

-- +migrate Down
DROP INDEX idx_notifications_contents_notifications_id;
DROP TABLE notification_contents;
DROP INDEX idx_notifications_seen;
DROP INDEX idx_notifications_id;
DROP INDEX idx_notifications_users_id;
DROP TABLE notifications;