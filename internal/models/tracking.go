package models

import (
	"time"
)

type TrackingTestRequest struct {
	TestID      string `json:"test_id" bson:"test_id"`
	TestName    string `json:"test_name" bson:"test_name"`
	NormalValue string `json:"normal_value" bson:"normal_value"`
	Unit        string `json:"unit" bson:"unit"`
	Order       int    `json:"order" bson:"order"`
}

type TrackingTestData struct {
	TestID      string `json:"test_id" bson:"test_id"`
	TestName    string `json:"test_name" bson:"test_name"`
	NormalValue string `json:"normal_value" bson:"normal_value"`
	Unit        string `json:"unit" bson:"unit"`
	Order       int    `json:"order" bson:"order"`
}

type Tracking struct {
	ID        string             `json:"id" bson:"_id,omitempty" db:"id"`
	Name      string             `json:"name" bson:"name" db:"name"`
	Tests     []TrackingTestData `json:"tests" bson:"tests" db:"-"` // Not directly mapped, handled by junction table
	CreatedAt time.Time          `json:"created_at" bson:"created_at" db:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at" db:"updated_at"`
}

type CreateTrackingRequest struct {
	Name  string                `json:"tracking_name"`
	Tests []TrackingTestRequest `json:"tests"`
}

type TrackingQueryOptions struct {
	Keyword string `json:"keyword,omitempty" bson:"keyword,omitempty"`
}

func NewTracking(name string, testRequests []TrackingTestRequest) Tracking {
	var tests []TrackingTestData
	for _, test := range testRequests {
		tests = append(tests, TrackingTestData(test))
	}

	now := time.Now()
	return Tracking{
		ID:        GenerateRandomID("tracking_"), // Use tracking name for ID
		Name:      name,
		Tests:     tests,
		CreatedAt: now,
		UpdatedAt: now,
	}
}
