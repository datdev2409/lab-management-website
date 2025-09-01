package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/a-h/templ"
	"github.com/datdev2409/lab-admin-go/internal/models"
	"github.com/datdev2409/lab-admin-go/internal/templates/pages"
	"github.com/datdev2409/lab-admin-go/internal/templates/partials"
	"github.com/go-chi/chi"
)

func (h *Handler) HandlePatientPage(w http.ResponseWriter, r *http.Request) error {
	return Render(r.Context(), w, pages.PatientsPage())
}

func (h *Handler) HandleCreatePatient(w http.ResponseWriter, r *http.Request) error {
	patient := models.NewPatient(
		r.FormValue("patient_name"),
		r.FormValue("patient_yob"),
		r.FormValue("patient_gender"),
		r.FormValue("patient_address"),
		r.FormValue("patient_phone"),
	)

	_, err := h.Store.InsertPatient(r.Context(), patient)
	if err != nil {
		return err
	}

	HTMXRedirect(w, "/danh-muc-benh-nhan")
	return nil
}

func (h *Handler) ListPatients(w http.ResponseWriter, r *http.Request) error {
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		page = 1
	}
	pageSize, err := strconv.Atoi(r.URL.Query().Get("page_size"))
	if err != nil {
		pageSize = 10
	}

	keyword := r.URL.Query().Get("patient_name")
	patients, pagination, err := h.Store.SearchPatientByNameOrPhone(r.Context(), models.PatientQueryOptions{Keyword: keyword}, models.GenericQueryOptions{Page: page, PageSize: pageSize})
	if err != nil {
		return err
	}

	target := r.Header.Get("HX-Target")

	switch target {
	case "patient-table":
		return RenderMultiComponents(r.Context(), w, []templ.Component{
			partials.PatientTable(patients),
			partials.Pagination(pagination, "patient-page"),
		})
	case "patient-autocomplete":
		return Render(r.Context(), w, partials.PatientAutocomplete(patients))
	case "tracking-patient-autocomplete":
		return Render(r.Context(), w, partials.TrackingPatientAutocomplete(patients))
	}

	return nil
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

func (h *Handler) UpdatePatient(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")

	update := models.PatientUpdate{
		Name:    models.GetStringPtr(r.FormValue("patient_name")),
		YOB:     models.GetStringPtr(r.FormValue("patient_yob")),
		Gender:  models.GetStringPtr(r.FormValue("patient_gender")),
		Address: models.GetStringPtr(r.FormValue("patient_address")),
		Phone:   models.GetStringPtr(r.FormValue("patient_phone")),
	}

	err := h.Store.UpdatePatientById(r.Context(), id, update)
	if err != nil {
		return err
	}

	patient, err := h.Store.GetPatientById(r.Context(), id)
	if err != nil {
		return err
	}

	return Render(r.Context(), w, partials.PatientRow(patient))
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

	keyword := r.URL.Query().Get("patient_name")
	patients, pagination, err := h.Store.SearchPatientByNameOrPhone(r.Context(), models.PatientQueryOptions{Keyword: keyword}, models.GenericQueryOptions{Page: page, PageSize: pageSize})
	if err != nil {
		return err
	}

	RespondJSON(w, http.StatusOK, map[string]interface{}{
		"patients":   patients,
		"pagination": pagination,
	})
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
