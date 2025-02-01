package models

type Combo struct {
	ID    string   `bson:"_id,omitempty"`
	Name  string   `json:"name" bson:"name"`
	Tests []string `json:"tests" bson:"tests"`
}
