package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type TestResult struct {
	ID         string `json:"id" bson:"_id,omitempty"`
	Test       Test   `json:"test" bson:"test"`
	Result     string `json:"result" bson:"result"`
	ResultText string `json:"result_text" bson:"result_text"`
}

type TestResultWithDetails struct {
	Test       Test   `json:"test" bson:"test"`
	TestID     string `json:"test_id" bson:"test_id"`
	Result     string `json:"result" bson:"result"`
	ResultText string `json:"result_text" bson:"result_text"`
}

type EmbeddedPatient struct {
	ID    string `json:"id" bson:"id"`
	Name  string `json:"name" bson:"name"`
	Phone string `json:"phone" bson:"phone"`
}

type Record struct {
	ID          string          `json:"id" bson:"_id"`
	ComboName   string          `json:"combo_name" bson:"combo_name"`
	Patient     EmbeddedPatient `json:"patient" bson:"patient"`
	TestResults []TestResult    `json:"test_results" bson:"test_results"`
	Status      string          `json:"status" bson:"status"`
	CreatedAt   time.Time       `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at" bson:"updated_at"`
}

type RecordWithDetails struct {
	ID          string          `json:"id" bson:"_id"`
	ComboName   string          `json:"combo_name" bson:"combo_name"`
	Patient     Patient         `json:"patient" bson:"patient"`
	TestResults []TestResult    `json:"test_results" bson:"test_results"`
	Status      string          `json:"status" bson:"status"`
	TestInfoMap map[string]Test `json:"test_info_map" bson:"test_info_map"`
	CreatedAt   time.Time       `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at" bson:"updated_at"`
}

type RecordQueryOptions struct {
	Keyword   string
	Status    string
	PatientID string
	StartDate *time.Time
	EndDate   *time.Time
}

type GenericQueryOptions struct {
	// Sorting
	SortBy    string
	SortOrder string

	// Pagination
	Page     int
	PageSize int
}

func (r *Record) MarshalBSON() ([]byte, error) {
	if r.TestResults == nil {
		r.TestResults = []TestResult{}
	}

	type my Record
	return bson.Marshal((*my)(r))
}
