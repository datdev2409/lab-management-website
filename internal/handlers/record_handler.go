package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/a-h/templ"
	"github.com/datdev2409/lab-admin-go/internal/models"
	"github.com/datdev2409/lab-admin-go/internal/sheets"
	"github.com/datdev2409/lab-admin-go/internal/templates/pages"
	"github.com/datdev2409/lab-admin-go/internal/templates/partials"
	"github.com/go-chi/chi"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func (h *Handler) HandleRecordPage(w http.ResponseWriter, r *http.Request) error {
	return Render(context.Background(), w, pages.RecordPage())
}

func (h *Handler) HandleCreateNewRecord(w http.ResponseWriter, r *http.Request) error {
	return Render(r.Context(), w, pages.RecordCreatePage())
}

func (h *Handler) HandleRecordDetailPage(w http.ResponseWriter, r *http.Request) error {
	recordId := chi.URLParam(r, "id")
	return Render(r.Context(), w, pages.RecordDetailsPage(recordId))
}

func (h *Handler) CreateRecord(w http.ResponseWriter, r *http.Request) error {
	var request models.CreateRecordRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return err
	}

	patient, err := h.Store.Patients().GetById(request.PatientID)
	if err != nil {
		return err
	}

	recordTestResults := []models.TestResult{}
	for _, testResult := range request.TestResults {
		testId, err := bson.ObjectIDFromHex(testResult.ID)
		if err != nil {
			return err
		}
		recordTestResults = append(recordTestResults, models.TestResult{
			ID:          testId,
			Name:        testResult.Name,
			Price:       testResult.Price,
			NormalValue: testResult.NormalValue,
			Unit:        testResult.Unit,
			LowerBound:  testResult.LowerBound,
			UpperBound:  testResult.UpperBound,
			Result:      testResult.Result,
			ResultText:  testResult.ResultText,
		})
	}

	record := models.Record{
		Patient:     *patient,
		ComboName:   request.ComboName,
		TestResults: recordTestResults,
		Status:      "pending",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	recordId, err := h.Store.Records().Insert(r.Context(), &record)
	if err != nil {
		return err
	}

	jsonResponse, err := json.Marshal(models.CreateRecordResponse{ID: recordId})

	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonResponse)
	return nil
}

func (h *Handler) ListRecords(w http.ResponseWriter, r *http.Request) error {
	patientId := r.URL.Query().Get("patient_id")
	keyword := r.URL.Query().Get("keyword")
	status := r.URL.Query().Get("status")

	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		page = 1
	}
	pageSize, err := strconv.Atoi(r.URL.Query().Get("page_size"))
	if err != nil {
		pageSize = 20
	}

	sortBy := r.URL.Query().Get("sort_by")
	if sortBy == "" {
		sortBy = "created_at"
	}
	sortOrder := r.URL.Query().Get("sort_order")
	if sortOrder == "" {
		sortOrder = "desc"
	}

	recordsQueryOptions := models.RecordQueryOptions{
		Keyword:   keyword,
		Status:    status,
		PatientID: patientId,
	}

	genericQueryOptions := models.GenericQueryOptions{
		SortBy:    sortBy,
		SortOrder: sortOrder,
		Page:      page,
		PageSize:  pageSize,
	}

	records, pagination, err := h.Store.Records().ListRecords(r.Context(), recordsQueryOptions, genericQueryOptions)
	if err != nil {
		return err
	}

	log.Println(records)

	RenderMultiComponents(r.Context(), w, []templ.Component{
		partials.RecordTable(*records),
		partials.Pagination(pagination, "record-page"),
	})
	return nil
}

func (h *Handler) UpdateRecord(w http.ResponseWriter, r *http.Request) error {
	recordId := chi.URLParam(r, "id")

	var request models.UpdateRecordRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Println("Error while decoding request", err)
		return err
	}

	err := h.Store.Records().UpdateTestResults(r.Context(), recordId, request.TestResults)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Record updated successfully"}`))
	return nil
}

func (h *Handler) GetRecord(w http.ResponseWriter, r *http.Request) error {
	recordId := chi.URLParam(r, "id")
	record, err := h.Store.Records().GetById(r.Context(), recordId)
	if err != nil {
		return err
	}

	jsonResponse, err := json.Marshal(record)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
	return nil
}

type FileExportResponse struct {
	PDFPath   string `json:"pdf_path"`
	ExcelPath string `json:"excel_path"`
}

func (h *Handler) ExportRecordBilling(w http.ResponseWriter, r *http.Request) error {
	recordId := chi.URLParam(r, "id")

	record, err := h.Store.Records().GetById(r.Context(), recordId)
	if err != nil {
		return err
	}

	filepath, err := sheets.CreateRecordBillingFile(*record)
	if err != nil {
		return err
	}

	pdfPath, err := sheets.ConvertExcelToPDF(filepath)
	if err != nil {
		return err
	}

	response := FileExportResponse{
		PDFPath:   pdfPath,
		ExcelPath: filepath,
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
	return nil
}

// func (h *Handler) HandleRecordDetailPage(w http.ResponseWriter, r *http.Request) error {
// 	recordId := chi.URLParam(r, "id")
// 	recordWithDetails, err := h.Store.Records().GetDetails(r.Context(), recordId)
// 	// record, err := h.Store.Records().GetById(r.Context(), recordId)
// 	if err != nil {
// 		log.Println(err)
// 		return nil
// 	}
// 	return Render(r.Context(), w, pages.RecordDetailPage(*recordWithDetails, ""))
// }

// func (h *Handler) UpdateRecordCombo(w http.ResponseWriter, r *http.Request) error {
// 	recordId := chi.URLParam(r, "id")
// 	comboId := r.FormValue("combo_id")

// 	if comboId == "" {
// 		return Render(r.Context(), w, pages.ComboSelect(recordId, ""))
// 	}

// 	log.Println("Updating record with combo", comboId)
// 	combo, err := h.Store.Combos().GetById(r.Context(), comboId)
// 	if err != nil {
// 		log.Println(err)
// 		return err
// 	}

// 	err = h.Store.Records().UpdateCombo(r.Context(), recordId, combo)
// 	if err != nil {
// 		log.Println(err)
// 	}

// 	tests, err := h.Store.Tests().GetByIds(r.Context(), combo.GetTestIDs())
// 	if err != nil {
// 		log.Println(err)
// 	}

// 	log.Println(tests)

// 	err = h.Store.Records().AddTests(r.Context(), recordId, tests)
// 	if err != nil {
// 		log.Println(err)
// 	}

// 	// record, err := h.Store.Records().GetById(r.Context(), recordId)
// 	// if err != nil {
// 	// 	log.Println(err)
// 	// }

// 	return RenderMultiComponents(r.Context(), w, []templ.Component{
// 		pages.ComboSelect(recordId, combo.Name),
// 		// pages.RecordPageTestTable(tests, record.TestResults, true),
// 	})
// }

// func (h *Handler) UpdateRecordPatient(w http.ResponseWriter, r *http.Request) error {
// 	recordId := chi.URLParam(r, "id")
// 	patientId := r.FormValue("patient_id")

// 	if patientId == "" {
// 		return Render(r.Context(), w, pages.PatientSelectForm(recordId, models.Patient{}))
// 	}

// 	log.Println("Updating record with patient", patientId)
// 	patient, err := h.Store.Patients().GetById(patientId)
// 	if err != nil {
// 		log.Println(err)
// 	}
// 	err = h.Store.Records().UpdatePatient(r.Context(), recordId, *patient)
// 	if err != nil {
// 		log.Println(err)
// 	}
// 	return Render(r.Context(), w, pages.PatientSelectForm(recordId, *patient))
// }

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

// func (h *Handler) AddTestToRecord(w http.ResponseWriter, r *http.Request) error {
// 	recordId := chi.URLParam(r, "id")
// 	testId := r.FormValue("test_id")

// 	test, err := h.Store.Tests().GetById(testId)
// 	if err != nil {
// 		log.Println(err)
// 		return nil
// 	}
// 	err = h.Store.Records().AddTest(r.Context(), recordId, test)
// 	if err != nil {
// 		log.Println(err)
// 		return nil
// 	}

// 	log.Println("Updating record with test", recordId)
// 	log.Println("Updating record with test", testId)

// 	return Render(r.Context(), w, pages.RecordPageTestRow(test, models.TestResult{}))
// }

// func (h *Handler) UpdateRecordTests(w http.ResponseWriter, r *http.Request) error {
// 	recordId := chi.URLParam(r, "id")
// 	log.Println(recordId)

// 	r.ParseForm()

// 	tests := map[string]*models.TestResult{}

// 	for key, values := range r.Form {
// 		field, testId := ParseInputName(key, "#")

// 		_, ok := tests[testId]
// 		if !ok {
// 			// tests[testId] = &models.TestResult{TestID: testId}
// 		}

// 		if field == "result" {
// 			tests[testId].Result = SafeAccessSliceIndex(values, 0)
// 		}

// 		if field == "result_text" {
// 			tests[testId].ResultText = SafeAccessSliceIndex(values, 0)
// 		}
// 	}

// 	testResults := make([]models.TestResult, 0, len(tests))
// 	for _, result := range tests {
// 		testResults = append(testResults, *result)
// 	}
// 	h.Store.Records().SaveTestResults(r.Context(), recordId, testResults)

// 	return nil
// }
