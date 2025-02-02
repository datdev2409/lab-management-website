package handlers

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/datdev2409/lab-admin-go/internal/models"
	"github.com/datdev2409/lab-admin-go/internal/templates/pages"
	"github.com/datdev2409/lab-admin-go/internal/templates/partials"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

func (h *Handler) HandleRecordPage(w http.ResponseWriter, r *http.Request) error {
	return Render(context.Background(), w, pages.RecordPage())
}

func (h *Handler) HandleCreateNewRecord(w http.ResponseWriter, r *http.Request) error {
	// Create empty record with status "new"
	generatedRecordId := "r-" + uuid.NewString()
	record := &models.Record{
		ID:     generatedRecordId,
		Status: "new",
	}

	h.Store.Records().Insert(r.Context(), record)

	// Redirect to record page with record id
	http.Redirect(w, r, "/phieu-xet-nghiem/"+generatedRecordId, http.StatusSeeOther)
	return nil
}

func (h *Handler) HandleRecordDetailPage(w http.ResponseWriter, r *http.Request) error {
	recordId := chi.URLParam(r, "id")
	record, err := h.Store.Records().GetById(r.Context(), recordId)
	if err != nil {
		log.Println(err)
		return nil
	}
	return Render(r.Context(), w, pages.RecordDetailPage(*record, ""))
}

func (h *Handler) UpdateRecordPatient(w http.ResponseWriter, r *http.Request) error {
	recordId := chi.URLParam(r, "id")
	patientId := r.FormValue("patient_id")
	comboId := r.FormValue("combo_id")

	if strings.TrimSpace(patientId) != "" {
		log.Println("Updating record with patient", patientId)
		h.Store.Records().UpdatePatient(r.Context(), recordId, patientId)
		http.Redirect(w, r, "/api/patients/"+patientId, http.StatusSeeOther)
		return nil
	}
	if comboId != "" {
		log.Println("Updating record with combo", comboId)
		combo, _ := h.Store.Combos().GetById(r.Context(), comboId)
		err := h.Store.Records().UpdateCombo(r.Context(), recordId, combo)
		if err != nil {
			log.Println(err)
		}
		http.Redirect(w, r, "/api/records/"+recordId+"/tests", http.StatusSeeOther)
		return nil
	}

	return nil
}

func (h *Handler) GetRecordTests(w http.ResponseWriter, r *http.Request) error {
	recordId := chi.URLParam(r, "id")

	record, err := h.Store.Records().GetById(r.Context(), recordId)
	if err != nil {
		return Render(r.Context(), w, partials.TestTable([]models.Test{}, "record-page"))
	}

	testIds := []string{}
	for _, result := range record.TestResults {
		testIds = append(testIds, result.TestID)
	}
	tests, _ := h.Store.Tests().GetByIds(r.Context(), testIds)
	return Render(r.Context(), w, partials.TestTable(*tests, "record-page"))
}
