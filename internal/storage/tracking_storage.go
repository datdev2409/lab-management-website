package storage

import (
	"context"

	"github.com/datdev2409/lab-admin-go/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func (m *MongoStorage) InsertTracking(ctx context.Context, tracking *models.Tracking) (string, error) {
	col := m.getCollection("trackings")
	return MongoInsert(ctx, col, tracking)
}

func (m *MongoStorage) ListTrackings(ctx context.Context, filterOpts models.TrackingQueryOptions, opts models.GenericQueryOptions) ([]*models.Tracking, *models.PaginationResponse, error) {
	filters := bson.D{}
	if filterOpts.Keyword != "" {
		filters = append(filters, bson.E{Key: "name", Value: bson.D{{Key: "$regex", Value: filterOpts.Keyword}, {Key: "$options", Value: "i"}}})
	}
	col := m.getCollection("trackings")
	return MongoList[models.Tracking](ctx, col, filters, opts)
}

func (m *MongoStorage) GetTrackingById(ctx context.Context, id string) (*models.Tracking, error) {
	col := m.getCollection("trackings")
	return MongoGetById[models.Tracking](ctx, col, id)
}

func (m *MongoStorage) DeleteTrackingById(ctx context.Context, id string) error {
	col := m.getCollection("trackings")
	return MongoDeleteById[models.Tracking](ctx, col, id)
}
