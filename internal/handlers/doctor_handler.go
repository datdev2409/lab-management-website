package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/datdev2409/lab-admin-go/internal/models"
	"github.com/datdev2409/lab-admin-go/internal/templates/pages"
	"github.com/go-chi/chi"
)

func (h *Handler) HandleDoctorPage(w http.ResponseWriter, r *http.Request) error {
	return Render(r.Context(), w, pages.DoctorsPage())
}

// ListDoctorsV1 handles GET /api/v1/doctors
func (h *Handler) ListDoctorsV1(w http.ResponseWriter, r *http.Request) error {
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		page = 1
	}
	pageSize, err := strconv.Atoi(r.URL.Query().Get("page_size"))
	if err != nil {
		pageSize = 10
	}

	keyword := r.URL.Query().Get("q")
	doctors, pagination, err := h.Store.SearchDoctorByNameOrPhone(r.Context(), models.DoctorQueryOptions{Keyword: keyword}, models.GenericQueryOptions{Page: page, PageSize: pageSize})
	if err != nil {
		return err
	}

	RespondJSONWithPagination(w, http.StatusOK, doctors, pagination)
	return nil
}

// CreateDoctorV1 handles POST /api/v1/doctors
func (h *Handler) CreateDoctorV1(w http.ResponseWriter, r *http.Request) error {
	var req struct {
		Name    string `json:"name"`
		Phone   string `json:"phone"`
		Address string `json:"address"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}

	existing, err := h.Store.FindDoctorByNameAndPhone(r.Context(), req.Name, req.Phone)
	if err != nil {
		return err
	}
	if existing != nil {
		return &AppError{http.StatusBadRequest, DUPLICATE_DOCTOR_ERROR}
	}

	doctor := models.NewDoctor(req.Name, req.Phone, req.Address)
	newDoctor, err := h.Store.InsertDoctor(r.Context(), doctor)
	if err != nil {
		return err
	}
	RespondJSON(w, http.StatusCreated, newDoctor)
	return nil
}

// GetDoctorV1 handles GET /api/v1/doctors/{id}
func (h *Handler) GetDoctorV1(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")
	doctor, err := h.Store.GetDoctorById(r.Context(), id)
	if err != nil {
		return err
	}
	RespondJSON(w, http.StatusOK, doctor)
	return nil
}

// UpdateDoctorV1 handles PUT /api/v1/doctors/{id}
func (h *Handler) UpdateDoctorV1(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")
	var req struct {
		Name    *string `json:"name"`
		Phone   *string `json:"phone"`
		Address *string `json:"address"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}
	
	// Check for duplicate if name or phone is being updated
	if req.Name != nil || req.Phone != nil {
		// Get current doctor to check if name/phone combination is changing
		currentDoctor, err := h.Store.GetDoctorById(r.Context(), id)
		if err != nil {
			return err
		}
		
		// Use current values if not being updated
		newName := currentDoctor.Name
		newPhone := currentDoctor.Phone
		if req.Name != nil {
			newName = *req.Name
		}
		if req.Phone != nil {
			newPhone = *req.Phone
		}
		
		// Check if this name+phone combination exists for a different doctor
		existing, err := h.Store.FindDoctorByNameAndPhone(r.Context(), newName, newPhone)
		if err != nil {
			return err
		}
		if existing != nil && existing.ID != id {
			return &AppError{http.StatusBadRequest, DUPLICATE_DOCTOR_ERROR}
		}
	}

	update := models.DoctorUpdate{
		Name:    req.Name,
		Phone:   req.Phone,
		Address: req.Address,
	}
	if err := h.Store.UpdateDoctorById(r.Context(), id, update); err != nil {
		return err
	}
	doctor, err := h.Store.GetDoctorById(r.Context(), id)
	if err != nil {
		return err
	}

	RespondJSON(w, http.StatusOK, doctor)
	return nil
}

// DeleteDoctorV1 handles DELETE /api/v1/doctors/{id}
func (h *Handler) DeleteDoctorV1(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")
	if err := h.Store.DeleteDoctorById(r.Context(), id); err != nil {
		return err
	}
	RespondJSON(w, http.StatusOK, map[string]string{"result": "deleted"})
	return nil
}