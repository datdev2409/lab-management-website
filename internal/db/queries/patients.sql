-- name: CreatePatient :one
INSERT INTO patients (name, yob, gender, address, phone, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
RETURNING *;

-- name: GetPatientByID :one
SELECT * FROM patients WHERE id = $1;

-- name: ListPatients :many
SELECT * FROM patients ORDER BY created_at DESC;

-- name: SearchPatientsByNameOrPhone :many
SELECT * FROM patients WHERE @keyword = '' OR name ILIKE '%@keyword%' OR phone ILIKE '%@keyword%' ORDER BY created_at DESC LIMIT @limit_arg OFFSET @offset_arg;

-- name: CountPatientsByNameOrPhone :one
SELECT COUNT(*) FROM patients WHERE @keyword = '' OR name ILIKE @keyword OR phone ILIKE @keyword;

-- name: IsPatientNameAndPhoneExists :one
SELECT EXISTS (SELECT 1 FROM patients WHERE name = $1 AND phone = $2 LIMIT 1);

-- name: UpdatePatientByID :exec
UPDATE patients
SET name = $2, yob = $3, address = $4, phone = $5, updated_at = NOW()
WHERE id = $1;

-- name: DeletePatientByID :exec
DELETE FROM patients WHERE id = $1;


