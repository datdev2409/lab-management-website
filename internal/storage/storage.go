package storage

import (
	"context"
	"errors"
	"math"

	"github.com/datdev2409/lab-admin-go/internal/db/sqlc"
	"github.com/datdev2409/lab-admin-go/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Storage interface {
	// User
	CreateUser(ctx context.Context, user *models.User) (string, error)
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
	// Patient
	InsertPatient(ctx context.Context, patient *models.Patient) (string, error)
	GetPatientById(ctx context.Context, id string) (*models.Patient, error)
	UpdatePatientById(ctx context.Context, id string, update models.PatientUpdate) error
	DeletePatientById(ctx context.Context, id string) error
	SearchPatientByNameOrPhone(ctx context.Context, filterOpts models.PatientQueryOptions, opts models.GenericQueryOptions) ([]*models.Patient, *models.PaginationResponse, error)
	IsPatientExists(ctx context.Context, name, phone string) (bool, error)
	// Doctor
	InsertDoctor(ctx context.Context, doctor *models.Doctor) (string, error)
	GetDoctorById(ctx context.Context, id string) (*models.Doctor, error)
	UpdateDoctorById(ctx context.Context, id string, update models.DoctorUpdate) error
	DeleteDoctorById(ctx context.Context, id string) error
	SearchDoctorByNameOrPhone(ctx context.Context, filterOpts models.DoctorQueryOptions, opts models.GenericQueryOptions) ([]*models.Doctor, *models.PaginationResponse, error)
	FindDoctorByNameAndPhone(ctx context.Context, name, phone string) (*models.Doctor, error)
	// Test
	InsertTest(ctx context.Context, test *models.Test) (string, error)
	ListTests(ctx context.Context, filterOpts models.TestQueryOptions, opts models.GenericQueryOptions) ([]*models.Test, *models.PaginationResponse, error)
	GetTestById(ctx context.Context, id string) (*models.Test, error)
	UpdateTestById(ctx context.Context, id string, update map[string]interface{}) error
	DeleteTestById(ctx context.Context, id string) error
	// Combo
	InsertCombo(ctx context.Context, combo *models.Combo) (string, error)
	ListCombos(ctx context.Context, filterOpts models.ComboQueryOptions, opts models.GenericQueryOptions) ([]*models.Combo, *models.PaginationResponse, error)
	GetComboById(ctx context.Context, id string) (*models.Combo, error)
	UpdateComboById(ctx context.Context, id string, update map[string]interface{}) error
	UpdateComboByIdAndReturn(ctx context.Context, id string, update map[string]interface{}) (*models.Combo, error)
	DeleteComboById(ctx context.Context, id string) error
	GetTestsInCombo(ctx context.Context, comboId string) (*models.Combo, []*models.Test, error)
	GetTestsByComboId(ctx context.Context, comboId string) ([]*models.Test, error)
	// Record
	InsertRecord(ctx context.Context, record *models.Record) (string, error)
	ListRecords(ctx context.Context, filters models.RecordQueryOptions, opts models.GenericQueryOptions) ([]*models.Record, *models.PaginationResponse, error)
	GetRecordById(ctx context.Context, id string) (*models.Record, error)
	GetRecordsByIds(ctx context.Context, ids []string) ([]*models.Record, error)
	GetRecordsByPatientId(ctx context.Context, patientId string) ([]*models.Record, error)
	UpdateRecord(ctx context.Context, recordId string, updateRequest models.UpdateRecordRequest) error
	DeleteRecord(ctx context.Context, recordId string) error
	GetRecordsWithRevenue(ctx context.Context, filters models.RecordQueryOptions, opts models.GenericQueryOptions) (*models.ReportResponse, error)
	// Tracking
	InsertTracking(ctx context.Context, tracking *models.Tracking) (string, error)
	ListTrackings(ctx context.Context, filterOpts models.TrackingQueryOptions, opts models.GenericQueryOptions) ([]*models.Tracking, *models.PaginationResponse, error)
	GetTrackingById(ctx context.Context, id string) (*models.Tracking, error)
	DeleteTrackingById(ctx context.Context, id string) error
}

type MongoStorage struct {
	db *mongo.Database
}

func NewMongoStorage(dbClient *mongo.Client) *MongoStorage {
	db := dbClient.Database("labadmin")
	return &MongoStorage{db: db}
}

type PgStorage struct {
	queries *sqlc.Queries
}

func (p PgStorage) CreateUser(ctx context.Context, user *models.User) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (p PgStorage) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	//TODO implement me
	panic("implement me")
}

func (p PgStorage) GetPatientById(ctx context.Context, id string) (*models.Patient, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	sqlcPatient, err := p.queries.GetPatientByID(ctx, uid)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	patient := &models.Patient{
		ID:        sqlcPatient.ID.String(),
		Name:      sqlcPatient.Name,
		YOB:       sqlcPatient.Yob,
		Gender:    sqlcPatient.Gender,
		Address:   sqlcPatient.Address,
		Phone:     sqlcPatient.Phone,
		CreatedAt: sqlcPatient.CreatedAt,
		UpdatedAt: sqlcPatient.UpdatedAt,
	}
	return patient, nil
}

func (p PgStorage) UpdatePatientById(ctx context.Context, id string, update models.PatientUpdate) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	// fetch existing patient to preserve fields not provided in update
	existing, err := p.queries.GetPatientByID(ctx, uid)
	if err != nil {
		return err
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

	return p.queries.UpdatePatientByID(ctx, sqlc.UpdatePatientByIDParams{
		ID:      uid,
		Name:    name,
		Yob:     yob,
		Address: address,
		Phone:   phone,
	})
}

func (p PgStorage) DeletePatientById(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return p.queries.DeletePatientByID(ctx, uid)
}

func (p PgStorage) SearchPatientByNameOrPhone(ctx context.Context, filterOpts models.PatientQueryOptions, opts models.GenericQueryOptions) ([]*models.Patient, *models.PaginationResponse, error) {
	var patients []*models.Patient
	sqlcPatients, err := p.queries.SearchPatientsByNameOrPhone(ctx, sqlc.SearchPatientsByNameOrPhoneParams{
		Keyword:   filterOpts.Keyword,
		LimitArg:  int32(opts.PageSize),
		OffsetArg: int32((opts.Page - 1) * opts.PageSize),
	})
	if err != nil {
		return nil, nil, err
	}

	totalPatients, err := p.queries.CountPatientsByNameOrPhone(ctx, filterOpts.Keyword)
	if err != nil {
		return nil, nil, err
	}

	pagination := &models.PaginationResponse{
		Total:     int(totalPatients),
		TotalPage: int(math.Ceil(float64(totalPatients) / float64(opts.PageSize))),
		Page:      opts.Page,
		PageSize:  opts.PageSize,
	}

	for _, p := range sqlcPatients {
		patient := &models.Patient{
			ID:      p.ID.String(),
			Name:    p.Name,
			YOB:     p.Yob,
			Gender:  p.Gender,
			Address: p.Address,
			Phone:   p.Phone,
		}
		patients = append(patients, patient)
	}

	return patients, pagination, nil
}

func (p PgStorage) InsertDoctor(ctx context.Context, doctor *models.Doctor) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (p PgStorage) GetDoctorById(ctx context.Context, id string) (*models.Doctor, error) {
	//TODO implement me
	panic("implement me")
}

func (p PgStorage) UpdateDoctorById(ctx context.Context, id string, update models.DoctorUpdate) error {
	//TODO implement me
	panic("implement me")
}

func (p PgStorage) DeleteDoctorById(ctx context.Context, id string) error {
	//TODO implement me
	panic("implement me")
}

func (p PgStorage) SearchDoctorByNameOrPhone(ctx context.Context, filterOpts models.DoctorQueryOptions, opts models.GenericQueryOptions) ([]*models.Doctor, *models.PaginationResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (p PgStorage) FindDoctorByNameAndPhone(ctx context.Context, name, phone string) (*models.Doctor, error) {
	//TODO implement me
	panic("implement me")
}

func (p PgStorage) InsertTest(ctx context.Context, test *models.Test) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (p PgStorage) ListTests(ctx context.Context, filterOpts models.TestQueryOptions, opts models.GenericQueryOptions) ([]*models.Test, *models.PaginationResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (p PgStorage) GetTestById(ctx context.Context, id string) (*models.Test, error) {
	//TODO implement me
	panic("implement me")
}

func (p PgStorage) UpdateTestById(ctx context.Context, id string, update map[string]interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (p PgStorage) DeleteTestById(ctx context.Context, id string) error {
	//TODO implement me
	panic("implement me")
}

func (p PgStorage) InsertCombo(ctx context.Context, combo *models.Combo) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (p PgStorage) ListCombos(ctx context.Context, filterOpts models.ComboQueryOptions, opts models.GenericQueryOptions) ([]*models.Combo, *models.PaginationResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (p PgStorage) GetComboById(ctx context.Context, id string) (*models.Combo, error) {
	//TODO implement me
	panic("implement me")
}

func (p PgStorage) UpdateComboById(ctx context.Context, id string, update map[string]interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (p PgStorage) UpdateComboByIdAndReturn(ctx context.Context, id string, update map[string]interface{}) (*models.Combo, error) {
	//TODO implement me
	panic("implement me")
}

func (p PgStorage) DeleteComboById(ctx context.Context, id string) error {
	//TODO implement me
	panic("implement me")
}

func (p PgStorage) GetTestsInCombo(ctx context.Context, comboId string) (*models.Combo, []*models.Test, error) {
	//TODO implement me
	panic("implement me")
}

func (p PgStorage) GetTestsByComboId(ctx context.Context, comboId string) ([]*models.Test, error) {
	//TODO implement me
	panic("implement me")
}

func (p PgStorage) InsertRecord(ctx context.Context, record *models.Record) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (p PgStorage) ListRecords(ctx context.Context, filters models.RecordQueryOptions, opts models.GenericQueryOptions) ([]*models.Record, *models.PaginationResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (p PgStorage) GetRecordById(ctx context.Context, id string) (*models.Record, error) {
	//TODO implement me
	panic("implement me")
}

func (p PgStorage) GetRecordsByIds(ctx context.Context, ids []string) ([]*models.Record, error) {
	//TODO implement me
	panic("implement me")
}

func (p PgStorage) GetRecordsByPatientId(ctx context.Context, patientId string) ([]*models.Record, error) {
	//TODO implement me
	panic("implement me")
}

func (p PgStorage) UpdateRecord(ctx context.Context, recordId string, updateRequest models.UpdateRecordRequest) error {
	//TODO implement me
	panic("implement me")
}

func (p PgStorage) DeleteRecord(ctx context.Context, recordId string) error {
	//TODO implement me
	panic("implement me")
}

func (p PgStorage) GetRecordsWithRevenue(ctx context.Context, filters models.RecordQueryOptions, opts models.GenericQueryOptions) (*models.ReportResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (p PgStorage) InsertTracking(ctx context.Context, tracking *models.Tracking) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (p PgStorage) ListTrackings(ctx context.Context, filterOpts models.TrackingQueryOptions, opts models.GenericQueryOptions) ([]*models.Tracking, *models.PaginationResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (p PgStorage) GetTrackingById(ctx context.Context, id string) (*models.Tracking, error) {
	//TODO implement me
	panic("implement me")
}

func (p PgStorage) DeleteTrackingById(ctx context.Context, id string) error {
	//TODO implement me
	panic("implement me")
}

func NewPgStorage(db *pgxpool.Pool) *PgStorage {
	queries := sqlc.New(db)
	return &PgStorage{queries: queries}
}
