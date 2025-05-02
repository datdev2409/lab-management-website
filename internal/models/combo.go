package models

import (
	"encoding/json"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Combo struct {
	ID        bson.ObjectID   `bson:"_id,omitempty"`
	Name      string          `json:"name" bson:"name"`
	TestIDs   []bson.ObjectID `json:"test_ids" bson:"test_ids"`
	CreatedAt time.Time       `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time       `json:"updated_at" bson:"updated_at"`
}

type ComboDetailsResponse struct {
	Combo *Combo  `json:"combo"`
	Tests []*Test `json:"tests"`
}

type ComboQueryOptions struct {
	Keyword string
}

func (c *Combo) GetTestIDs() []string {
	ids := []string{}
	for _, id := range c.TestIDs {
		ids = append(ids, id.Hex())
	}
	return ids
}

func (c *Combo) GetID() string {
	return c.ID.Hex()
}

func ConvertIDsToObjectIDs(ids []string) ([]bson.ObjectID, error) {
	objectIDs := []bson.ObjectID{}
	for _, id := range ids {
		oid, err := bson.ObjectIDFromHex(id)
		if err != nil {
			return nil, err
		}
		objectIDs = append(objectIDs, oid)
	}
	return objectIDs, nil
}

// When response to json, convert TestIDs to string
func (c *Combo) MarshalJSON() ([]byte, error) {
	type ComboJSON struct {
		ID      string   `json:"id"`
		Name    string   `json:"name"`
		TestIDs []string `json:"test_ids"`
	}

	return json.Marshal(ComboJSON{
		ID:      c.ID.Hex(),
		Name:    c.Name,
		TestIDs: c.GetTestIDs(),
	})
}
