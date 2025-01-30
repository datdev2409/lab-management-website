package main

import (
	"context"
	"log"

	"github.com/datdev2409/lab-admin-go/internal/db"
)

func main() {
	app := &Application{
		Config: &Config{
			Env:  "development",
			Port: ":8081",
			DB: &DBConfig{
				Addr: "mongodb://root:password123@localhost:27017/",
			},
		},
	}

	mongoClient := db.NewMongoClient(app.Config.DB.Addr)
	defer func() {
		if err := mongoClient.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
	log.Println("DB connection is established")

	app.Store = mongoClient.Database("labadmin")

	router := app.NewRouter()

	log.Println("Server is running on port", app.Config.Port)
	log.Fatal(app.Run(&router))
}
