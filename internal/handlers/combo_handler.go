package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

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
	testIds, err := models.ConvertIDsToObjectIDs(strings.Split(r.FormValue("test_ids"), ","))
	if err != nil {
		return err
	}
	combo := &models.Combo{
		Name:      r.FormValue("combo_name"),
		TestIDs:   testIds,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err = h.Store.Combos().Insert(combo)
	if err != nil {
		log.Println(err)
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
	combos, pagination, err := h.Store.Combos().ListCombos(r.Context(), models.ComboQueryOptions{Keyword: keyword}, models.GenericQueryOptions{Page: page, PageSize: pageSize})
	if err != nil {
		return err
	}

	target := r.Header.Get("HX-Target")
	switch target {
	case "combo-autocomplete":
		return Render(r.Context(), w, partials.ComboAutocomplete(combos))
	}

	return RenderMultiComponents(r.Context(), w, []templ.Component{
		partials.ComboTable(combos),
		partials.Pagination(pagination, "combo-page"),
	})
}

func (h *Handler) GetComboDetails(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")
	combo, tests, err := h.Store.Combos().GetTestsInCombo(r.Context(), id)
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

	// // return Render(r.Context(), w, partials.RecordTestTable(tests))
	// return RenderMultiComponents(r.Context(), w, []templ.Component{
	// 	partials.RecordTestTable(tests),
	// 	partials.RecordComboInfo(combo),
	// })
	response := models.ComboDetailsResponse{
		Combo: combo,
		Tests: tests,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
	return nil
}
