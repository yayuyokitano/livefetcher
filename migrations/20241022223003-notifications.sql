-- +migrate Up
CREATE TABLE notifications (
	id BIGSERIAL PRIMARY KEY,
	lives_id BIGINT,
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	FOREIGN KEY (lives_id) REFERENCES lives(id) ON DELETE SET NULL
);
CREATE INDEX idx_notifications_id ON notifications(id);

CREATE TABLE usernotifications (
	notifications_id BIGINT NOT NULL,
	users_id BIGINT NOT NULL,
	seen BOOLEAN NOT NULL DEFAULT FALSE,
	FOREIGN KEY (notifications_id) REFERENCES notifications(id) ON DELETE CASCADE,
	FOREIGN KEY (users_id) REFERENCES users(id) ON DELETE CASCADE
);
CREATE INDEX idx_usernotifications_users_id ON usernotifications(users_id);
CREATE INDEX idx_notifications_seen ON usernotifications(seen);

CREATE TABLE notification_fields (
	notifications_id BIGINT NOT NULL,
	field_type SMALLINT NOT NULL,
	old_value TEXT NOT NULL,
	new_value TEXT NOT NULL,
	FOREIGN KEY (notifications_id) REFERENCES notifications(id) ON DELETE CASCADE
);
CREATE INDEX idx_notification_fields_notifications_id ON notification_fields(notifications_id);

-- +migrate Down
DROP INDEX idx_notification_fields_notifications_id;
DROP TABLE notification_fields;
DROP INDEX idx_usernotifications_users_id;
DROP INDEX idx_notifications_seen;
DROP TABLE usernotifications;
DROP INDEX idx_notifications_id;
DROP TABLE notifications;