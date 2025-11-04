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

func TestHandler_DeleteTest(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockStore := new(storage.MockStorage)
		handler := &Handler{Store: mockStore}

		mockStore.On("DeleteTestById", mock.Anything, "test_123").
			Return(nil)

		req := httptest.NewRequest(http.MethodDelete, "/api/v1/tests/test_123", nil)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "test_123")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
		
		w := httptest.NewRecorder()

		err := handler.DeleteTestV1(w, req)

		require.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, w.Code)
		mockStore.AssertExpectations(t)
	})
}

func TestHandler_ListTestsV1(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockStore := new(storage.MockStorage)
		handler := &Handler{Store: mockStore}

		expectedTests := []*models.Test{
			{ID: "test_1", Name: "Blood Sugar", Price: 50000},
			{ID: "test_2", Name: "Cholesterol", Price: 60000},
		}
		expectedPagination := &models.PaginationResponse{
			Total:     2,
			TotalPage: 1,
			Page:      1,
			PageSize:  10,
		}

		mockStore.On("ListTests", 
			mock.Anything, 
			models.TestQueryOptions{Keyword: "Blood"},
			models.GenericQueryOptions{Page: 1, PageSize: 10}).
			Return(expectedTests, expectedPagination, nil)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/tests?q=Blood&page=1&page_size=10", nil)
		w := httptest.NewRecorder()

		err := handler.ListTestsV1(w, req)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err = json.NewDecoder(w.Body).Decode(&response)
		require.NoError(t, err)
		assert.Equal(t, "success", response["status"])

		mockStore.AssertExpectations(t)
	})
}

func TestHandler_DeleteCombo(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockStore := new(storage.MockStorage)
		handler := &Handler{Store: mockStore}

		mockStore.On("DeleteComboById", mock.Anything, "combo_123").
			Return(nil)

		req := httptest.NewRequest(http.MethodDelete, "/api/v1/combos/combo_123", nil)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "combo_123")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
		
		w := httptest.NewRecorder()

		err := handler.DeleteComboV1(w, req)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, w.Code)
		mockStore.AssertExpectations(t)
	})

	t.Run("error - delete fails", func(t *testing.T) {
		mockStore := new(storage.MockStorage)
		handler := &Handler{Store: mockStore}

		mockStore.On("DeleteComboById", mock.Anything, "combo_123").
			Return(errors.New("database error"))

		req := httptest.NewRequest(http.MethodDelete, "/api/v1/combos/combo_123", nil)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "combo_123")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
		
		w := httptest.NewRecorder()

		err := handler.DeleteComboV1(w, req)

		require.Error(t, err)
		mockStore.AssertExpectations(t)
	})
}

func TestHandler_ListCombosV1(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockStore := new(storage.MockStorage)
		handler := &Handler{Store: mockStore}

		expectedCombos := []*models.Combo{
			{ID: "combo_1", Name: "Basic Health Check", TestIDs: []string{"test_1", "test_2"}},
		}
		expectedPagination := &models.PaginationResponse{
			Total:     1,
			TotalPage: 1,
			Page:      1,
			PageSize:  10,
		}

		mockStore.On("ListCombos", 
			mock.Anything, 
			models.ComboQueryOptions{Keyword: ""},
			models.GenericQueryOptions{Page: 1, PageSize: 10}).
			Return(expectedCombos, expectedPagination, nil)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/combos", nil)
		w := httptest.NewRecorder()

		err := handler.ListCombosV1(w, req)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, w.Code)

		mockStore.AssertExpectations(t)
	})
}

func TestHandler_GetComboDetail(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockStore := new(storage.MockStorage)
		handler := &Handler{Store: mockStore}

		expectedTests := []*models.Test{
			{ID: "test_1", Name: "Blood Sugar"},
			{ID: "test_2", Name: "Cholesterol"},
		}

		mockStore.On("GetTestsByComboId", mock.Anything, "combo_123").
			Return(expectedTests, nil)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/combos/combo_123/tests", nil)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "combo_123")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
		
		w := httptest.NewRecorder()

		err := handler.GetComboTestsV1(w, req)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, w.Code)

		mockStore.AssertExpectations(t)
	})

	t.Run("error - combo not found", func(t *testing.T) {
		mockStore := new(storage.MockStorage)
		handler := &Handler{Store: mockStore}

		mockStore.On("GetTestsByComboId", mock.Anything, "nonexistent").
			Return(nil, errors.New("combo not found"))

		req := httptest.NewRequest(http.MethodGet, "/api/v1/combos/nonexistent/tests", nil)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "nonexistent")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
		
		w := httptest.NewRecorder()

		err := handler.GetComboTestsV1(w, req)

		require.Error(t, err)
		mockStore.AssertExpectations(t)
	})
}
