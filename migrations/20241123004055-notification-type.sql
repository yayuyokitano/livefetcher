-- +migrate Up
ALTER TABLE notifications ADD COLUMN notification_type SMALLINT NOT NULL DEFAULT 1;


-- +migrate Down
ALTER TABLE notifications DROP COLUMN notification_type;