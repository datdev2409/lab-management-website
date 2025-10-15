package models

import (
	"time"
)

type TestResult struct {
	ID                     string  `json:"id" bson:"_id" db:"id"`
	RecordID               string  `json:"record_id" bson:"record_id" db:"record_id"` // Foreign key to records table
	TestID                 *string `json:"test_id" bson:"test_id" db:"test_id"`       // Foreign key to tests table (nullable)
	Name                   string  `json:"name" bson:"name" db:"name"`
	Price                  int     `json:"price" bson:"price" db:"price"`
	NormalValue            string  `json:"normal_value" bson:"normal_value" db:"normal_value"`
	Unit                   string  `json:"unit" bson:"unit" db:"unit"`
	LowerBound             float64 `json:"lower_bound" bson:"lower_bound" db:"lower_bound"`
	UpperBound             float64 `json:"upper_bound" bson:"upper_bound" db:"upper_bound"`
	Result                 string  `json:"result" bson:"result" db:"result"`
	ResultText             string  `json:"result_text" bson:"result_text" db:"result_text"`
	Abnormal               bool    `json:"abnormal" bson:"abnormal" db:"abnormal"`
	ManualAbnormalOverride bool    `json:"manual_abnormal_override" bson:"manual_abnormal_override" db:"manual_abnormal_override"`
}

type Record struct {
	ID          string       `json:"id" bson:"_id,omitempty" db:"id"`
	PatientID   string       `json:"patient_id" bson:"patient_id" db:"patient_id"`                  // Foreign key to patients table
	DoctorID    *string      `json:"doctor_id,omitempty" bson:"doctor_id,omitempty" db:"doctor_id"` // Foreign key to doctors table (nullable)
	ComboName   string       `json:"combo_name" bson:"combo_name" db:"combo_name"`
	DoctorName  string       `json:"doctor_name,omitempty" bson:"doctor_name,omitempty" db:"doctor_name"`
	Patient     Patient      `json:"patient" bson:"patient" db:"-"`           // Not directly mapped, populated via join
	TestResults []TestResult `json:"test_results" bson:"test_results" db:"-"` // Not directly mapped, handled via separate table
	Status      string       `json:"status" bson:"status" db:"status"`
	CreatedAt   time.Time    `json:"created_at" bson:"created_at" db:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at" bson:"updated_at" db:"updated_at"`
}

func NewRecord(patient Patient, comboName string, testResults []TestResult) Record {
	now := time.Now()
	record := Record{
		ID:          GenerateRandomID("record_"),
		PatientID:   patient.ID,
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
	if doctorID != "" {
		record.DoctorID = &doctorID
	}
	record.DoctorName = doctorName
	return record
}

// HasDoctor returns true if the record has doctor information
func (r *Record) HasDoctor() bool {
	return r.DoctorID != nil && *r.DoctorID != "" && r.DoctorName != ""
}

// SetDoctor sets the doctor information for the record
func (r *Record) SetDoctor(doctorID, doctorName string) {
	if doctorID != "" {
		r.DoctorID = &doctorID
	} else {
		r.DoctorID = nil
	}
	r.DoctorName = doctorName
}

// ClearDoctor removes the doctor information from the record
func (r *Record) ClearDoctor() {
	r.DoctorID = nil
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
