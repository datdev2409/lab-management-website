package storage

import (
	"context"
	"log"

	"github.com/datdev2409/lab-admin-go/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type RecordStorage interface {
	Insert(ctx context.Context, record models.Record) error
	ListRecords(ctx context.Context, filters models.RecordQueryOptions, opts models.GenericQueryOptions) ([]*models.Record, *models.PaginationResponse, error)
	GetById(ctx context.Context, id string) (*models.Record, error)
	GetByIds(ctx context.Context, ids []string) ([]*models.Record, error)
	UpdateRecord(ctx context.Context, recordId string, updateRequest models.UpdateRecordRequest) error
}

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

func (m *MongoRecordStorage) UpdateRecord(ctx context.Context, recordId string, updateRequest models.UpdateRecordRequest) error {

	update := bson.D{}
	if updateRequest.ComboName != "" {
		update = append(update, bson.E{Key: "combo_name", Value: updateRequest.ComboName})
	}

	if updateRequest.Patient != nil {
		patientBSON, err := models.ToBSONDocument(updateRequest.Patient)
		if err != nil {
			return err
		}
		update = append(update, bson.E{Key: "patient", Value: patientBSON})
	}

	updatedTestResults := []models.TestResult{}
	for _, testResult := range updateRequest.TestResults {
		testId, err := bson.ObjectIDFromHex(testResult.ID)
		if err != nil {
			log.Println("Error while converting test id to object id", err)
			return err
		}
		updatedTestResults = append(updatedTestResults, models.TestResult{
			ID:          testId,
			Name:        testResult.Name,
			Price:       testResult.Price,
			NormalValue: testResult.NormalValue,
			Unit:        testResult.Unit,
			LowerBound:  testResult.LowerBound,
			UpperBound:  testResult.UpperBound,
			Result:      testResult.Result,
			ResultText:  testResult.ResultText,
		})
	}
	update = append(update, bson.E{Key: "test_results", Value: updatedTestResults})

	return m.UpdateById(ctx, recordId, bson.M{"$set": update})
}
