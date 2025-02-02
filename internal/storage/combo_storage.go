package storage

import (
	"context"
	"strconv"

	"github.com/datdev2409/lab-admin-go/internal/models"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func (m MongoComboStorage) Insert(combo *models.Combo) error {
	_, err := m.col.InsertOne(context.Background(), combo)
	return err
}

func (m MongoComboStorage) SearchByKeyword(ctx context.Context, keyword string, opts map[string]string) (*[]models.Combo, error) {
	combos := []models.Combo{}
	filter := BuildMongoFilter(map[string]FilterCondition{
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

	cursor, err := m.col.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	err = cursor.All(ctx, &combos)
	if err != nil {
		return &combos, err
	}

	return &combos, nil
}

func (m MongoComboStorage) GetById(ctx context.Context, id string) (*models.Combo, error) {
	var combo models.Combo
	err := m.col.FindOne(ctx, map[string]string{"_id": id}).Decode(&combo)
	return &combo, err
}
