-- name: CreateCombo :one
INSERT INTO combos (name, created_at, updated_at)
VALUES ($1, NOW(), NOW())
RETURNING *;

-- name: GetComboByID :one
SELECT * FROM combos WHERE id = $1;

-- name: GetComboByIDWithTests :one
SELECT c.id, c.name, c.created_at, c.updated_at
FROM combos c
WHERE c.id = $1;

-- name: SearchCombosByName :many
SELECT * FROM combos WHERE @keyword = '' OR name ILIKE '%' || @keyword || '%' ORDER BY name ASC LIMIT @limit_arg OFFSET @offset_arg;

-- name: CountCombosByName :one
SELECT COUNT(*) FROM combos WHERE @keyword = '' OR name ILIKE '%' || @keyword || '%';

-- name: IsComboNameExists :one
SELECT EXISTS (SELECT 1 FROM combos WHERE name = $1 LIMIT 1);

-- name: UpdateComboByID :one
UPDATE combos
SET name = $2, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteComboByID :exec
DELETE FROM combos WHERE id = $1;

-- name: ListAllCombos :many
SELECT id, name, created_at, updated_at FROM combos ORDER BY name ASC;

-- name: AddTestToCombo :exec
INSERT INTO combo_tests (combo_id, test_id, test_order)
VALUES ($1, $2, $3)
ON CONFLICT DO NOTHING;

-- name: RemoveTestFromCombo :exec
DELETE FROM combo_tests WHERE combo_id = $1 AND test_id = $2;

-- name: GetComboTests :many
SELECT t.id, t.name, t.price, t.imported_price, t.normal_value, t.unit, t.lower_bound, t.upper_bound, t.created_at, t.updated_at
FROM tests t
INNER JOIN combo_tests ct ON t.id = ct.test_id
WHERE ct.combo_id = $1
ORDER BY ct.test_order ASC, t.name ASC;

-- name: GetComboTestCount :one
SELECT COUNT(*) FROM combo_tests WHERE combo_id = $1;

-- name: RemoveAllTestsFromCombo :exec
DELETE FROM combo_tests WHERE combo_id = $1;
