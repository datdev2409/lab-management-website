package storage

import (
	"context"

	"github.com/datdev2409/lab-admin-go/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func (m *MongoStorage) InsertTest(ctx context.Context, test *models.Test) (string, error) {
	col := m.getCollection("tests")
	return MongoInsert(ctx, col, test)
}

func (m *MongoStorage) ListTests(ctx context.Context, filterOpts models.TestQueryOptions, opts models.GenericQueryOptions) ([]*models.Test, *models.PaginationResponse, error) {
	filters := bson.D{}
	if filterOpts.Keyword != "" {
		filters = append(filters, bson.E{Key: "name", Value: bson.D{{Key: "$regex", Value: filterOpts.Keyword}, {Key: "$options", Value: "i"}}})
	}
	col := m.getCollection("tests")
	return MongoList[models.Test](ctx, col, filters, opts)
}

func (m *MongoStorage) GetTestById(ctx context.Context, id string) (*models.Test, error) {
	col := m.getCollection("tests")
	return MongoGetById[models.Test](ctx, col, id)
}

func (m *MongoStorage) UpdateTestById(ctx context.Context, id string, update map[string]interface{}) error {
	col := m.getCollection("tests")
	return MongoUpdateById[models.Test](ctx, col, id, bson.M{"$set": update})
}

func (m *MongoStorage) DeleteTestById(ctx context.Context, id string) error {
	col := m.getCollection("tests")
	return MongoDeleteById[models.Test](ctx, col, id)
}
