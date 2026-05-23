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
	DoctorID    string       `json:"doctor_id,omitempty" bson:"doctor_id,omitempty"`
	DoctorName  string       `json:"doctor_name,omitempty" bson:"doctor_name,omitempty"`
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

// NewRecordWithDoctor creates a new record with optional doctor information
func NewRecordWithDoctor(patient Patient, comboName string, testResults []TestResult, doctorID, doctorName string) Record {
	record := NewRecord(patient, comboName, testResults)
	record.DoctorID = doctorID
	record.DoctorName = doctorName
	return record
}

// HasDoctor returns true if the record has doctor information
func (r *Record) HasDoctor() bool {
	return r.DoctorID != "" && r.DoctorName != ""
}

// SetDoctor sets the doctor information for the record
func (r *Record) SetDoctor(doctorID, doctorName string) {
	r.DoctorID = doctorID
	r.DoctorName = doctorName
}

// ClearDoctor removes the doctor information from the record
func (r *Record) ClearDoctor() {
	r.DoctorID = ""
	r.DoctorName = ""
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
	DoctorID    string              `json:"doctor_id,omitempty"`
	DoctorName  string              `json:"doctor_name,omitempty"`
	TestResults []TestResultRequest `json:"test_results"`
}

type UpdateRecordRequest struct {
	ComboName   string              `json:"combo_name" bson:"combo_name"`
	PatientID   string              `json:"patient_id" bson:"patient_id"`
	Patient     *Patient            `json:"patient,omitempty" bson:"patient,omitempty"`
	DoctorID    string              `json:"doctor_id,omitempty" bson:"doctor_id,omitempty"`
	DoctorName  string              `json:"doctor_name,omitempty" bson:"doctor_name,omitempty"`
	TestResults []TestResultRequest `json:"test_results" bson:"test_results"`
}

type RecordQueryOptions struct {
	Keyword   string
	Status    string
	PatientID string
	DoctorID  string
	TestID    string // Filter by test_id
	TestName  string // For display in report summary
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
	RevenueReport           ReportType = "bao_cao_doanh_thu"
)

// RecordForRevenueReport represents a minimal record for revenue reports (lightweight response)
type RecordForRevenueReport struct {
	ID             string    `json:"id"`
	PatientName    string    `json:"patient_name"`
	PatientPhone   string    `json:"patient_phone"`
	PatientAddress string    `json:"patient_address,omitempty"`
	ComboName      string    `json:"combo_name"`
	DoctorName     string    `json:"doctor_name,omitempty"`
	DoctorID       string    `json:"doctor_id,omitempty"`
	Status         string    `json:"status"`
	TotalPrice     int       `json:"total_price"`
	CreatedAt      time.Time `json:"created_at"`
}

// ReportSummary represents aggregated data for revenue reports
type ReportSummary struct {
	TotalRecords     int        `json:"total_records"`
	TotalRevenue     int        `json:"total_revenue"`
	StartDate        *time.Time `json:"start_date"`
	EndDate          *time.Time `json:"end_date"`
	FilteredTestName string     `json:"filtered_test_name,omitempty"`
	TestCount        int        `json:"test_count,omitempty"`
}

// ReportResponse represents the complete response for revenue reports
type ReportResponse struct {
	Records    []*RecordForRevenueReport `json:"records"`
	Pagination *PaginationResponse       `json:"pagination"`
	Summary    *ReportSummary            `json:"summary"`
}
