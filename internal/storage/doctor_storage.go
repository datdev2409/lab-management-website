package storage

import (
	"context"

	"github.com/datdev2409/lab-admin-go/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// FindDoctorByNameAndPhone returns a doctor with the given name and phone, or nil if not found
func (m *MongoStorage) FindDoctorByNameAndPhone(ctx context.Context, name, phone string) (*models.Doctor, error) {
	col := m.getCollection("doctors")
	filter := bson.M{"name": name, "phone": phone}
	var result models.Doctor
	err := col.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			return nil, nil
		}
		return nil, err
	}
	return &result, nil
}

func (m *MongoStorage) InsertDoctor(ctx context.Context, doctor *models.Doctor) (string, error) {
	col := m.getCollection("doctors")
	return MongoInsert(ctx, col, doctor)
}

func (m *MongoStorage) GetDoctorById(ctx context.Context, id string) (*models.Doctor, error) {
	col := m.getCollection("doctors")
	return MongoGetById[models.Doctor](ctx, col, id)
}

func (m *MongoStorage) UpdateDoctorById(ctx context.Context, id string, update models.DoctorUpdate) error {
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
	col := m.getCollection("doctors")
	return MongoUpdateById[models.Doctor](ctx, col, id, bson.M{"$set": updateBSON})
}

func (m *MongoStorage) DeleteDoctorById(ctx context.Context, id string) error {
	col := m.getCollection("doctors")
	return MongoDeleteById[models.Doctor](ctx, col, id)
}

func (m *MongoStorage) SearchDoctorByNameOrPhone(ctx context.Context, filterOpts models.DoctorQueryOptions, opts models.GenericQueryOptions) ([]*models.Doctor, *models.PaginationResponse, error) {
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
	col := m.getCollection("doctors")
	return MongoList[models.Doctor](ctx, col, filters, opts)
}