package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/datdev2409/lab-admin-go/internal/db"
	"github.com/datdev2409/lab-admin-go/internal/handlers"
	"github.com/datdev2409/lab-admin-go/internal/storage"
	"github.com/joho/godotenv"
)

type DBConfig struct {
	Addr string
	Port int
}

type Config struct {
	Env  string
	Port string
	DB   *DBConfig
}

type Application struct {
	Config  *Config
	Handler http.Handler
}

func (app *Application) Init(config *Config, handler http.Handler) {
	app.Config = config
	app.Handler = handler
}

func (app *Application) Start() error {
	err := http.ListenAndServe(app.Config.Port, app.Handler)
	return err
}

func main() {
	env := os.Getenv("GO_ENV")
	if env == "" {
		env = "local"
	}
	err := godotenv.Load(".env." + env)
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	app := &Application{}

	// Init config
	config := &Config{
		Env:  os.Getenv("ENV"),
		Port: os.Getenv("SERVER_PORT"),
		DB: &DBConfig{
			Addr: os.Getenv("MONGODB_URI"),
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

	app.Init(config, handler.Router)

	log.Println("Server is running on port", app.Config.Port)
	log.Fatal(app.Start())
}
