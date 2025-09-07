package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/datdev2409/lab-admin-go/internal/models"
	"github.com/datdev2409/lab-admin-go/internal/sheets"
	"github.com/datdev2409/lab-admin-go/internal/templates/pages"
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

	patient, err := h.Store.GetPatientById(r.Context(), request.PatientID)
	if err != nil {
		return err
	}

	recordTestResults := []models.TestResult{}
	for _, testResult := range request.TestResults {
		recordTestResults = append(recordTestResults, models.TestResult(testResult))
	}

	record := models.NewRecord(*patient, request.ComboName, recordTestResults)

	_, err = h.Store.InsertRecord(r.Context(), &record)
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusCreated, map[string]string{"id": record.ID})
}

func (h *Handler) UpdateRecord(w http.ResponseWriter, r *http.Request) error {
	recordId := chi.URLParam(r, "id")

	var request models.UpdateRecordRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return err
	}

	if request.PatientID != "" {
		patient, err := h.Store.GetPatientById(r.Context(), request.PatientID)
		if err != nil {
			return err
		}
		request.Patient = patient
	}

	err := h.Store.UpdateRecord(r.Context(), recordId, request)
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, map[string]string{
		"message": "Record updated successfully",
	})
}

func (h *Handler) GetRecord(w http.ResponseWriter, r *http.Request) error {
	recordId := chi.URLParam(r, "id")
	record, err := h.Store.GetRecordById(r.Context(), recordId)
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

	record, err := h.Store.GetRecordById(r.Context(), req.RecordId)
	if err != nil {
		return err
	}

	var filePath string
	var pdfPath string

	switch models.ReportType(req.ReportType) {
	case models.BillingReport:
		filePath, err = sheets.CreateRecordBillingFile(r.Context(), record)
	case models.ResultsReport:
		filePath, err = sheets.CreateRecordResultFile(r.Context(), record)
	case models.ResultsWithSignature:
		filePath, err = sheets.CreateRecordResultWithSignatureFile(r.Context(), record)
	case models.ResultsWithSignaturePDF:
		filePath, err = sheets.CreateRecordResultPDF(r.Context(), record)
		if err != nil {
			return err
		}
		pdfPath, err = sheets.ConvertExcelToPDF(r.Context(), filePath)
		if err != nil {
			return err
		}
		return WriteJSON(w, http.StatusOK, map[string]string{
			"pdf_path": pdfPath,
		})
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

// ListRecordsV1 handles GET /api/v1/records
func (h *Handler) ListRecordsV1(w http.ResponseWriter, r *http.Request) error {
	keyword := r.URL.Query().Get("q")
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
		Keyword: keyword,
		Status:  status,
	}

	genericQueryOptions := models.GenericQueryOptions{
		SortBy:    sortBy,
		SortOrder: sortOrder,
		Page:      page,
		PageSize:  pageSize,
	}

	records, pagination, err := h.Store.ListRecords(r.Context(), recordsQueryOptions, genericQueryOptions)
	if err != nil {
		return err
	}

	RespondJSONWithPagination(w, http.StatusOK, records, pagination)
	return nil
}

// CreateRecordV1 handles POST /api/v1/records
func (h *Handler) CreateRecordV1(w http.ResponseWriter, r *http.Request) error {
	var request models.CreateRecordRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return BadRequestError("invalid request body")
	}
	if request.PatientID == "" || len(request.TestResults) == 0 {
		return BadRequestError("patient_id and test_results are required")
	}
	patient, err := h.Store.GetPatientById(r.Context(), request.PatientID)
	if err != nil {
		return err
	}
	recordTestResults := []models.TestResult{}
	for _, tr := range request.TestResults {
		recordTestResults = append(recordTestResults, models.TestResult(tr))
	}
	record := models.NewRecord(*patient, request.ComboName, recordTestResults)
	id, err := h.Store.InsertRecord(r.Context(), &record)
	if err != nil {
		return err
	}
	RespondJSON(w, http.StatusCreated, map[string]string{"id": id})
	return nil
}

// GetRecordV1 handles GET /api/v1/records/{id}
func (h *Handler) GetRecordV1(w http.ResponseWriter, r *http.Request) error {
	recordId := chi.URLParam(r, "id")
	record, err := h.Store.GetRecordById(r.Context(), recordId)
	if err != nil {
		return NotFoundError("record not found")
	}
	RespondJSON(w, http.StatusOK, record)
	return nil
}

// UpdateRecordV1 handles PUT /api/v1/records/{id}
func (h *Handler) UpdateRecordV1(w http.ResponseWriter, r *http.Request) error {
	recordId := chi.URLParam(r, "id")
	var request models.UpdateRecordRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return BadRequestError("invalid request body")
	}
	if request.PatientID != "" {
		patient, err := h.Store.GetPatientById(r.Context(), request.PatientID)
		if err != nil {
			return err
		}
		request.Patient = patient
	}
	if err := h.Store.UpdateRecord(r.Context(), recordId, request); err != nil {
		return err
	}
	record, err := h.Store.GetRecordById(r.Context(), recordId)
	if err != nil {
		return err
	}
	RespondJSON(w, http.StatusOK, record)
	return nil
}

// DeleteRecordV1 handles DELETE /api/v1/records/{id}
func (h *Handler) DeleteRecordV1(w http.ResponseWriter, r *http.Request) error {
	recordId := chi.URLParam(r, "id")
	if err := h.Store.DeleteRecord(r.Context(), recordId); err != nil {
		return err
	}
	RespondJSON(w, http.StatusOK, map[string]string{"result": "deleted"})
	return nil
}
