package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/a-h/templ"
	"github.com/datdev2409/lab-admin-go/internal/models"
	"github.com/datdev2409/lab-admin-go/internal/templates/pages"
	"github.com/datdev2409/lab-admin-go/internal/templates/partials"
	"github.com/go-chi/chi"
)

func (h *Handler) HandlePatientPage(w http.ResponseWriter, r *http.Request) error {
	return Render(r.Context(), w, pages.PatientsPage())
}

func (h *Handler) HandleCreatePatient(w http.ResponseWriter, r *http.Request) error {
	// patient := models.Patient{
	// 	Name:    r.FormValue("patient_name"),
	// 	YOB:     r.FormValue("patient_yob"),
	// 	Gender:  r.FormValue("patient_gender"),
	// 	Address: r.FormValue("patient_address"),
	// 	Phone:   r.FormValue("patient_phone"),
	// }
	patient := models.NewPatient(
		r.FormValue("patient_name"),
		r.FormValue("patient_yob"),
		r.FormValue("patient_gender"),
		r.FormValue("patient_address"),
		r.FormValue("patient_phone"),
	)

	_, err := h.StoreV2.InsertPatient(r.Context(), patient)
	if err != nil {
		log.Println(err)
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
	patients, pagination, err := h.StoreV2.SearchPatientByNameOrPhone(r.Context(), models.PatientQueryOptions{Keyword: keyword}, models.GenericQueryOptions{Page: page, PageSize: pageSize})
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
	case "tracking-patient-autocomplete":
		return Render(r.Context(), w, partials.TrackingPatientAutocomplete(patients))
	}

	return nil
}

func (h *Handler) GetPatient(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")
	patient, err := h.StoreV2.GetPatientById(r.Context(), id)
	if err != nil {
		return err
	}

	jsonResponse, err := json.Marshal(patient)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
	return nil
}

func (h *Handler) UpdatePatient(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")

	update := models.PatientUpdate{
		Name:    models.GetStringPtr(r.FormValue("patient_name")),
		YOB:     models.GetStringPtr(r.FormValue("patient_yob")),
		Gender:  models.GetStringPtr(r.FormValue("patient_gender")),
		Address: models.GetStringPtr(r.FormValue("patient_address")),
		Phone:   models.GetStringPtr(r.FormValue("patient_phone")),
	}

	err := h.StoreV2.UpdatePatientById(r.Context(), id, update)
	if err != nil {
		log.Println(err)
		return err
	}

	patient, err := h.StoreV2.GetPatientById(r.Context(), id)
	if err != nil {
		return err
	}

	return Render(r.Context(), w, partials.PatientRow(patient))
}

func (h *Handler) DeletePatient(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")
	err := h.StoreV2.DeletePatientById(r.Context(), id)
	if err != nil {
		log.Println(err)
		return err
	}

	w.WriteHeader(http.StatusOK)
	return nil
}
