package storage

import (
	"github.com/datdev2409/lab-admin-go/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type FilterCondition struct {
	Operator string
	Value    string
	Option   string // optional field
}

func BuildMongoFilter(filterCriterias map[string]FilterCondition) bson.D {
	filters := bson.D{}

	for key, condition := range filterCriterias {
		filter := bson.D{{Key: key, Value: bson.D{{Key: condition.Operator, Value: condition.Value}}}}
		if condition.Option != "" {
			filter[0].Value = append(filter[0].Value.(bson.D), bson.E{Key: "$options", Value: condition.Option})
		}
		filters = append(filters, filter...)
	}

	return filters
}

func BuildMongoSortAndPaginationOptions(opts models.GenericQueryOptions) *options.FindOptionsBuilder {
	findOptions := options.Find()

	if opts.Page > 0 && opts.PageSize > 0 {
		skip := (opts.Page - 1) * opts.PageSize // Skip items for the current page
		limit := opts.PageSize                  // Limit per page

		findOptions.SetSkip(int64(skip))
		findOptions.SetLimit(int64(limit))
	}

	// Sorting: Apply if specified
	if opts.SortBy != "" {
		sortDirection := 1 // ascending
		if opts.SortOrder == "desc" {
			sortDirection = -1
		}
		findOptions.SetSort(map[string]int{
			opts.SortBy: sortDirection,
		})
	}

	return findOptions
}
