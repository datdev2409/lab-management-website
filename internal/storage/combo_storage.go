package storage

import (
	"context"
	"math"
	"strconv"

	"github.com/datdev2409/lab-admin-go/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func (m MongoComboStorage) Insert(combo *models.Combo) error {
	_, err := m.col.InsertOne(context.Background(), combo)
	return err
}

func (m MongoComboStorage) ListCombos(ctx context.Context, filterOpts models.ComboQueryOptions, opts models.GenericQueryOptions) ([]*models.Combo, *models.PaginationResponse, error) {
	combos := []*models.Combo{}

	filters := bson.D{}
	if filterOpts.Keyword != "" {
		filters = append(filters, bson.E{Key: "name", Value: bson.D{{Key: "$regex", Value: filterOpts.Keyword}}})
	}

	findOpts := BuildMongoSortAndPaginationOptions(opts)

	cursor, err := m.col.Find(ctx, filters, findOpts)
	if err != nil {
		return nil, nil, err
	}

	count, err := m.col.CountDocuments(ctx, filters)
	if err != nil {
		return nil, nil, err
	}

	pagniation := &models.PaginationResponse{
		Total:     int(count),
		Page:      opts.Page,
		PageSize:  opts.PageSize,
		TotalPage: int(math.Ceil(float64(count) / float64(opts.PageSize))),
	}

	if err = cursor.All(ctx, &combos); err != nil {
		return nil, nil, err
	}

	return combos, pagniation, nil
}

func (m MongoComboStorage) SearchByKeyword(ctx context.Context, keyword string, opts map[string]string) ([]*models.Combo, error) {
	combos := []*models.Combo{}
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
		return nil, err
	}

	return combos, nil
}

func (m MongoComboStorage) GetById(ctx context.Context, id string) (*models.Combo, error) {
	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var combo models.Combo
	err = m.col.FindOne(ctx, bson.D{{Key: "_id", Value: oid}}).Decode(&combo)
	return &combo, err
}

func (m MongoComboStorage) GetTestsInCombo(ctx context.Context, comboId string) (*models.Combo, []*models.Test, error) {
	oid, err := bson.ObjectIDFromHex(comboId)
	if err != nil {
		return nil, nil, err
	}

	var combo *models.Combo
	err = m.col.FindOne(ctx, bson.D{{Key: "_id", Value: oid}}).Decode(&combo)
	if err != nil {
		return nil, nil, err
	}

	cursor, err := m.db.Collection("tests").Find(ctx, bson.D{{Key: "_id", Value: bson.D{{Key: "$in", Value: combo.TestIDs}}}})
	if err != nil {
		return nil, nil, err
	}

	tests := []*models.Test{}
	if err = cursor.All(ctx, &tests); err != nil {
		return combo, nil, err
	}

	return combo, tests, nil
}
