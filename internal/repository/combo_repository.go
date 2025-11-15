package repository

import (
	"context"
	"fmt"
	"math"

	"github.com/datdev2409/lab-admin-go/internal/db/sqlc"
	"github.com/datdev2409/lab-admin-go/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ComboRepository interface {
	InsertCombo(ctx context.Context, name string, testIDs []string) (string, error)
	IsComboNameExists(ctx context.Context, name string) (bool, error)
	SearchCombosByName(ctx context.Context, keyword string, page, pageSize int) ([]*models.Combo, *models.PaginationResponse, error)
	GetComboById(ctx context.Context, id string) (*models.Combo, error)
	UpdateComboById(ctx context.Context, id string, update models.ComboUpdate) (*models.Combo, error)
	DeleteComboById(ctx context.Context, id string) error
	ListAllCombos(ctx context.Context) ([]*models.Combo, error)
	GetComboTests(ctx context.Context, comboId string) ([]*models.Test, error)
}

type PgComboRepository struct {
	queries *sqlc.Queries
	pool    *pgxpool.Pool
}

func NewPgComboRepository(queries *sqlc.Queries, pool *pgxpool.Pool) *PgComboRepository {
	return &PgComboRepository{queries: queries, pool: pool}
}

// InsertCombo inserts a combo and its related combo_tests rows inside a single transaction.
// Returns the created combo id as string.
func (r *PgComboRepository) InsertCombo(ctx context.Context, name string, testIDs []string) (string, error) {
	// start transaction
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to begin tx: %w", err)
	}
	defer func() {
		// if still in progress, rollback (safe to call after commit)
		_ = tx.Rollback(ctx)
	}()

	q := r.queries.WithTx(tx)

	// create combo
	combo, err := q.CreateCombo(ctx, name)
	if err != nil {
		return "", fmt.Errorf("failed to create combo: %w", err)
	}

	// insert combo_tests entries with ordering
	for i, tid := range testIDs {
		uid, err := uuid.Parse(tid)
		if err != nil {
			return "", fmt.Errorf("invalid test id %q: %w", tid, err)
		}
		// Use AddTestToCombo sqlc function instead of raw SQL
		if err := q.AddTestToCombo(ctx, sqlc.AddTestToComboParams{
			ComboID:   combo.ID,
			TestID:    uid,
			TestOrder: int32(i),
		}); err != nil {
			return "", fmt.Errorf("failed to add test %s to combo: %w", tid, err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return "", fmt.Errorf("failed to commit tx: %w", err)
	}

	return combo.ID.String(), nil
}

func (r *PgComboRepository) IsComboNameExists(ctx context.Context, name string) (bool, error) {
	exists, err := r.queries.IsComboNameExists(ctx, name)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (r *PgComboRepository) SearchCombosByName(ctx context.Context, keyword string, page, pageSize int) ([]*models.Combo, *models.PaginationResponse, error) {
	combos, err := r.queries.SearchCombosByName(ctx, sqlc.SearchCombosByNameParams{
		Keyword:   keyword,
		LimitArg:  int32(pageSize),
		OffsetArg: int32((page - 1) * pageSize),
	})
	if err != nil {
		return nil, nil, err
	}

	total, err := r.queries.CountCombosByName(ctx, keyword)
	if err != nil {
		return nil, nil, err
	}

	pagination := &models.PaginationResponse{
		Total:     int(total),
		TotalPage: int(math.Ceil(float64(total) / float64(pageSize))),
		Page:      page,
		PageSize:  pageSize,
	}

	var domainCombos []*models.Combo
	for _, c := range combos {
		domainCombos = append(domainCombos, r.ToDomainCombo(c))
	}

	return domainCombos, pagination, nil
}

func (r *PgComboRepository) GetComboById(ctx context.Context, id string) (*models.Combo, error) {
	uid, err := parseUUID(id)
	if err != nil {
		return nil, err
	}

	combo, err := r.queries.GetComboByID(ctx, uid)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("combo not found")
		}
		return nil, err
	}

	return r.ToDomainCombo(combo), nil
}

func (r *PgComboRepository) UpdateComboById(ctx context.Context, id string, update models.ComboUpdate) (*models.Combo, error) {
	uid, err := parseUUID(id)
	if err != nil {
		return nil, err
	}

	// Fetch existing combo to preserve fields not provided in update
	existing, err := r.queries.GetComboByID(ctx, uid)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("combo not found")
		}
		return nil, err
	}

	name := existing.Name
	if update.Name != nil {
		name = *update.Name
	}

	// If tests are being updated, use transaction to update both combo and tests
	if len(update.TestIDs) > 0 {
		return r.updateComboWithTests(ctx, uid, name, update.TestIDs)
	}

	// Just update the name if no tests are provided
	updated, err := r.queries.UpdateComboByID(ctx, sqlc.UpdateComboByIDParams{
		ID:   uid,
		Name: name,
	})
	if err != nil {
		return nil, err
	}

	return r.ToDomainCombo(updated), nil
}

// updateComboWithTests handles updating combo with test list in a transaction
func (r *PgComboRepository) updateComboWithTests(ctx context.Context, comboID uuid.UUID, name string, testIDs []string) (*models.Combo, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin tx: %w", err)
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	q := r.queries.WithTx(tx)

	// Update combo name
	updated, err := q.UpdateComboByID(ctx, sqlc.UpdateComboByIDParams{
		ID:   comboID,
		Name: name,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update combo: %w", err)
	}

	// Remove all existing tests for this combo
	err = q.RemoveAllTestsFromCombo(ctx, comboID)
	if err != nil {
		return nil, fmt.Errorf("failed to remove existing tests: %w", err)
	}

	// Add new tests with order
	for order, testID := range testIDs {
		testUUID, err := parseUUID(testID)
		if err != nil {
			return nil, fmt.Errorf("invalid test id: %w", err)
		}

		err = q.AddTestToCombo(ctx, sqlc.AddTestToComboParams{
			ComboID:   comboID,
			TestID:    testUUID,
			TestOrder: int32(order),
		})
		if err != nil {
			return nil, fmt.Errorf("failed to add test to combo: %w", err)
		}
	}

	// Commit transaction
	err = tx.Commit(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to commit tx: %w", err)
	}

	return r.ToDomainCombo(updated), nil
}

func (r *PgComboRepository) DeleteComboById(ctx context.Context, id string) error {
	uid, err := parseUUID(id)
	if err != nil {
		return err
	}

	return r.queries.DeleteComboByID(ctx, uid)
}

func (r *PgComboRepository) ListAllCombos(ctx context.Context) ([]*models.Combo, error) {
	combos, err := r.queries.ListAllCombos(ctx)
	if err != nil {
		return nil, err
	}

	var domainCombos []*models.Combo
	for _, c := range combos {
		domainCombos = append(domainCombos, r.ToDomainCombo(c))
	}

	return domainCombos, nil
}

func (r *PgComboRepository) GetComboTests(ctx context.Context, comboId string) ([]*models.Test, error) {
	uid, err := parseUUID(comboId)
	if err != nil {
		return nil, err
	}

	tests, err := r.queries.GetComboTests(ctx, uid)
	if err != nil {
		return nil, err
	}

	var domainTests []*models.Test
	for _, t := range tests {
		domainTests = append(domainTests, r.testSqlcToDomain(t))
	}

	return domainTests, nil
}

func (r *PgComboRepository) ToDomainCombo(dbCombo sqlc.Combo) *models.Combo {
	return &models.Combo{
		ID:        dbCombo.ID.String(),
		Name:      dbCombo.Name,
		CreatedAt: dbCombo.CreatedAt,
		UpdatedAt: dbCombo.UpdatedAt,
	}
}

// testSqlcToDomain converts sqlc Test to domain Test model
func (r *PgComboRepository) testSqlcToDomain(dbTest sqlc.Test) *models.Test {
	return &models.Test{
		ID:            dbTest.ID.String(),
		Name:          dbTest.Name,
		Price:         int(dbTest.Price),
		ImportedPrice: int(dbTest.ImportedPrice),
		NormalValue:   dbTest.NormalValue,
		Unit:          dbTest.Unit,
		LowerBound:    fromFloat8Optional(dbTest.LowerBound),
		UpperBound:    fromFloat8Optional(dbTest.UpperBound),
		CreatedAt:     dbTest.CreatedAt,
		UpdatedAt:     dbTest.UpdatedAt,
	}
}
