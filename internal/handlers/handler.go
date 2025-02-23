package handlers

import (
	"log"
	"net/http"
	"strings"

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
		r.Get("/search", Make(h.SearchRecordsByPatientNameOrPhone))
		r.Get("/", Make(h.ListRecords))
		r.Patch("/{id}", Make(h.UpdateRecordPatient))
		r.Patch("/{id}/patients", Make(h.UpdateRecordPatient))
		r.Patch("/{id}/combos", Make(h.UpdateRecordCombo))
		r.Post("/{id}/tests", Make(h.AddTestToRecord))
		r.Patch("/{id}/tests", Make(h.UpdateRecordTests))
		r.Post("/{id}/export/billing", Make(h.ExportRecordBilling))
		r.Post("/{id}/export/results", Make(h.ExportRecordResults))
	})

	r.Route("/api/tracking", func(r chi.Router) {
		r.Post("/export", Make(h.ExportTracking))
	})

	return h
}

func (h *Handler) HandleTrackingPage(w http.ResponseWriter, r *http.Request) error {
	return Render(r.Context(), w, pages.TrackingPage(""))
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

func (h *Handler) ListRecords(w http.ResponseWriter, r *http.Request) error {
	patientId := r.URL.Query().Get("patient_id")
	records, err := h.Store.Records().ListByPatientId(r.Context(), patientId)
	if err != nil {
		return nil
	}

	log.Println(patientId)
	view.RecordList(*records).Render(w)
	return nil
	// records, err := h.Store.Records().GetAll(r.Context())
	// if err != nil {
	// 	return err
	// }

	// return view.RecordTable(*records).Render(w)
}

func (h *Handler) ExportRecordResults(w http.ResponseWriter, r *http.Request) error {
	recordId := chi.URLParam(r, "id")

	record, err := h.Store.Records().GetDetails(r.Context(), recordId)
	if err != nil {
		return err
	}

	filepath, err := sheets.CreateRecordResultFile(*record)
	if err != nil {
		return err
	}

	HTMXRedirect(w, "/"+filepath)
	return nil
}

func (h *Handler) ExportRecordBilling(w http.ResponseWriter, r *http.Request) error {
	recordId := chi.URLParam(r, "id")

	record, err := h.Store.Records().GetDetails(r.Context(), recordId)
	if err != nil {
		return err
	}

	filepath, err := sheets.CreateRecordBillingFile(*record)
	if err != nil {
		return err
	}

	HTMXRedirect(w, "/"+filepath)
	return nil
}
