-- +goose Up
-- +goose StatementBegin
CREATE TABLE tests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    price INTEGER NOT NULL DEFAULT 0,
    imported_price INTEGER NOT NULL DEFAULT 0,
    normal_value VARCHAR(255) NOT NULL DEFAULT '',
    unit VARCHAR(50) NOT NULL DEFAULT '',
    lower_bound DOUBLE PRECISION,
    upper_bound DOUBLE PRECISION,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_tests_name ON tests(name);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS tests;
-- +goose StatementEnd
