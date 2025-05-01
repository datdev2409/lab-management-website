package storage

import (
	"context"
	"math"
	"strconv"

	"github.com/datdev2409/lab-admin-go/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func (m *MongoPatientStorage) Insert(patient *models.Patient) error {
	_, err := m.col.InsertOne(context.Background(), patient)
	return err
}

func (m *MongoPatientStorage) GetById(id string) (*models.Patient, error) {
	var patient models.Patient
	filter := bson.D{{Key: "_id", Value: id}}
	err := m.col.FindOne(context.Background(), filter).Decode(&patient)
	return &patient, err
}

func (m *MongoPatientStorage) SearchByKeyword(ctx context.Context, keyword string, opts map[string]string) (*[]models.Patient, error) {
	var patients []models.Patient

	filters := bson.D{{}}
	if keyword != "" {
		filters = bson.D{{Key: "name", Value: bson.D{{Key: "$regex", Value: keyword}, {Key: "$options", Value: "i"}}}}
	}

	findOpts := options.Find()
	if val, ok := opts["limit"]; ok {
		limit, err := strconv.Atoi(val)
		if err != nil {
			limit = 5
		}
		findOpts.SetLimit(int64(limit))
	}

	cursor, err := m.col.Find(context.Background(), filters, findOpts)
	if err != nil {
		return nil, err
	}

	if err := cursor.All(context.Background(), &patients); err != nil {
		patients = []models.Patient{}
	}

	return &patients, nil
}

func (m *MongoPatientStorage) ListPatients(ctx context.Context, filterOpts models.PatientQueryOptions, opts models.GenericQueryOptions) ([]*models.Patient, *models.PaginationResponse, error) {
	patients := []*models.Patient{}
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

	findOpts := BuildMongoSortAndPaginationOptions(opts)

	cursor, err := m.col.Find(ctx, filters, findOpts)
	if err != nil {
		return nil, nil, err
	}

	if err := cursor.All(ctx, &patients); err != nil {
		return nil, nil, err
	}

	// Get total count for pagination
	total, err := m.col.CountDocuments(ctx, filters)
	if err != nil {
		return nil, nil, err
	}

	totalPage := int(math.Ceil(float64(total) / float64(opts.PageSize)))
	pagination := &models.PaginationResponse{
		Total:     int(total),
		TotalPage: totalPage,
		Page:      opts.Page,
		PageSize:  opts.PageSize,
	}

	return patients, pagination, nil
}

func (m *MongoPatientStorage) UpdateById(ctx context.Context, id string, patient *models.Patient) error {
	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = m.col.UpdateOne(context.Background(), bson.D{{Key: "_id", Value: oid}}, bson.D{{Key: "$set", Value: patient}})
	return err
}

func (m *MongoPatientStorage) Delete(id string) error {
	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = m.col.DeleteOne(context.Background(), bson.D{{Key: "_id", Value: oid}})
	return err
}
