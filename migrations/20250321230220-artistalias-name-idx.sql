-- +migrate Up
CREATE INDEX idx_artistaliases_artistsnames ON artistaliases(artists_name);
-- +migrate Down
DROP INDEX idx_artistaliases_artistsnames;