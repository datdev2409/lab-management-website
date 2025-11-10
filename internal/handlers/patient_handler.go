package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/datdev2409/lab-admin-go/internal/logger"
	"github.com/datdev2409/lab-admin-go/internal/models"
	"github.com/datdev2409/lab-admin-go/internal/service"
	"github.com/datdev2409/lab-admin-go/internal/sheets"
	"github.com/datdev2409/lab-admin-go/internal/templates/pages"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

type PatientHandler struct {
	patientService *service.PatientService
	validator      *validator.Validate
}

func NewPatientHandler(patientService *service.PatientService, validator *validator.Validate) *PatientHandler {
	return &PatientHandler{
		patientService: patientService,
		validator:      validator,
	}
}

func (h *PatientHandler) CreatePatient(w http.ResponseWriter, r *http.Request) error {
	var input models.CreatePatientInput

	if err := BindAndValidate(r, h.validator, &input); err != nil {
		return err
	}

	patient, err := h.patientService.CreatePatient(r.Context(), &input)
	if err != nil {
		if errors.Is(err, service.ErrPatientAlreadyExists) {
			return &AppError{StatusCode: http.StatusConflict, Message: PATIENT_ALREADY_EXISTS}
		}
		return err
	}

	RespondJSON(w, http.StatusCreated, patient)
	return nil
}

func (h *PatientHandler) SearchPatientsByKeyword(w http.ResponseWriter, r *http.Request) error {
	queryOpts := ParseListParams(r, 10)

	keyword := r.URL.Query().Get("q")
	patients, pagination, err := h.patientService.SearchPatientsByKeyword(r.Context(), keyword, queryOpts.Page, queryOpts.PageSize)
	if err != nil {
		return err
	}

	RespondJSONWithPagination(w, http.StatusOK, patients, pagination)
	return nil
}

func (h *PatientHandler) GetPatient(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")
	patient, err := h.patientService.GetPatientById(r.Context(), id)
	if err != nil {
		return err
	}

	RespondJSON(w, http.StatusOK, patient)
	return nil
}

func (h *PatientHandler) UpdatePatient(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")

	var req models.PatientUpdate
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return &AppError{StatusCode: http.StatusBadRequest, Message: "Invalid request body"}
	}

	patient, err := h.patientService.UpdatePatientById(r.Context(), id, req)
	if err != nil {
		return err
	}

	RespondJSON(w, http.StatusOK, patient)
	return nil
}

func (h *PatientHandler) DeletePatient(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")
	err := h.patientService.DeletePatientById(r.Context(), id)
	if err != nil {
		return err
	}

	RespondJSON(w, http.StatusOK, map[string]string{"result": "deleted"})
	return nil
}

func (h *PatientHandler) HandlePatientPage(w http.ResponseWriter, r *http.Request) error {
	return Render(r.Context(), w, pages.PatientsPage())
}

// Keep this on Handler for now - it will be removed once all pages are migrated
func (h *Handler) HandlePatientPage(w http.ResponseWriter, r *http.Request) error {
	return Render(r.Context(), w, pages.PatientsPage())
}

// UpdatePatientV1 handles PUT /api/v1/patients/{id}
// NOTE: This method stays on Handler temporarily because it updates records (cross-domain logic)
func (h *Handler) UpdatePatientV1(w http.ResponseWriter, r *http.Request) error {
	log := logger.FromCtx(r.Context())
	id := chi.URLParam(r, "id")
	var req struct {
		Name    *string `json:"name"`
		YOB     *string `json:"yob"`
		Gender  *string `json:"gender"`
		Address *string `json:"address"`
		Phone   *string `json:"phone"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}
	update := models.PatientUpdate{
		Name:    req.Name,
		YOB:     req.YOB,
		Gender:  req.Gender,
		Address: req.Address,
		Phone:   req.Phone,
	}
	if err := h.Store.UpdatePatientById(r.Context(), id, update); err != nil {
		return err
	}
	patient, err := h.Store.GetPatientById(r.Context(), id)
	if err != nil {
		return err
	}
	// When update patient info, also update all records with the new patient info
	records, err := h.Store.GetRecordsByPatientId(r.Context(), id)
	if err == nil {
		patientUpdate := models.UpdateRecordRequest{
			Patient: patient,
		}
		for _, record := range records {
			if err := h.Store.UpdateRecord(r.Context(), record.ID, patientUpdate); err != nil {
				log.Warn("failed to update patient info in the record",
					zap.String("record_id", record.ID),
					zap.String("patient_id", patient.ID),
					zap.Error(err),
				)
			}
		}
	}

	RespondJSON(w, http.StatusOK, patient)
	return nil
}

// GetPatientRecordsV1 handles GET /api/v1/patients/{id}/records
// NOTE: This method stays on Handler temporarily because it accesses record storage
func (h *Handler) GetPatientRecordsV1(w http.ResponseWriter, r *http.Request) error {
	patientId := chi.URLParam(r, "id")

	// Verify patient exists
	_, err := h.Store.GetPatientById(r.Context(), patientId)
	if err != nil {
		return err
	}

	// Parse query parameters for pagination and filtering
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		page = 1
	}
	pageSize, err := strconv.Atoi(r.URL.Query().Get("page_size"))
	if err != nil {
		pageSize = 20
	}

	status := r.URL.Query().Get("status")
	keyword := r.URL.Query().Get("keyword")

	sortBy := r.URL.Query().Get("sort_by")
	if sortBy == "" {
		sortBy = "created_at"
	}
	sortOrder := r.URL.Query().Get("sort_order")
	if sortOrder == "" {
		sortOrder = "desc"
	}

	// Query records for this patient
	recordsQueryOptions := models.RecordQueryOptions{
		PatientID: patientId,
		Status:    status,
		Keyword:   keyword,
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

// ComparePatientRecordsV1 handles POST /api/v1/patients/{id}/records/compare
// Creates an Excel comparison report for the specified patient's records.
// Supports optional tracking template for ordered test comparison.
func (h *Handler) ComparePatientRecordsV1(w http.ResponseWriter, r *http.Request) error {
	patientId := chi.URLParam(r, "id")

	// Verify patient exists
	patient, err := h.Store.GetPatientById(r.Context(), patientId)
	if err != nil {
		return NotFoundError("patient not found")
	}

	// Parse request body
	var req struct {
		RecordIDs  []string `json:"record_ids"`
		TrackingID string   `json:"tracking_id,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return BadRequestError("invalid request body")
	}

	if len(req.RecordIDs) == 0 {
		return BadRequestError("record_ids are required")
	}

	// Fetch records and verify they belong to this patient
	var records []*models.Record
	for _, recordID := range req.RecordIDs {
		record, err := h.Store.GetRecordById(r.Context(), recordID)
		if err != nil {
			return BadRequestError("record " + recordID + " not found")
		}
		if record.Patient.ID != patient.ID {
			return BadRequestError("record " + recordID + " does not belong to this patient")
		}
		records = append(records, record)
	}

	// Build ordered test list for comparison
	testList, err := h.buildOrderedTestList(r.Context(), records, req.TrackingID)
	if err != nil {
		return BadRequestError("failed to build test list: " + err.Error())
	}

	// Generate tracking report using new strategy pattern
	reportGenerator, err := sheets.NewReportGenerator(r.Context(), models.TrackingReport)
	if err != nil {
		return InternalServerError("failed to create report generator")
	}

	trackingData := &sheets.TrackingReportData{
		Records:  records,
		TestList: testList,
	}

	reader, err := reportGenerator.Generate(r.Context(), trackingData)
	if err != nil {
		return InternalServerError("failed to generate tracking report")
	}

	storer := sheets.LocalFileStoreStrategy{
		BaseDir: "./reports",
	}

	fileName := fmt.Sprintf("%s-%s-tracking.xlsx",
		time.Now().Format("20060102"),
		strings.ReplaceAll(records[0].Patient.Name, " ", "_"))

	filePath, err := storer.Store(r.Context(), reader, fileName)
	if err != nil {
		return InternalServerError("failed to store tracking report")
	}

	RespondJSON(w, http.StatusOK, map[string]string{
		"excel_file_path": filePath,
	})
	return nil
}

// buildOrderedTestList creates an ordered list of tests for comparison.
// If trackingID is provided, it uses the tracking template with proper ordering.
// Otherwise, it extracts unique tests from records maintaining order of appearance.
func (h *Handler) buildOrderedTestList(ctx context.Context, records []*models.Record, trackingID string) ([]models.TestInfo, error) {
	if trackingID != "" {
		return h.buildTestListFromTracking(ctx, trackingID)
	}
	return h.buildTestListFromRecords(records), nil
}

// buildTestListFromTracking creates a test list from a tracking template, sorted by Order field
func (h *Handler) buildTestListFromTracking(ctx context.Context, trackingID string) ([]models.TestInfo, error) {
	tracking, err := h.Store.GetTrackingById(ctx, trackingID)
	if err != nil {
		return nil, err
	}

	// Sort tracking tests by Order field using Go's built-in sort
	tests := make([]models.TrackingTestData, len(tracking.Tests))
	copy(tests, tracking.Tests)

	sort.Slice(tests, func(i, j int) bool {
		return tests[i].Order < tests[j].Order
	})

	// Convert to TestInfo
	testList := make([]models.TestInfo, 0, len(tests))
	for _, test := range tests {
		testList = append(testList, models.TestInfo{
			Name:        test.TestName,
			NormalValue: test.NormalValue,
			Unit:        test.Unit,
			Order:       test.Order,
		})
	}

	return testList, nil
}

// buildTestListFromRecords extracts unique tests from records maintaining order of appearance
func (h *Handler) buildTestListFromRecords(records []*models.Record) []models.TestInfo {
	seen := make(map[string]bool)
	var testList []models.TestInfo

	for _, record := range records {
		for _, test := range record.TestResults {
			if !seen[test.Name] {
				seen[test.Name] = true
				testList = append(testList, models.TestInfo{
					Name:        test.Name,
					NormalValue: test.NormalValue,
					Unit:        test.Unit,
					Order:       0, // Default order when no tracking template
				})
			}
		}
	}

	return testList
}
