package models

import (
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
