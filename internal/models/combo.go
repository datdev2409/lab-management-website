package models

type Combo struct {
	ID      string   `bson:"_id,omitempty"`
	Name    string   `json:"name" bson:"name"`
	TestIDs []string `json:"test_ids" bson:"test_ids"`
}
