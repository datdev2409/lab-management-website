package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/datdev2409/lab-admin-go/internal/db"
	"github.com/datdev2409/lab-admin-go/internal/handlers"
	"github.com/datdev2409/lab-admin-go/internal/logger"
	"github.com/datdev2409/lab-admin-go/internal/storage"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
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

func GetEnv(key, defaultValue string) string {
	if value, found := os.LookupEnv(key); found {
		return value
	}
	return defaultValue
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	log := logger.Init()
	defer log.Sync()

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
			log.Error("Mongo disconnect error", zap.Error(err))
		}
	}()

	store := storage.NewMongoStorage(mongoClient)
	handler := handlers.NewHandler(store, log)

	app.Init(config, handler.Router)

	log.Info("Server is running", zap.String("port", app.Config.Port))
	err = app.Start()
	if err != nil {
		log.Error("Server error", zap.Error(err))
	}
}
