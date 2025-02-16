package storage

import (
	"context"
	"log"
	"strconv"

	"github.com/datdev2409/lab-admin-go/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
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

func (m *MongoRecordStorage) UpdatePatient(ctx context.Context, recordId string, patient models.Patient) error {
	filter := bson.D{{Key: "_id", Value: recordId}}

	patientDoc := bson.D{{Key: "id", Value: patient.ID}, {Key: "name", Value: patient.Name}, {Key: "phone", Value: patient.Phone}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "patient", Value: patientDoc}}}}
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

func (m *MongoRecordStorage) AddTest(ctx context.Context, recordId string, testId string) error {
	testResult := models.TestResult{TestID: testId, Result: "", ResultText: ""}
	update := bson.D{{Key: "$push", Value: bson.D{{Key: "test_results", Value: testResult}}}}
	_, err := m.col.UpdateByID(ctx, recordId, update)
	return err
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

func (m *MongoRecordStorage) GetDetails(ctx context.Context, id string) (*models.RecordWithDetails, error) {
	recordFilterStage := bson.D{{Key: "$match", Value: bson.D{{Key: "_id", Value: id}}}}
	patientLookupStage := bson.D{{Key: "$lookup", Value: bson.D{
		{Key: "from", Value: "patients"},
		{Key: "localField", Value: "patient.id"},
		{Key: "foreignField", Value: "_id"},
		{Key: "as", Value: "patient"},
	}}}
	patientUnwindStage := bson.D{{Key: "$unwind", Value: "$patient"}}

	testsLookupStage := bson.D{{Key: "$lookup", Value: bson.D{
		{Key: "from", Value: "tests"},
		{Key: "localField", Value: "test_results.test_id"},
		{Key: "foreignField", Value: "_id"},
		{Key: "as", Value: "tests"},
	}}}
	testsArrayToObjectStage := bson.D{{Key: "$addFields", Value: bson.D{
		{Key: "test_info_map", Value: bson.D{
			{Key: "$arrayToObject", Value: bson.D{
				{Key: "$map", Value: bson.D{
					{Key: "input", Value: "$tests"},
					{Key: "as", Value: "test"},
					{Key: "in", Value: bson.D{
						{Key: "k", Value: "$$test._id"},
						{Key: "v", Value: "$$test"},
					}},
				}},
			}},
		}},
	}}}

	lookupPipeline := bson.A{
		recordFilterStage,
		patientLookupStage,
		patientUnwindStage,
		testsLookupStage,
		testsArrayToObjectStage,
	}

	cursor, err := m.col.Aggregate(ctx, lookupPipeline)
	if err != nil {
		log.Println("Error while aggregating record with details", err)
		return nil, err
	}

	var result models.RecordWithDetails
	if !cursor.Next(ctx) {
		log.Println("Record not found")
		return nil, mongo.ErrNoDocuments
	}

	if err = cursor.Decode(&result); err != nil {
		log.Println("Error while decoding record with details", err)
		return nil, err
	}

	return &result, nil
}
