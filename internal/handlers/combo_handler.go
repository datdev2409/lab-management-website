package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/a-h/templ"
	"github.com/datdev2409/lab-admin-go/internal/models"
	"github.com/datdev2409/lab-admin-go/internal/templates/pages"
	"github.com/datdev2409/lab-admin-go/internal/templates/partials"
	"github.com/go-chi/chi"
)

func (h *Handler) HandleComboPage(w http.ResponseWriter, r *http.Request) error {
	return Render(r.Context(), w, pages.ComboPage(""))
}

func (h *Handler) HandleComboCreatePage(w http.ResponseWriter, r *http.Request) error {
	return Render(r.Context(), w, pages.ComboCreatePage())
}

func (h *Handler) CreateCombo(w http.ResponseWriter, r *http.Request) error {
	if r.FormValue("test_ids") == "" {
		return errors.New("test_ids is required")
	}
	testIds := strings.Split(r.FormValue("test_ids"), ",")
	combo := models.NewCombo(r.FormValue("combo_name"), testIds)
	_, err := h.Store.InsertCombo(r.Context(), combo)
	if err != nil {
		return err
	}

	HTMXRedirect(w, "/danh-muc-goi-xet-nghiem")
	return nil
}

func (h *Handler) ListCombos(w http.ResponseWriter, r *http.Request) error {
	keyword := r.URL.Query().Get("combo_name")
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		page = 1
	}
	pageSize, err := strconv.Atoi(r.URL.Query().Get("page_size"))
	if err != nil {
		pageSize = 10
	}
	combos, pagination, err := h.Store.ListCombos(r.Context(), models.ComboQueryOptions{Keyword: keyword}, models.GenericQueryOptions{Page: page, PageSize: pageSize})
	if err != nil {
		return err
	}

	return RenderMultiComponents(r.Context(), w, []templ.Component{
		partials.ComboTable(combos),
		partials.Pagination(pagination, "combo-page"),
	})
}

func (h *Handler) GetComboDetails(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")
	combo, tests, err := h.Store.GetTestsInCombo(r.Context(), id)
	if err != nil {
		log.Println(err)
		return nil
	}

	target := r.Header.Get("HX-Target")

	if strings.HasPrefix(target, "combo-tests") {
		test_names := []string{}
		for _, test := range tests {
			test_names = append(test_names, test.Name)
		}
		w.Write([]byte(strings.Join(test_names, ", ")))
		return nil
	}

	response := models.ComboDetailsResponse{
		Combo: combo,
		Tests: tests,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
	return nil
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
	combo, tests, err := h.Store.GetTestsInCombo(r.Context(), id)
	if err != nil {
		return NotFoundError("combo not found")
	}
	RespondJSON(w, http.StatusOK, models.ComboDetailsResponse{Combo: combo, Tests: tests})
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
	if err := h.Store.UpdateComboById(r.Context(), id, update); err != nil {
		return err
	}
	combo, tests, err := h.Store.GetTestsInCombo(r.Context(), id)
	if err != nil {
		return NotFoundError("combo not found")
	}
	RespondJSON(w, http.StatusOK, models.ComboDetailsResponse{Combo: combo, Tests: tests})
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
