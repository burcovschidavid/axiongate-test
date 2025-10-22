CREATE TABLE IF NOT EXISTS shipment_records (
    id VARCHAR(36) PRIMARY KEY,
    provider VARCHAR(50) NOT NULL,
    generic_payload JSONB NOT NULL,
    transformed_payload JSONB NOT NULL,
    provider_response JSONB NOT NULL,
    success BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_shipment_records_provider ON shipment_records(provider);
CREATE INDEX idx_shipment_records_created_at ON shipment_records(created_at DESC);
CREATE INDEX idx_shipment_records_success ON shipment_records(success);
