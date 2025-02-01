package handlers

import (
	"github.com/datdev2409/lab-admin-go/internal/models"
	"github.com/datdev2409/lab-admin-go/internal/storage"
	"github.com/datdev2409/lab-admin-go/internal/templates/pages"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"net/http"
	"strings"
)

type Handler struct {
	Router http.Handler
	Store  storage.AppStorage
}

func (h *Handler) HandleComboPage(w http.ResponseWriter, r *http.Request) error {
	return Render(r.Context(), w, pages.ComboPage(""))
}

func (h *Handler) CreateCombo(w http.ResponseWriter, r *http.Request) error {
	combo := &models.Combo{
		ID:    "c" + uuid.NewString(),
		Name:  r.FormValue("combo_name"),
		Tests: strings.Split(r.FormValue("test_ids"), ","),
	}
	return h.Store.Combos().Insert(combo)
}

func NewHandler(store storage.AppStorage) *Handler {
	r := chi.NewRouter()
	h := &Handler{Router: r, Store: store}

	r.Use(middleware.Logger)

	// Handle pages
	r.Route("/", func(r chi.Router) {
		r.Get("/", Make(HandleRecordPage))
		r.Get("/phieu-xet-nghiem", Make(HandleRecordPage))
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
	})

	return h
}
