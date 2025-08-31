package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/a-h/templ"
	"github.com/datdev2409/lab-admin-go/internal/models"
	"github.com/datdev2409/lab-admin-go/internal/sheets"
	"github.com/datdev2409/lab-admin-go/internal/templates/pages"
	"github.com/datdev2409/lab-admin-go/internal/templates/partials"
)

func (h *Handler) HandleTrackingPage(w http.ResponseWriter, r *http.Request) error {
	return Render(r.Context(), w, pages.TrackingPage())
}

func (h *Handler) HandleTrackingListPage(w http.ResponseWriter, r *http.Request) error {
	return Render(r.Context(), w, pages.TrackingListPage())
}

func (h *Handler) HandleCreateTrackingListPage(w http.ResponseWriter, r *http.Request) error {
	return Render(r.Context(), w, pages.TrackingCreatePage())
}

func (h *Handler) ListTrackings(w http.ResponseWriter, r *http.Request) error {
	keyword := r.URL.Query().Get("keyword")
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		page = 1
	}
	pageSize, err := strconv.Atoi(r.URL.Query().Get("page_size"))
	if err != nil {
		pageSize = 10
	}
	trackings, pagination, err := h.Store.ListTrackings(r.Context(), models.TrackingQueryOptions{Keyword: keyword}, models.GenericQueryOptions{Page: page, PageSize: pageSize})
	if err != nil {
		return err
	}

	target := r.Header.Get("HX-Target")
	slog.Debug("HTMX Target:", slog.String("target", target))

	switch target {
	case "tracking-autocomplete":
		return Render(r.Context(), w, partials.TrackingAutocomplete(trackings))
	default:
		return RenderMultiComponents(r.Context(), w, []templ.Component{
			partials.TrackingTable(trackings),
			partials.Pagination(pagination, "tracking-page"),
		})
	}

}

func (h *Handler) ListTestsForTracking(w http.ResponseWriter, r *http.Request) error {
	r.ParseForm()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("failed to read request body: %v", err)
	}
	log.Println(string(body))

	records, err := h.Store.GetRecordsByIds(r.Context(), r.Form["record_ids"])
	if err != nil {
		return err
	}

	// Collect all unique tests
	tests := []*models.TestInfo{}
	existingTests := make(map[string]bool)

	for _, record := range records {
		for _, test := range record.TestResults {
			if _, exists := existingTests[test.Name]; !exists {
				existingTests[test.Name] = true
				tests = append(tests, &models.TestInfo{
					Name:        test.Name,
					NormalValue: test.NormalValue,
					Unit:        test.Unit,
				})
			}
		}
	}

	return Render(r.Context(), w, partials.TrackingTestTable(tests))
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
		record, err := h.Store.GetRecordById(r.Context(), id)
		if err != nil {
			return fmt.Errorf("failed to fetch record %s: %v", id, err)
		}
		records = append(records, record)
	}

	// Handle custom tracking ID
	testMap := make(map[string]models.TestInfo)
	trackingId := r.FormValue("tracking_id")
	if trackingId == "" {
		for _, record := range records {
			for _, test := range record.TestResults {
				testMap[test.Name] = models.TestInfo{
					Name:        test.Name,
					NormalValue: test.NormalValue,
					Unit:        test.Unit,
				}
			}
		}
	} else {
		tracking, err := h.Store.GetTrackingById(r.Context(), trackingId)
		if err != nil {
			return err
		}

		for _, test := range tracking.Tests {
			testMap[test.TestName] = models.TestInfo{
				Name:        test.TestName,
				NormalValue: test.NormalValue,
				Unit:        test.Unit,
			}
		}

	}

	filename, err := sheets.CreateRecordTrackingFile(records, testMap)
	if err != nil {
		return err
	}

	HTMXRedirect(w, filename)
	return nil
}

func (h *Handler) CreateTracking(w http.ResponseWriter, r *http.Request) error {
	var request models.CreateTrackingRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return err
	}

	tracking := models.NewTracking(request.Name, request.Tests)

	_, err := h.Store.InsertTracking(r.Context(), &tracking)
	if err != nil {
		return err
	}

	return nil
}
