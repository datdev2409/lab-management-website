package models

import "go.mongodb.org/mongo-driver/v2/bson"

type TrackingTest struct {
	TestID      bson.ObjectID `json:"test_id" bson:"test_id"`
	TestName    string        `json:"test_name" bson:"test_name"`
	NormalValue string        `json:"normal_value" bson:"normal_value"`
	Order       int           `json:"order" bson:"order"`
}

type Tracking struct {
	ID    bson.ObjectID  `json:"id" bson:"_id,omitempty"`
	Name  string         `json:"name" bson:"name"`
	Tests []TrackingTest `json:"tests" bson:"tests"`
}
