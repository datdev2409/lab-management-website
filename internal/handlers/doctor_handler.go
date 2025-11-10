package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/datdev2409/lab-admin-go/internal/models"
	"github.com/datdev2409/lab-admin-go/internal/service"
	"github.com/datdev2409/lab-admin-go/internal/templates/pages"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type DoctorHandler struct {
	doctorService *service.DoctorService
	validator     *validator.Validate
}

func NewDoctorHandler(doctorService *service.DoctorService, validator *validator.Validate) *DoctorHandler {
	return &DoctorHandler{
		doctorService: doctorService,
		validator:     validator,
	}
}

// HandleDoctorPage handles GET /danh-muc-bac-si
func (h *DoctorHandler) HandleDoctorPage(w http.ResponseWriter, r *http.Request) error {
	return Render(r.Context(), w, pages.DoctorsPage())
}

// CreateDoctor handles POST /api/v1/doctors
func (h *DoctorHandler) CreateDoctor(w http.ResponseWriter, r *http.Request) error {
	var input models.CreateDoctorInput

	if err := BindAndValidate(r, h.validator, &input); err != nil {
		return err
	}

	doctor, err := h.doctorService.CreateDoctor(r.Context(), &input)
	if err != nil {
		if errors.Is(err, service.ErrDoctorAlreadyExists) {
			return &AppError{
				StatusCode: http.StatusConflict,
				Message:    DOCTOR_ALREADY_EXISTS,
			}
		}
		return err
	}

	RespondJSON(w, http.StatusCreated, doctor)
	return nil
}

// SearchDoctorsByKeyword handles GET /api/v1/doctors?q=keyword&page=1&page_size=10
func (h *DoctorHandler) SearchDoctorsByKeyword(w http.ResponseWriter, r *http.Request) error {
	queryOpts := ParseListParams(r, 10) // default page size 10

	keyword := r.URL.Query().Get("q")
	doctors, pagination, err := h.doctorService.SearchDoctorsByKeyword(
		r.Context(),
		keyword,
		queryOpts.Page,
		queryOpts.PageSize,
	)
	if err != nil {
		return err
	}

	RespondJSONWithPagination(w, http.StatusOK, doctors, pagination)
	return nil
}

// GetDoctor handles GET /api/v1/doctors/{id}
func (h *DoctorHandler) GetDoctor(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")

	doctor, err := h.doctorService.GetDoctorById(r.Context(), id)
	if err != nil {
		return err
	}

	RespondJSON(w, http.StatusOK, doctor)
	return nil
}

// UpdateDoctor handles PATCH /api/v1/doctors/{id}
func (h *DoctorHandler) UpdateDoctor(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")

	var update models.DoctorUpdate
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		return &AppError{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid request body",
		}
	}

	doctor, err := h.doctorService.UpdateDoctorById(r.Context(), id, update)
	if err != nil {
		return err
	}

	RespondJSON(w, http.StatusOK, doctor)
	return nil
}

// DeleteDoctor handles DELETE /api/v1/doctors/{id}
func (h *DoctorHandler) DeleteDoctor(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")

	err := h.doctorService.DeleteDoctorById(r.Context(), id)
	if err != nil {
		return err
	}

	RespondJSON(w, http.StatusOK, map[string]string{"result": "deleted"})
	return nil
}
