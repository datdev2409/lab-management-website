package handlers

import (
	"log"
	"net/http"

	"github.com/datdev2409/lab-admin-go/internal/models"
	"github.com/datdev2409/lab-admin-go/internal/templates/pages"
	"github.com/datdev2409/lab-admin-go/internal/templates/partials"
	"github.com/datdev2409/lab-admin-go/internal/view"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
	. "maragu.dev/gomponents"
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

	log.Println(keyword)

	patients, err := h.Store.Patients().SearchByKeyword(r.Context(), keyword, opts)
	if err != nil {
		patients = &[]models.Patient{}
	}

	target := r.Header.Get("HX-Target")
	log.Println(target)

	switch target {
	case "patient-table":
		return Render(r.Context(), w, partials.PatientTable(*patients))
	case "patient-suggestion-list":
		recordId := r.URL.Query().Get("record_id")
		log.Println(recordId)
		return Render(r.Context(), w, pages.PatientSuggestionList(*patients, recordId))
	case "cp_patient-suggestion-list":
		return view.PatientSuggestionList(*patients, false).Render(w)
	}

	return nil
}

func (h *Handler) GetPatient(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")
	patient, err := h.Store.Patients().GetById(id)
	if err != nil {
		log.Println(err)
		// return Render(r.Context(), w, pages.PatientInfo(models.Patient{}))
		return nil
	}

	return RenderOOB(r.Context(), w, []Node{
		view.PatientSelectInput(patient.Name, patient.ID),
		view.PatientInfo(patient, true),
		view.PatientSuggestionList([]models.Patient{}, true),
	})
	// return view.PatientInfo(patient, false).Render(w)
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
