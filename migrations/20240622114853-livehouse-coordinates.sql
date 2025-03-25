-- +migrate Up
ALTER TABLE livehouses ADD latitude REAL;
ALTER TABLE livehouses ADD longitude REAL;
-- +migrate Down
ALTER TABLE livehouses DROP COLUMN latitude;
ALTER TABLE livehouses DROP COLUMN longitude;