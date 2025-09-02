package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/datdev2409/lab-admin-go/internal/models"
	"github.com/datdev2409/lab-admin-go/internal/templates/pages"
	"github.com/go-chi/chi"
)

func (h *Handler) HandleComboPage(w http.ResponseWriter, r *http.Request) error {
	return Render(r.Context(), w, pages.ComboPage(""))
}

func (h *Handler) HandleComboCreatePage(w http.ResponseWriter, r *http.Request) error {
	return Render(r.Context(), w, pages.ComboCreatePage(""))
}

func (h *Handler) HandleComboEditPage(w http.ResponseWriter, r *http.Request) error {
	comboId := chi.URLParam(r, "id")
	return Render(r.Context(), w, pages.ComboCreatePage(comboId))
}

// ListCombosV1 handles GET /api/v1/combos
func (h *Handler) ListCombosV1(w http.ResponseWriter, r *http.Request) error {
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		page = 1
	}
	pageSize, err := strconv.Atoi(r.URL.Query().Get("page_size"))
	if err != nil {
		pageSize = 10
	}
	keyword := r.URL.Query().Get("q")
	combos, pagination, err := h.Store.ListCombos(r.Context(), models.ComboQueryOptions{Keyword: keyword}, models.GenericQueryOptions{Page: page, PageSize: pageSize})
	if err != nil {
		return err
	}
	RespondJSONWithPagination(w, http.StatusOK, combos, pagination)
	return nil
}

// CreateComboV1 handles POST /api/v1/combos
func (h *Handler) CreateComboV1(w http.ResponseWriter, r *http.Request) error {
	var req struct {
		Name    string   `json:"name"`
		TestIDs []string `json:"test_ids"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return BadRequestError("invalid request body")
	}
	if req.Name == "" || len(req.TestIDs) == 0 {
		return BadRequestError("name and test_ids are required")
	}
	combo := models.NewCombo(req.Name, req.TestIDs)
	id, err := h.Store.InsertCombo(r.Context(), combo)
	if err != nil {
		return err
	}
	RespondJSON(w, http.StatusCreated, map[string]string{"id": id})
	return nil
}

// GetComboV1 handles GET /api/v1/combos/{id}
func (h *Handler) GetComboV1(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")
	combo, err := h.Store.GetComboById(r.Context(), id)
	if err != nil {
		return NotFoundError("combo not found")
	}
	RespondJSON(w, http.StatusOK, combo)
	return nil
}

// GetComboTestsV1 handles GET /api/v1/combos/{id}/tests
func (h *Handler) GetComboTestsV1(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")
	tests, err := h.Store.GetTestsByComboId(r.Context(), id)
	if err != nil {
		return NotFoundError("combo not found")
	}
	RespondJSON(w, http.StatusOK, tests)
	return nil
}

// UpdateComboV1 handles PUT /api/v1/combos/{id}
func (h *Handler) UpdateComboV1(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")
	var req struct {
		Name    *string  `json:"name"`
		TestIDs []string `json:"test_ids"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return BadRequestError("invalid request body")
	}
	update := map[string]interface{}{}
	if req.Name != nil {
		update["name"] = *req.Name
	}
	if req.TestIDs != nil {
		update["test_ids"] = req.TestIDs
	}
	if len(update) == 0 {
		RespondJSON(w, http.StatusNotModified, map[string]string{"message": "no fields to update"})
		return nil
	}
	combo, err := h.Store.UpdateComboByIdAndReturn(r.Context(), id, update)
	if err != nil {
		return err
	}
	RespondJSON(w, http.StatusOK, combo)
	return nil
}

// DeleteComboV1 handles DELETE /api/v1/combos/{id}
func (h *Handler) DeleteComboV1(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")
	if err := h.Store.DeleteComboById(r.Context(), id); err != nil {
		return err
	}
	RespondJSON(w, http.StatusOK, map[string]string{"result": "deleted"})
	return nil
}
