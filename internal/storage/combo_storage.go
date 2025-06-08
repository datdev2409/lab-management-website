package storage

import (
	"context"

	"github.com/datdev2409/lab-admin-go/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func (m MongoComboStorage) ListCombos(ctx context.Context, filterOpts models.ComboQueryOptions, opts models.GenericQueryOptions) ([]*models.Combo, *models.PaginationResponse, error) {
	filters := bson.D{}
	if filterOpts.Keyword != "" {
		filters = append(filters, bson.E{Key: "name", Value: bson.D{{Key: "$regex", Value: filterOpts.Keyword}}})
	}

	return m.List(ctx, filters, opts)

}

func (m MongoComboStorage) GetTestsInCombo(ctx context.Context, comboId string) (*models.Combo, []*models.Test, error) {
	combo, err := m.GetById(ctx, comboId)
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
