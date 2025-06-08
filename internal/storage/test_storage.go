package storage

import (
	"context"

	"github.com/datdev2409/lab-admin-go/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func (t *MongoTestStorage) ListTests(ctx context.Context, filterOpts models.TestQueryOptions, opts models.GenericQueryOptions) ([]*models.Test, *models.PaginationResponse, error) {
	filters := bson.D{}

	if filterOpts.Keyword != "" {
		filters = append(filters, bson.E{Key: "name", Value: bson.D{{Key: "$regex", Value: filterOpts.Keyword}}})
	}

	return t.List(ctx, filters, opts)
}
