-- name: CreateDoctor :one
INSERT INTO doctors (name, phone, address, created_at, updated_at)
VALUES ($1, $2, $3, NOW(), NOW())
RETURNING *;

-- name: GetDoctorByID :one
SELECT * FROM doctors WHERE id = $1;

-- name: SearchDoctorsByNameOrPhone :many
SELECT * FROM doctors 
WHERE @keyword = '' OR name ILIKE @keyword OR phone ILIKE @keyword 
ORDER BY created_at DESC 
LIMIT @limit_arg OFFSET @offset_arg;

-- name: CountDoctorsByNameOrPhone :one
SELECT COUNT(*) FROM doctors 
WHERE @keyword = '' OR name ILIKE @keyword OR phone ILIKE @keyword;

-- name: IsDoctorNameAndPhoneExists :one
SELECT EXISTS (
    SELECT 1 FROM doctors 
    WHERE name = $1 AND phone = $2 
    LIMIT 1
);

-- name: UpdateDoctorByID :exec
UPDATE doctors
SET name = $2, phone = $3, address = $4, updated_at = NOW()
WHERE id = $1;

-- name: DeleteDoctorByID :exec
DELETE FROM doctors WHERE id = $1;
