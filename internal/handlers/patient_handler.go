package handlers

import (
	"github.com/datdev2409/lab-admin-go/internal/models"
	"github.com/datdev2409/lab-admin-go/internal/templates/pages"
	"github.com/datdev2409/lab-admin-go/internal/templates/partials"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"log"
	"net/http"
)

func (h *Handler) HandlePatientPage(w http.ResponseWriter, r *http.Request) error {
	patients, err := h.Store.Patients().SearchByKeyword(r.Context(), "", map[string]string{"limit": "10"})
	if err != nil {
		patients = &[]models.Patient{}
	}
	return Render(r.Context(), w, pages.PatientsPage(*patients))
}

func (h *Handler) HandleCreatePatient(w http.ResponseWriter, r *http.Request) error {
	patient := models.Patient{
		ID:      "p-" + uuid.NewString(),
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
	opts := make(map[string]string)
	if limit := r.URL.Query().Get("limit"); limit != "" {
		opts["limit"] = limit
	}

	keyword := r.URL.Query().Get("patient_name")

	patients, err := h.Store.Patients().SearchByKeyword(r.Context(), keyword, opts)
	if err != nil {
		patients = &[]models.Patient{}
	}

	target := r.Header.Get("HX-Target")

	switch target {
	case "patient-table":
		return Render(r.Context(), w, partials.PatientTable(*patients))
	case "patient-suggestion-list":
		return Render(r.Context(), w, partials.PatientSuggestionList(*patients))
	}
	return nil
}

func (h *Handler) GetPatient(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")
	patient, err := h.Store.Patients().GetById(id)
	if err != nil {
		log.Println(err)
		return Render(r.Context(), w, partials.SelectUserForm(models.Patient{}, false))
	}

	return Render(r.Context(), w, partials.SelectUserForm(*patient, false))
}

func (h *Handler) UpdatePatient(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")

	patient := models.Patient{
		ID:      id,
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
	return Render(r.Context(), w, partials.PatientRow(patient))
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
