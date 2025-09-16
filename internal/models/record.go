package models

import (
	"time"
)

type TestResult struct {
	ID                     string  `json:"id" bson:"_id"`
	Name                   string  `json:"name" bson:"name"`
	Price                  int     `json:"price" bson:"price"`
	NormalValue            string  `json:"normal_value" bson:"normal_value"`
	Unit                   string  `json:"unit" bson:"unit"`
	LowerBound             float64 `json:"lower_bound" bson:"lower_bound"`
	UpperBound             float64 `json:"upper_bound" bson:"upper_bound"`
	Result                 string  `json:"result" bson:"result"`
	ResultText             string  `json:"result_text" bson:"result_text"`
	Abnormal               bool    `json:"abnormal" bson:"abnormal"`
	ManualAbnormalOverride bool    `json:"manual_abnormal_override" bson:"manual_abnormal_override"`
}

type Record struct {
	ID          string       `json:"id" bson:"_id,omitempty"`
	ComboName   string       `json:"combo_name" bson:"combo_name"`
	Patient     Patient      `json:"patient" bson:"patient"`
	TestResults []TestResult `json:"test_results" bson:"test_results"`
	Status      string       `json:"status" bson:"status"`
	CreatedAt   time.Time    `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at" bson:"updated_at"`
}

func NewRecord(patient Patient, comboName string, testResults []TestResult) Record {
	now := time.Now()
	record := Record{
		ID:          GenerateRandomID("record_"),
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
	ID                     string  `json:"id"`
	Name                   string  `json:"name" bson:"name"`
	Price                  int     `json:"price" bson:"price"`
	NormalValue            string  `json:"normal_value" bson:"normal_value"`
	Unit                   string  `json:"unit" bson:"unit"`
	LowerBound             float64 `json:"lower_bound" bson:"lower_bound"`
	UpperBound             float64 `json:"upper_bound" bson:"upper_bound"`
	Result                 string  `json:"result"`
	ResultText             string  `json:"result_text"`
	Abnormal               bool    `json:"abnormal"`
	ManualAbnormalOverride bool    `json:"manual_abnormal_override"`
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
	ComboName   string              `json:"combo_name" bson:"combo_name"`
	PatientID   string              `json:"patient_id" bson:"patient_id"`
	Patient     *Patient            `json:"patient,omitempty" bson:"patient,omitempty"`
	TestResults []TestResultRequest `json:"test_results" bson:"test_results"`
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
	BillingReport           ReportType = "phieu_thu"
	ResultsReport           ReportType = "phieu_ket_qua"
	ResultsWithSignature    ReportType = "phieu_ket_qua_chu_ky"
	ResultsWithSignaturePDF ReportType = "phieu_ket_qua_chu_ky_pdf"
	TrackingReport          ReportType = "phieu_theo_doi"
)

// RecordWithTotal represents a record with its calculated total price
type RecordWithTotal struct {
	*Record
	TotalPrice int `json:"total_price"`
}

// ReportSummary represents aggregated data for revenue reports
type ReportSummary struct {
	TotalRecords int        `json:"total_records"`
	TotalRevenue int        `json:"total_revenue"`
	StartDate    *time.Time `json:"start_date"`
	EndDate      *time.Time `json:"end_date"`
}

// ReportResponse represents the complete response for revenue reports
type ReportResponse struct {
	Records    []*RecordWithTotal  `json:"records"`
	Pagination *PaginationResponse `json:"pagination"`
	Summary    *ReportSummary      `json:"summary"`
}
