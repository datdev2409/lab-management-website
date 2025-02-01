package storage

import (
	"context"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"strconv"

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

func (t *MongoTestStorage) GetAll() ([]*models.Test, error) {
	var tests []*models.Test
	cursor, err := t.col.Find(context.Background(), map[string]string{})
	if err != nil {
		return nil, err
	}

	for cursor.Next(context.Background()) {
		var test models.Test
		err := cursor.Decode(&test)
		if err != nil {
			return nil, err
		}

		tests = append(tests, &test)
	}

	return tests, nil
}

func (t *MongoTestStorage) SearchByKeyword(ctx context.Context, keyword string, opts map[string]string) (*[]models.Test, error) {
	tests := []models.Test{}

	filters := bson.D{{}}
	if keyword != "" {
		filters = bson.D{{Key: "name", Value: bson.D{{Key: "$regex", Value: keyword}, {Key: "$options", Value: "i"}}}}
	}

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
		return &tests, err
	}

	if err = cursor.All(ctx, &tests); err != nil {
		return &tests, err
	}

	return &tests, nil
}

func (t *MongoTestStorage) Update(test *models.Test) error {
	_, err := t.col.UpdateOne(context.Background(), map[string]string{"id": test.ID}, test)
	return err
}

func (t *MongoTestStorage) Delete(id string) error {
	_, err := t.col.DeleteOne(context.Background(), map[string]string{"id": id})
	return err
}
