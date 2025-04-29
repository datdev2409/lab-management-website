package handlers

import (
	"log"
	"math"
	"net/http"
	"strconv"

	"github.com/datdev2409/lab-admin-go/internal/models"
	"github.com/datdev2409/lab-admin-go/internal/templates/pages"
	"github.com/datdev2409/lab-admin-go/internal/templates/partials"
	"github.com/google/uuid"
)

func (h *Handler) HandleTestPage(w http.ResponseWriter, r *http.Request) error {
	tests, err := h.Store.Tests().SearchByKeyword(r.Context(), "", map[string]string{"limit": "10"})
	if err != nil {
		log.Println(err)
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
		log.Println(err)
		errorMessage := `<div class="alert alert-danger" role="alert">Đã có lỗi xảy ra khi thêm xét nghiệm.</div>`
		w.Write([]byte(errorMessage))
	}
	testUpperBound, err := strconv.ParseFloat(r.FormValue("test_upper_bound"), 32)
	if err != nil {
		log.Println(err)
		errorMessage := `<div class="alert alert-danger" role="alert">Đã có lỗi xảy ra khi thêm xét nghiệm.</div>`
		w.Write([]byte(errorMessage))
	}
	testPrice, err := strconv.Atoi(r.FormValue("test_price"))
	if err != nil {
		log.Println(err)
		errorMessage := `<div class="alert alert-danger" role="alert">Đã có lỗi xảy ra khi thêm xét nghiệm.</div>`
		w.Write([]byte(errorMessage))
	}
	test := models.Test{
		ID:          "t-" + uuid.NewString(),
		Name:        r.FormValue("test_name"),
		NormalValue: r.FormValue("test_normal_value"),
		Unit:        r.FormValue("test_unit"),
		LowerBound:  math.Round(testLowerBound*100) / 100,
		UpperBound:  math.Round(testUpperBound*100) / 100,
		Price:       testPrice,
	}

	err = h.Store.Tests().Insert(&test)
	if err != nil {
		log.Println(err)
		errorMessage := `<div class="alert alert-danger" role="alert">Đã có lỗi xảy ra khi thêm xét nghiệm.</div>`
		w.Write([]byte(errorMessage))
	}

	SetFlashCookie(w, "test:create:success")
	HTMXRedirect(w, "/danh-muc-xet-nghiem")
	return nil
}

func (h *Handler) SearchTestsByKeyword(w http.ResponseWriter, r *http.Request) error {
	keyword := r.URL.Query().Get("test_name")
	tests, err := h.Store.Tests().SearchByKeyword(r.Context(), keyword, map[string]string{"limit": "5"})
	if err != nil {
		return err
	}

	log.Println(r.Header.Get("HX-Target"))
	target := r.Header.Get("HX-Target")
	switch target {
	case "test-autocomplete":
		recordId := r.URL.Query().Get("record_id")
		return Render(r.Context(), w, pages.TestAutocomplete(tests, recordId))
	case "test-table":
		return Render(r.Context(), w, partials.TestTable(tests, "test-page", false))
	case "test-search-result-combo-page":
		return Render(r.Context(), w, partials.TestAutocomplete(tests, "combo-page"))
	case "test-search-result-record-page":
		return Render(r.Context(), w, partials.TestAutocomplete(tests, "record-page"))
	}
	return nil
}

//func (h *Handler) SelectTest(w http.ResponseWriter, r *http.Request) error {
//	id := chi.URLParam(r, "id")
//
//	test, err := h.Store.Tests().GetById(id)
//	if err != nil {
//		test = &models.Test{}
//		return Render(r.Context(), w, pages.TestSearchRow(*test))
//	}
//
//	log.Println(test)
//	return Render(r.Context(), w, pages.TestSearchRow(*test))
//}
