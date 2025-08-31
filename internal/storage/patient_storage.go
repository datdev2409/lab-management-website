package storage

import (
	"context"

	"github.com/datdev2409/lab-admin-go/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func (m *MongoStorage) InsertPatient(ctx context.Context, patient *models.Patient) (string, error) {
	col := m.getCollection("patients")
	return MongoInsert(ctx, col, patient)
}

func (m *MongoStorage) GetPatientById(ctx context.Context, id string) (*models.Patient, error) {
	col := m.getCollection("patients")
	return MongoGetById[models.Patient](ctx, col, id)
}

func (m *MongoStorage) UpdatePatientById(ctx context.Context, id string, update models.PatientUpdate) error {
	updateBSON := bson.M{}
	if update.Name != nil {
		updateBSON["name"] = *update.Name
	}
	if update.Phone != nil {
		updateBSON["phone"] = *update.Phone
	}
	if update.Address != nil {
		updateBSON["address"] = *update.Address
	}
	if update.YOB != nil {
		updateBSON["yob"] = *update.YOB
	}
	if update.Gender != nil {
		updateBSON["gender"] = *update.Gender
	}
	col := m.getCollection("patients")
	return MongoUpdateById[models.Patient](ctx, col, id, bson.M{"$set": updateBSON})
}

func (m *MongoStorage) DeletePatientById(ctx context.Context, id string) error {
	col := m.getCollection("patients")
	return MongoDeleteById[models.Patient](ctx, col, id)
}

func (m *MongoStorage) SearchPatientByNameOrPhone(ctx context.Context, filterOpts models.PatientQueryOptions, opts models.GenericQueryOptions) ([]*models.Patient, *models.PaginationResponse, error) {
	filters := bson.D{}
	if filterOpts.Keyword != "" {
		regexPattern := bson.D{{Key: "$regex", Value: filterOpts.Keyword}, {Key: "$options", Value: "i"}}
		filters = append(filters, bson.E{
			Key: "$or",
			Value: bson.A{
				bson.D{{Key: "name", Value: regexPattern}},
				bson.D{{Key: "phone", Value: regexPattern}},
				bson.D{{Key: "address", Value: regexPattern}},
			},
		})
	}
	col := m.getCollection("patients")
	return MongoList[models.Patient](ctx, col, filters, opts)
}
