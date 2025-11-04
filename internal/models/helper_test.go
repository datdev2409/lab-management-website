package models

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateRandomID(t *testing.T) {
	tests := []struct {
		name   string
		prefix string
	}{
		{
			name:   "patient prefix",
			prefix: "patient_",
		},
		{
			name:   "test prefix",
			prefix: "test_",
		},
		{
			name:   "record prefix",
			prefix: "record_",
		},
		{
			name:   "empty prefix",
			prefix: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id := GenerateRandomID(tt.prefix)
			
			assert.Greater(t, len(id), len(tt.prefix), "ID should be longer than prefix")
			
			if tt.prefix != "" {
				assert.True(t, strings.HasPrefix(id, tt.prefix), "ID should start with prefix")
			}
		})
	}

	// Test uniqueness
	t.Run("generates unique IDs", func(t *testing.T) {
		id1 := GenerateRandomID("test_")
		id2 := GenerateRandomID("test_")
		assert.NotEqual(t, id1, id2, "Generated IDs should be unique")
	})
}

func TestNormalizeString(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "basic text",
			input: "Hello World",
			want:  "helloworld",
		},
		{
			name:  "multiple spaces",
			input: "Hello   World   Test",
			want:  "helloworldtest",
		},
		{
			name:  "tabs and newlines",
			input: "Hello\tWorld\nTest",
			want:  "helloworldtest",
		},
		{
			name:  "mixed case",
			input: "HeLLo WoRLd",
			want:  "helloworld",
		},
		{
			name:  "empty string",
			input: "",
			want:  "",
		},
		{
			name:  "only spaces",
			input: "   ",
			want:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NormalizeString(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestToBSONDocument(t *testing.T) {
	t.Run("convert struct to BSON", func(t *testing.T) {
		type testStruct struct {
			Name  string `bson:"name"`
			Age   int    `bson:"age"`
			Email string `bson:"email"`
		}

		data := testStruct{
			Name:  "John Doe",
			Age:   30,
			Email: "john@example.com",
		}

		result, err := ToBSONDocument(data)
		require.NoError(t, err)
		
		assert.Equal(t, "John Doe", result["name"])
		assert.Equal(t, int32(30), result["age"])
		assert.Equal(t, "john@example.com", result["email"])
	})

	t.Run("convert map to BSON", func(t *testing.T) {
		data := map[string]interface{}{
			"key1": "value1",
			"key2": 123,
			"key3": true,
		}

		result, err := ToBSONDocument(data)
		require.NoError(t, err)
		
		assert.Equal(t, "value1", result["key1"])
	})

	t.Run("handle nil", func(t *testing.T) {
		// nil cannot be marshaled to BSON
		_, err := ToBSONDocument(nil)
		assert.Error(t, err, "Should return error for nil input")
	})

	t.Run("handle invalid type", func(t *testing.T) {
		// channels cannot be marshaled to BSON
		ch := make(chan int)
		_, err := ToBSONDocument(ch)
		assert.Error(t, err, "Should return error for channel type")
	})
}
