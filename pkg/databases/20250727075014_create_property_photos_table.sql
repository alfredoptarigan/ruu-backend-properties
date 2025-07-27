-- +goose Up
-- +goose StatementBegin
CREATE TABLE property_photos (
   uuid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
   property_uuid UUID REFERENCES properties(uuid),
   photo_url TEXT NOT NULL,
   is_primary BOOLEAN NOT NULL DEFAULT FALSE,
   created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
   updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
   deleted_at TIMESTAMP NOT NULL DEFAULT NULL
);
CREATE INDEX idx_property_photos_property_uuid ON property_photos(property_uuid);
CREATE INDEX idx_property_photos_photo_url ON property_photos(photo_url);
CREATE INDEX idx_property_photos_deleted_at ON property_photos(deleted_at);
CREATE INDEX idx_property_photos_is_primary ON property_photos(is_primary);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS property_photos;
DROP INDEX IF EXISTS idx_property_photos_property_uuid;
DROP INDEX IF EXISTS idx_property_photos_photo_url;
DROP INDEX IF EXISTS idx_property_photos_deleted_at;
DROP INDEX IF EXISTS idx_property_photos_is_primary;
-- +goose StatementEnd
