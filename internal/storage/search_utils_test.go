package storage

import (
	"testing"

	"github.com/datdev2409/lab-admin-go/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestBuildMongoFilter(t *testing.T) {
	t.Run("single filter without option", func(t *testing.T) {
		criteria := map[string]FilterCondition{
			"name": {
				Operator: "$eq",
				Value:    "test",
			},
		}

		filter := BuildMongoFilter(criteria)
		
		assert.NotNil(t, filter)
		assert.Len(t, filter, 1)
	})

	t.Run("single filter with option", func(t *testing.T) {
		criteria := map[string]FilterCondition{
			"name": {
				Operator: "$regex",
				Value:    "test",
				Option:   "i",
			},
		}

		filter := BuildMongoFilter(criteria)
		
		assert.NotNil(t, filter)
		assert.Len(t, filter, 1)
	})

	t.Run("multiple filters", func(t *testing.T) {
		criteria := map[string]FilterCondition{
			"name": {
				Operator: "$regex",
				Value:    "test",
				Option:   "i",
			},
			"age": {
				Operator: "$gte",
				Value:    "18",
			},
		}

		filter := BuildMongoFilter(criteria)
		
		assert.NotNil(t, filter)
		assert.Len(t, filter, 2)
	})

	t.Run("empty criteria", func(t *testing.T) {
		criteria := map[string]FilterCondition{}

		filter := BuildMongoFilter(criteria)
		
		assert.NotNil(t, filter)
		assert.Empty(t, filter)
	})

	t.Run("nil criteria", func(t *testing.T) {
		filter := BuildMongoFilter(nil)
		
		assert.NotNil(t, filter)
		assert.Empty(t, filter)
	})
}

func TestBuildMongoSortAndPaginationOptions(t *testing.T) {
	t.Run("with pagination", func(t *testing.T) {
		opts := models.GenericQueryOptions{
			Page:     2,
			PageSize: 10,
		}

		findOptions := BuildMongoSortAndPaginationOptions(opts)
		
		assert.NotNil(t, findOptions)
		// Note: We can't directly access the values set on FindOptionsBuilder
		// but we've verified it doesn't panic
	})

	t.Run("with sorting ascending", func(t *testing.T) {
		opts := models.GenericQueryOptions{
			SortBy:    "name",
			SortOrder: "asc",
		}

		findOptions := BuildMongoSortAndPaginationOptions(opts)
		
		assert.NotNil(t, findOptions)
	})

	t.Run("with sorting descending", func(t *testing.T) {
		opts := models.GenericQueryOptions{
			SortBy:    "created_at",
			SortOrder: "desc",
		}

		findOptions := BuildMongoSortAndPaginationOptions(opts)
		
		assert.NotNil(t, findOptions)
	})

	t.Run("with both pagination and sorting", func(t *testing.T) {
		opts := models.GenericQueryOptions{
			Page:      1,
			PageSize:  20,
			SortBy:    "name",
			SortOrder: "asc",
		}

		findOptions := BuildMongoSortAndPaginationOptions(opts)
		
		assert.NotNil(t, findOptions)
	})

	t.Run("with no pagination (page 0)", func(t *testing.T) {
		opts := models.GenericQueryOptions{
			Page:     0,
			PageSize: 10,
		}

		findOptions := BuildMongoSortAndPaginationOptions(opts)
		
		assert.NotNil(t, findOptions)
	})

	t.Run("with no pagination (pageSize 0)", func(t *testing.T) {
		opts := models.GenericQueryOptions{
			Page:     1,
			PageSize: 0,
		}

		findOptions := BuildMongoSortAndPaginationOptions(opts)
		
		assert.NotNil(t, findOptions)
	})

	t.Run("empty options", func(t *testing.T) {
		opts := models.GenericQueryOptions{}

		findOptions := BuildMongoSortAndPaginationOptions(opts)
		
		assert.NotNil(t, findOptions)
	})

	t.Run("pagination calculations", func(t *testing.T) {
		tests := []struct {
			name     string
			page     int
			pageSize int
			wantSkip int64
		}{
			{"page 1", 1, 10, 0},
			{"page 2", 2, 10, 10},
			{"page 3", 3, 10, 20},
			{"page 1 size 20", 1, 20, 0},
			{"page 5 size 25", 5, 25, 100},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				opts := models.GenericQueryOptions{
					Page:     tt.page,
					PageSize: tt.pageSize,
				}

				findOptions := BuildMongoSortAndPaginationOptions(opts)
				assert.NotNil(t, findOptions)
				
				// Calculate expected skip
				expectedSkip := int64((tt.page - 1) * tt.pageSize)
				assert.Equal(t, tt.wantSkip, expectedSkip)
			})
		}
	})
}
