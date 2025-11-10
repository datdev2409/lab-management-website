package main

import (
	"log"
	"net/http"
	"os"

	"github.com/datdev2409/lab-admin-go/internal/db"
	"github.com/datdev2409/lab-admin-go/internal/db/sqlc"
	"github.com/datdev2409/lab-admin-go/internal/handlers"
	"github.com/datdev2409/lab-admin-go/internal/logger"
	"github.com/datdev2409/lab-admin-go/internal/repository"
	"github.com/datdev2409/lab-admin-go/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func GetEnv(key, defaultValue string) string {
	if value, found := os.LookupEnv(key); found {
		return value
	}
	return defaultValue
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("No .env file found")
	}

	log := logger.Init()
	defer log.Sync()

	port := GetEnv("SERVER_PORT", "8080")

	pgPool, err := db.NewPostgresPool(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Failed to connect to Postgres", zap.Error(err))
	}
	defer pgPool.Close()

	v := validator.New()
	queries := sqlc.New(pgPool)
	patientRepository := repository.NewPgPatientRepository(queries)
	patientService := service.NewPatientService(patientRepository)
	patientHandler := handlers.NewPatientHandler(patientService, v)

	r := chi.NewRouter()

	r.Route("/api/v1/patients", func(r chi.Router) {
		r.Get("/", handlers.Make(patientHandler.SearchPatientsByKeyword))
		r.Post("/", handlers.Make(patientHandler.CreatePatient))
		r.Get("/{id}", handlers.Make(patientHandler.GetPatient))
		r.Patch("/{id}", handlers.Make(patientHandler.UpdatePatient))
		r.Delete("/{id}", handlers.Make(patientHandler.DeletePatient))
	})

	// Legacy routes for compatibility with old API paths
	r.Route("/api/patients", func(r chi.Router) {
		r.Get("/{id}", handlers.Make(patientHandler.GetPatient))
		r.Delete("/{id}", handlers.Make(patientHandler.DeletePatient))
	})

	log.Info("Server is running", zap.String("port", port))
	err = http.ListenAndServe(":"+port, r)
	if err != nil {
		log.Error("Server error", zap.Error(err))
	}
}
