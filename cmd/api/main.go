package main

import (
	"context"
	"log"
	"log/slog"
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

func getLogger(env string, logLevelStr string) *slog.Logger {
	var logLevel slog.Level
	switch logLevelStr {
	case "debug":
		logLevel = slog.LevelDebug
	case "info":
		logLevel = slog.LevelInfo
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}
	opts := &slog.HandlerOptions{
		Level:     logLevel,
		AddSource: true,
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, opts))
	if env == "production" {
		logger = slog.New(slog.NewJSONHandler(os.Stdout, opts))
	}

	return logger
}

func GetEnv(key, defaultValue string) string {
	if value, found := os.LookupEnv(key); found {
		return value
	}
	return defaultValue
}

func main() {
	env := GetEnv("GO_ENV", "local")
	// Load environment variables from .env file
	err := godotenv.Load(".env." + env)
	if err != nil {
		log.Fatalf("Error loading .env.%s file", env)
	}

	logLevel := GetEnv("LOG_LEVEL", "debug")
	slog.SetDefault(getLogger(env, logLevel))

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
			slog.Error("Mongo disconnect error", slog.String("error", err.Error()))
		}
	}()

	store := storage.NewMongoStorage(mongoClient)
	handler := handlers.NewHandler(store)

	app.Init(config, handler.Router)

	slog.Info("Server is running", slog.String("port", app.Config.Port))
	err = app.Start()
	if err != nil {
		slog.Error("Server error", slog.String("error", err.Error()))
	}
}
