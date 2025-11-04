package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/datdev2409/lab-admin-go/internal/models"
	"github.com/datdev2409/lab-admin-go/internal/storage"
	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestHandler_GetPatient(t *testing.T) {
	t.Run("success - patient found", func(t *testing.T) {
		mockStore := new(storage.MockStorage)
		handler := &Handler{Store: mockStore}

		expectedPatient := &models.Patient{
			ID:      "patient_123",
			Name:    "John Doe",
			YOB:     "1990",
			Gender:  "Male",
			Address: "123 Main St",
			Phone:   "555-1234",
		}

		mockStore.On("GetPatientById", mock.Anything, "patient_123").
			Return(expectedPatient, nil)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/patients/patient_123", nil)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "patient_123")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
		
		w := httptest.NewRecorder()

		err := handler.GetPatient(w, req)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var patient models.Patient
		err = json.NewDecoder(w.Body).Decode(&patient)
		require.NoError(t, err)
		assert.Equal(t, expectedPatient.ID, patient.ID)
		assert.Equal(t, expectedPatient.Name, patient.Name)

		mockStore.AssertExpectations(t)
	})

	t.Run("error - patient not found", func(t *testing.T) {
		mockStore := new(storage.MockStorage)
		handler := &Handler{Store: mockStore}

		mockStore.On("GetPatientById", mock.Anything, "nonexistent").
			Return(nil, errors.New("patient not found"))

		req := httptest.NewRequest(http.MethodGet, "/api/v1/patients/nonexistent", nil)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "nonexistent")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
		
		w := httptest.NewRecorder()

		err := handler.GetPatient(w, req)

		require.Error(t, err)
		mockStore.AssertExpectations(t)
	})
}

func TestHandler_DeletePatient(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockStore := new(storage.MockStorage)
		handler := &Handler{Store: mockStore}

		mockStore.On("DeletePatientById", mock.Anything, "patient_123").
			Return(nil)

		req := httptest.NewRequest(http.MethodDelete, "/api/v1/patients/patient_123", nil)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "patient_123")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
		
		w := httptest.NewRecorder()

		err := handler.DeletePatient(w, req)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, w.Code)
		mockStore.AssertExpectations(t)
	})

	t.Run("error - delete fails", func(t *testing.T) {
		mockStore := new(storage.MockStorage)
		handler := &Handler{Store: mockStore}

		mockStore.On("DeletePatientById", mock.Anything, "patient_123").
			Return(errors.New("database error"))

		req := httptest.NewRequest(http.MethodDelete, "/api/v1/patients/patient_123", nil)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "patient_123")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
		
		w := httptest.NewRecorder()

		err := handler.DeletePatient(w, req)

		require.Error(t, err)
		mockStore.AssertExpectations(t)
	})
}

func TestHandler_ListPatientsV1(t *testing.T) {
	t.Run("success - with results", func(t *testing.T) {
		mockStore := new(storage.MockStorage)
		handler := &Handler{Store: mockStore}

		expectedPatients := []*models.Patient{
			{ID: "patient_1", Name: "John Doe", Phone: "555-1234"},
			{ID: "patient_2", Name: "Jane Smith", Phone: "555-5678"},
		}
		expectedPagination := &models.PaginationResponse{
			Total:     2,
			TotalPage: 1,
			Page:      1,
			PageSize:  10,
		}

		mockStore.On("SearchPatientByNameOrPhone", 
			mock.Anything, 
			models.PatientQueryOptions{Keyword: "John"},
			models.GenericQueryOptions{Page: 1, PageSize: 10}).
			Return(expectedPatients, expectedPagination, nil)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/patients?q=John&page=1&page_size=10", nil)
		w := httptest.NewRecorder()

		err := handler.ListPatientsV1(w, req)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err = json.NewDecoder(w.Body).Decode(&response)
		require.NoError(t, err)
		assert.Equal(t, "success", response["status"])
		assert.NotNil(t, response["data"])
		assert.NotNil(t, response["pagination"])

		mockStore.AssertExpectations(t)
	})

	t.Run("success - default pagination", func(t *testing.T) {
		mockStore := new(storage.MockStorage)
		handler := &Handler{Store: mockStore}

		expectedPatients := []*models.Patient{}
		expectedPagination := &models.PaginationResponse{
			Total:     0,
			TotalPage: 0,
			Page:      1,
			PageSize:  10,
		}

		mockStore.On("SearchPatientByNameOrPhone", 
			mock.Anything, 
			models.PatientQueryOptions{Keyword: ""},
			models.GenericQueryOptions{Page: 1, PageSize: 10}).
			Return(expectedPatients, expectedPagination, nil)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/patients", nil)
		w := httptest.NewRecorder()

		err := handler.ListPatientsV1(w, req)

		require.NoError(t, err)
		mockStore.AssertExpectations(t)
	})

	t.Run("error - storage failure", func(t *testing.T) {
		mockStore := new(storage.MockStorage)
		handler := &Handler{Store: mockStore}

		mockStore.On("SearchPatientByNameOrPhone", 
			mock.Anything, 
			mock.Anything,
			mock.Anything).
			Return(nil, nil, errors.New("database connection error"))

		req := httptest.NewRequest(http.MethodGet, "/api/v1/patients", nil)
		w := httptest.NewRecorder()

		err := handler.ListPatientsV1(w, req)

		require.Error(t, err)
		mockStore.AssertExpectations(t)
	})
}
