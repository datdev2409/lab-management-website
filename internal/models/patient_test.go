package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPatient(t *testing.T) {
	name := "John Doe"
	yob := "1990"
	gender := "Male"
	address := "123 Main St"
	phone := "555-1234"

	patient := NewPatient(name, yob, gender, address, phone)

	assert.NotNil(t, patient)
	assert.NotEmpty(t, patient.ID)
	assert.Contains(t, patient.ID, "patient_")
	assert.Equal(t, name, patient.Name)
	assert.Equal(t, yob, patient.YOB)
	assert.Equal(t, gender, patient.Gender)
	assert.Equal(t, address, patient.Address)
	assert.Equal(t, phone, patient.Phone)
	assert.False(t, patient.CreatedAt.IsZero())
	assert.False(t, patient.UpdatedAt.IsZero())
	assert.True(t, patient.CreatedAt.Equal(patient.UpdatedAt))
}

func TestGetStringPtr(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  *string
	}{
		{
			name:  "non-empty string",
			input: "test",
			want:  stringPtr("test"),
		},
		{
			name:  "empty string returns nil",
			input: "",
			want:  nil,
		},
		{
			name:  "whitespace string",
			input: "  ",
			want:  stringPtr("  "),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetStringPtr(tt.input)
			
			if tt.want == nil {
				assert.Nil(t, got)
			} else {
				assert.NotNil(t, got)
				assert.Equal(t, *tt.want, *got)
			}
		})
	}
}

// Helper functions
func stringPtr(s string) *string {
	return &s
}
