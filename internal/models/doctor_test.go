package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDoctor(t *testing.T) {
	name := "Dr. Smith"
	phone := "555-9999"
	address := "456 Medical Center"

	doctor := NewDoctor(name, phone, address)

	assert.NotNil(t, doctor)
	assert.NotEmpty(t, doctor.ID)
	assert.Contains(t, doctor.ID, "doctor_")
	assert.Equal(t, name, doctor.Name)
	assert.Equal(t, phone, doctor.Phone)
	assert.Equal(t, address, doctor.Address)
	assert.False(t, doctor.CreatedAt.IsZero())
	assert.False(t, doctor.UpdatedAt.IsZero())
}

func TestDoctor_GetID(t *testing.T) {
	doctor := NewDoctor("Dr. Jones", "555-8888", "789 Hospital Rd")
	
	got := doctor.GetID()
	assert.Equal(t, doctor.ID, got)
	assert.NotEmpty(t, got)
}
