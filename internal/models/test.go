package models

import (
	"time"
)

type Test struct {
	ID            string    `json:"id" bson:"_id,omitempty"`
	Name          string    `json:"name" bson:"name"`
	Price         int       `json:"price" bson:"price"`
	ImportedPrice int       `json:"imported_price" bson:"imported_price"`
	NormalValue   string    `json:"normal_value" bson:"normal_value"`
	Unit          string    `json:"unit" bson:"unit"`
	LowerBound    *float64  `json:"lower_bound,omitempty" bson:"lower_bound"`
	UpperBound    *float64  `json:"upper_bound,omitempty" bson:"upper_bound"`
	CreatedAt     time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" bson:"updated_at"`
}

type CreateTestInput struct {
	Name          string   `json:"name" validate:"required"`
	Price         int      `json:"price"`
	ImportedPrice int      `json:"imported_price"`
	NormalValue   string   `json:"normal_value"`
	Unit          string   `json:"unit"`
	LowerBound    *float64 `json:"lower_bound"`
	UpperBound    *float64 `json:"upper_bound"`
}

type TestUpdate struct {
	Name          *string  `json:"name,omitempty"`
	Price         *int     `json:"price,omitempty"`
	ImportedPrice *int     `json:"imported_price,omitempty"`
	NormalValue   *string  `json:"normal_value,omitempty"`
	Unit          *string  `json:"unit,omitempty"`
	LowerBound    *float64 `json:"lower_bound,omitempty"`
	UpperBound    *float64 `json:"upper_bound,omitempty"`
}

type BulkCreateTestResult struct {
	Success int                   `json:"success"`
	Failure int                   `json:"failure"`
	Errors  []BulkCreateTestError `json:"errors,omitempty"`
}

type BulkCreateTestError struct {
	Index   int    `json:"index"`
	Name    string `json:"name,omitempty"`
	Message string `json:"message"`
}

type CreateTestRequest struct {
	Name        string  `json:"name"`
	Price       int     `json:"price"`
	NormalValue string  `json:"normal_value"`
	Unit        string  `json:"unit"`
	LowerBound  float64 `json:"lower_bound"`
	UpperBound  float64 `json:"upper_bound"`
}

type TestInfo struct {
	Name        string
	NormalValue string
	Unit        string
	Order       int
}

type TestQueryOptions struct {
	Keyword string `json:"keyword"`
}

func NewTest(name string, price int, normalValue, unit string, lowerBound, upperBound float64) *Test {
	now := time.Now()
	lb := lowerBound
	ub := upperBound
	return &Test{
		ID:          GenerateRandomID("test_"),
		Name:        name,
		Price:       price,
		NormalValue: normalValue,
		Unit:        unit,
		LowerBound:  &lb,
		UpperBound:  &ub,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

func (t Test) GetID() string {
	return t.ID
}
