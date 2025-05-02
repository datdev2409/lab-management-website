package storage

import (
	"context"

	"github.com/datdev2409/lab-admin-go/internal/models"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type PatientStorage interface {
	Insert(patient *models.Patient) error
	GetById(id string) (*models.Patient, error)
	ListPatients(ctx context.Context, filterOpts models.PatientQueryOptions, opts models.GenericQueryOptions) ([]*models.Patient, *models.PaginationResponse, error)
	SearchByKeyword(ctx context.Context, keyword string, opts map[string]string) (*[]models.Patient, error)
	UpdateById(ctx context.Context, id string, patient *models.Patient) error
	Delete(id string) error
}

type TestStorage interface {
	Insert(test *models.Test) error
	ListTests(ctx context.Context, filterOpts models.TestQueryOptions, opts models.GenericQueryOptions) ([]*models.Test, *models.PaginationResponse, error)
	GetById(id string) (*models.Test, error)
	GetByIds(ctx context.Context, ids []string) ([]*models.Test, error)
	SearchByKeyword(ctx context.Context, keyword string, opts map[string]string) ([]*models.Test, error)
	Update(test *models.Test) error
	Delete(id string) error
}

type ComboStorage interface {
	Insert(combo *models.Combo) error
	ListCombos(ctx context.Context, filterOpts models.ComboQueryOptions, opts models.GenericQueryOptions) ([]*models.Combo, *models.PaginationResponse, error)
	GetById(ctx context.Context, id string) (*models.Combo, error)
	SearchByKeyword(ctx context.Context, keyword string, opts map[string]string) ([]*models.Combo, error)
	GetTestsInCombo(ctx context.Context, comboId string) (*models.Combo, []*models.Test, error)
}

type RecordStorage interface {
	Insert(ctx context.Context, record *models.Record) (string, error)
	// GetById(ctx context.Context, id string) (*models.Record, error)
	// GetDetails(ctx context.Context, id string) (*models.RecordWithDetails, error)
	// ListByPatientId(ctx context.Context, patientId string) (*[]models.Record, error)
	ListRecords(ctx context.Context, filters models.RecordQueryOptions, opts models.GenericQueryOptions) (*[]models.Record, *models.PaginationResponse, error)
	// UpdatePatient(ctx context.Context, recordId string, patient models.Patient) error
	// UpdateCombo(ctx context.Context, recordId string, combo *models.Combo) error
	// AddTest(ctx context.Context, recordId string, test *models.Test) error
	// AddTests(ctx context.Context, recordId string, tests []*models.Test) error
	// SaveTestResults(ctx context.Context, recordId string, testResults []models.TestResult) error
}

type AppStorage interface {
	Patients() PatientStorage
	Tests() TestStorage
	Combos() ComboStorage
	Records() RecordStorage
}

func NewMongoStorage(client *mongo.Client) *MongoStorage {
	return &MongoStorage{db: client.Database("labadmin")}
}

type MongoStorage struct {
	db *mongo.Database
}

func (m *MongoStorage) Patients() PatientStorage {
	return &MongoPatientStorage{db: m.db, col: m.db.Collection("patients")}
}

func (m *MongoStorage) Tests() TestStorage {
	return &MongoTestStorage{db: m.db, col: m.db.Collection("tests")}
}

func (m *MongoStorage) Combos() ComboStorage {
	return &MongoComboStorage{db: m.db, col: m.db.Collection("combos")}
}

func (m *MongoStorage) Records() RecordStorage {
	return &MongoRecordStorage{db: m.db, col: m.db.Collection("records")}
}

type MongoPatientStorage struct {
	db  *mongo.Database
	col *mongo.Collection
}

type MongoTestStorage struct {
	db  *mongo.Database
	col *mongo.Collection
}

type MongoComboStorage struct {
	db  *mongo.Database
	col *mongo.Collection
}

type MongoRecordStorage struct {
	db  *mongo.Database
	col *mongo.Collection
}
