-- +migrate Up
ALTER TABLE lives ADD price_en TEXT;
-- +migrate Down
ALTER TABLE lives DROP COLUMN price_en;