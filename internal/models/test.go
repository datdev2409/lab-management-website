package models

import (
	"time"

	"github.com/google/uuid"
)

type Test struct {
	ID          string    `json:"id" bson:"_id,omitempty" db:"id"`
	Name        string    `json:"name" bson:"name" db:"name"`
	Price       int       `json:"price" bson:"price" db:"price"`
	NormalValue string    `json:"normal_value" bson:"normal_value" db:"normal_value"`
	Unit        string    `json:"unit" bson:"unit" db:"unit"`
	LowerBound  float64   `json:"lower_bound" bson:"lower_bound" db:"lower_bound"`
	UpperBound  float64   `json:"upper_bound" bson:"upper_bound" db:"upper_bound"`
	CreatedAt   time.Time `json:"created_at" bson:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" bson:"updated_at" db:"updated_at"`
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
	return &Test{
		ID:          uuid.New().String(),
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

func (t Test) GetID() string {
	return t.ID
}
