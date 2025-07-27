-- +goose Up
-- +goose StatementBegin
CREATE TABLE property_features (
    property_uuid UUID REFERENCES properties(uuid),
    feature_uuid UUID REFERENCES features(uuid),
    PRIMARY KEY (property_uuid, feature_uuid)
);
CREATE INDEX idx_property_features_property_uuid ON property_features(property_uuid);
CREATE INDEX idx_property_features_feature_uuid ON property_features(feature_uuid);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS property_features;
DROP INDEX IF EXISTS idx_property_features_property_uuid;
DROP INDEX IF EXISTS idx_property_features_feature_uuid;
-- +goose StatementEnd
