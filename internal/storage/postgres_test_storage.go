package storage

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/datdev2409/lab-admin-go/internal/models"
)

// InsertTest inserts a new test record
func (p *PostgresStorage) InsertTest(ctx context.Context, test *models.Test) (string, error) {
	// Generate UUID for the new test
	testID := uuid.New().String()

	query := `
		INSERT INTO tests (id, name, price, normal_value, unit, lower_bound, upper_bound, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id`

	now := time.Now()

	rows, _ := p.pool.Query(ctx, query,
		testID,
		test.Name,
		test.Price,
		test.NormalValue,
		test.Unit,
		test.LowerBound,
		test.UpperBound,
		now,
		now,
	)
	defer rows.Close()

	returnedID, err := pgx.CollectExactlyOneRow(rows, pgx.RowTo[string])
	if err != nil {
		return "", fmt.Errorf("failed to insert test: %w", err)
	}

	return returnedID, nil
}

// GetTestById retrieves a test by ID
func (p *PostgresStorage) GetTestById(ctx context.Context, id string) (*models.Test, error) {
	query := `
		SELECT id, name, price, normal_value, unit, lower_bound, upper_bound, created_at, updated_at
		FROM tests 
		WHERE id = $1`

	rows, _ := p.pool.Query(ctx, query, id)
	defer rows.Close()

	test, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[models.Test])
	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("test not found with id: %s", id)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get test by id: %w", err)
	}

	return &test, nil
}

// UpdateTestById updates a test by ID
func (p *PostgresStorage) UpdateTestById(ctx context.Context, id string, update map[string]interface{}) error {
	if len(update) == 0 {
		return fmt.Errorf("no update fields provided")
	}

	// Build dynamic SET clause with named args
	setParts := []string{}
	args := pgx.NamedArgs{
		"id":         id,
		"updated_at": time.Now(),
	}

	for field, value := range update {
		switch field {
		case "name", "price", "normal_value", "unit", "lower_bound", "upper_bound":
			setParts = append(setParts, fmt.Sprintf("%s = @%s", field, field))
			args[field] = value
		default:
			return fmt.Errorf("invalid update field: %s", field)
		}
	}

	// Always update updated_at
	setParts = append(setParts, "updated_at = @updated_at")

	query := fmt.Sprintf(`
		UPDATE tests 
		SET %s 
		WHERE id = @id`,
		strings.Join(setParts, ", "),
	)

	result, err := p.pool.Exec(ctx, query, args)
	if err != nil {
		return fmt.Errorf("failed to update test: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("test not found with id: %s", id)
	}

	return nil
}

// DeleteTestById deletes a test by ID
func (p *PostgresStorage) DeleteTestById(ctx context.Context, id string) error {
	query := `DELETE FROM tests WHERE id = $1`

	result, err := p.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete test: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("test not found with id: %s", id)
	}

	return nil
}

// ListTests retrieves a list of tests with filtering and pagination
func (p *PostgresStorage) ListTests(ctx context.Context, filterOpts models.TestQueryOptions, opts models.GenericQueryOptions) ([]*models.Test, *models.PaginationResponse, error) {
	pageSize := opts.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}
	offset := (opts.Page - 1) * pageSize
	if offset < 0 {
		offset = 0
	}

	keyword := strings.ToLower(filterOpts.Keyword)
	args := pgx.NamedArgs{
		"keyword": "%" + keyword + "%",
	}

	whereClauses := ""
	if keyword != "" {
		whereClauses = "WHERE LOWER(name) LIKE @keyword"
	}

	// Count total
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM tests %s", whereClauses)
	var total int64
	if err := p.pool.QueryRow(ctx, countQuery, args).Scan(&total); err != nil {
		return nil, nil, fmt.Errorf("failed to count tests: %w", err)
	}

	// Select with pagination
	args["limit"] = pageSize
	args["offset"] = offset
	selectQuery := fmt.Sprintf(`
		SELECT id, name, price, normal_value, unit, lower_bound, upper_bound, created_at, updated_at
		FROM tests
		%s
		ORDER BY created_at DESC
		LIMIT @limit OFFSET @offset
	`, whereClauses)

	rows, err := p.pool.Query(ctx, selectQuery, args)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to list tests: %w", err)
	}
	defer rows.Close()

	tests, err := pgx.CollectRows(rows, pgx.RowToStructByName[*models.Test])
	if err != nil {
		return nil, nil, fmt.Errorf("failed to collect tests: %w", err)
	}

	pagination := &models.PaginationResponse{
		Total:     int(total),
		Page:      opts.Page,
		PageSize:  pageSize,
		TotalPage: int((total + int64(pageSize) - 1) / int64(pageSize)),
	}

	if pageSize == 0 {
		pagination.TotalPage = 1
		pagination.Page = 1
	}

	return tests, pagination, nil
}
