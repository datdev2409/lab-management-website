package storage

import (
	"context"
	"regexp"

	"github.com/datdev2409/lab-admin-go/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func (m *MongoStorage) InsertCombo(ctx context.Context, combo *models.Combo) (string, error) {
	col := m.getCollection("combos")
	return MongoInsert(ctx, col, combo)
}

func (m *MongoStorage) ListCombos(ctx context.Context, filterOpts models.ComboQueryOptions, opts models.GenericQueryOptions) ([]*models.Combo, *models.PaginationResponse, error) {
	filters := bson.D{}
	if filterOpts.Keyword != "" {
		filters = append(filters, bson.E{Key: "name", Value: bson.D{{Key: "$regex", Value: regexp.QuoteMeta(filterOpts.Keyword)}, {Key: "$options", Value: "i"}}})
	}
	col := m.getCollection("combos")
	return MongoList[models.Combo](ctx, col, filters, opts)
}

func (m *MongoStorage) GetComboById(ctx context.Context, id string) (*models.Combo, error) {
	col := m.getCollection("combos")
	return MongoGetById[models.Combo](ctx, col, id)
}

func (m *MongoStorage) UpdateComboById(ctx context.Context, id string, update map[string]interface{}) error {
	col := m.getCollection("combos")
	return MongoUpdateById[models.Combo](ctx, col, id, bson.M{"$set": update})
}

func (m *MongoStorage) DeleteComboById(ctx context.Context, id string) error {
	col := m.getCollection("combos")
	return MongoDeleteById[models.Combo](ctx, col, id)
}

func (m *MongoStorage) GetTestsInCombo(ctx context.Context, comboId string) (*models.Combo, []*models.Test, error) {
	combo, err := m.GetComboById(ctx, comboId)
	if err != nil {
		return nil, nil, err
	}
	col := m.getCollection("tests")
	tests, err := MongoGetByIds[models.Test](ctx, col, combo.TestIDs)
	if err != nil {
		return combo, nil, err
	}
	return combo, tests, nil
}
