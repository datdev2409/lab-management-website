package repository

import (
	"context"
	"math"

	"github.com/datdev2409/lab-admin-go/internal/db/sqlc"
	"github.com/datdev2409/lab-admin-go/internal/models"
)

type DoctorRepository interface {
	InsertDoctor(ctx context.Context, input *models.CreateDoctorInput) (*models.Doctor, error)
	IsDoctorExists(ctx context.Context, name, phone string) (bool, error)
	SearchDoctorsByKeyword(ctx context.Context, keyword string, page, pageSize int) ([]*models.Doctor, *models.PaginationResponse, error)
	GetDoctorById(ctx context.Context, id string) (*models.Doctor, error)
	UpdateDoctorById(ctx context.Context, id string, update models.DoctorUpdate) (*models.Doctor, error)
	DeleteDoctorById(ctx context.Context, id string) error
}

type PgDoctorRepository struct {
	queries *sqlc.Queries
}

func NewPgDoctorRepository(queries *sqlc.Queries) *PgDoctorRepository {
	return &PgDoctorRepository{
		queries: queries,
	}
}

func (r *PgDoctorRepository) InsertDoctor(ctx context.Context, input *models.CreateDoctorInput) (*models.Doctor, error) {
	doctor, err := r.queries.CreateDoctor(ctx, sqlc.CreateDoctorParams{
		Name:    input.Name,
		Phone:   input.Phone,
		Address: input.Address,
	})
	if err != nil {
		return nil, err
	}

	return r.ToDomainDoctor(doctor), nil
}

func (r *PgDoctorRepository) IsDoctorExists(ctx context.Context, name, phone string) (bool, error) {
	exists, err := r.queries.IsDoctorNameAndPhoneExists(ctx, sqlc.IsDoctorNameAndPhoneExistsParams{
		Name:  name,
		Phone: phone,
	})
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (r *PgDoctorRepository) SearchDoctorsByKeyword(ctx context.Context, keyword string, page, pageSize int) ([]*models.Doctor, *models.PaginationResponse, error) {
	doctors, err := r.queries.SearchDoctorsByNameOrPhone(ctx, sqlc.SearchDoctorsByNameOrPhoneParams{
		Keyword:   keyword,
		LimitArg:  int32(pageSize),
		OffsetArg: int32((page - 1) * pageSize),
	})
	if err != nil {
		return nil, nil, err
	}

	total, err := r.queries.CountDoctorsByNameOrPhone(ctx, keyword)
	if err != nil {
		return nil, nil, err
	}

	pagination := &models.PaginationResponse{
		Total:     int(total),
		TotalPage: int(math.Ceil(float64(total) / float64(pageSize))),
		Page:      page,
		PageSize:  pageSize,
	}

	var domainDoctors []*models.Doctor
	for _, d := range doctors {
		domainDoctors = append(domainDoctors, r.ToDomainDoctor(d))
	}

	return domainDoctors, pagination, nil
}

func (r *PgDoctorRepository) GetDoctorById(ctx context.Context, id string) (*models.Doctor, error) {
	uid, err := parseUUID(id)
	if err != nil {
		return nil, err
	}

	doctor, err := r.queries.GetDoctorByID(ctx, uid)
	if err != nil {
		return nil, err
	}

	return r.ToDomainDoctor(doctor), nil
}

func (r *PgDoctorRepository) UpdateDoctorById(ctx context.Context, id string, update models.DoctorUpdate) (*models.Doctor, error) {
	uid, err := parseUUID(id)
	if err != nil {
		return nil, err
	}

	// Fetch existing doctor to preserve fields not provided in update
	existing, err := r.queries.GetDoctorByID(ctx, uid)
	if err != nil {
		return nil, err
	}

	name := existing.Name
	phone := existing.Phone
	address := existing.Address

	if update.Name != nil {
		name = *update.Name
	}
	if update.Phone != nil {
		phone = *update.Phone
	}
	if update.Address != nil {
		address = *update.Address
	}

	err = r.queries.UpdateDoctorByID(ctx, sqlc.UpdateDoctorByIDParams{
		ID:      uid,
		Name:    name,
		Phone:   phone,
		Address: address,
	})
	if err != nil {
		return nil, err
	}

	// Return updated doctor
	updated, err := r.queries.GetDoctorByID(ctx, uid)
	if err != nil {
		return nil, err
	}

	return r.ToDomainDoctor(updated), nil
}

func (r *PgDoctorRepository) DeleteDoctorById(ctx context.Context, id string) error {
	uid, err := parseUUID(id)
	if err != nil {
		return err
	}

	return r.queries.DeleteDoctorByID(ctx, uid)
}

func (r *PgDoctorRepository) ToDomainDoctor(dbDoctor sqlc.Doctor) *models.Doctor {
	return &models.Doctor{
		ID:        dbDoctor.ID.String(),
		Name:      dbDoctor.Name,
		Phone:     dbDoctor.Phone,
		Address:   dbDoctor.Address,
		CreatedAt: dbDoctor.CreatedAt,
		UpdatedAt: dbDoctor.UpdatedAt,
	}
}
