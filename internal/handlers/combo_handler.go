package handlers

import (
	"log"
	"net/http"
	"strings"

	"github.com/datdev2409/lab-admin-go/internal/models"
	"github.com/datdev2409/lab-admin-go/internal/templates/pages"
	"github.com/datdev2409/lab-admin-go/internal/templates/partials"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

func (h *Handler) HandleComboPage(w http.ResponseWriter, r *http.Request) error {
	return Render(r.Context(), w, pages.ComboPage(""))
}

func (h *Handler) CreateCombo(w http.ResponseWriter, r *http.Request) error {
	combo := &models.Combo{
		ID:    "c" + uuid.NewString(),
		Name:  r.FormValue("combo_name"),
		Tests: strings.Split(r.FormValue("test_ids"), ","),
	}
	return h.Store.Combos().Insert(combo)
}

func (h *Handler) SearchCombosByKeyword(w http.ResponseWriter, r *http.Request) error {
	keyword := r.URL.Query().Get("combo_name")
	recordId := r.URL.Query().Get("record_id")
	log.Println(recordId)
	combos, err := h.Store.Combos().SearchByKeyword(r.Context(), keyword, map[string]string{"limit": "5"})
	if err != nil {
		return err
	}

	return Render(r.Context(), w, pages.ComboSuggestionList(*combos, recordId))
}

func (h *Handler) GetCombo(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")
	combo, err := h.Store.Combos().GetById(r.Context(), id)

	if err != nil {
		log.Println(err)
		return nil
	}

	tests, err := h.Store.Tests().GetByIds(r.Context(), combo.Tests)
	if err != nil {
		log.Println(err)
		return nil
	}

	log.Println(tests)

	return Render(r.Context(), w, partials.TestTable(tests, "record-page", false))
}
