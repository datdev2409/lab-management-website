package handlers

import (
	"errors"
	"net/http"

	"github.com/datdev2409/lab-admin-go/internal/models"
	"github.com/datdev2409/lab-admin-go/internal/service"
	"github.com/datdev2409/lab-admin-go/internal/templates/pages"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type ComboHandler struct {
	comboService *service.ComboService
	validator    *validator.Validate
}

func NewComboHandler(comboService *service.ComboService, validator *validator.Validate) *ComboHandler {
	return &ComboHandler{
		comboService: comboService,
		validator:    validator,
	}
}

// HandleComboPage handles GET /danh-muc-goi-xet-nghiem
func (h *ComboHandler) HandleComboPage(w http.ResponseWriter, r *http.Request) error {
	return Render(r.Context(), w, pages.ComboPage(""))
}

// HandleComboCreatePage handles GET /danh-muc-goi-xet-nghiem/new
func (h *ComboHandler) HandleComboCreatePage(w http.ResponseWriter, r *http.Request) error {
	return Render(r.Context(), w, pages.ComboCreatePage(""))
}

// HandleComboEditPage handles GET /danh-muc-goi-xet-nghiem/{id}/edit
func (h *ComboHandler) HandleComboEditPage(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")
	return Render(r.Context(), w, pages.ComboCreatePage(id))
}

// CreateCombo handles POST /api/v1/combos
func (h *ComboHandler) CreateCombo(w http.ResponseWriter, r *http.Request) error {
	var input models.CreateComboInput

	if err := BindAndValidate(r, h.validator, &input); err != nil {
		return err
	}

	comboId, err := h.comboService.CreateCombo(r.Context(), &input)
	if err != nil {
		if errors.Is(err, service.ErrComboAlreadyExists) {
			return &AppError{StatusCode: http.StatusConflict, Message: COMBO_ALREADY_EXISTS}
		}
		if errors.Is(err, service.ErrInvalidComboTests) {
			return &AppError{StatusCode: http.StatusBadRequest, Message: INVALID_COMBO_TESTS}
		}
		return err
	}

	RespondJSON(w, http.StatusCreated, map[string]string{"id": comboId})
	return nil
}

// SearchCombosByName handles GET /api/v1/combos?q=keyword&page=1&page_size=10
func (h *ComboHandler) SearchCombosByName(w http.ResponseWriter, r *http.Request) error {
	queryOpts := ParseListParams(r, 10) // default page size 10

	keyword := r.URL.Query().Get("q")
	combos, pagination, err := h.comboService.SearchCombosByName(
		r.Context(),
		keyword,
		queryOpts.Page,
		queryOpts.PageSize,
	)
	if err != nil {
		return err
	}

	RespondJSONWithPagination(w, http.StatusOK, combos, pagination)
	return nil
}

// GetCombo handles GET /api/v1/combos/{id}
func (h *ComboHandler) GetCombo(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")

	// Delegate fetching both combo and its tests to the service for simpler handler logic
	details, err := h.comboService.GetComboDetails(r.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrComboNotFound) {
			return &AppError{StatusCode: http.StatusNotFound, Message: COMBO_NOT_FOUND}
		}
		return err
	}

	RespondJSON(w, http.StatusOK, details)
	return nil
}

// UpdateCombo handles PATCH /api/v1/combos/{id}
func (h *ComboHandler) UpdateCombo(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")

	var update models.ComboUpdate
	if err := BindAndValidate(r, h.validator, &update); err != nil {
		return err
	}

	combo, err := h.comboService.UpdateComboById(r.Context(), id, update)
	if err != nil {
		if errors.Is(err, service.ErrComboNotFound) {
			return &AppError{StatusCode: http.StatusNotFound, Message: COMBO_NOT_FOUND}
		}
		return err
	}

	RespondJSON(w, http.StatusOK, combo)
	return nil
}

// DeleteCombo handles DELETE /api/v1/combos/{id}
func (h *ComboHandler) DeleteCombo(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")

	err := h.comboService.DeleteComboById(r.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrComboNotFound) {
			return &AppError{StatusCode: http.StatusNotFound, Message: COMBO_NOT_FOUND}
		}
		return err
	}

	RespondJSON(w, http.StatusOK, map[string]string{"result": "deleted"})
	return nil
}

// GetComboTests handles GET /api/v1/combos/{id}/tests
func (h *ComboHandler) GetComboTests(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")

	tests, err := h.comboService.GetComboTests(r.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrComboNotFound) {
			return &AppError{StatusCode: http.StatusNotFound, Message: COMBO_NOT_FOUND}
		}
		return err
	}

	RespondJSON(w, http.StatusOK, tests)
	return nil
}

// ListAllCombos handles GET /api/v1/combos/all
func (h *ComboHandler) ListAllCombos(w http.ResponseWriter, r *http.Request) error {
	combos, err := h.comboService.ListAllCombos(r.Context())
	if err != nil {
		return err
	}

	RespondJSON(w, http.StatusOK, combos)
	return nil
}
