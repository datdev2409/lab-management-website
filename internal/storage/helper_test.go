package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetCurrentTime(t *testing.T) {
	// Test that GetCurrentTime returns a valid time
	currentTime := GetCurrentTime()
	
	assert.False(t, currentTime.IsZero(), "GetCurrentTime should not return zero time")
	
	// Test that the location is set (should be Asia/Ho_Chi_Minh)
	location := currentTime.Location()
	assert.NotNil(t, location)
	
	// The location name should be "Asia/Ho_Chi_Minh" or fall back to local time
	locationName := location.String()
	assert.NotEmpty(t, locationName)
}
