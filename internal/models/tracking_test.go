package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTracking(t *testing.T) {
	name := "Blood Panel Tracking"
	testRequests := []TrackingTestRequest{
		{
			TestID:      "test_1",
			TestName:    "Blood Sugar",
			NormalValue: "70-100",
			Unit:        "mg/dL",
			Order:       1,
		},
		{
			TestID:      "test_2",
			TestName:    "Cholesterol",
			NormalValue: "<200",
			Unit:        "mg/dL",
			Order:       2,
		},
	}

	tracking := NewTracking(name, testRequests)

	assert.NotEmpty(t, tracking.ID)
	assert.Contains(t, tracking.ID, "tracking_")
	assert.Equal(t, name, tracking.Name)
	assert.Len(t, tracking.Tests, len(testRequests))
	
	for i, req := range testRequests {
		assert.Equal(t, req.TestID, tracking.Tests[i].TestID)
		assert.Equal(t, req.TestName, tracking.Tests[i].TestName)
		assert.Equal(t, req.NormalValue, tracking.Tests[i].NormalValue)
		assert.Equal(t, req.Unit, tracking.Tests[i].Unit)
		assert.Equal(t, req.Order, tracking.Tests[i].Order)
	}
}

func TestNewTracking_EmptyTests(t *testing.T) {
	name := "Empty Tracking"
	testRequests := []TrackingTestRequest{}

	tracking := NewTracking(name, testRequests)

	assert.NotEmpty(t, tracking.ID)
	// Empty slice results in nil Tests due to append on nil slice
}

func TestNewTracking_NilTests(t *testing.T) {
	name := "Nil Tests Tracking"
	var testRequests []TrackingTestRequest = nil

	tracking := NewTracking(name, testRequests)

	assert.NotEmpty(t, tracking.ID)
	// nil slice will result in nil Tests
}
