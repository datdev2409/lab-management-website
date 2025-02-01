package main

import (
	"context"
	"github.com/datdev2409/lab-admin-go/internal/db"
	"github.com/datdev2409/lab-admin-go/internal/handlers"
	"github.com/datdev2409/lab-admin-go/internal/storage"
	"log"
)

func main() {
	app := &Application{}

	// Init config
	config := &Config{
		Env:  "development",
		Port: ":8081",
		DB: &DBConfig{
			Addr: "mongodb://root:password123@localhost:27017/",
		},
	}

	// Init storage
	mongoClient := db.NewMongoClient(config.DB.Addr)
	defer func() {
		if err := mongoClient.Disconnect(context.TODO()); err != nil {
			log.Fatalln(err)
		}
	}()

	store := storage.NewMongoStorage(mongoClient)
	handler := handlers.NewHandler(store)

	app.Init(config, store, handler.Router)

	log.Println("Server is running on port", app.Config.Port)
	log.Fatal(app.Start())
}
