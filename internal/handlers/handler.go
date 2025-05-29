package handlers

import (
	"net/http"
	"time"

	"github.com/datdev2409/lab-admin-go/internal/storage"
	"github.com/datdev2409/lab-admin-go/internal/view"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5/middleware"
)

type Handler struct {
	Router http.Handler
	Store  storage.AppStorage
}

func NewHandler(store storage.AppStorage) *Handler {
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
		r.Get("/danh-muc-so-sanh", func(w http.ResponseWriter, r *http.Request) {
			view.CompareResultsPage(view.PageProps{}).Render(w)
		})
	})

	// Handle patients
	r.Route("/api/patients", func(r chi.Router) {
		r.Post("/", Make(h.HandleCreatePatient))
		r.Get("/", Make(h.ListPatients))
		r.Get("/{id}", Make(h.GetPatient))
		r.Put("/{id}", Make(h.UpdatePatient))
		r.Delete("/{id}", Make(h.DeletePatient))
	})

	// Handle tests
	r.Route("/api/tests", func(r chi.Router) {
		r.Get("/", Make(h.ListTests))
		r.Get("/search", Make(h.SearchTestsByKeyword))
		r.Post("/", Make(h.HandleCreateTest))
	})

	r.Route("/api/combos", func(r chi.Router) {
		r.Post("/", Make(h.CreateCombo))
		r.Get("/", Make(h.ListCombos))
		r.Get("/{id}", Make(h.GetComboDetails))
	})

	r.Route("/api/records", func(r chi.Router) {
		r.Post("/", Make(h.CreateRecord))
		r.Get("/", Make(h.ListRecords))
		r.Get("/{id}", Make(h.GetRecord))
		r.Patch("/{id}", Make(h.UpdateRecord))
	})

	r.Route("/api/reports", func(r chi.Router) {
		r.Post("/export", Make(h.ExportRecord))
	})

	r.Route("/api/tracking", func(r chi.Router) {
		r.Post("/", Make(h.CreateTrackingReport))
	})

	return h
}
