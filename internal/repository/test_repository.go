package repository

import (
	"context"
	"math"

	"github.com/datdev2409/lab-admin-go/internal/db/sqlc"
	"github.com/datdev2409/lab-admin-go/internal/models"
)

type TestRepository interface {
	InsertTest(ctx context.Context, test *models.CreateTestInput) (*models.Test, error)
	IsTestNameExists(ctx context.Context, name string) (bool, error)
	SearchTestsByName(ctx context.Context, keyword string, page, pageSize int) ([]*models.Test, *models.PaginationResponse, error)
	GetTestById(ctx context.Context, id string) (*models.Test, error)
	UpdateTestById(ctx context.Context, id string, update models.TestUpdate) (*models.Test, error)
	DeleteTestById(ctx context.Context, id string) error
	ListAllTests(ctx context.Context) ([]*models.Test, error)
}

type PgTestRepository struct {
	queries *sqlc.Queries
}

func NewPgTestRepository(queries *sqlc.Queries) *PgTestRepository {
	return &PgTestRepository{
		queries: queries,
	}
}

func (r *PgTestRepository) InsertTest(ctx context.Context, input *models.CreateTestInput) (*models.Test, error) {
	test, err := r.queries.CreateTest(ctx, sqlc.CreateTestParams{
		Name:          input.Name,
		Price:         int32(input.Price),
		ImportedPrice: int32(input.ImportedPrice),
		NormalValue:   input.NormalValue,
		Unit:          input.Unit,
		LowerBound:    toFloat8Optional(input.LowerBound),
		UpperBound:    toFloat8Optional(input.UpperBound),
	})
	if err != nil {
		return nil, err
	}

	return r.ToDomainTest(test), nil
}

func (r *PgTestRepository) IsTestNameExists(ctx context.Context, name string) (bool, error) {
	exists, err := r.queries.IsTestNameExists(ctx, name)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (r *PgTestRepository) SearchTestsByName(ctx context.Context, keyword string, page, pageSize int) ([]*models.Test, *models.PaginationResponse, error) {
	tests, err := r.queries.SearchTestsByName(ctx, sqlc.SearchTestsByNameParams{
		Keyword:   keyword,
		LimitArg:  int32(pageSize),
		OffsetArg: int32((page - 1) * pageSize),
	})
	if err != nil {
		return nil, nil, err
	}

	total, err := r.queries.CountTestsByName(ctx, keyword)
	if err != nil {
		return nil, nil, err
	}

	pagination := &models.PaginationResponse{
		Total:     int(total),
		TotalPage: int(math.Ceil(float64(total) / float64(pageSize))),
		Page:      page,
		PageSize:  pageSize,
	}

	var domainTests []*models.Test
	for _, t := range tests {
		domainTests = append(domainTests, r.ToDomainTest(t))
	}

	return domainTests, pagination, nil
}

func (r *PgTestRepository) GetTestById(ctx context.Context, id string) (*models.Test, error) {
	uid, err := parseUUID(id)
	if err != nil {
		return nil, err
	}

	test, err := r.queries.GetTestByID(ctx, uid)
	if err != nil {
		return nil, err
	}

	return r.ToDomainTest(test), nil
}

func (r *PgTestRepository) UpdateTestById(ctx context.Context, id string, update models.TestUpdate) (*models.Test, error) {
	uid, err := parseUUID(id)
	if err != nil {
		return nil, err
	}

	// Fetch existing test to preserve fields not provided in update
	existing, err := r.queries.GetTestByID(ctx, uid)
	if err != nil {
		return nil, err
	}

	name := existing.Name
	price := existing.Price
	importedPrice := existing.ImportedPrice
	normalValue := existing.NormalValue
	unit := existing.Unit
	lowerBound := existing.LowerBound
	upperBound := existing.UpperBound

	if update.Name != nil {
		name = *update.Name
	}
	if update.Price != nil {
		price = int32(*update.Price)
	}
	if update.ImportedPrice != nil {
		importedPrice = int32(*update.ImportedPrice)
	}
	if update.NormalValue != nil {
		normalValue = *update.NormalValue
	}
	if update.Unit != nil {
		unit = *update.Unit
	}
	if update.LowerBound != nil {
		lowerBound = toFloat8Optional(update.LowerBound)
	}
	if update.UpperBound != nil {
		upperBound = toFloat8Optional(update.UpperBound)
	}

	updated, err := r.queries.UpdateTestByID(ctx, sqlc.UpdateTestByIDParams{
		ID:            uid,
		Name:          name,
		Price:         price,
		ImportedPrice: importedPrice,
		NormalValue:   normalValue,
		Unit:          unit,
		LowerBound:    lowerBound,
		UpperBound:    upperBound,
	})
	if err != nil {
		return nil, err
	}

	return r.ToDomainTest(updated), nil
}

func (r *PgTestRepository) DeleteTestById(ctx context.Context, id string) error {
	uid, err := parseUUID(id)
	if err != nil {
		return err
	}

	return r.queries.DeleteTestByID(ctx, uid)
}

func (r *PgTestRepository) ListAllTests(ctx context.Context) ([]*models.Test, error) {
	tests, err := r.queries.ListAllTests(ctx)
	if err != nil {
		return nil, err
	}

	var domainTests []*models.Test
	for _, t := range tests {
		domainTests = append(domainTests, r.ToDomainTest(t))
	}

	return domainTests, nil
}

func (r *PgTestRepository) ToDomainTest(dbTest sqlc.Test) *models.Test {
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
