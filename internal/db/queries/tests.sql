-- name: CreateTest :one
INSERT INTO tests (name, price, imported_price, normal_value, unit, lower_bound, upper_bound, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())
RETURNING *;

-- name: GetTestByID :one
SELECT * FROM tests WHERE id = $1;

-- name: SearchTestsByName :many
SELECT * FROM tests WHERE @keyword = '' OR name ILIKE '%' || @keyword || '%' ORDER BY name ASC LIMIT @limit_arg OFFSET @offset_arg;

-- name: CountTestsByName :one
SELECT COUNT(*) FROM tests WHERE @keyword = '' OR name ILIKE '%' || @keyword || '%';

-- name: IsTestNameExists :one
SELECT EXISTS (SELECT 1 FROM tests WHERE name = $1 LIMIT 1);

-- name: UpdateTestByID :one
UPDATE tests
SET name = $2, price = $3, imported_price = $4, normal_value = $5, unit = $6, lower_bound = $7, upper_bound = $8, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteTestByID :exec
DELETE FROM tests WHERE id = $1;

-- name: ListAllTests :many
SELECT id, name, price, imported_price, normal_value, unit, lower_bound, upper_bound, created_at, updated_at FROM tests ORDER BY name ASC;
