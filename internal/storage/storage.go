package storage

import (
	"context"

	"github.com/datdev2409/lab-admin-go/internal/models"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type TestStorage interface {
	Insert(ctx context.Context, test models.Test) error
	ListTests(ctx context.Context, filterOpts models.TestQueryOptions, opts models.GenericQueryOptions) ([]*models.Test, *models.PaginationResponse, error)
	GetById(ctx context.Context, id string) (*models.Test, error)
	UpdateById(ctx context.Context, id string, update interface{}) error
	DeleteById(ctx context.Context, id string) error
}

type ComboStorage interface {
	Insert(ctx context.Context, combo models.Combo) error
	ListCombos(ctx context.Context, filterOpts models.ComboQueryOptions, opts models.GenericQueryOptions) ([]*models.Combo, *models.PaginationResponse, error)
	GetById(ctx context.Context, id string) (*models.Combo, error)
	GetTestsInCombo(ctx context.Context, comboId string) (*models.Combo, []*models.Test, error)
	UpdateById(ctx context.Context, id string, update interface{}) error
	DeleteById(ctx context.Context, id string) error
}

type TrackingStorage interface {
	Insert(ctx context.Context, tracking models.Tracking) error
	ListTrackings(ctx context.Context, filterOpts models.TrackingQueryOptions, opts models.GenericQueryOptions) ([]*models.Tracking, *models.PaginationResponse, error)
	GetById(ctx context.Context, id string) (*models.Tracking, error)
}

type AppStorage interface {
	Patients() PatientStorage
	Tests() TestStorage
	Combos() ComboStorage
	Records() RecordStorage
	Trackings() TrackingStorage
}

func NewMongoStorage(client *mongo.Client) *MongoStorage {
	return &MongoStorage{db: client.Database("labadmin")}
}

type MongoStorage struct {
	db *mongo.Database
}

type MongoPatientStorage struct {
	*BaseRepository[models.Patient]
}

type MongoTestStorage struct {
	*BaseRepository[models.Test]
}

type MongoComboStorage struct {
	*BaseRepository[models.Combo]
}

type MongoRecordStorage struct {
	*BaseRepository[models.Record]
}

type MongoTrackingStorage struct {
	*BaseRepository[models.Tracking]
}

func (m *MongoStorage) Patients() PatientStorage {
	return &MongoPatientStorage{
		BaseRepository: NewBaseRepository[models.Patient](m.db, "patients"),
	}
}

func (m *MongoStorage) Tests() TestStorage {
	return &MongoTestStorage{
		BaseRepository: NewBaseRepository[models.Test](m.db, "tests"),
	}
}

func (m *MongoStorage) Combos() ComboStorage {
	return &MongoComboStorage{
		BaseRepository: NewBaseRepository[models.Combo](m.db, "combos"),
	}
}

func (m *MongoStorage) Records() RecordStorage {
	return &MongoRecordStorage{
		BaseRepository: NewBaseRepository[models.Record](m.db, "records"),
	}
}

func (m *MongoStorage) Trackings() TrackingStorage {
	return &MongoTrackingStorage{
		BaseRepository: NewBaseRepository[models.Tracking](m.db, "trackings"),
	}
}
