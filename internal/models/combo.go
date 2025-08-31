package models

import (
	"encoding/json"
	"time"
)

type Combo struct {
	ID        string    `bson:"_id,omitempty"`
	Name      string    `json:"name" bson:"name"`
	TestIDs   []string  `json:"test_ids" bson:"test_ids"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

type ComboDetailsResponse struct {
	Combo *Combo  `json:"combo"`
	Tests []*Test `json:"tests"`
}

type ComboQueryOptions struct {
	Keyword string
}

func (c *Combo) GetTestIDs() []string {
	return c.TestIDs
}

func (c *Combo) GetID() string {
	return c.ID
}

func ConvertIDsToObjectIDs(ids []string) ([]string, error) {
	return ids, nil // No conversion needed
}

// When response to json, convert TestIDs to string
func (c *Combo) MarshalJSON() ([]byte, error) {
	type ComboJSON struct {
		ID      string   `json:"id"`
		Name    string   `json:"name"`
		TestIDs []string `json:"test_ids"`
	}

	return json.Marshal(ComboJSON{
		ID:      c.ID,
		Name:    c.Name,
		TestIDs: c.GetTestIDs(),
	})
}

func NewCombo(name string, testIDs []string) *Combo {
	comboId := GenerateComboID(name)
	now := time.Now()
	return &Combo{
		ID:        comboId,
		Name:      name,
		TestIDs:   testIDs,
		CreatedAt: now,
		UpdatedAt: now,
	}
}
