package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

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

	record := models.NewRecord(*patient, recordTestResults)

	recordId, err := h.Store.Records().Insert(r.Context(), record)
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusCreated, map[string]string{"id": recordId})
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

	return WriteJSON(w, http.StatusOK, map[string]string{
		"message": "Record updated successfully",
	})
}

func (h *Handler) GetRecord(w http.ResponseWriter, r *http.Request) error {
	recordId := chi.URLParam(r, "id")
	record, err := h.Store.Records().GetById(r.Context(), recordId)
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, record)
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

	filePath, err := sheets.CreateRecordBillingFile(record)
	if err != nil {
		return err
	}

	pdfPath, err := sheets.ConvertExcelToPDF(filePath)
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, map[string]string{
		"pdf_path":   pdfPath,
		"excel_path": filePath,
	})
}

func (h *Handler) ExportRecordResults(w http.ResponseWriter, r *http.Request) error {
	recordId := chi.URLParam(r, "id")

	log.Println("Create the record result pdf file", recordId)

	record, err := h.Store.Records().GetById(r.Context(), recordId)
	if err != nil {
		return err
	}

	filePath, err := sheets.CreateRecordResultFile(record)
	if err != nil {
		return err
	}

	pdfPath, err := sheets.ConvertExcelToPDF(filePath)
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, map[string]string{
		"pdf_path":   pdfPath,
		"excel_path": filePath,
	})
}
