package db

import (
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func NewMongoClient(conn string) *mongo.Client {
	client, err := mongo.Connect(options.Client().ApplyURI(conn))

	if err != nil {
		panic(err)
	}

	return client
}

type Patient struct {
	ID      string `json:"id" bson:"id,omitempty"`
	Name    string `json:"name" bson:"name"`
	YOB     string `json:"yob" bson:"yob"`
	Gender  string `json:"gender" bson:"gender"`
	Address string `json:"address" bson:"address"`
	Phone   string `json:"phone" bson:"phone"`
}

type Test struct {
	ID          string  `json:"id" bson:"id,omitempty"`
	Name        string  `json:"name" bson:"name"`
	Price       int     `json:"price" bson:"price"`
	NormalValue string  `json:"normal_value" bson:"normal_value"`
	Unit        string  `json:"unit" bson:"unit"`
	LowerBound  float64 `json:"lower_bound" bson:"lower_bound"`
	UpperBound  float64 `json:"upper_bound" bson:"upper_bound"`
}
