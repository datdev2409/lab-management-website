package handlers

import (
	"encoding/json"
	"log"
	"log/slog"
	"math"
	"net/http"
	"strconv"

	"github.com/a-h/templ"
	"github.com/datdev2409/lab-admin-go/internal/models"
	"github.com/datdev2409/lab-admin-go/internal/templates/pages"
	"github.com/datdev2409/lab-admin-go/internal/templates/partials"
	"github.com/go-chi/chi"
)

func (h *Handler) HandleTestPage(w http.ResponseWriter, r *http.Request) error {
	tests, _, err := h.Store.ListTests(r.Context(), models.TestQueryOptions{}, models.GenericQueryOptions{Page: 1, PageSize: 10})
	if err != nil {
		return err
	}
	messages := map[string]string{
		"test:create:success": "Thêm xét nghiệm thành công",
	}
	redirectCode := GetAndDeleteFlashCookie(w, r)
	return Render(r.Context(), w, pages.TestPage(tests, messages[redirectCode]))
}

func (h *Handler) HandleCreateTest(w http.ResponseWriter, r *http.Request) error {
	testLowerBound, err := strconv.ParseFloat(r.FormValue("test_lower_bound"), 32)
	if err != nil {
		errorMessage := `<div class="alert alert-danger" role="alert">Đã có lỗi xảy ra khi thêm xét nghiệm.</div>`
		w.Write([]byte(errorMessage))
	}
	testUpperBound, err := strconv.ParseFloat(r.FormValue("test_upper_bound"), 32)
	if err != nil {
		errorMessage := `<div class="alert alert-danger" role="alert">Đã có lỗi xảy ra khi thêm xét nghiệm.</div>`
		w.Write([]byte(errorMessage))
	}
	testPrice, err := strconv.Atoi(r.FormValue("test_price"))
	if err != nil {
		errorMessage := `<div class="alert alert-danger" role="alert">Đã có lỗi xảy ra khi thêm xét nghiệm.</div>`
		w.Write([]byte(errorMessage))
	}

	test := models.NewTest(
		r.FormValue("test_name"),
		testPrice,
		r.FormValue("test_normal_value"),
		r.FormValue("test_unit"),
		math.Round(testLowerBound*100)/100,
		math.Round(testUpperBound*100)/100,
	)

	_, err = h.Store.InsertTest(r.Context(), test)
	if err != nil {
		errorMessage := `<div class="alert alert-danger" role="alert">Đã có lỗi xảy ra khi thêm xét nghiệm.</div>`
		w.Write([]byte(errorMessage))
	}

	SetFlashCookie(w, "test:create:success")
	HTMXRedirect(w, "/danh-muc-xet-nghiem")
	return nil
}

func (h *Handler) SearchTestsByKeyword(w http.ResponseWriter, r *http.Request) error {
	keyword := r.URL.Query().Get("test_name")
	tests, _, err := h.Store.ListTests(r.Context(), models.TestQueryOptions{Keyword: keyword}, models.GenericQueryOptions{Page: 1, PageSize: 5})
	if err != nil {
		return err
	}

	log.Println(r.Header.Get("HX-Target"))
	target := r.Header.Get("HX-Target")
	switch target {
	case "test-autocomplete":
		return nil
	case "test-table":
		return Render(r.Context(), w, partials.TestTable(tests, "test-page", false))
	case "test-search-result-combo-page":
		return Render(r.Context(), w, partials.TestAutocomplete(tests, "combo-page"))
	case "test-search-result-record-page":
		return Render(r.Context(), w, partials.TestAutocomplete(tests, "record-page"))
	}
	return nil
}

func (h *Handler) ListTests(w http.ResponseWriter, r *http.Request) error {
	keyword := r.URL.Query().Get("test_name")
	slog.Debug("Listing tests with keyword", "keyword", keyword)
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		page = 1
	}
	tests, pagination, err := h.Store.ListTests(r.Context(), models.TestQueryOptions{Keyword: keyword}, models.GenericQueryOptions{Page: page, PageSize: 10})
	if err != nil {
		return err
	}
	target := r.Header.Get("HX-Target")
	switch target {
	case "test-table":
		return RenderMultiComponents(r.Context(), w, []templ.Component{
			partials.TestTable(tests, "test-page", false),
			partials.Pagination(pagination, "test-page"),
		})
	case "test-search-result-combo-page":
		return Render(r.Context(), w, partials.TestAutocomplete(tests, "combo-page"))
	case "test-autocomplete":
		return Render(r.Context(), w, partials.RecordTestAutocomplete(tests))
	}
	return nil
}

func (h *Handler) DeleteTest(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")
	err := h.Store.DeleteTestById(r.Context(), id)
	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusOK)
	return nil
}

// ListTestsV1 handles GET /api/v1/tests
func (h *Handler) ListTestsV1(w http.ResponseWriter, r *http.Request) error {
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		page = 1
	}
	pageSize, err := strconv.Atoi(r.URL.Query().Get("page_size"))
	if err != nil {
		pageSize = 10
	}
	keyword := r.URL.Query().Get("test_name")
	tests, pagination, err := h.Store.ListTests(r.Context(), models.TestQueryOptions{Keyword: keyword}, models.GenericQueryOptions{Page: page, PageSize: pageSize})
	if err != nil {
		return err
	}
	RespondJSON(w, http.StatusOK, map[string]interface{}{
		"tests":      tests,
		"pagination": pagination,
	})
	return nil
}

// CreateTestV1 handles POST /api/v1/tests
func (h *Handler) CreateTestV1(w http.ResponseWriter, r *http.Request) error {
	var req struct {
		Name        string  `json:"name"`
		Price       int     `json:"price"`
		NormalValue string  `json:"normal_value"`
		Unit        string  `json:"unit"`
		LowerBound  float64 `json:"lower_bound"`
		UpperBound  float64 `json:"upper_bound"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return nil
	}
	test := models.NewTest(
		req.Name,
		req.Price,
		req.NormalValue,
		req.Unit,
		math.Round(req.LowerBound*100)/100,
		math.Round(req.UpperBound*100)/100,
	)
	newTest, err := h.Store.InsertTest(r.Context(), test)
	if err != nil {
		RespondJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to create test"})
		return nil
	}
	RespondJSON(w, http.StatusCreated, newTest)
	return nil
}

// GetTestV1 handles GET /api/v1/tests/{id}
func (h *Handler) GetTestV1(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")
	test, err := h.Store.GetTestById(r.Context(), id)
	if err != nil {
		RespondJSON(w, http.StatusNotFound, map[string]string{"error": "test not found"})
		return nil
	}
	RespondJSON(w, http.StatusOK, test)
	return nil
}

// UpdateTestV1 handles PUT /api/v1/tests/{id}
func (h *Handler) UpdateTestV1(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")
	var req struct {
		Name        *string  `json:"name"`
		Price       *int     `json:"price"`
		NormalValue *string  `json:"normal_value"`
		Unit        *string  `json:"unit"`
		LowerBound  *float64 `json:"lower_bound"`
		UpperBound  *float64 `json:"upper_bound"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return nil
	}
	update := make(map[string]interface{})
	if req.Name != nil {
		update["name"] = *req.Name
	}
	if req.Price != nil {
		update["price"] = *req.Price
	}
	if req.NormalValue != nil {
		update["normal_value"] = *req.NormalValue
	}
	if req.Unit != nil {
		update["unit"] = *req.Unit
	}
	if req.LowerBound != nil {
		update["lower_bound"] = *req.LowerBound
	}
	if req.UpperBound != nil {
		update["upper_bound"] = *req.UpperBound
	}
	if len(update) == 0 {
		RespondJSON(w, http.StatusBadRequest, map[string]string{"error": "no fields to update"})
		return nil
	}
	if err := h.Store.UpdateTestById(r.Context(), id, update); err != nil {
		RespondJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to update test"})
		return nil
	}
	test, err := h.Store.GetTestById(r.Context(), id)
	if err != nil {
		RespondJSON(w, http.StatusNotFound, map[string]string{"error": "test not found"})
		return nil
	}
	RespondJSON(w, http.StatusOK, test)
	return nil
}

// DeleteTestV1 handles DELETE /api/v1/tests/{id}
func (h *Handler) DeleteTestV1(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")
	if err := h.Store.DeleteTestById(r.Context(), id); err != nil {
		RespondJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to delete test"})
		return nil
	}
	RespondJSON(w, http.StatusOK, map[string]string{"result": "deleted"})
	return nil
}
