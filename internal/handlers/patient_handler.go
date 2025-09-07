package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/datdev2409/lab-admin-go/internal/models"
	"github.com/datdev2409/lab-admin-go/internal/sheets"
	"github.com/datdev2409/lab-admin-go/internal/templates/pages"
	"github.com/go-chi/chi"
)

func (h *Handler) HandlePatientPage(w http.ResponseWriter, r *http.Request) error {
	return Render(r.Context(), w, pages.PatientsPage())
}

func (h *Handler) GetPatient(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")
	patient, err := h.Store.GetPatientById(r.Context(), id)
	if err != nil {
		return err
	}

	jsonResponse, err := json.Marshal(patient)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
	return nil
}

func (h *Handler) DeletePatient(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")
	err := h.Store.DeletePatientById(r.Context(), id)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	return nil
}

// ListPatientsV1 handles GET /api/v1/patients
func (h *Handler) ListPatientsV1(w http.ResponseWriter, r *http.Request) error {
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		page = 1
	}
	pageSize, err := strconv.Atoi(r.URL.Query().Get("page_size"))
	if err != nil {
		pageSize = 10
	}

	keyword := r.URL.Query().Get("q")
	patients, pagination, err := h.Store.SearchPatientByNameOrPhone(r.Context(), models.PatientQueryOptions{Keyword: keyword}, models.GenericQueryOptions{Page: page, PageSize: pageSize})
	if err != nil {
		return err
	}

	RespondJSONWithPagination(w, http.StatusOK, patients, pagination)
	return nil
}

// CreatePatientV1 handles POST /api/v1/patients
func (h *Handler) CreatePatientV1(w http.ResponseWriter, r *http.Request) error {
	var req struct {
		Name    string `json:"name"`
		YOB     string `json:"yob"`
		Gender  string `json:"gender"`
		Address string `json:"address"`
		Phone   string `json:"phone"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}
	patient := models.NewPatient(req.Name, req.YOB, req.Gender, req.Address, req.Phone)
	newPatient, err := h.Store.InsertPatient(r.Context(), patient)
	if err != nil {
		return err
	}
	RespondJSON(w, http.StatusCreated, newPatient)
	return nil
}

// GetPatientV1 handles GET /api/v1/patients/{id}
func (h *Handler) GetPatientV1(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")
	patient, err := h.Store.GetPatientById(r.Context(), id)
	if err != nil {
		return err
	}
	RespondJSON(w, http.StatusOK, patient)
	return nil
}

// UpdatePatientV1 handles PUT /api/v1/patients/{id}
func (h *Handler) UpdatePatientV1(w http.ResponseWriter, r *http.Request) error {
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
	RespondJSON(w, http.StatusOK, patient)
	return nil
}

// DeletePatientV1 handles DELETE /api/v1/patients/{id}
func (h *Handler) DeletePatientV1(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")
	if err := h.Store.DeletePatientById(r.Context(), id); err != nil {
		return err
	}
	RespondJSON(w, http.StatusOK, map[string]string{"result": "deleted"})
	return nil
}

// GetPatientRecordsV1 handles GET /api/v1/patients/{id}/records
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

	// Build test map for comparison
	testMap := make(map[string]models.TestInfo)
	if req.TrackingID != "" {
		// Use tracking template
		tracking, err := h.Store.GetTrackingById(r.Context(), req.TrackingID)
		if err != nil {
			return BadRequestError("tracking template not found")
		}
		for _, test := range tracking.Tests {
			testMap[test.TestName] = models.TestInfo{
				Name:        test.TestName,
				NormalValue: test.NormalValue,
				Unit:        test.Unit,
			}
		}
	} else {
		// Use all tests from the records
		for _, record := range records {
			for _, test := range record.TestResults {
				testMap[test.Name] = models.TestInfo{
					Name:        test.Name,
					NormalValue: test.NormalValue,
					Unit:        test.Unit,
				}
			}
		}
	}

	// Generate Excel file
	filePath, err := sheets.CreateRecordTrackingFile(r.Context(), records, testMap)
	if err != nil {
		return InternalServerError("failed to create comparison file")
	}

	RespondJSON(w, http.StatusOK, map[string]string{
		"excel_file_path": filePath,
	})
	return nil
}
