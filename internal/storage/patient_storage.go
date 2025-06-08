package storage

import (
	"context"

	"github.com/datdev2409/lab-admin-go/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type PatientStorage interface {
	Insert(ctx context.Context, patient models.Patient) error
	SearchByNameOrPhone(ctx context.Context, filterOpts models.PatientQueryOptions, opts models.GenericQueryOptions) ([]*models.Patient, *models.PaginationResponse, error)
	GetById(ctx context.Context, id string) (*models.Patient, error)
	UpdatePatientById(ctx context.Context, id string, update models.PatientUpdate) error
	UpdateById(ctx context.Context, id string, update interface{}) error
	DeleteById(ctx context.Context, id string) error
}

func (m *MongoPatientStorage) SearchByNameOrPhone(ctx context.Context, filterOpts models.PatientQueryOptions, opts models.GenericQueryOptions) ([]*models.Patient, *models.PaginationResponse, error) {
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
	return m.List(ctx, filters, opts)
}

func (m *MongoPatientStorage) UpdatePatientById(ctx context.Context, id string, update models.PatientUpdate) error {
	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	updateBSON := bson.D{}
	if update.Name != nil {
		updateBSON = append(updateBSON, bson.E{Key: "name", Value: *update.Name})
	}
	if update.Phone != nil {
		updateBSON = append(updateBSON, bson.E{Key: "phone", Value: *update.Phone})
	}
	if update.Address != nil {
		updateBSON = append(updateBSON, bson.E{Key: "address", Value: *update.Address})
	}
	if update.YOB != nil {
		updateBSON = append(updateBSON, bson.E{Key: "yob", Value: *update.YOB})
	}
	if update.Gender != nil {
		updateBSON = append(updateBSON, bson.E{Key: "gender", Value: *update.Gender})
	}
	return m.UpdateById(ctx, oid.Hex(), bson.M{"$set": updateBSON})
}
