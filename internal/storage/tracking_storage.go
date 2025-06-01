package storage

import (
	"context"
	"log"
	"math"

	"github.com/datdev2409/lab-admin-go/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func (m *MongoTrackingStorage) Insert(ctx context.Context, tracking *models.Tracking) (string, error) {
	result, err := m.col.InsertOne(ctx, tracking)
	if err != nil {
		log.Printf("Error inserting tracking: %v", err)
		return "", err
	}

	return result.InsertedID.(bson.ObjectID).Hex(), nil
}

func (m *MongoTrackingStorage) GetById(ctx context.Context, id string) (*models.Tracking, error) {
	var tracking models.Tracking
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		log.Printf("Error converting ID to ObjectID: %v", err)
		return nil, err
	}

	err = m.col.FindOne(ctx, bson.M{"_id": objectID}).Decode(&tracking)
	if err != nil {
		log.Printf("Error finding tracking by ID: %v", err)
		return nil, err
	}

	return &tracking, nil
}

func (m *MongoTrackingStorage) ListTrackings(ctx context.Context, filterOpts models.TrackingQueryOptions, opts models.GenericQueryOptions) ([]*models.Tracking, *models.PaginationResponse, error) {
	filters := bson.D{}
	if filterOpts.Keyword != "" {
		filters = append(filters, bson.E{Key: "name", Value: bson.D{{Key: "$regex", Value: filterOpts.Keyword}, {Key: "$options", Value: "i"}}})
	}

	findOpts := BuildMongoSortAndPaginationOptions(opts)

	cursor, err := m.col.Find(ctx, filters, findOpts)
	if err != nil {
		return nil, nil, err
	}
	defer cursor.Close(ctx)

	count, err := m.col.CountDocuments(ctx, filters)
	if err != nil {
		return nil, nil, err
	}

	pagination := &models.PaginationResponse{
		Total:     int(count),
		Page:      opts.Page,
		PageSize:  opts.PageSize,
		TotalPage: int(math.Ceil(float64(count) / float64(opts.PageSize))),
	}

	var trackings []*models.Tracking
	if err := cursor.All(ctx, &trackings); err != nil {
		return nil, nil, err
	}

	return trackings, pagination, nil
}
