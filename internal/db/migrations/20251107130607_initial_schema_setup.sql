-- +goose Up
-- +goose StatementBegin
CREATE TABLE patients (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    yob VARCHAR(25) NOT NULL DEFAULT '', 
    gender VARCHAR(50) NOT NULL DEFAULT '',
    address VARCHAR(255) NOT NULL DEFAULT '',
    phone VARCHAR(50) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS patients;
-- +goose StatementEnd
