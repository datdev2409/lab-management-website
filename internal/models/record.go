package models

import (
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

func NewRecord(patient Patient, comboName string, testResults []TestResult) Record {
	now := time.Now()
	record := Record{
		Patient:     patient,
		ComboName:   comboName,
		TestResults: testResults,
		Status:      "pending",
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	return record
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

type ReportType string

const (
	BillingReport        ReportType = "phieu_thu"
	ResultsReport        ReportType = "phieu_ket_qua"
	ResultsWithSignature ReportType = "phieu_ket_qua_chu_ky"
	TrackingReport       ReportType = "phieu_theo_doi"
)
