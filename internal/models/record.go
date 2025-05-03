package models

import (
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type TestResult struct {
	ID          bson.ObjectID `json:"id" bson:"_id"`
	Name        string        `json:"name" bson:"name"`
	Price       int           `json:"price" bson:"price"`
	NormalValue string        `json:"normal_value" bson:"normal_value"`
	Unit        string        `json:"unit" bson:"unit"`
	LowerBound  float64       `json:"lower_bound" bson:"lower_bound"`
	UpperBound  float64       `json:"upper_bound" bson:"upper_bound"`
	Result      string        `json:"result" bson:"result"`
	ResultText  string        `json:"result_text" bson:"result_text"`
}

type Record struct {
	ID          bson.ObjectID `json:"id" bson:"_id,omitempty"`
	ComboName   string        `json:"combo_name" bson:"combo_name"`
	Patient     Patient       `json:"patient" bson:"patient"`
	TestResults []TestResult  `json:"test_results" bson:"test_results"`
	Status      string        `json:"status" bson:"status"`
	CreatedAt   time.Time     `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at" bson:"updated_at"`
}

type TestResultRequest struct {
	ID          string  `json:"id"`
	Name        string  `json:"name" bson:"name"`
	Price       int     `json:"price" bson:"price"`
	NormalValue string  `json:"normal_value" bson:"normal_value"`
	Unit        string  `json:"unit" bson:"unit"`
	LowerBound  float64 `json:"lower_bound" bson:"lower_bound"`
	UpperBound  float64 `json:"upper_bound" bson:"upper_bound"`
	Result      string  `json:"result"`
	ResultText  string  `json:"result_text"`
}

type CreateRecordResponse struct {
	ID string `json:"id"`
}

type CreateRecordRequest struct {
	PatientID   string              `json:"patient_id"`
	ComboName   string              `json:"combo_name"`
	TestResults []TestResultRequest `json:"test_results"`
}

type UpdateRecordRequest struct {
	TestResults []TestResultRequest `json:"test_results"`
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

type PaginationResponse struct {
	Total     int `json:"total"`
	TotalPage int `json:"total_page"`
	Page      int `json:"page"`
	PageSize  int `json:"page_size"`
}

// Convert frontend JSON data to TestResult
func ConvertFrontendTestResult(data map[string]interface{}) (TestResult, error) {
	// Convert string ID to bson.ObjectID
	idStr, ok := data["id"].(string)
	if !ok {
		return TestResult{}, fmt.Errorf("invalid id format")
	}

	id, err := bson.ObjectIDFromHex(idStr)
	if err != nil {
		return TestResult{}, fmt.Errorf("invalid id format: %v", err)
	}

	// Convert other fields
	name, _ := data["name"].(string)
	price, _ := data["price"].(float64) // JSON numbers are float64 by default
	normalValue, _ := data["normal_value"].(string)
	unit, _ := data["unit"].(string)
	lowerBound, _ := data["lower_bound"].(float64)
	upperBound, _ := data["upper_bound"].(float64)
	result, _ := data["result"].(string)
	resultText, _ := data["result_text"].(string)

	return TestResult{
		ID:          id,
		Name:        name,
		Price:       int(price),
		NormalValue: normalValue,
		Unit:        unit,
		LowerBound:  lowerBound,
		UpperBound:  upperBound,
		Result:      result,
		ResultText:  resultText,
	}, nil
}

// func (r *Record) MarshalBSON() ([]byte, error) {
// 	if r.TestResults == nil {
// 		r.TestResults = []TestResult{}
// 	}

// 	type my Record
// 	return bson.Marshal((*my)(r))
// }

// Update TestId from string to bson.ObjectID when unmarshal
