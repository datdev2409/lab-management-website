package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCombo(t *testing.T) {
	name := "Basic Health Check"
	testIDs := []string{"test_1", "test_2", "test_3"}

	combo := NewCombo(name, testIDs)

	assert.NotNil(t, combo)
	assert.NotEmpty(t, combo.ID)
	assert.Contains(t, combo.ID, "combo_")
	assert.Equal(t, name, combo.Name)
	assert.Len(t, combo.TestIDs, len(testIDs))
	
	for i, id := range testIDs {
		assert.Equal(t, id, combo.TestIDs[i])
	}
	
	assert.False(t, combo.CreatedAt.IsZero())
	assert.False(t, combo.UpdatedAt.IsZero())
}

func TestNewCombo_EmptyTestIDs(t *testing.T) {
	name := "Empty Combo"
	testIDs := []string{}

	combo := NewCombo(name, testIDs)

	assert.NotNil(t, combo)
	assert.Empty(t, combo.TestIDs)
}

func TestNewCombo_NilTestIDs(t *testing.T) {
	name := "Nil Test IDs"
	var testIDs []string = nil

	combo := NewCombo(name, testIDs)

	assert.NotNil(t, combo)
	// nil slice is preserved
}
