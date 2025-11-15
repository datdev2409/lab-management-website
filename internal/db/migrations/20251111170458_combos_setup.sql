-- +goose Up
-- +goose StatementBegin
CREATE TABLE combos (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_combos_name ON combos(name);

-- Join table for many-to-many relationship between combos and tests
CREATE TABLE combo_tests (
    combo_id UUID NOT NULL REFERENCES combos(id) ON DELETE CASCADE,
    test_id UUID NOT NULL REFERENCES tests(id) ON DELETE CASCADE,
    test_order INTEGER NOT NULL DEFAULT 0,
    PRIMARY KEY (combo_id, test_id)
);

CREATE INDEX idx_combo_tests_combo_id ON combo_tests(combo_id);
CREATE INDEX idx_combo_tests_test_id ON combo_tests(test_id);
CREATE INDEX idx_combo_tests_combo_order ON combo_tests(combo_id, test_order);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS combo_tests;
DROP TABLE IF EXISTS combos;
-- +goose StatementEnd
