-- +goose Up
-- +goose StatementBegin
CREATE TABLE features (
    uuid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NOT NULL DEFAULT NULL
);
CREATE INDEX idx_features_name ON features(name);
CREATE INDEX idx_features_deleted_at ON features(deleted_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS features;
DROP INDEX IF EXISTS idx_features_name;
DROP INDEX IF EXISTS idx_features_deleted_at;
-- +goose StatementEnd
