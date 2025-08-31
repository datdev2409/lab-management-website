package models

import (
	"time"
)

type Test struct {
	ID          string    `json:"id" bson:"_id,omitempty"`
	Name        string    `json:"name" bson:"name"`
	Price       int       `json:"price" bson:"price"`
	NormalValue string    `json:"normal_value" bson:"normal_value"`
	Unit        string    `json:"unit" bson:"unit"`
	LowerBound  float64   `json:"lower_bound" bson:"lower_bound"`
	UpperBound  float64   `json:"upper_bound" bson:"upper_bound"`
	CreatedAt   time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" bson:"updated_at"`
}

type TestInfo struct {
	Name        string
	NormalValue string
	Unit        string
}

type TestQueryOptions struct {
	Keyword string `json:"keyword"`
}

func NewTest(name string, price int, normalValue, unit string, lowerBound, upperBound float64) *Test {
	now := time.Now()
	return &Test{
		ID:          GenerateRandomID("test_"),
		Name:        name,
		Price:       price,
		NormalValue: normalValue,
		Unit:        unit,
		LowerBound:  lowerBound,
		UpperBound:  upperBound,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}
