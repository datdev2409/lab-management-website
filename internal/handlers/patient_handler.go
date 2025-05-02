package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/a-h/templ"
	"github.com/datdev2409/lab-admin-go/internal/models"
	"github.com/datdev2409/lab-admin-go/internal/templates/pages"
	"github.com/datdev2409/lab-admin-go/internal/templates/partials"
	"github.com/datdev2409/lab-admin-go/internal/view"
	"github.com/go-chi/chi"
	. "maragu.dev/gomponents"
)

func (h *Handler) HandlePatientPage(w http.ResponseWriter, r *http.Request) error {
	return Render(r.Context(), w, pages.PatientsPage())
}

func (h *Handler) HandleCreatePatient(w http.ResponseWriter, r *http.Request) error {
	patient := models.Patient{
		Name:    r.FormValue("patient_name"),
		YOB:     r.FormValue("patient_yob"),
		Gender:  r.FormValue("patient_gender"),
		Address: r.FormValue("patient_address"),
		Phone:   r.FormValue("patient_phone"),
	}

	err := h.Store.Patients().Insert(&patient)
	if err != nil {
		return err
	}

	SetFlashCookie(w, "patient:create:success")
	HTMXRedirect(w, "/phieu-xet-nghiem")
	return nil
}

func (h *Handler) ListPatients(w http.ResponseWriter, r *http.Request) error {
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		page = 1
	}
	pageSize, err := strconv.Atoi(r.URL.Query().Get("page_size"))
	if err != nil {
		pageSize = 10
	}

	keyword := r.URL.Query().Get("patient_name")
	patients, pagination, err := h.Store.Patients().ListPatients(r.Context(), models.PatientQueryOptions{Keyword: keyword}, models.GenericQueryOptions{Page: page, PageSize: pageSize})
	if err != nil {
		return err
	}

	target := r.Header.Get("HX-Target")

	switch target {
	case "patient-table":
		return RenderMultiComponents(r.Context(), w, []templ.Component{
			partials.PatientTable(patients),
			partials.Pagination(pagination, "patient-page"),
		})
	case "patient-autocomplete":
		return Render(r.Context(), w, partials.PatientAutocomplete(patients))
	case "cp_patient-suggestion-list":
		return view.PatientSuggestionList(patients, false).Render(w)
	}

	return nil
}

func (h *Handler) GetPatient(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")
	patient, err := h.Store.Patients().GetById(id)
	if err != nil {
		return err
	}

	target := r.Header.Get("HX-Target")
	switch target {
	case "patient-info":
		return Render(r.Context(), w, partials.PatientInfo(patient))
	}

	return RenderOOB(r.Context(), w, []Node{
		view.PatientSelectInput(patient.Name, patient.ID.Hex()),
		view.PatientInfo(patient, true),
		// view.PatientSuggestionList([]models.Patient{}, true),
	})
	// return view.PatientInfo(patient, false).Render(w)
}

func (h *Handler) UpdatePatient(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")

	patient := models.Patient{
		Name:    r.FormValue("patient_name"),
		YOB:     r.FormValue("patient_yob"),
		Gender:  r.FormValue("patient_gender"),
		Address: r.FormValue("patient_address"),
		Phone:   r.FormValue("patient_phone"),
	}

	err := h.Store.Patients().UpdateById(r.Context(), id, &patient)
	if err != nil {
		log.Println(err)
		return err
	}
	return Render(r.Context(), w, partials.PatientRow(&patient))
}

func (h *Handler) DeletePatient(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")
	err := h.Store.Patients().Delete(id)
	if err != nil {
		log.Println(err)
		return err
	}

	w.WriteHeader(http.StatusOK)
	return nil
}
