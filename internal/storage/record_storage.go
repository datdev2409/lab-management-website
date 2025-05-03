package storage

import (
	"context"
	"log"
	"math"

	"github.com/datdev2409/lab-admin-go/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// func (m *MongoRecordStorage) UpdateCombo(ctx context.Context, recordId string, combo *models.Combo) error {
// 	filter := bson.D{{Key: "_id", Value: recordId}}

// 	updateComboName := bson.D{{Key: "$set", Value: bson.D{{Key: "combo_name", Value: combo.Name}}}}

// 	_, err := m.col.UpdateOne(context.TODO(), filter, updateComboName)
// 	return err
// }

// func (m *MongoRecordStorage) UpdatePatient(ctx context.Context, recordId string, patient models.Patient) error {
// 	filter := bson.D{{Key: "_id", Value: recordId}}

// 	patientDoc := bson.D{{Key: "id", Value: patient.ID}, {Key: "name", Value: patient.Name}, {Key: "phone", Value: patient.Phone}}
// 	update := bson.D{
// 		{Key: "$set", Value: bson.D{{Key: "patient", Value: patientDoc}}},
// 		{Key: "$set", Value: bson.D{{Key: "updated_at", Value: GetCurrentTime()}}},
// 	}
// 	_, err := m.col.UpdateOne(context.TODO(), filter, update)

// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

func (m *MongoRecordStorage) GetById(ctx context.Context, id string) (*models.Record, error) {
	recordId, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var record models.Record
	err = m.col.FindOne(ctx, bson.D{{Key: "_id", Value: recordId}}).Decode(&record)
	return &record, err
}

func (m *MongoRecordStorage) Insert(ctx context.Context, record *models.Record) (string, error) {
	result, err := m.col.InsertOne(ctx, record)
	if err != nil {
		return "", err
	}
	return result.InsertedID.(bson.ObjectID).Hex(), err
}

// func (m *MongoRecordStorage) AddTest(ctx context.Context, recordId string, test *models.Test) error {
// 	testResult := models.TestResult{Test: *test, Result: "", ResultText: ""}
// 	update := bson.D{
// 		{Key: "$push", Value: bson.D{{Key: "test_results", Value: testResult}}},
// 		{Key: "$set", Value: bson.D{{Key: "updated_at", Value: GetCurrentTime()}}},
// 	}
// 	_, err := m.col.UpdateByID(ctx, recordId, update)
// 	return err
// }

// func (m *MongoRecordStorage) AddTests(ctx context.Context, recordId string, tests []*models.Test) error {
// 	update := bson.D{}
// 	for _, test := range tests {
// 		testResult := models.TestResult{Test: *test, Result: "", ResultText: ""}
// 		update = append(update, bson.E{Key: "$push", Value: bson.A{bson.D{{Key: "test_results", Value: testResult}}}})
// 	}

// 	update = append(update, bson.E{Key: "$set", Value: bson.D{{Key: "updated_at", Value: GetCurrentTime()}}})
// 	_, err := m.col.UpdateByID(ctx, recordId, update)
// 	return err
// }

// // SearchRecords implements RecordStorage.
func (m *MongoRecordStorage) ListRecords(ctx context.Context, filterOpts models.RecordQueryOptions, opts models.GenericQueryOptions) (*[]models.Record, *models.PaginationResponse, error) {
	records := []models.Record{}

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
		filters = append(filters, bson.E{Key: "patient.id", Value: filterOpts.PatientID})
	}

	if filterOpts.Status != "" {
		filters = append(filters, bson.E{Key: "status", Value: filterOpts.Status})
	}

	// TODO: Add date range filters if provided
	findOpts := BuildMongoSortAndPaginationOptions(opts)
	cursor, err := m.col.Find(ctx, filters, findOpts)
	if err != nil {
		return nil, nil, err
	}

	count, err := m.col.CountDocuments(ctx, filters)
	if err != nil {
		return nil, nil, err
	}

	pagination := models.PaginationResponse{
		Total:     int(count),
		TotalPage: int(math.Ceil(float64(count) / float64(opts.PageSize))),
		Page:      opts.Page,
		PageSize:  opts.PageSize,
	}

	if err = cursor.All(ctx, &records); err != nil {
		return nil, nil, err
	}

	return &records, &pagination, nil
}

// func (m *MongoRecordStorage) GetDetails(ctx context.Context, id string) (*models.RecordWithDetails, error) {
// 	recordFilterStage := bson.D{{Key: "$match", Value: bson.D{{Key: "_id", Value: id}}}}
// 	patientLookupStage := bson.D{{Key: "$lookup", Value: bson.D{
// 		{Key: "from", Value: "patients"},
// 		{Key: "localField", Value: "patient.id"},
// 		{Key: "foreignField", Value: "_id"},
// 		{Key: "as", Value: "patient"},
// 	}}}
// 	patientUnwindStage := bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$patient"}, {Key: "preserveNullAndEmptyArrays", Value: true}}}}

// 	testsLookupStage := bson.D{{Key: "$lookup", Value: bson.D{
// 		{Key: "from", Value: "tests"},
// 		{Key: "localField", Value: "test_results.test_id"},
// 		{Key: "foreignField", Value: "_id"},
// 		{Key: "as", Value: "tests"},
// 	}}}
// 	testsArrayToObjectStage := bson.D{{Key: "$addFields", Value: bson.D{
// 		{Key: "test_info_map", Value: bson.D{
// 			{Key: "$arrayToObject", Value: bson.D{
// 				{Key: "$map", Value: bson.D{
// 					{Key: "input", Value: "$tests"},
// 					{Key: "as", Value: "test"},
// 					{Key: "in", Value: bson.D{
// 						{Key: "k", Value: "$$test._id"},
// 						{Key: "v", Value: "$$test"},
// 					}},
// 				}},
// 			}},
// 		}},
// 	}}}

// 	lookupPipeline := bson.A{
// 		recordFilterStage,
// 		patientLookupStage,
// 		patientUnwindStage,
// 		testsLookupStage,
// 		testsArrayToObjectStage,
// 	}

// 	cursor, err := m.col.Aggregate(ctx, lookupPipeline)
// 	if err != nil {
// 		log.Println("Error while aggregating record with details", err)
// 		return nil, err
// 	}

// 	var result models.RecordWithDetails
// 	if !cursor.Next(ctx) {
// 		log.Println("Record not found")
// 		return nil, mongo.ErrNoDocuments
// 	}

// 	if err = cursor.Decode(&result); err != nil {
// 		log.Println("Error while decoding record with details", err)
// 		return nil, err
// 	}

// 	return &result, nil
// }

func (m *MongoRecordStorage) UpdateTestResults(ctx context.Context, recordId string, testResults []models.TestResultRequest) error {
	recordOId, err := bson.ObjectIDFromHex(recordId)
	if err != nil {
		log.Println("Error while converting record id to object id", err)
		return err
	}

	updatedTestResults := []models.TestResult{}
	for _, testResult := range testResults {
		log.Println("Test result", testResult.ID)
		testId, err := bson.ObjectIDFromHex(testResult.ID)
		if err != nil {
			log.Println("Error while converting test id to object idd", err)
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

	update := bson.D{{Key: "$set", Value: bson.D{{Key: "test_results", Value: updatedTestResults}}}}
	_, err = m.col.UpdateOne(ctx, bson.D{{Key: "_id", Value: recordOId}}, update)
	return err
}

// func (m *MongoRecordStorage) ListByPatientId(ctx context.Context, patientId string) (*[]models.Record, error) {
// 	records := []models.Record{}

// 	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})
// 	cursor, err := m.col.Find(ctx, map[string]string{"patient.id": patientId}, opts)
// 	if err != nil {
// 		return &records, err
// 	}

// 	if err = cursor.All(ctx, &records); err != nil {
// 		return &records, err
// 	}

// 	return &records, nil
// }
