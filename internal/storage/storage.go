package storage

import (
	"context"
	"fmt"
	"strconv"

	"github.com/datdev2409/lab-admin-go/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type PatientStorage interface {
	Insert(patient *models.Patient) error
	GetById(id string) (*models.Patient, error)
	SearchByKeyword(ctx context.Context, keyword string, opts map[string]string) (*[]models.Patient, error)
	UpdateById(ctx context.Context, id string, patient *models.Patient) error
	Delete(id string) error
}

type TestStorage interface {
	Insert(test *models.Test) error
	GetById(id string) (*models.Test, error)
	GetByIds(ctx context.Context, ids []string) (*[]models.Test, error)
	SearchByKeyword(ctx context.Context, keyword string, opts map[string]string) (*[]models.Test, error)
	Update(test *models.Test) error
	Delete(id string) error
}

type ComboStorage interface {
	Insert(combo *models.Combo) error
	GetById(ctx context.Context, id string) (*models.Combo, error)
	SearchByKeyword(ctx context.Context, keyword string, opts map[string]string) (*[]models.Combo, error)
}

type RecordStorage interface {
	Insert(ctx context.Context, record *models.Record) error
	GetById(ctx context.Context, id string) (*models.Record, error)
	GetAllInformation(ctx context.Context, id string) (*models.Record, *models.Patient, *[]models.Test, error)
	SearchByKeyword(ctx context.Context, keyword string, opts map[string]string) (*[]models.Record, error)
	UpdatePatient(ctx context.Context, recordId string, patient models.Patient) error
	UpdateCombo(ctx context.Context, recordId string, combo *models.Combo) error
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
	return &MongoPatientStorage{col: m.db.Collection("patients")}
}

func (m *MongoStorage) Tests() TestStorage {
	return &MongoTestStorage{col: m.db.Collection("tests")}
}

func (m *MongoStorage) Combos() ComboStorage {
	return &MongoComboStorage{col: m.db.Collection("combos")}
}

func (m *MongoStorage) Records() RecordStorage {
	return &MongoRecordStorage{storage: m, col: m.db.Collection("records")}
}

type MongoPatientStorage struct {
	col *mongo.Collection
}

type MongoTestStorage struct {
	col *mongo.Collection
}

type MongoComboStorage struct {
	col *mongo.Collection
}

type MongoRecordStorage struct {
	storage *MongoStorage
	col     *mongo.Collection
}

// SearchByKeyword implements RecordStorage.
func (m *MongoRecordStorage) SearchByKeyword(ctx context.Context, keyword string, opts map[string]string) (*[]models.Record, error) {
	records := []models.Record{}

	// Support filter by patient name and phone
	filters := bson.D{
		{Key: "$or", Value: bson.A{
			bson.D{{Key: "patient.name", Value: bson.D{{Key: "$regex", Value: keyword}, {Key: "$options", Value: "i"}}}},
			bson.D{{Key: "patient.phone", Value: bson.D{{Key: "$regex", Value: keyword}, {Key: "$options", Value: "i"}}}},
		}},
	}

	findOpts := options.Find()
	if val, ok := opts["limit"]; ok {
		limit, err := strconv.Atoi(val)
		if err != nil {
			limit = 5
		}
		findOpts.SetLimit(int64(limit))
	}

	cursor, err := m.col.Find(ctx, filters, findOpts)
	if err != nil {
		return &records, err
	}

	if err = cursor.All(ctx, &records); err != nil {
		return &records, err
	}

	return &records, nil
}

func (m *MongoRecordStorage) GetAllInformation(ctx context.Context, id string) (*models.Record, *models.Patient, *[]models.Test, error) {
	lookupPipeline := bson.A{
		bson.D{{Key: "$match", Value: bson.D{{Key: "_id", Value: id}}}},
		bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "patients"},
			{Key: "localField", Value: "patient.id"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "patient"},
		}},
		},
	}

	// m.storage.db.Aggregate(ctx, lookupPipeline)
	cursor, err := m.col.Aggregate(ctx, lookupPipeline)
	if err != nil {
		return nil, nil, nil, err
	}

	var results []bson.M
	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}
	for _, result := range results {
		fmt.Println(result)
	}
	return nil, nil, nil, nil
	// record, err := m.GetById(ctx, id)
	// if err != nil {
	// 	return nil, nil, nil, err
	// }

	// patient, err := m.getPatient(ctx, record.PatientID)
	// if err != nil {
	// 	return nil, nil, nil, err
	// }

	// tests, err := m.getTests(ctx, record.TestIDs)
	// if err != nil {
	// 	return nil, nil, nil, err
	// }

	// return record, patient, tests, nil
}
