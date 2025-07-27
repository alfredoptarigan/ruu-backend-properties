-- +goose Up
-- +goose StatementBegin
CREATE TABLE properties (
   uuid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
   owner_client_uuid UUID REFERENCES clients(uuid),
   agent_user_uuid UUID REFERENCES users(uuid),
   title VARCHAR(255) NOT NULL,
   description TEXT,
   status VARCHAR(20) CHECK (status IN ('Available', 'Sold', 'Rented', 'Under Offer')),
   listing_type VARCHAR(20) CHECK (listing_type IN ('For Sale', 'For Rent')),
   property_type VARCHAR(20) CHECK (property_type IN ('House', 'Apartment', 'Land', 'Villa')),
   price DECIMAL(10, 2) NOT NULL,
   address TEXT NOT NULL,
   city VARCHAR(100) NOT NULL,
   size_land INTEGER,
   size_building INTEGER,
   bedrooms INTEGER,
   bathrooms INTEGER,
   year_built INTEGER,
   created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
   updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
   deleted_at TIMESTAMP NOT NULL DEFAULT NULL
);

CREATE INDEX idx_properties_owner_client_uuid ON properties(owner_client_uuid);
CREATE INDEX idx_properties_agent_user_uuid ON properties(agent_user_uuid);
CREATE INDEX idx_properties_status ON properties(status);
CREATE INDEX idx_properties_listing_type ON properties(listing_type);
CREATE INDEX idx_properties_property_type ON properties(property_type);
CREATE INDEX idx_properties_city ON properties(city);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS properties;
DROP INDEX IF EXISTS idx_properties_owner_client_uuid;
DROP INDEX IF EXISTS idx_properties_agent_user_uuid;
DROP INDEX IF EXISTS idx_properties_status;
DROP INDEX IF EXISTS idx_properties_listing_type;
DROP INDEX IF EXISTS idx_properties_property_type;
DROP INDEX IF EXISTS idx_properties_city;
-- +goose StatementEnd
