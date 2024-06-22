-- +migrate Up
ALTER TABLE livehouses ADD location geography(POINT,4326);
CREATE INDEX livehouse_loc_idx ON livehouses USING GIST ( location );
-- +migrate Down
DROP INDEX livehouse_loc_idx;
ALTER TABLE lives DROP COLUMN location;