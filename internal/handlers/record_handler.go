package handlers

import (
	"context"
	"log"
	"net/http"

	"github.com/a-h/templ"
	"github.com/datdev2409/lab-admin-go/internal/models"
	"github.com/datdev2409/lab-admin-go/internal/templates/pages"
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
	recordWithDetails, err := h.Store.Records().GetDetails(r.Context(), recordId)
	// record, err := h.Store.Records().GetById(r.Context(), recordId)
	if err != nil {
		log.Println(err)
		return nil
	}
	return Render(r.Context(), w, pages.RecordDetailPage(*recordWithDetails, ""))
}

func (h *Handler) UpdateRecordCombo(w http.ResponseWriter, r *http.Request) error {
	recordId := chi.URLParam(r, "id")
	comboId := r.FormValue("combo_id")

	if comboId == "" {
		return Render(r.Context(), w, pages.ComboSelect(recordId, ""))
	}

	log.Println("Updating record with combo", comboId)
	combo, err := h.Store.Combos().GetById(r.Context(), comboId)
	if err != nil {
		log.Println(err)
		return err
	}

	err = h.Store.Records().UpdateCombo(r.Context(), recordId, combo)
	if err != nil {
		log.Println(err)
	}

	tests, err := h.Store.Tests().GetByIds(r.Context(), combo.Tests)
	if err != nil {
		log.Println(err)
	}

	log.Println(tests)

	err = h.Store.Records().AddTests(r.Context(), recordId, tests)
	if err != nil {
		log.Println(err)
	}

	// record, err := h.Store.Records().GetById(r.Context(), recordId)
	// if err != nil {
	// 	log.Println(err)
	// }

	return RenderMultiComponents(r.Context(), w, []templ.Component{
		pages.ComboSelect(recordId, combo.Name),
		// pages.RecordPageTestTable(tests, record.TestResults, true),
	})
}

func (h *Handler) UpdateRecordPatient(w http.ResponseWriter, r *http.Request) error {
	recordId := chi.URLParam(r, "id")
	patientId := r.FormValue("patient_id")

	if patientId == "" {
		return Render(r.Context(), w, pages.PatientSelectForm(recordId, models.Patient{}))
	}

	log.Println("Updating record with patient", patientId)
	patient, err := h.Store.Patients().GetById(patientId)
	if err != nil {
		log.Println(err)
	}
	err = h.Store.Records().UpdatePatient(r.Context(), recordId, *patient)
	if err != nil {
		log.Println(err)
	}
	return Render(r.Context(), w, pages.PatientSelectForm(recordId, *patient))
}

// func (h *Handler) GetRecordTests(w http.ResponseWriter, r *http.Request) error {
// 	recordId := chi.URLParam(r, "id")

// 	record, err := h.Store.Records().GetById(r.Context(), recordId)
// 	if err != nil {
// 		return Render(r.Context(), w, partials.TestTable([]models.Test{}, "record-page"))
// 	}

// 	testIds := []string{}
// 	for _, result := range record.TestResults {
// 		testIds = append(testIds, result.TestID)
// 	}
// 	tests, _ := h.Store.Tests().GetByIds(r.Context(), testIds)
// 	return Render(r.Context(), w, partials.TestTable(*tests, "record-page"))
// }

// func (h *Handler) SearchRecordsByPatientNameOrPhone(w http.ResponseWriter, r *http.Request) error {
// 	keyword := r.FormValue("keyword")
// 	status := r.FormValue("status")
// 	limit := 20

// 	log.Println(keyword, status, limit)

// 	filter := models.RecordSearchFilter{
// 		Keyword: &keyword,
// 		Status:  &status,
// 		Limit:   &limit,
// 	}

// 	records, err := h.Store.Records().SearchRecords(r.Context(), filter)
// 	log.Println(records)
// 	if err != nil {
// 		return err
// 	}

// 	return Render(r.Context(), w, partials.RecordTable(*records))
// }

func (h *Handler) AddTestToRecord(w http.ResponseWriter, r *http.Request) error {
	recordId := chi.URLParam(r, "id")
	testId := r.FormValue("test_id")

	test, err := h.Store.Tests().GetById(testId)
	if err != nil {
		log.Println(err)
		return nil
	}
	err = h.Store.Records().AddTest(r.Context(), recordId, test)
	if err != nil {
		log.Println(err)
		return nil
	}

	log.Println("Updating record with test", recordId)
	log.Println("Updating record with test", testId)

	return Render(r.Context(), w, pages.RecordPageTestRow(test, models.TestResult{}))
}

func (h *Handler) UpdateRecordTests(w http.ResponseWriter, r *http.Request) error {
	recordId := chi.URLParam(r, "id")
	log.Println(recordId)

	r.ParseForm()

	tests := map[string]*models.TestResult{}

	for key, values := range r.Form {
		field, testId := ParseInputName(key, "#")

		_, ok := tests[testId]
		if !ok {
			// tests[testId] = &models.TestResult{TestID: testId}
		}

		if field == "result" {
			tests[testId].Result = SafeAccessSliceIndex(values, 0)
		}

		if field == "result_text" {
			tests[testId].ResultText = SafeAccessSliceIndex(values, 0)
		}
	}

	testResults := make([]models.TestResult, 0, len(tests))
	for _, result := range tests {
		testResults = append(testResults, *result)
	}
	h.Store.Records().SaveTestResults(r.Context(), recordId, testResults)

	return nil
}
