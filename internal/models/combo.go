package models

import (
	"time"
)

type Combo struct {
	ID        string    `json:"id" bson:"_id,omitempty" db:"id"`
	Name      string    `json:"name" bson:"name" db:"name"`
	TestIDs   []string  `json:"test_ids" bson:"test_ids" db:"-"` // Not directly mapped, handled by junction table
	CreatedAt time.Time `json:"created_at" bson:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at" db:"updated_at"`
}

type ComboDetailsResponse struct {
	Combo *Combo  `json:"combo"`
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
		TestIDs:   testIDs,
		CreatedAt: now,
		UpdatedAt: now,
	}
}
