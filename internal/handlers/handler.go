package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/datdev2409/lab-admin-go/internal/storage"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/xuri/excelize/v2"
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
		r.Get("/", Make(h.SearchRecordsByPatientNameOrPhone))
		r.Patch("/{id}", Make(h.UpdateRecordPatient))
		r.Patch("/{id}/patients", Make(h.UpdateRecordPatient))
		r.Patch("/{id}/combos", Make(h.UpdateRecordCombo))
		r.Patch("/{id}/tests", Make(h.UpdateRecordTest))
		r.Post("/{id}/export/billing", Make(h.ExportRecordBilling))
	})

	return h
}

func (h *Handler) ExportRecordBilling(w http.ResponseWriter, r *http.Request) error {
	recordId := chi.URLParam(r, "id")

	f, err := excelize.OpenFile("report-templates/billing.xlsx")
	if err != nil {
		return err
	}
	defer f.Close()

	record, err := h.Store.Records().GetDetails(r.Context(), recordId)
	if err != nil {
		return err
	}

	borderStyle, _ := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
	})

	f.SetCellValue("Sheet1", "B6", record.Patient.Name)
	f.SetCellValue("Sheet1", "B7", record.Patient.Address)
	f.SetCellValue("Sheet1", "D6", record.Patient.YOB)

	startTestRow := 10
	for i, testResult := range record.TestResults {
		row := startTestRow + i
		testInfo := record.TestInfoMap[testResult.TestID]
		f.SetCellValue("Sheet1", fmt.Sprintf("A%d", row), i+1)
		f.SetCellValue("Sheet1", fmt.Sprintf("B%d", row), testInfo.Name)
		f.SetCellValue("Sheet1", fmt.Sprintf("C%d", row), 1)
		f.SetCellValue("Sheet1", fmt.Sprintf("D%d", row), testInfo.Price)
		f.SetCellValue("Sheet1", fmt.Sprintf("E%d", row), testInfo.Price)
	}

	endTestRow := startTestRow + len(record.TestResults) - 1
	startTestCell := fmt.Sprintf("A%d", startTestRow)
	endTestCell := fmt.Sprintf("E%d", endTestRow)
	f.SetCellStyle("Sheet1", startTestCell, endTestCell, borderStyle)

	now := time.Now().Format("2006-01-02")
	filename := fmt.Sprintf("reports/%s-billing-%s.xlsx", recordId, now)
	if err := f.SaveAs(filename); err != nil {
		return err
	}
	return nil
}
