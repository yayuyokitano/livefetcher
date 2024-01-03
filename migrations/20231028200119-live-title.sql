-- +migrate Up
ALTER TABLE lives ADD title TEXT;
-- +migrate Down
ALTER TABLE lives DROP COLUMN title;