package db

import (
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"log"
)

func NewMongoClient(conn string) *mongo.Client {
	client, err := mongo.Connect(options.Client().ApplyURI(conn))

	if err != nil {
		log.Println("error here")
		log.Fatalln(err)
	}

	log.Println("Connected to MongoDB!")

	return client
}
