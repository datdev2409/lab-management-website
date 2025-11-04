package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRecord(t *testing.T) {
	patient := *NewPatient("John Doe", "1990", "Male", "123 Main St", "555-1234")
	comboName := "Basic Health Check"
	testResults := []TestResult{
		{
			ID:          "test_1",
			Name:        "Blood Sugar",
			Price:       50000,
			NormalValue: "70-100",
			Unit:        "mg/dL",
			LowerBound:  70.0,
			UpperBound:  100.0,
			Result:      "85",
			Abnormal:    false,
		},
	}

	record := NewRecord(patient, comboName, testResults)

	assert.NotEmpty(t, record.ID)
	assert.Contains(t, record.ID, "record_")
	assert.Equal(t, patient.Name, record.Patient.Name)
	assert.Equal(t, comboName, record.ComboName)
	assert.Len(t, record.TestResults, len(testResults))
	assert.Equal(t, "pending", record.Status)
	assert.False(t, record.CreatedAt.IsZero())
	assert.False(t, record.UpdatedAt.IsZero())
	// Default record should not have doctor info
	assert.Empty(t, record.DoctorID)
	assert.Empty(t, record.DoctorName)
}

func TestNewRecordWithDoctor(t *testing.T) {
	patient := *NewPatient("John Doe", "1990", "Male", "123 Main St", "555-1234")
	comboName := "Basic Health Check"
	testResults := []TestResult{}
	doctorID := "doctor_123"
	doctorName := "Dr. Smith"

	record := NewRecordWithDoctor(patient, comboName, testResults, doctorID, doctorName)

	assert.NotEmpty(t, record.ID)
	assert.Contains(t, record.ID, "record_")
	assert.Equal(t, doctorID, record.DoctorID)
	assert.Equal(t, doctorName, record.DoctorName)
	assert.Equal(t, "pending", record.Status)
}

func TestRecord_HasDoctor(t *testing.T) {
	tests := []struct {
		name       string
		doctorID   string
		doctorName string
		want       bool
	}{
		{
			name:       "has both doctor ID and name",
			doctorID:   "doctor_123",
			doctorName: "Dr. Smith",
			want:       true,
		},
		{
			name:       "missing doctor ID",
			doctorID:   "",
			doctorName: "Dr. Smith",
			want:       false,
		},
		{
			name:       "missing doctor name",
			doctorID:   "doctor_123",
			doctorName: "",
			want:       false,
		},
		{
			name:       "missing both",
			doctorID:   "",
			doctorName: "",
			want:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			patient := *NewPatient("John", "1990", "Male", "Address", "Phone")
			record := NewRecordWithDoctor(patient, "Combo", []TestResult{}, tt.doctorID, tt.doctorName)
			
			got := record.HasDoctor()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestRecord_SetDoctor(t *testing.T) {
	patient := *NewPatient("John", "1990", "Male", "Address", "Phone")
	record := NewRecord(patient, "Combo", []TestResult{})

	doctorID := "doctor_456"
	doctorName := "Dr. Jones"

	record.SetDoctor(doctorID, doctorName)

	assert.Equal(t, doctorID, record.DoctorID)
	assert.Equal(t, doctorName, record.DoctorName)
	assert.True(t, record.HasDoctor())
}

func TestRecord_ClearDoctor(t *testing.T) {
	patient := *NewPatient("John", "1990", "Male", "Address", "Phone")
	record := NewRecordWithDoctor(patient, "Combo", []TestResult{}, "doctor_123", "Dr. Smith")

	assert.True(t, record.HasDoctor(), "Initial record should have doctor")

	record.ClearDoctor()

	assert.Empty(t, record.DoctorID)
	assert.Empty(t, record.DoctorName)
	assert.False(t, record.HasDoctor())
}
