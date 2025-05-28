package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/datdev2409/lab-admin-go/internal/models"
	"github.com/datdev2409/lab-admin-go/internal/sheets"
	"github.com/datdev2409/lab-admin-go/internal/storage"
	"github.com/datdev2409/lab-admin-go/internal/templates/pages"
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
		r.Get("/danh-muc-so-sanh", Make(h.HandleTrackingPage))
		r.Get("/so-sanh-ket-qua", func(w http.ResponseWriter, r *http.Request) {
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
		r.Post("/", Make(h.CreateTracking))
		r.Post("/export", Make(h.ExportTracking))
	})

	return h
}

func (h *Handler) HandleTrackingPage(w http.ResponseWriter, r *http.Request) error {
	return Render(r.Context(), w, pages.TrackingPage())
}

func (h *Handler) ExportTracking(w http.ResponseWriter, r *http.Request) error {
	r.ParseForm()
	var compareRecordIds []string
	for key, values := range r.Form {
		if strings.HasPrefix(key, "record_id_") {
			compareRecordIds = append(compareRecordIds, values[0])
		}
	}
	log.Printf("Selected records: %v", compareRecordIds)
	return nil
}

func (h *Handler) CreateTracking(w http.ResponseWriter, r *http.Request) error {
	r.ParseForm()

	recordIds := r.Form["record_ids"]

	if len(recordIds) == 0 {
		return fmt.Errorf("no record IDs provided for comparison")
	}

	// Fetch records from storage
	var records []*models.Record
	for _, id := range recordIds {
		record, err := h.Store.Records().GetById(r.Context(), id)
		if err != nil {
			return fmt.Errorf("failed to fetch record %s: %v", id, err)
		}
		records = append(records, record)
	}

	// Collect all unique tests
	testMap := make(map[string]models.TestInfo)
	for _, record := range records {
		for _, test := range record.TestResults {
			testMap[test.Name] = models.TestInfo{
				Name:        test.Name,
				NormalValue: test.NormalValue,
			}
		}
	}

	log.Println(testMap)

	filename, err := sheets.CreateRecordTrackingFile(records, testMap)
	if err != nil {
		log.Println("Error creating tracking file:", err)
		return err
	}

	HTMXRedirect(w, filename)
	return nil
}
