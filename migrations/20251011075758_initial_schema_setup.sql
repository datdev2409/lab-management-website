-- +goose Up
-- +goose StatementBegin

CREATE TABLE users (
    id UUID PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL DEFAULT 'user',
    active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_users_username ON users (username);
CREATE INDEX idx_users_active ON users (active);

CREATE TABLE patients (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    yob VARCHAR(25), 
    gender VARCHAR(50),
    address VARCHAR(255),
    phone VARCHAR(50),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE doctors (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    phone VARCHAR(50),
    address VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE tests (
    id UUID PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL,
    price INTEGER NOT NULL,
    normal_value VARCHAR(255),
    unit VARCHAR(50),
    lower_bound NUMERIC(10, 3), -- Exact precision: XXXXXXX.YYY
    upper_bound NUMERIC(10, 3), -- Exact precision: XXXXXXX.YYY
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE combos (
    id UUID PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Junction Table: Combo-Test Relationship
CREATE TABLE combo_tests (
    combo_id UUID NOT NULL REFERENCES combos(id) ON DELETE CASCADE,
    test_id UUID NOT NULL REFERENCES tests(id) ON DELETE RESTRICT,
    PRIMARY KEY (combo_id, test_id)
);

CREATE TABLE trackings (
    id UUID PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Junction Table: Tracking-Test Relationship (Includes Order)
CREATE TABLE tracking_tests (
    tracking_id UUID NOT NULL REFERENCES trackings(id) ON DELETE CASCADE,
    test_id UUID NOT NULL REFERENCES tests(id) ON DELETE RESTRICT,
    test_order INTEGER NOT NULL, -- Field for display sequence
    PRIMARY KEY (tracking_id, test_id)
);

CREATE TABLE records (
    id UUID PRIMARY KEY,
    
    patient_id UUID NOT NULL REFERENCES patients(id) ON DELETE RESTRICT,
    doctor_id UUID REFERENCES doctors(id) ON DELETE SET NULL, 
    
    combo_name VARCHAR(255), 
    doctor_name VARCHAR(255), 
    
    status VARCHAR(50) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_records_patient_id ON records (patient_id);
CREATE INDEX idx_records_doctor_id ON records (doctor_id);

CREATE TABLE test_results (
    id UUID PRIMARY KEY,
    
    record_id UUID NOT NULL REFERENCES records(id) ON DELETE CASCADE, 
    test_id UUID REFERENCES tests(id) ON DELETE SET NULL,            
    
    -- Historical Snapshot Data (Denormalized)
    name VARCHAR(255) NOT NULL,
    price INTEGER NOT NULL,
    normal_value VARCHAR(255),
    unit VARCHAR(50),
    lower_bound NUMERIC(10, 3),
    upper_bound NUMERIC(10, 3),
    
    -- Result Data
    result TEXT NOT NULL,
    result_text TEXT,
    abnormal BOOLEAN NOT NULL DEFAULT FALSE,
    manual_abnormal_override BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE INDEX idx_test_results_record_id ON test_results (record_id);

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
-- Drop tables in reverse order to respect Foreign Key constraints

DROP TABLE IF EXISTS test_results;
DROP TABLE IF EXISTS records;
DROP TABLE IF EXISTS combo_tests;
DROP TABLE IF EXISTS combos;
DROP TABLE IF EXISTS tracking_tests;
DROP TABLE IF EXISTS trackings;
DROP TABLE IF EXISTS tests;
DROP TABLE IF EXISTS doctors;
DROP TABLE IF EXISTS patients;
DROP TABLE IF EXISTS users;

-- +goose StatementEnd