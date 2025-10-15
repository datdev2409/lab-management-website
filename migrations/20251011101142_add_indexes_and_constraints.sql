-- +goose Up
-- +goose StatementBegin

-- Enable pg_trgm extension for trigram-based text search
CREATE EXTENSION IF NOT EXISTS pg_trgm;

-- Add unique constraints for duplicate prevention
ALTER TABLE users ADD CONSTRAINT unique_user_username UNIQUE (username);
ALTER TABLE patients ADD CONSTRAINT unique_patient_name_phone UNIQUE (name, phone);
ALTER TABLE doctors ADD CONSTRAINT unique_doctor_name_phone UNIQUE (name, phone);

-- Add performance indexes
-- User indexes
CREATE INDEX idx_users_username ON users (username);
CREATE INDEX idx_users_active ON users (active);

-- Patient indexes for regex/LIKE pattern search using GIN trigram
CREATE INDEX idx_patients_name_gin ON patients USING GIN (LOWER(name) gin_trgm_ops);
CREATE INDEX idx_patients_phone_gin ON patients USING GIN (phone gin_trgm_ops);
CREATE INDEX idx_patients_created_at ON patients (created_at DESC);

-- Doctor indexes for regex/LIKE pattern search using GIN trigram
CREATE INDEX idx_doctors_name_gin ON doctors USING GIN (LOWER(name) gin_trgm_ops);
CREATE INDEX idx_doctors_phone_gin ON doctors USING GIN (phone gin_trgm_ops);
CREATE INDEX idx_doctors_created_at ON doctors (created_at DESC);

-- Test indexes for regex/LIKE pattern search using GIN trigram
CREATE INDEX idx_tests_name_gin ON tests USING GIN (LOWER(name) gin_trgm_ops);
CREATE INDEX idx_tests_created_at ON tests (created_at DESC);

-- Combo indexes for regex/LIKE pattern search using GIN trigram
CREATE INDEX idx_combos_name_gin ON combos USING GIN (LOWER(name) gin_trgm_ops);
CREATE INDEX idx_combos_created_at ON combos (created_at DESC);

-- Tracking indexes for regex/LIKE pattern search using GIN trigram
CREATE INDEX idx_trackings_name_gin ON trackings USING GIN (LOWER(name) gin_trgm_ops);
CREATE INDEX idx_trackings_created_at ON trackings (created_at DESC);

-- Record indexes (already exist in initial schema, but adding for completeness)
-- CREATE INDEX idx_records_patient_id ON records (patient_id); -- Already exists
-- CREATE INDEX idx_records_doctor_id ON records (doctor_id);   -- Already exists
CREATE INDEX idx_records_status ON records (status);
CREATE INDEX idx_records_created_at ON records (created_at DESC);

-- Test results indexes (already exists in initial schema)  
-- CREATE INDEX idx_test_results_record_id ON test_results (record_id); -- Already exists

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Drop indexes (in reverse order)
DROP INDEX IF EXISTS idx_test_results_record_id;
DROP INDEX IF EXISTS idx_records_created_at;
DROP INDEX IF EXISTS idx_records_status;
DROP INDEX IF EXISTS idx_records_doctor_id;
DROP INDEX IF EXISTS idx_records_patient_id;

-- Drop tracking indexes
DROP INDEX IF EXISTS idx_trackings_created_at;
DROP INDEX IF EXISTS idx_trackings_name_gin;

-- Drop combo indexes
DROP INDEX IF EXISTS idx_combos_created_at;
DROP INDEX IF EXISTS idx_combos_name_gin;

-- Drop test indexes
DROP INDEX IF EXISTS idx_tests_created_at;
DROP INDEX IF EXISTS idx_tests_name_gin;

-- Drop doctor indexes
DROP INDEX IF EXISTS idx_doctors_created_at;
DROP INDEX IF EXISTS idx_doctors_phone_gin;
DROP INDEX IF EXISTS idx_doctors_name_gin;

-- Drop patient indexes
DROP INDEX IF EXISTS idx_patients_created_at;
DROP INDEX IF EXISTS idx_patients_address_gin;
DROP INDEX IF EXISTS idx_patients_phone_gin;
DROP INDEX IF EXISTS idx_patients_name_gin;

-- Drop user indexes
DROP INDEX IF EXISTS idx_users_active;
DROP INDEX IF EXISTS idx_users_username;

-- Drop unique constraints
ALTER TABLE doctors DROP CONSTRAINT IF EXISTS unique_doctor_name_phone;
ALTER TABLE patients DROP CONSTRAINT IF EXISTS unique_patient_name_phone;
ALTER TABLE users DROP CONSTRAINT IF EXISTS unique_user_username;

-- Drop pg_trgm extension (only if no other objects depend on it)
-- DROP EXTENSION IF EXISTS pg_trgm;

-- +goose StatementEnd
