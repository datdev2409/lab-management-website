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
	ID    string         `json:"id"`
	Name  string         `json:"name"`
	Tests []*TestInCombo `json:"tests"`
}

// TestInCombo represents a test as returned inside a combo details response.
// It intentionally excludes timestamp fields to keep the payload small.
type TestInCombo struct {
	ID            string   `json:"id"`
	Name          string   `json:"name"`
	Price         int      `json:"price"`
	ImportedPrice int      `json:"imported_price"`
	NormalValue   string   `json:"normal_value"`
	Unit          string   `json:"unit"`
	LowerBound    *float64 `json:"lower_bound"`
	UpperBound    *float64 `json:"upper_bound"`
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
