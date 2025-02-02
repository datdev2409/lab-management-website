package storage

import (
	"go.mongodb.org/mongo-driver/v2/bson"
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
