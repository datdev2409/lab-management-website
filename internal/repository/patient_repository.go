package repository

import (
	"context"
	"fmt"
	"math"

	"github.com/datdev2409/lab-admin-go/internal/db/sqlc"
	"github.com/datdev2409/lab-admin-go/internal/models"
	"github.com/google/uuid"
)

type PatientRepository interface {
	InsertPatient(ctx context.Context, patient *models.CreatePatientInput) (*models.Patient, error)
	IsPatientExists(ctx context.Context, name, phone string) (bool, error)
	SearchPatientsByKeyword(ctx context.Context, keyword string, page, pageSize int) ([]*models.Patient, *models.PaginationResponse, error)
	GetPatientById(ctx context.Context, id string) (*models.Patient, error)
	UpdatePatientById(ctx context.Context, id string, update models.PatientUpdate) (*models.Patient, error)
	DeletePatientById(ctx context.Context, id string) error
}

type PgPatientRepository struct {
	queries *sqlc.Queries
}

func NewPgPatientRepository(queries *sqlc.Queries) *PgPatientRepository {
	return &PgPatientRepository{
		queries: queries,
	}
}

func parseUUID(id string) (uuid.UUID, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid UUID: %w", err)
	}
	return uid, nil
}

func (r *PgPatientRepository) InsertPatient(ctx context.Context, input *models.CreatePatientInput) (*models.Patient, error) {
	patient, err := r.queries.CreatePatient(ctx, sqlc.CreatePatientParams{
		Name:    input.Name,
		Yob:     input.YOB,
		Gender:  input.Gender,
		Address: input.Address,
		Phone:   input.Phone,
	})
	if err != nil {
		return nil, err
	}

	return r.ToDomainPatient(patient), nil
}

func (r *PgPatientRepository) IsPatientExists(ctx context.Context, name, phone string) (bool, error) {
	exists, err := r.queries.IsPatientNameAndPhoneExists(ctx, sqlc.IsPatientNameAndPhoneExistsParams{
		Name:  name,
		Phone: phone,
	})
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (r *PgPatientRepository) SearchPatientsByKeyword(ctx context.Context, keyword string, page, pageSize int) ([]*models.Patient, *models.PaginationResponse, error) {
	patients, err := r.queries.SearchPatientsByNameOrPhone(ctx, sqlc.SearchPatientsByNameOrPhoneParams{
		Keyword:   keyword,
		LimitArg:  int32(pageSize),
		OffsetArg: int32((page - 1) * pageSize),
	})
	if err != nil {
		return nil, nil, err
	}
	total, err := r.queries.CountPatientsByNameOrPhone(ctx, keyword)
	if err != nil {
		return nil, nil, err
	}

	pagination := &models.PaginationResponse{
		Total:     int(total),
		TotalPage: int(math.Ceil(float64(total) / float64(pageSize))),
		Page:      page,
		PageSize:  pageSize,
	}

	var domainPatients []*models.Patient
	for _, p := range patients {
		domainPatients = append(domainPatients, r.ToDomainPatient(p))
	}

	return domainPatients, pagination, nil
}

func (r *PgPatientRepository) GetPatientById(ctx context.Context, id string) (*models.Patient, error) {
	uid, err := parseUUID(id)
	if err != nil {
		return nil, err
	}

	patient, err := r.queries.GetPatientByID(ctx, uid)
	if err != nil {
		return nil, err
	}

	return r.ToDomainPatient(patient), nil
}

func (r *PgPatientRepository) UpdatePatientById(ctx context.Context, id string, update models.PatientUpdate) (*models.Patient, error) {
	uid, err := parseUUID(id)
	if err != nil {
		return nil, err
	}

	// Fetch existing patient to preserve fields not provided in update
	existing, err := r.queries.GetPatientByID(ctx, uid)
	if err != nil {
		return nil, err
	}

	name := existing.Name
	yob := existing.Yob
	address := existing.Address
	phone := existing.Phone

	if update.Name != nil {
		name = *update.Name
	}
	if update.YOB != nil {
		yob = *update.YOB
	}
	if update.Address != nil {
		address = *update.Address
	}
	if update.Phone != nil {
		phone = *update.Phone
	}

	err = r.queries.UpdatePatientByID(ctx, sqlc.UpdatePatientByIDParams{
		ID:      uid,
		Name:    name,
		Yob:     yob,
		Address: address,
		Phone:   phone,
	})
	if err != nil {
		return nil, err
	}

	// Return updated patient
	updated, err := r.queries.GetPatientByID(ctx, uid)
	if err != nil {
		return nil, err
	}

	return r.ToDomainPatient(updated), nil
}

func (r *PgPatientRepository) DeletePatientById(ctx context.Context, id string) error {
	uid, err := parseUUID(id)
	if err != nil {
		return err
	}

	return r.queries.DeletePatientByID(ctx, uid)
}

func (r *PgPatientRepository) ToDomainPatient(dbPatient sqlc.Patient) *models.Patient {
	return &models.Patient{
		ID:        dbPatient.ID.String(),
		Name:      dbPatient.Name,
		YOB:       dbPatient.Yob,
		Gender:    dbPatient.Gender,
		Address:   dbPatient.Address,
		Phone:     dbPatient.Phone,
		CreatedAt: dbPatient.CreatedAt,
		UpdatedAt: dbPatient.UpdatedAt,
	}
}
