package storage

import (
	"context"
	"github.com/datdev2409/lab-admin-go/internal/models"
)

func (m MongoComboStorage) Insert(combo *models.Combo) error {
	_, err := m.col.InsertOne(context.Background(), combo)
	return err
}
