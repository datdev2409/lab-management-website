package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/a-h/templ"
	"github.com/datdev2409/lab-admin-go/internal/models"
	"github.com/datdev2409/lab-admin-go/internal/sheets"
	"github.com/datdev2409/lab-admin-go/internal/templates/pages"
	"github.com/datdev2409/lab-admin-go/internal/templates/partials"
	"github.com/go-chi/chi"
)

func (h *Handler) HandleRecordPage(w http.ResponseWriter, r *http.Request) error {
	return Render(context.Background(), w, pages.RecordPage())
}

func (h *Handler) HandleCreateNewRecord(w http.ResponseWriter, r *http.Request) error {
	return Render(r.Context(), w, pages.RecordCreatePage())
}

func (h *Handler) HandleRecordDetailPage(w http.ResponseWriter, r *http.Request) error {
	recordId := chi.URLParam(r, "id")
	return Render(r.Context(), w, pages.RecordDetailsPage(recordId))
}

func (h *Handler) CreateRecord(w http.ResponseWriter, r *http.Request) error {
	var request models.CreateRecordRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return err
	}

	patient, err := h.StoreV2.GetPatientById(r.Context(), request.PatientID)
	if err != nil {
		return err
	}

	recordTestResults := []models.TestResult{}
	for _, testResult := range request.TestResults {
		recordTestResults = append(recordTestResults, models.TestResult{
			ID:          testResult.ID,
			Name:        testResult.Name,
			Price:       testResult.Price,
			NormalValue: testResult.NormalValue,
			Unit:        testResult.Unit,
			LowerBound:  testResult.LowerBound,
			UpperBound:  testResult.UpperBound,
			Result:      testResult.Result,
			ResultText:  testResult.ResultText,
		})
	}

	record := models.NewRecord(*patient, request.ComboName, recordTestResults)

	_, err = h.StoreV2.InsertRecord(r.Context(), &record)
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusCreated, map[string]string{"id": record.ID})
}

func (h *Handler) ListRecords(w http.ResponseWriter, r *http.Request) error {
	patientId := r.URL.Query().Get("patient_id")
	keyword := r.URL.Query().Get("keyword")
	status := r.URL.Query().Get("status")

	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		page = 1
	}
	pageSize, err := strconv.Atoi(r.URL.Query().Get("page_size"))
	if err != nil {
		pageSize = 20
	}

	sortBy := r.URL.Query().Get("sort_by")
	if sortBy == "" {
		sortBy = "created_at"
	}
	sortOrder := r.URL.Query().Get("sort_order")
	if sortOrder == "" {
		sortOrder = "desc"
	}

	recordsQueryOptions := models.RecordQueryOptions{
		Keyword:   keyword,
		Status:    status,
		PatientID: patientId,
	}

	genericQueryOptions := models.GenericQueryOptions{
		SortBy:    sortBy,
		SortOrder: sortOrder,
		Page:      page,
		PageSize:  pageSize,
	}

	records, pagination, err := h.StoreV2.ListRecords(r.Context(), recordsQueryOptions, genericQueryOptions)
	if err != nil {
		return err
	}

	target := r.Header.Get("HX-Target")
	switch target {
	case "tracking-record-list":
		return Render(r.Context(), w, partials.TrackingRecordTable(records))
	default:
		RenderMultiComponents(r.Context(), w, []templ.Component{
			partials.RecordTable(records),
			partials.Pagination(pagination, "record-page"),
		})
	}
	return nil
}

func (h *Handler) UpdateRecord(w http.ResponseWriter, r *http.Request) error {
	recordId := chi.URLParam(r, "id")

	var request models.UpdateRecordRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Println("Error while decoding request", err)
		return err
	}

	if request.PatientID != "" {
		patient, err := h.StoreV2.GetPatientById(r.Context(), request.PatientID)
		if err != nil {
			return err
		}
		request.Patient = patient
	}

	err := h.StoreV2.UpdateRecord(r.Context(), recordId, request)
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, map[string]string{
		"message": "Record updated successfully",
	})
}

func (h *Handler) GetRecord(w http.ResponseWriter, r *http.Request) error {
	recordId := chi.URLParam(r, "id")
	record, err := h.StoreV2.GetRecordById(r.Context(), recordId)
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, record)
}

func (h *Handler) ExportRecord(w http.ResponseWriter, r *http.Request) error {
	var req struct {
		RecordId   string `json:"record_id"`
		ReportType string `json:"report_type"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}

	record, err := h.StoreV2.GetRecordById(r.Context(), req.RecordId)
	if err != nil {
		return err
	}

	var filePath string

	switch models.ReportType(req.ReportType) {
	case models.BillingReport:
		filePath, err = sheets.CreateRecordBillingFile(record)
	case models.ResultsReport:
		filePath, err = sheets.CreateRecordResultFile(record)
	case models.ResultsWithSignature:
		filePath, err = sheets.CreateRecordResultWithSignatureFile(record)
	default:
		return errors.New("invalid export type")
	}

	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, map[string]string{
		"excel_path": filePath,
	})
}
