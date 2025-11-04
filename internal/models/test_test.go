package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTest(t *testing.T) {
	name := "Blood Sugar"
	price := 50000
	normalValue := "70-100"
	unit := "mg/dL"
	lowerBound := 70.0
	upperBound := 100.0

	test := NewTest(name, price, normalValue, unit, lowerBound, upperBound)

	assert.NotNil(t, test)
	assert.NotEmpty(t, test.ID)
	assert.Contains(t, test.ID, "test_")
	assert.Equal(t, name, test.Name)
	assert.Equal(t, price, test.Price)
	assert.Equal(t, normalValue, test.NormalValue)
	assert.Equal(t, unit, test.Unit)
	assert.Equal(t, lowerBound, test.LowerBound)
	assert.Equal(t, upperBound, test.UpperBound)
	assert.False(t, test.CreatedAt.IsZero())
	assert.False(t, test.UpdatedAt.IsZero())
}

func TestTest_GetID(t *testing.T) {
	test := NewTest("Test Name", 1000, "normal", "unit", 0.0, 100.0)
	
	got := test.GetID()
	assert.Equal(t, test.ID, got)
	assert.NotEmpty(t, got)
}
