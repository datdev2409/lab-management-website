package handlers

import (
	"net/http"
	"time"

	"github.com/datdev2409/lab-admin-go/internal/storage"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

type Handler struct {
	Router         http.Handler
	Validator      *validator.Validate
	Store          storage.Storage
	patientHandler *PatientHandler
}

func NewHandler(store storage.Storage, log *zap.Logger, patientHandler *PatientHandler) *Handler {
	r := chi.NewRouter()

	r.Use(middleware.Timeout(60 * time.Second))

	v := validator.New()
	h := &Handler{Router: r, Store: store, Validator: v, patientHandler: patientHandler}

	r.Use(middleware.RequestID)
	// r.Use(requestLogger)
	r.Use(LoggingMiddleware(log))
	r.Use(HTTPLogger)

	// Handle static files
	r.Get("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))).ServeHTTP)
	r.Get("/reports/*", http.StripPrefix("/reports/", http.FileServer(http.Dir("reports"))).ServeHTTP)

	// Health check endpoint
	r.Get("/health", Make(h.HandleHealthCheck))

	r.Get("/login", Make(h.HandleLoginPage))
	r.Get("/register", Make(h.HandleRegisterPage))

	// Handle pages
	r.Route("/", func(r chi.Router) {
		// r.Use(JWTAuthWebEndpoint)
		r.Get("/", Make(h.HandleRecordPage))
		r.Get("/phieu-xet-nghiem", Make(h.HandleRecordPage))
		r.Get("/phieu-xet-nghiem/new", Make(h.HandleCreateNewRecord))
		r.Get("/phieu-xet-nghiem/{id}", Make(h.HandleRecordDetailPage))
		r.Get("/danh-muc-benh-nhan", Make(h.HandlePatientPage))
		// Doctor page now handled in main.go
		// r.Get("/danh-muc-bac-si", Make(h.HandleDoctorPage))
		r.Get("/danh-muc-xet-nghiem", Make(h.HandleTestPage))
		r.Get("/danh-muc-goi-xet-nghiem", Make(h.HandleComboPage))
		r.Get("/danh-muc-goi-xet-nghiem/new", Make(h.HandleComboCreatePage))
		r.Get("/danh-muc-goi-xet-nghiem/{id}/edit", Make(h.HandleComboEditPage))
		r.Get("/so-sanh-ket-qua", Make(h.HandleTrackingPage))
		r.Get("/danh-muc-so-sanh", Make(h.HandleTrackingListPage))
		r.Get("/danh-muc-so-sanh/new", Make(h.HandleCreateTrackingListPage))
		r.Get("/bao-cao-thong-ke-doanh-so", Make(h.HandleReportPage))
	})

	// Handle patients
	r.Route("/api/patients", func(r chi.Router) {
		r.Get("/{id}", Make(h.patientHandler.GetPatient))
		r.Delete("/{id}", Make(h.patientHandler.DeletePatient))
	})

	r.Route("/api/v1/auth", func(r chi.Router) {
		r.Post("/register", Make(h.RegisterHandler))
		r.Post("/login", Make(h.LoginHandler))
		r.Post("/logout", Make(h.LogoutHandler))
	})

	// API v1 routes
	r.Route("/api/v1", func(r chi.Router) {
		// r.Use(JWTAuthAPIEndpoint)
		r.Route("/patients", func(r chi.Router) {
			r.Get("/", Make(h.patientHandler.SearchPatientsByKeyword))
			r.Post("/", Make(h.patientHandler.CreatePatient))
			r.Get("/{id}", Make(h.patientHandler.GetPatient))
			r.Patch("/{id}", Make(h.UpdatePatientV1))
			r.Delete("/{id}", Make(h.patientHandler.DeletePatient))
			r.Get("/{id}/records", Make(h.GetPatientRecordsV1))
			r.Post("/{id}/records/compare", Make(h.ComparePatientRecordsV1))
		})

		// Doctor routes now handled in main.go
		// r.Route("/doctors", func(r chi.Router) {
		// 	r.Get("/", Make(h.ListDoctorsV1))
		// 	r.Post("/", Make(h.CreateDoctorV1))
		// 	r.Get("/{id}", Make(h.GetDoctorV1))
		// 	r.Put("/{id}", Make(h.UpdateDoctorV1))
		// 	r.Delete("/{id}", Make(h.DeleteDoctorV1))
		// })

		r.Route("/tests", func(r chi.Router) {
			r.Get("/", Make(h.ListTestsV1))
			r.Post("/", Make(h.CreateTestV1))
			r.Get("/{id}", Make(h.GetTestV1))
			r.Put("/{id}", Make(h.UpdateTestV1))
			r.Delete("/{id}", Make(h.DeleteTestV1))
		})

		r.Route("/combos", func(r chi.Router) {
			r.Get("/", Make(h.ListCombosV1))
			r.Post("/", Make(h.CreateComboV1))
			r.Get("/{id}", Make(h.GetComboV1))
			r.Get("/{id}/tests", Make(h.GetComboTestsV1))
			r.Put("/{id}", Make(h.UpdateComboV1))
			r.Delete("/{id}", Make(h.DeleteComboV1))
		})

		r.Route("/trackings", func(r chi.Router) {
			r.Get("/", Make(h.ListTrackingsV1))
			r.Post("/", Make(h.CreateTrackingV1))
			r.Get("/{id}", Make(h.GetTrackingV1))
			r.Delete("/{id}", Make(h.DeleteTrackingV1))
		})

		r.Route("/records", func(r chi.Router) {
			r.Get("/", Make(h.ListRecordsV1))
			r.Post("/", Make(h.CreateRecordV1))
			r.Get("/{id}", Make(h.GetRecordV1))
			r.Put("/{id}", Make(h.UpdateRecordV1))
			r.Delete("/{id}", Make(h.DeleteRecordV1))
			r.Get("/revenue", Make(h.GetRecordsWithRevenue))
			r.Get("/revenue/export", Make(h.ExportRecordsRevenue))
		})
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

// HandleHealthCheck provides a health check endpoint
func (h *Handler) HandleHealthCheck(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(`{"status":"ok","service":"lab-admin-go"}`))
	return err
}
