package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Test struct {
	ID          bson.ObjectID `json:"id" bson:"_id,omitempty"`
	Name        string        `json:"name" bson:"name"`
	Price       int           `json:"price" bson:"price"`
	NormalValue string        `json:"normal_value" bson:"normal_value"`
	Unit        string        `json:"unit" bson:"unit"`
	LowerBound  float64       `json:"lower_bound" bson:"lower_bound"`
	UpperBound  float64       `json:"upper_bound" bson:"upper_bound"`
	CreatedAt   time.Time     `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at" bson:"updated_at"`
}

type TestInfo struct {
	Name        string
	NormalValue string
	Unit        string
}

type TestQueryOptions struct {
	Keyword string `json:"keyword"`
}
