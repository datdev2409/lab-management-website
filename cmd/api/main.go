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

	// Initialize Patient
	patientRepository := repository.NewPgPatientRepository(queries)
	patientService := service.NewPatientService(patientRepository)
	patientHandler := handlers.NewPatientHandler(patientService, v)

	// Initialize Doctor
	doctorRepository := repository.NewPgDoctorRepository(queries)
	doctorService := service.NewDoctorService(doctorRepository)
	doctorHandler := handlers.NewDoctorHandler(doctorService, v)

	// Initialize Test
	testRepository := repository.NewPgTestRepository(queries)
	testService := service.NewTestService(testRepository)
	testHandler := handlers.NewTestHandler(testService, v)

	// Initialize Combo
	comboRepository := repository.NewPgComboRepository(queries, pgPool)
	comboService := service.NewComboService(comboRepository)
	comboHandler := handlers.NewComboHandler(comboService, v)

	r := chi.NewRouter()

	// Patient routes
	r.Route("/api/v1/patients", func(r chi.Router) {
		r.Get("/", handlers.Make(patientHandler.SearchPatientsByKeyword))
		r.Post("/", handlers.Make(patientHandler.CreatePatient))
		r.Get("/{id}", handlers.Make(patientHandler.GetPatient))
		r.Patch("/{id}", handlers.Make(patientHandler.UpdatePatient))
		r.Delete("/{id}", handlers.Make(patientHandler.DeletePatient))
	})

	// Legacy patient routes for compatibility
	r.Route("/api/patients", func(r chi.Router) {
		r.Get("/{id}", handlers.Make(patientHandler.GetPatient))
		r.Delete("/{id}", handlers.Make(patientHandler.DeletePatient))
	})

	// Doctor routes
	r.Route("/api/v1/doctors", func(r chi.Router) {
		r.Get("/", handlers.Make(doctorHandler.SearchDoctorsByKeyword))
		r.Post("/", handlers.Make(doctorHandler.CreateDoctor))
		r.Get("/{id}", handlers.Make(doctorHandler.GetDoctor))
		r.Patch("/{id}", handlers.Make(doctorHandler.UpdateDoctor))
		r.Delete("/{id}", handlers.Make(doctorHandler.DeleteDoctor))
	})

	// Doctor page route
	r.Get("/danh-muc-bac-si", handlers.Make(doctorHandler.HandleDoctorPage))

	// Test routes
	r.Route("/api/v1/tests", func(r chi.Router) {
		r.Get("/", handlers.Make(testHandler.SearchTestsByName))
		r.Post("/", handlers.Make(testHandler.CreateTest))
		r.Post("/bulk", handlers.Make(testHandler.BulkCreateTests))
		r.Get("/{id}", handlers.Make(testHandler.GetTest))
		r.Patch("/{id}", handlers.Make(testHandler.UpdateTest))
		r.Delete("/{id}", handlers.Make(testHandler.DeleteTest))
	})

	// Test page route
	r.Get("/danh-muc-xet-nghiem", handlers.Make(testHandler.HandleTestPage))

	// Combo routes
	r.Route("/api/v1/combos", func(r chi.Router) {
		r.Get("/", handlers.Make(comboHandler.SearchCombosByName))
		r.Get("/all", handlers.Make(comboHandler.ListAllCombos))
		r.Post("/", handlers.Make(comboHandler.CreateCombo))
		r.Get("/{id}", handlers.Make(comboHandler.GetCombo))
		r.Put("/{id}", handlers.Make(comboHandler.UpdateCombo))
		r.Delete("/{id}", handlers.Make(comboHandler.DeleteCombo))
		r.Get("/{id}/tests", handlers.Make(comboHandler.GetComboTests))
	})

	r.Route("/danh-muc-goi-xet-nghiem", func(r chi.Router) {
		r.Get("/new", handlers.Make(comboHandler.HandleComboCreatePage))
		r.Get("/{id}/edit", handlers.Make(comboHandler.HandleComboEditPage))
		r.Get("/", handlers.Make(comboHandler.HandleComboPage))
	})

	log.Info("Server is running", zap.String("port", port))
	err = http.ListenAndServe(":"+port, r)
	if err != nil {
		log.Error("Server error", zap.Error(err))
	}
}
