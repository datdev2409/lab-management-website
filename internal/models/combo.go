package models

import (
	"time"
)

type Combo struct {
	ID        string    `json:"id" bson:"_id,omitempty"`
	Name      string    `json:"name" bson:"name"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

type CreateComboInput struct {
	Name    string   `json:"name" validate:"required"`
	TestIDs []string `json:"test_ids" validate:"required,min=1"`
}

type ComboUpdate struct {
	Name    *string  `json:"name,omitempty"`
	TestIDs []string `json:"test_ids,omitempty"`
}

type ComboDetailsResponse struct {
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	Tests []*Test `json:"tests"`
}

type ComboQueryOptions struct {
	Keyword string
}

func NewCombo(name string, testIDs []string) *Combo {
	comboId := GenerateRandomID("combo_")
	now := time.Now()
	return &Combo{
		ID:        comboId,
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
	}
}
