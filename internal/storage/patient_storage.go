package storage

import (
	"context"
	"github.com/datdev2409/lab-admin-go/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"strconv"
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

func (m *MongoPatientStorage) UpdateById(ctx context.Context, id string, patient *models.Patient) error {
	_, err := m.col.UpdateOne(context.Background(), map[string]string{"_id": id}, bson.D{{Key: "$set", Value: patient}})
	return err
}

func (m *MongoPatientStorage) Delete(id string) error {
	_, err := m.col.DeleteOne(context.Background(), map[string]string{"_id": id})
	return err
}
