package handlers

import (
	"net/http"

	"github.com/datdev2409/lab-admin-go/internal/storage"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5/middleware"
)

type Handler struct {
	Router http.Handler
	Store  storage.AppStorage
}

func NewHandler(store storage.AppStorage) *Handler {
	r := chi.NewRouter()
	h := &Handler{Router: r, Store: store}

	r.Use(middleware.Logger)

	// Handle pages
	r.Route("/", func(r chi.Router) {
		r.Get("/", Make(h.HandleRecordPage))
		r.Get("/phieu-xet-nghiem", Make(h.HandleRecordPage))
		r.Get("/phieu-xet-nghiem/new", Make(h.HandleCreateNewRecord))
		r.Get("/phieu-xet-nghiem/{id}", Make(h.HandleRecordDetailPage))
		r.Get("/danh-muc-benh-nhan", Make(h.HandlePatientPage))
		r.Get("/danh-muc-xet-nghiem", Make(h.HandleTestPage))
		r.Get("/danh-muc-goi-xet-nghiem", Make(h.HandleComboPage))
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
		r.Get("/", Make(h.SearchTestsByKeyword))
		//r.Get("/{id}", Make(h.SelectTest))
		r.Post("/", Make(h.HandleCreateTest))
	})

	r.Route("/api/combos", func(r chi.Router) {
		r.Post("/", Make(h.CreateCombo))
		r.Get("/", Make(h.SearchCombosByKeyword))
		r.Get("/{id}", Make(h.GetCombo))
	})

	r.Route("/api/records", func(r chi.Router) {
		r.Patch("/{id}", Make(h.UpdateRecordPatient))
		r.Get("/{id}/tests", Make(h.GetRecordTests))
	})

	return h
}
