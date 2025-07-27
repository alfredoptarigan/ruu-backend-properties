-- +goose Up
-- +goose StatementBegin
CREATE TABLE clients(
    uuid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    phone_number VARCHAR(15) NOT NULL,
    address TEXT NOT NULL,
    contact_person VARCHAR(100) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE DEFAULT NULL
);

CREATE INDEX idx_clients_email ON clients(email);
CREATE INDEX idx_clients_phone_number ON clients(phone_number);
CREATE INDEX idx_clients_created_at ON clients(created_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS clients;
DROP INDEX IF EXISTS idx_clients_phone_number;
DROP INDEX IF EXISTS idx_clients_created_at;
-- +goose StatementEnd
