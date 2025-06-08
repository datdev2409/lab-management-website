package storage

import (
	"context"
	"log"

	"github.com/datdev2409/lab-admin-go/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// func (m *MongoRecordStorage) GetById(ctx context.Context, id string) (*models.Record, error) {
// 	recordId, err := bson.ObjectIDFromHex(id)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var record models.Record
// 	err = m.col.FindOne(ctx, bson.D{{Key: "_id", Value: recordId}}).Decode(&record)
// 	return &record, err
// }

// func (m *MongoRecordStorage) Insert(ctx context.Context, record *models.Record) (string, error) {
// 	result, err := m.col.InsertOne(ctx, record)
// 	if err != nil {
// 		return "", err
// 	}
// 	return result.InsertedID.(bson.ObjectID).Hex(), err
// }

func (m *MongoRecordStorage) ListRecords(ctx context.Context, filterOpts models.RecordQueryOptions, opts models.GenericQueryOptions) ([]*models.Record, *models.PaginationResponse, error) {
	filters := bson.D{}
	if filterOpts.Keyword != "" {
		filters = append(filters, bson.E{
			Key: "$or",
			Value: bson.A{
				bson.D{{Key: "patient.name", Value: bson.D{{Key: "$regex", Value: filterOpts.Keyword}, {Key: "$options", Value: "i"}}}},
				bson.D{{Key: "patient.phone", Value: bson.D{{Key: "$regex", Value: filterOpts.Keyword}, {Key: "$options", Value: "i"}}}},
			},
		})
	}

	if filterOpts.PatientID != "" {
		patientOID, err := bson.ObjectIDFromHex(filterOpts.PatientID)
		if err != nil {
			log.Println("Error while converting patient ID to ObjectID:", err)
			return nil, nil, err
		}
		filters = append(filters, bson.E{Key: "patient._id", Value: patientOID})
	}

	if filterOpts.Status != "" {
		filters = append(filters, bson.E{Key: "status", Value: filterOpts.Status})
	}

	return m.List(ctx, filters, opts)

}

// func (m *MongoRecordStorage) UpdateTestResults(ctx context.Context, recordId string, testResults []models.TestResultRequest) error {
// 	recordOId, err := bson.ObjectIDFromHex(recordId)
// 	if err != nil {
// 		log.Println("Error while converting record id to object id", err)
// 		return err
// 	}

// 	updatedTestResults := []models.TestResult{}
// 	for _, testResult := range testResults {
// 		testId, err := bson.ObjectIDFromHex(testResult.ID)
// 		if err != nil {
// 			log.Println("Error while converting test id to object idd", err)
// 			return err
// 		}
// 		updatedTestResults = append(updatedTestResults, models.TestResult{
// 			ID:          testId,
// 			Name:        testResult.Name,
// 			Price:       testResult.Price,
// 			NormalValue: testResult.NormalValue,
// 			Unit:        testResult.Unit,
// 			LowerBound:  testResult.LowerBound,
// 			UpperBound:  testResult.UpperBound,
// 			Result:      testResult.Result,
// 			ResultText:  testResult.ResultText,
// 		})
// 	}

// 	update := bson.D{{Key: "$set", Value: bson.D{{Key: "test_results", Value: updatedTestResults}}}}
// 	_, err = m.col.UpdateOne(ctx, bson.D{{Key: "_id", Value: recordOId}}, update)
// 	return err
// }
