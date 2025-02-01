package models

type Test struct {
	ID          string  `json:"id" bson:"_id,omitempty"`
	Name        string  `json:"name" bson:"name"`
	Price       int     `json:"price" bson:"price"`
	NormalValue string  `json:"normal_value" bson:"normal_value"`
	Unit        string  `json:"unit" bson:"unit"`
	LowerBound  float64 `json:"lower_bound" bson:"lower_bound"`
	UpperBound  float64 `json:"upper_bound" bson:"upper_bound"`
}
