package storage

import (
	"context"

	"github.com/datdev2409/lab-admin-go/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func (m *MongoTrackingStorage) ListTrackings(ctx context.Context, filterOpts models.TrackingQueryOptions, opts models.GenericQueryOptions) ([]*models.Tracking, *models.PaginationResponse, error) {
	filters := bson.D{}
	if filterOpts.Keyword != "" {
		filters = append(filters, bson.E{Key: "name", Value: bson.D{{Key: "$regex", Value: filterOpts.Keyword}, {Key: "$options", Value: "i"}}})
	}

	return m.List(ctx, filters, opts)
}
