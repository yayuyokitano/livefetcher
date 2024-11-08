-- +migrate Up
ALTER TABLE notifications ADD COLUMN deleted BOOLEAN DEFAULT FALSE;

-- +migrate Down
ALTER TABLE notifications DROP COLUMN deleted;