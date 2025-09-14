package storage

import (
	"context"

	"github.com/datdev2409/lab-admin-go/internal/models"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Storage interface {
	// Patient
	InsertPatient(ctx context.Context, patient *models.Patient) (string, error)
	GetPatientById(ctx context.Context, id string) (*models.Patient, error)
	UpdatePatientById(ctx context.Context, id string, update models.PatientUpdate) error
	DeletePatientById(ctx context.Context, id string) error
	SearchPatientByNameOrPhone(ctx context.Context, filterOpts models.PatientQueryOptions, opts models.GenericQueryOptions) ([]*models.Patient, *models.PaginationResponse, error)
	FindPatientByNameAndPhone(ctx context.Context, name, phone string) (*models.Patient, error)
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
