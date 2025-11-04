package handlers

import (
	"context"
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

func TestHandler_GetRecordById(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockStore := new(storage.MockStorage)
		handler := &Handler{Store: mockStore}

		patient := models.Patient{
			ID:   "patient_123",
			Name: "John Doe",
		}
		expectedRecord := &models.Record{
			ID:        "record_123",
			Patient:   patient,
			ComboName: "Basic Health Check",
			Status:    "pending",
		}

		mockStore.On("GetRecordById", mock.Anything, "record_123").
			Return(expectedRecord, nil)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/records/record_123", nil)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "record_123")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
		
		w := httptest.NewRecorder()

		err := handler.GetRecordV1(w, req)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, w.Code)
		mockStore.AssertExpectations(t)
	})

	t.Run("error - record not found", func(t *testing.T) {
		mockStore := new(storage.MockStorage)
		handler := &Handler{Store: mockStore}

		mockStore.On("GetRecordById", mock.Anything, "nonexistent").
			Return(nil, errors.New("record not found"))

		req := httptest.NewRequest(http.MethodGet, "/api/v1/records/nonexistent", nil)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "nonexistent")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
		
		w := httptest.NewRecorder()

		err := handler.GetRecordV1(w, req)

		require.Error(t, err)
		mockStore.AssertExpectations(t)
	})
}

func TestHandler_DeleteRecord(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockStore := new(storage.MockStorage)
		handler := &Handler{Store: mockStore}

		mockStore.On("DeleteRecord", mock.Anything, "record_123").
			Return(nil)

		req := httptest.NewRequest(http.MethodDelete, "/api/v1/records/record_123", nil)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "record_123")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
		
		w := httptest.NewRecorder()

		err := handler.DeleteRecordV1(w, req)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, w.Code)
		mockStore.AssertExpectations(t)
	})

	t.Run("error - delete fails", func(t *testing.T) {
		mockStore := new(storage.MockStorage)
		handler := &Handler{Store: mockStore}

		mockStore.On("DeleteRecord", mock.Anything, "record_123").
			Return(errors.New("cannot delete record"))

		req := httptest.NewRequest(http.MethodDelete, "/api/v1/records/record_123", nil)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "record_123")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
		
		w := httptest.NewRecorder()

		err := handler.DeleteRecordV1(w, req)

		require.Error(t, err)
		mockStore.AssertExpectations(t)
	})
}

func TestHandler_ListRecordsV1(t *testing.T) {
	t.Run("success - with filters", func(t *testing.T) {
		mockStore := new(storage.MockStorage)
		handler := &Handler{Store: mockStore}

		patient := models.Patient{ID: "patient_1", Name: "John Doe"}
		expectedRecords := []*models.Record{
			{ID: "record_1", Patient: patient, Status: "completed"},
		}
		expectedPagination := &models.PaginationResponse{
			Total:     1,
			TotalPage: 1,
			Page:      1,
			PageSize:  10,
		}

		mockStore.On("ListRecords", 
			mock.Anything, 
			mock.MatchedBy(func(filters models.RecordQueryOptions) bool {
				return filters.Status == "completed"
			}),
			models.GenericQueryOptions{SortBy: "created_at", SortOrder: "desc", Page: 1, PageSize: 10}).
			Return(expectedRecords, expectedPagination, nil)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/records?status=completed&page=1&page_size=10", nil)
		w := httptest.NewRecorder()

		err := handler.ListRecordsV1(w, req)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, w.Code)
		mockStore.AssertExpectations(t)
	})

	t.Run("success - no filters", func(t *testing.T) {
		mockStore := new(storage.MockStorage)
		handler := &Handler{Store: mockStore}

		expectedRecords := []*models.Record{}
		expectedPagination := &models.PaginationResponse{
			Total:     0,
			TotalPage: 0,
			Page:      1,
			PageSize:  20,
		}

		mockStore.On("ListRecords", 
			mock.Anything, 
			mock.Anything,
			models.GenericQueryOptions{SortBy: "created_at", SortOrder: "desc", Page: 1, PageSize: 20}).
			Return(expectedRecords, expectedPagination, nil)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/records", nil)
		w := httptest.NewRecorder()

		err := handler.ListRecordsV1(w, req)

		require.NoError(t, err)
		mockStore.AssertExpectations(t)
	})
}
