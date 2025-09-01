package handlers

import (
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
