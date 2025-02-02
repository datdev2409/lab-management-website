package storage

import (
	"context"

	"github.com/datdev2409/lab-admin-go/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func (m *MongoRecordStorage) UpdateCombo(ctx context.Context, recordId string, combo *models.Combo) error {
	filter := bson.D{{Key: "_id", Value: recordId}}

	update := bson.A{}
	updateComboName := bson.D{{Key: "$set", Value: bson.D{{Key: "combo_name", Value: combo.Name}}}}

	tests := bson.A{}
	for _, testId := range combo.Tests {
		testDocument := bson.D{{Key: "test_id", Value: testId}, {Key: "result", Value: ""}, {Key: "result_text", Value: ""}}
		tests = append(tests, testDocument)
	}

	updateTest := bson.D{{Key: "$set", Value: bson.D{{Key: "test_results", Value: tests}}}}
	update = append(update, updateComboName, updateTest)

	_, err := m.col.UpdateOne(context.TODO(), filter, update)
	return err
}

func (m *MongoRecordStorage) UpdatePatient(ctx context.Context, recordId string, patientId string) error {
	filter := bson.D{{Key: "_id", Value: recordId}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "patient_id", Value: patientId}}}}
	_, err := m.col.UpdateOne(context.TODO(), filter, update)

	if err != nil {
		return err
	}

	return nil
}

func (m *MongoRecordStorage) GetById(ctx context.Context, id string) (*models.Record, error) {
	var record models.Record
	err := m.col.FindOne(ctx, map[string]string{"_id": id}).Decode(&record)
	return &record, err
}

func (m *MongoRecordStorage) Insert(ctx context.Context, record *models.Record) error {
	_, err := m.col.InsertOne(ctx, record)
	return err
}
