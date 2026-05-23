package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	_ "time/tzdata"

	"github.com/datdev2409/lab-admin-go/internal/logger"
	"github.com/datdev2409/lab-admin-go/internal/models"
	"github.com/datdev2409/lab-admin-go/internal/sheets"
	"github.com/datdev2409/lab-admin-go/internal/templates/pages"
	"github.com/go-chi/chi"
	"go.uber.org/zap"
)

// validateAndSetDoctor validates doctor ID and sets doctor name in the request
func (h *Handler) validateAndSetDoctor(ctx context.Context, doctorID string) (string, error) {
	if doctorID == "" {
		return "", nil
	}

	doctor, err := h.Store.GetDoctorById(ctx, doctorID)
	if err != nil {
		logger.FromCtx(ctx).Error("Doctor not found", zap.String("doctor_id", doctorID), zap.Error(err))
		return "", BadRequestError(DOCTOR_NOT_FOUND_ERROR)
	}

	// Return the doctor's name from the database to ensure consistency
	return doctor.Name, nil
}

// parseRevenueReportFilters parses date range and query options from request for revenue reports
func parseRevenueReportFilters(r *http.Request) (models.RecordQueryOptions, models.GenericQueryOptions, error) {
	ctx := r.Context()
	log := logger.FromCtx(ctx)

	filters := models.RecordQueryOptions{}

	// Parse start_date
	if startDateStr := r.URL.Query().Get("start_date"); startDateStr != "" {
		startDate, err := ParseStartOfDayInVietnamTimezone(startDateStr)
		if err != nil {
			log.Warn("Invalid start_date format", zap.String("start_date", startDateStr), zap.Error(err))
		} else {
			filters.StartDate = startDate
		}
	}

	// Parse end_date
	if endDateStr := r.URL.Query().Get("end_date"); endDateStr != "" {
		endDate, err := ParseEndOfDayInVietnamTimezone(endDateStr)
		if err != nil {
			log.Warn("Invalid end_date format", zap.String("end_date", endDateStr), zap.Error(err))
		} else {
			filters.EndDate = endDate
		}
	}

	// Parse doctor_id (optional filter)
	if doctorID := r.URL.Query().Get("doctor_id"); doctorID != "" {
		filters.DoctorID = doctorID
	}

	// Parse test filters
	if testID := r.URL.Query().Get("test_id"); testID != "" {
		filters.TestID = testID
	}

	// Parse optional sorting parameters
	opts := models.GenericQueryOptions{
		Page:      1,
		PageSize:  0, // No pagination for reports
		SortBy:    r.URL.Query().Get("sort_by"),
		SortOrder: r.URL.Query().Get("sort_order"),
	}

	// Default sorting by created_at desc if not specified
	if opts.SortBy == "" {
		opts.SortBy = "created_at"
		opts.SortOrder = "desc"
	}

	return filters, opts, nil
}

// getRevenueReportData is a helper function to fetch revenue report data with common error handling
func (h *Handler) getRevenueReportData(ctx context.Context, filters models.RecordQueryOptions, opts models.GenericQueryOptions) (*models.ReportResponse, error) {
	log := logger.FromCtx(ctx)

	reportData, err := h.Store.GetRecordsWithRevenue(ctx, filters, opts)
	if err != nil {
		log.Error("Failed to get records with revenue", zap.Error(err))
		return nil, err
	}

	log.Info("Revenue report data retrieved successfully",
		zap.Int("record_count", len(reportData.Records)),
		zap.Int("total_revenue", reportData.Summary.TotalRevenue))

	return reportData, nil
}

// generateRevenueReportFilename creates a filename for revenue report based on date filters
func generateRevenueReportFilename(filters models.RecordQueryOptions) string {
	baseFilename := "revenue_report"

	if filters.StartDate != nil && filters.EndDate != nil {
		startDateStr := filters.StartDate.Format("2006-01-02")
		endDateStr := filters.EndDate.Format("2006-01-02")
		return fmt.Sprintf("%s_%s_to_%s.xlsx", baseFilename, startDateStr, endDateStr)
	} else if filters.StartDate != nil {
		startDateStr := filters.StartDate.Format("2006-01-02")
		return fmt.Sprintf("%s_from_%s.xlsx", baseFilename, startDateStr)
	} else if filters.EndDate != nil {
		endDateStr := filters.EndDate.Format("2006-01-02")
		return fmt.Sprintf("%s_to_%s.xlsx", baseFilename, endDateStr)
	}

	return baseFilename + ".xlsx"
}

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

	reportGenerator, err := sheets.NewReportGenerator(r.Context(), models.ReportType(req.ReportType))
	if err != nil {
		return err
	}

	reader, err := reportGenerator.Generate(r.Context(), record)
	if err != nil {
		return err
	}

	storer := sheets.LocalFileStoreStrategy{
		BaseDir: "./reports",
	}

	fileName := fmt.Sprintf("%s-%s-%s.xlsx", time.Now().Format("20060102-150405"), strings.ReplaceAll(record.Patient.Name, " ", "_"), req.ReportType)
	filePath, err := storer.Store(r.Context(), reader, fileName)
	if err != nil {
		return err
	}

	if models.ReportType(req.ReportType) == models.ResultsWithSignaturePDF {
		pdfPath, err := sheets.ConvertExcelToPDF(r.Context(), filePath)
		if err != nil {
			return err
		}

		return WriteJSON(w, http.StatusOK, map[string]string{
			"pdf_path": pdfPath,
		})
	}

	return WriteJSON(w, http.StatusOK, map[string]string{
		"excel_path": filePath,
	})
}

// ListRecordsV1 handles GET /api/v1/records
func (h *Handler) ListRecordsV1(w http.ResponseWriter, r *http.Request) error {
	keyword := r.URL.Query().Get("q")
	status := r.URL.Query().Get("status")
	doctorID := r.URL.Query().Get("doctor_id")

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
		Keyword:  keyword,
		Status:   status,
		DoctorID: doctorID,
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

	// Validate doctor if provided
	doctorName, err := h.validateAndSetDoctor(r.Context(), request.DoctorID)
	if err != nil {
		return err
	}
	request.DoctorName = doctorName

	recordTestResults := []models.TestResult{}
	for _, tr := range request.TestResults {
		recordTestResults = append(recordTestResults, models.TestResult(tr))
	}

	var record models.Record
	if request.DoctorID != "" && request.DoctorName != "" {
		record = models.NewRecordWithDoctor(*patient, request.ComboName, recordTestResults, request.DoctorID, request.DoctorName)
	} else {
		record = models.NewRecord(*patient, request.ComboName, recordTestResults)
	}

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

	// Validate doctor if provided
	doctorName, err := h.validateAndSetDoctor(r.Context(), request.DoctorID)
	if err != nil {
		return err
	}
	request.DoctorName = doctorName

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

// GetRecordsWithRevenue returns revenue report data based on date range filters
func (h *Handler) GetRecordsWithRevenue(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	log := logger.FromCtx(ctx)

	// Parse query parameters
	filters, opts, err := parseRevenueReportFilters(r)
	if err != nil {
		log.Error("Failed to parse revenue report filters", zap.Error(err))
		return WriteJSON(w, http.StatusBadRequest, map[string]string{
			"error": "Invalid query parameters",
		})
	}

	// Get report data using helper
	reportData, err := h.getRevenueReportData(ctx, filters, opts)
	if err != nil {
		return WriteJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "Failed to get records with revenue",
		})
	}

	return WriteJSON(w, http.StatusOK, reportData)
}

// ExportRecordsRevenue exports revenue report data to Excel format
func (h *Handler) ExportRecordsRevenue(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	log := logger.FromCtx(ctx)

	// Parse query parameters using shared helper
	filters, opts, err := parseRevenueReportFilters(r)
	if err != nil {
		log.Error("Failed to parse revenue report filters", zap.Error(err))
		return WriteJSON(w, http.StatusBadRequest, map[string]string{
			"error": "Invalid query parameters",
		})
	}

	// Get report data using helper
	reportData, err := h.getRevenueReportData(ctx, filters, opts)
	if err != nil {
		return WriteJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "Failed to export revenue report",
		})
	}

	// Generate Excel report
	generator, err := sheets.NewReportGenerator(ctx, models.RevenueReport)
	if err != nil {
		log.Error("Failed to create revenue export report generator", zap.Error(err))
		return WriteJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "Failed to generate export",
		})
	}

	excelFile, err := generator.Generate(ctx, reportData)
	if err != nil {
		log.Error("Failed to generate Excel file", zap.Error(err))
		return WriteJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "Failed to generate Excel file",
		})
	}

	// Generate filename using helper
	filename := generateRevenueReportFilename(filters)

	// Set response headers for file download
	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))

	// Stream the file to response
	_, err = io.Copy(w, excelFile)
	if err != nil {
		log.Error("Failed to write Excel file to response", zap.Error(err))
		return err
	}

	return nil
}
