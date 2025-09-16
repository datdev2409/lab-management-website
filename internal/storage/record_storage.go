package storage

import (
	"context"

	"github.com/datdev2409/lab-admin-go/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func (m *MongoStorage) InsertRecord(ctx context.Context, record *models.Record) (string, error) {
	col := m.getCollection("records")
	return MongoInsert(ctx, col, record)
}

func (m *MongoStorage) ListRecords(ctx context.Context, filters models.RecordQueryOptions, opts models.GenericQueryOptions) ([]*models.Record, *models.PaginationResponse, error) {
	mongoFilters := bson.D{}
	if filters.Keyword != "" {
		mongoFilters = append(mongoFilters, bson.E{
			Key: "$or",
			Value: bson.A{
				bson.D{{Key: "patient.name", Value: bson.D{{Key: "$regex", Value: filters.Keyword}, {Key: "$options", Value: "i"}}}},
				bson.D{{Key: "patient.phone", Value: bson.D{{Key: "$regex", Value: filters.Keyword}, {Key: "$options", Value: "i"}}}},
			},
		})
	}
	if filters.PatientID != "" {
		mongoFilters = append(mongoFilters, bson.E{Key: "patient._id", Value: filters.PatientID})
	}
	if filters.Status != "" {
		mongoFilters = append(mongoFilters, bson.E{Key: "status", Value: filters.Status})
	}
	// Add date range filtering
	if filters.StartDate != nil {
		mongoFilters = append(mongoFilters, bson.E{Key: "created_at", Value: bson.D{{Key: "$gte", Value: *filters.StartDate}}})
	}
	if filters.EndDate != nil {
		mongoFilters = append(mongoFilters, bson.E{Key: "created_at", Value: bson.D{{Key: "$lte", Value: *filters.EndDate}}})
	}
	col := m.getCollection("records")
	return MongoList[models.Record](ctx, col, mongoFilters, opts)
}

func (m *MongoStorage) GetRecordById(ctx context.Context, id string) (*models.Record, error) {
	col := m.getCollection("records")
	return MongoGetById[models.Record](ctx, col, id)
}

func (m *MongoStorage) GetRecordsByIds(ctx context.Context, ids []string) ([]*models.Record, error) {
	col := m.getCollection("records")
	return MongoGetByIds[models.Record](ctx, col, ids)
}

func (m *MongoStorage) GetRecordsByPatientId(ctx context.Context, patientId string) ([]*models.Record, error) {
	col := m.getCollection("records")
	filter := bson.M{"patient._id": patientId}
	cursor, err := col.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var records []*models.Record
	if err := cursor.All(ctx, &records); err != nil {
		return nil, err
	}
	return records, nil
}

func (m *MongoStorage) UpdateRecord(ctx context.Context, recordId string, updateRequest models.UpdateRecordRequest) error {
	update := bson.M{}
	if updateRequest.ComboName != "" {
		update["combo_name"] = updateRequest.ComboName
	}
	if updateRequest.Patient != nil {
		patientBSON, err := models.ToBSONDocument(updateRequest.Patient)
		if err != nil {
			return err
		}
		update["patient"] = patientBSON
	}
	updatedTestResults := []models.TestResult{}
	for _, testResult := range updateRequest.TestResults {
		updatedTestResults = append(updatedTestResults, models.TestResult(testResult))
	}
	update["test_results"] = updatedTestResults
	col := m.getCollection("records")
	return MongoUpdateById[models.Record](ctx, col, recordId, bson.M{"$set": update})
}

func (m *MongoStorage) DeleteRecord(ctx context.Context, recordId string) error {
	col := m.getCollection("records")
	return MongoDeleteById[models.Record](ctx, col, recordId)
}

// GetRecordsWithRevenue returns records with calculated total prices for revenue reporting
func (m *MongoStorage) GetRecordsWithRevenue(ctx context.Context, filters models.RecordQueryOptions, opts models.GenericQueryOptions) (*models.ReportResponse, error) {
	// Get all records without pagination for the report
	allRecordsOpts := models.GenericQueryOptions{
		Page:      1,
		PageSize:  0, // 0 means no limit, get all records
		SortBy:    opts.SortBy,
		SortOrder: opts.SortOrder,
	}
	records, _, err := m.ListRecords(ctx, filters, allRecordsOpts)
	if err != nil {
		return nil, err
	}

	// Convert records to RecordWithTotal and calculate totals
	recordsWithTotal := make([]*models.RecordWithTotal, 0, len(records))
	totalRevenue := 0

	for _, record := range records {
		// Calculate total price for this record
		totalPrice := 0
		for _, testResult := range record.TestResults {
			totalPrice += testResult.Price
		}

		recordWithTotal := &models.RecordWithTotal{
			Record:     record,
			TotalPrice: totalPrice,
		}
		recordsWithTotal = append(recordsWithTotal, recordWithTotal)
		totalRevenue += totalPrice
	}

	// Create summary
	summary := &models.ReportSummary{
		TotalRecords: len(records),
		TotalRevenue: totalRevenue,
		StartDate:    filters.StartDate,
		EndDate:      filters.EndDate,
	}

	return &models.ReportResponse{
		Records:    recordsWithTotal,
		Pagination: nil, // No pagination for reports
		Summary:    summary,
	}, nil
}
