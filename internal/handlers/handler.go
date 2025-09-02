package handlers

import (
	"net/http"
	"time"

	"github.com/datdev2409/lab-admin-go/internal/storage"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5/middleware"
)

type Handler struct {
	Router http.Handler
	Store  storage.Storage
}

func NewHandler(store storage.Storage) *Handler {
	r := chi.NewRouter()

	r.Use(middleware.Timeout(60 * time.Second))

	h := &Handler{Router: r, Store: store}

	r.Use(middleware.Logger)

	// Handle static files
	r.Get("/reports/*", http.StripPrefix("/reports/", http.FileServer(http.Dir("reports"))).ServeHTTP)

	// Handle pages
	r.Route("/", func(r chi.Router) {
		r.Get("/", Make(h.HandleRecordPage))
		r.Get("/phieu-xet-nghiem", Make(h.HandleRecordPage))
		r.Get("/phieu-xet-nghiem/new", Make(h.HandleCreateNewRecord))
		r.Get("/phieu-xet-nghiem/{id}", Make(h.HandleRecordDetailPage))
		r.Get("/danh-muc-benh-nhan", Make(h.HandlePatientPage))
		r.Get("/danh-muc-xet-nghiem", Make(h.HandleTestPage))
		r.Get("/danh-muc-goi-xet-nghiem", Make(h.HandleComboPage))
		r.Get("/danh-muc-goi-xet-nghiem/new", Make(h.HandleComboCreatePage))
		r.Get("/so-sanh-ket-qua", Make(h.HandleTrackingPage))
		r.Get("/danh-muc-so-sanh", Make(h.HandleTrackingListPage))
		r.Get("/danh-muc-so-sanh/new", Make(h.HandleCreateTrackingListPage))
	})

	// Handle patients
	r.Route("/api/patients", func(r chi.Router) {
		r.Get("/{id}", Make(h.GetPatient))
		r.Delete("/{id}", Make(h.DeletePatient))
	})

	// API v1 routes
	r.Route("/api/v1/patients", func(r chi.Router) {
		r.Get("/", Make(h.ListPatientsV1))
		r.Post("/", Make(h.CreatePatientV1))
		r.Get("/{id}", Make(h.GetPatientV1))
		r.Put("/{id}", Make(h.UpdatePatientV1))
		r.Delete("/{id}", Make(h.DeletePatientV1))
		r.Get("/{id}/records", Make(h.GetPatientRecordsV1))
		r.Post("/{id}/records/compare", Make(h.ComparePatientRecordsV1))
	})

	r.Route("/api/v1/tests", func(r chi.Router) {
		r.Get("/", Make(h.ListTestsV1))
		r.Post("/", Make(h.CreateTestV1))
		r.Get("/{id}", Make(h.GetTestV1))
		r.Put("/{id}", Make(h.UpdateTestV1))
		r.Delete("/{id}", Make(h.DeleteTestV1))
	})

	// New: v1 combo routes
	r.Route("/api/v1/combos", func(r chi.Router) {
		r.Get("/", Make(h.ListCombosV1))
		r.Post("/", Make(h.CreateComboV1))
		r.Get("/{id}", Make(h.GetComboV1))
		r.Get("/{id}/tests", Make(h.GetComboTestsV1))
		r.Put("/{id}", Make(h.UpdateComboV1))
		r.Delete("/{id}", Make(h.DeleteComboV1))
	})

	// New: v1 tracking routes
	r.Route("/api/v1/trackings", func(r chi.Router) {
		r.Get("/", Make(h.ListTrackingsV1))
		r.Post("/", Make(h.CreateTrackingV1))
		r.Get("/{id}", Make(h.GetTrackingV1))
		r.Delete("/{id}", Make(h.DeleteTrackingV1))
	})

	// New: v1 record routes
	r.Route("/api/v1/records", func(r chi.Router) {
		r.Get("/", Make(h.ListRecordsV1))
		r.Post("/", Make(h.CreateRecordV1))
		r.Get("/{id}", Make(h.GetRecordV1))
		r.Put("/{id}", Make(h.UpdateRecordV1))
		r.Delete("/{id}", Make(h.DeleteRecordV1))
	})

	// Handle tests
	r.Route("/api/tests", func(r chi.Router) {
		r.Get("/", Make(h.ListTests))
		r.Post("/", Make(h.HandleCreateTest))
		r.Delete("/{id}", Make(h.DeleteTest))
	})

	r.Route("/api/combos", func(r chi.Router) {
		r.Post("/", Make(h.CreateCombo))
		r.Get("/", Make(h.ListCombos))
		r.Get("/{id}", Make(h.GetComboDetails))
	})

	r.Route("/api/records", func(r chi.Router) {
		r.Post("/", Make(h.CreateRecord))
		r.Get("/{id}", Make(h.GetRecord))
		r.Patch("/{id}", Make(h.UpdateRecord))
	})

	r.Route("/api/reports", func(r chi.Router) {
		r.Post("/export", Make(h.ExportRecord))
	})

	r.Route("/api/tracking", func(r chi.Router) {
		r.Get("/", Make(h.ListTrackings))
		r.Post("/", Make(h.CreateTracking))
		r.Post("/tests", Make(h.ListTestsForTracking))
		r.Post("/export", Make(h.CreateTrackingReport))
	})

	return h
}
