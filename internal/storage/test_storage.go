package storage

import (
	"context"
	"math"
	"strconv"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"github.com/datdev2409/lab-admin-go/internal/models"
)

func (t *MongoTestStorage) Insert(test *models.Test) error {
	_, err := t.col.InsertOne(context.Background(), test)
	return err
}

func (t *MongoTestStorage) ListTests(ctx context.Context, filterOpts models.TestQueryOptions, opts models.GenericQueryOptions) ([]*models.Test, *models.PaginationResponse, error) {
	tests := []*models.Test{}
	filters := bson.D{}

	if filterOpts.Keyword != "" {
		filters = append(filters, bson.E{Key: "name", Value: bson.D{{Key: "$regex", Value: filterOpts.Keyword}}})
	}

	findOpts := BuildMongoSortAndPaginationOptions(opts)

	cursor, err := t.col.Find(ctx, filters, findOpts)
	if err != nil {
		return nil, nil, err
	}

	count, err := t.col.CountDocuments(ctx, filters)
	if err != nil {
		return nil, nil, err
	}

	pagniation := &models.PaginationResponse{
		Total:     int(count),
		Page:      opts.Page,
		PageSize:  opts.PageSize,
		TotalPage: int(math.Ceil(float64(count) / float64(opts.PageSize))),
	}

	if err = cursor.All(ctx, &tests); err != nil {
		return nil, nil, err
	}

	return tests, pagniation, nil
}

func (t *MongoTestStorage) GetById(id string) (*models.Test, error) {
	var test models.Test
	err := t.col.FindOne(context.Background(), map[string]string{"_id": id}).Decode(&test)
	return &test, err
}

func (t *MongoTestStorage) GetByIds(ctx context.Context, ids []string) ([]*models.Test, error) {
	tests := []*models.Test{}
	filter := bson.D{{Key: "_id", Value: bson.D{{Key: "$in", Value: ids}}}}

	cursor, err := t.col.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	if err = cursor.All(ctx, &tests); err != nil {
		return tests, err
	}

	return tests, nil
}

func (t *MongoTestStorage) SearchByKeyword(ctx context.Context, keyword string, opts map[string]string) ([]*models.Test, error) {
	tests := []*models.Test{}

	filters := BuildMongoFilter(map[string]FilterCondition{
		"name": {
			Operator: "$regex",
			Value:    keyword,
			Option:   "i",
		},
	})

	findOpts := options.Find()
	if val, ok := opts["limit"]; ok {
		limit, err := strconv.Atoi(val)
		if err != nil {
			limit = 5
		}
		findOpts.SetLimit(int64(limit))
	}

	cursor, err := t.col.Find(ctx, filters, findOpts)
	if err != nil {
		return tests, err
	}

	if err = cursor.All(ctx, &tests); err != nil {
		return tests, err
	}

	return tests, nil
}

func (t *MongoTestStorage) Update(test *models.Test) error {
	_, err := t.col.UpdateOne(context.Background(), map[string]string{"_id": test.ID.Hex()}, test)
	return err
}

func (t *MongoTestStorage) Delete(id string) error {
	_, err := t.col.DeleteOne(context.Background(), map[string]string{"id": id})
	return err
}
