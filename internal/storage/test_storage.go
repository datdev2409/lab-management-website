package storage

import (
	"context"
	"strconv"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"github.com/datdev2409/lab-admin-go/internal/models"
)

func (t *MongoTestStorage) Insert(test *models.Test) error {
	_, err := t.col.InsertOne(context.Background(), test)
	return err
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
	_, err := t.col.UpdateOne(context.Background(), map[string]string{"id": test.ID}, test)
	return err
}

func (t *MongoTestStorage) Delete(id string) error {
	_, err := t.col.DeleteOne(context.Background(), map[string]string{"id": id})
	return err
}
