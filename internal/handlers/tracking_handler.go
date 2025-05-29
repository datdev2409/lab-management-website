package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/datdev2409/lab-admin-go/internal/models"
	"github.com/datdev2409/lab-admin-go/internal/sheets"
	"github.com/datdev2409/lab-admin-go/internal/templates/pages"
)

func (h *Handler) HandleTrackingPage(w http.ResponseWriter, r *http.Request) error {
	return Render(r.Context(), w, pages.TrackingPage())
}

func (h *Handler) CreateTrackingReport(w http.ResponseWriter, r *http.Request) error {
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

	filename, err := sheets.CreateRecordTrackingFile(records, testMap)
	if err != nil {
		log.Println("Error creating tracking file:", err)
		return err
	}

	HTMXRedirect(w, filename)
	return nil
}
