package main

import (
	"encoding/json"
	"log"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/datdev2409/lab-admin-go/internal/db"
	"github.com/datdev2409/lab-admin-go/internal/templates/pages"
	"github.com/datdev2409/lab-admin-go/internal/templates/partials"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
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
	Config *Config
	Store  *mongo.Database
}

func (app *Application) Run(mux *http.Handler) error {
	server := &http.Server{
		Addr:    app.Config.Port,
		Handler: *mux,
	}

	err := server.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}

func writeJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

func (app *Application) NewRouter() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(6 * time.Second))

	fs := http.FileServer(http.Dir("public"))

	r.Put("/api/patients/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		patient := db.Patient{
			ID:      id,
			Name:    r.FormValue("patient_name"),
			YOB:     r.FormValue("patient_yob"),
			Gender:  r.FormValue("patient_gender"),
			Address: r.FormValue("patient_address"),
			Phone:   r.FormValue("patient_phone"),
		}

		filter := bson.D{{Key: "id", Value: id}}
		update := bson.D{{Key: "$set", Value: patient}}
		_, err := app.Store.Collection("patients").UpdateOne(r.Context(), filter, update)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		partials.PatientRow(patient).Render(r.Context(), w)
	})

	r.Get("/api/patients/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		filter := bson.D{{Key: "id", Value: id}}
		var patient db.Patient
		err := app.Store.Collection("patients").FindOne(r.Context(), filter).Decode(&patient)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusNotFound)
			return
		}

		partials.SelectUserForm(patient, false).Render(r.Context(), w)
	})

	r.Get("/api/patients", func(w http.ResponseWriter, r *http.Request) {
		limit, err := strconv.ParseInt(r.URL.Query().Get("limit"), 10, 64)
		if err != nil {
			limit = 5
		}

		keyword := r.URL.Query().Get("patient_name")
		filter := bson.D{{}}
		if keyword != "" {
			filter = bson.D{{Key: "name", Value: bson.D{{Key: "$regex", Value: keyword}, {Key: "$options", Value: "i"}}}}
		}
		cursor, err := app.Store.Collection("patients").Find(r.Context(), filter, options.Find().SetLimit(limit))

		var patients []db.Patient
		if err := cursor.All(r.Context(), &patients); err != nil {
			log.Println(err)
		}

		if err != nil {
			patients = []db.Patient{}
		}

		target := r.Header.Get("HX-Target")

		switch target {
		case "patient-table":
			partials.PatientTable(patients).Render(r.Context(), w)
		case "patient-suggestion-list":
			partials.PatientSuggestionList(patients).Render(r.Context(), w)
		}

	})

	r.Delete("/api/patients/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		filter := bson.D{{Key: "id", Value: id}}
		_, err := app.Store.Collection("patients").DeleteOne(r.Context(), filter)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	})

	r.Post("/api/patients", func(w http.ResponseWriter, r *http.Request) {
		patient := db.Patient{
			ID:      "p-" + uuid.NewString(),
			Name:    r.FormValue("patient_name"),
			YOB:     r.FormValue("patient_yob"),
			Gender:  r.FormValue("patient_gender"),
			Address: r.FormValue("patient_address"),
			Phone:   r.FormValue("patient_phone"),
		}

		_, err := app.Store.Collection("patients").InsertOne(r.Context(), patient)
		if err != nil {
			log.Println(err)
			errorMessage := `<div class="alert alert-danger" role="alert">Đã có lỗi xảy ra khi thêm bệnh nhân.</div>`
			w.Write([]byte(errorMessage))
		}

		w.WriteHeader(http.StatusCreated)
		successMessage := `<div class="alert alert-success" role="alert">Thêm bệnh nhân thành công.</div>`
		w.Write([]byte(successMessage))
	})

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/phieu-xet-nghiem", http.StatusSeeOther)
	})

	r.Get("/phieu-xet-nghiem", func(w http.ResponseWriter, r *http.Request) {
		// templates.().Render(r.Context(), w)
		// temp.PhieuXetNghiemPage().Render(r.Context(), w)
		pages.PhieuXetNghiemPage().Render(r.Context(), w)
	})

	r.Get("/danh-muc-benh-nhan", func(w http.ResponseWriter, r *http.Request) {
		filter := bson.D{{}}
		var patients []db.Patient
		cursor, err := app.Store.Collection("patients").Find(r.Context(), filter)
		if err != nil {
			patients = []db.Patient{}
		}

		if err := cursor.All(r.Context(), &patients); err != nil {
			patients = []db.Patient{}
		}
		pages.PatientsPage(patients).Render(r.Context(), w)
	})

	r.Get("/public/*", http.StripPrefix("/public", fs).ServeHTTP)

	r.Get("/v1/api/health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	// Test
	r.Get("/danh-muc-xet-nghiem", func(w http.ResponseWriter, r *http.Request) {
		pages.TestPage().Render(r.Context(), w)
	})

	r.Post("/api/tests", func(w http.ResponseWriter, r *http.Request) {
		testLowerBound, err := strconv.ParseFloat(r.FormValue("test_lower_bound"), 32)
		if err != nil {
			log.Println(err)
			errorMessage := `<div class="alert alert-danger" role="alert">Đã có lỗi xảy ra khi thêm xét nghiệm.</div>`
			w.Write([]byte(errorMessage))
		}
		testUpperBound, err := strconv.ParseFloat(r.FormValue("test_upper_bound"), 32)
		if err != nil {
			log.Println(err)
			errorMessage := `<div class="alert alert-danger" role="alert">Đã có lỗi xảy ra khi thêm xét nghiệm.</div>`
			w.Write([]byte(errorMessage))
		}
		testPrice, err := strconv.Atoi(r.FormValue("test_price"))
		if err != nil {
			log.Println(err)
			errorMessage := `<div class="alert alert-danger" role="alert">Đã có lỗi xảy ra khi thêm xét nghiệm.</div>`
			w.Write([]byte(errorMessage))
		}
		test := db.Test{
			ID:          "t-" + uuid.NewString(),
			Name:        r.FormValue("test_name"),
			NormalValue: r.FormValue("test_normal_value"),
			Unit:        r.FormValue("test_unit"),
			LowerBound:  math.Round(testLowerBound*100) / 100,
			UpperBound:  math.Round(testUpperBound*100) / 100,
			Price:       testPrice,
		}

		log.Println(test)

		_, err = app.Store.Collection("tests").InsertOne(r.Context(), test)
		if err != nil {
			log.Println(err)
			errorMessage := `<div class="alert alert-danger" role="alert">Đã có lỗi xảy ra khi thêm xét nghiệm.</div>`
			w.Write([]byte(errorMessage))
		}

		w.WriteHeader(http.StatusCreated)
		successMessage := `<div class="alert alert-success" role="alert">Thêm xét nghiệm thành công.</div>`
		w.Write([]byte(successMessage))
	})

	return r
}
