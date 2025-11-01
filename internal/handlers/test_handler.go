package handlers

import (
	"encoding/json"
	"math"
	"net/http"
	"strconv"

	"github.com/datdev2409/lab-admin-go/internal/logger"
	"github.com/datdev2409/lab-admin-go/internal/models"
	"github.com/datdev2409/lab-admin-go/internal/templates/pages"
	"github.com/go-chi/chi"
	"go.uber.org/zap"
)

func (h *Handler) HandleTestPage(w http.ResponseWriter, r *http.Request) error {
	return Render(r.Context(), w, pages.TestPage())
}

func (h *Handler) HandleCreateTest(w http.ResponseWriter, r *http.Request) error {
	testLowerBound, err := strconv.ParseFloat(r.FormValue("test_lower_bound"), 32)
	if err != nil {
		errorMessage := `<div class="alert alert-danger" role="alert">Đã có lỗi xảy ra khi thêm xét nghiệm.</div>`
		w.Write([]byte(errorMessage))
	}
	testUpperBound, err := strconv.ParseFloat(r.FormValue("test_upper_bound"), 32)
	if err != nil {
		errorMessage := `<div class="alert alert-danger" role="alert">Đã có lỗi xảy ra khi thêm xét nghiệm.</div>`
		w.Write([]byte(errorMessage))
	}
	testPrice, err := strconv.Atoi(r.FormValue("test_price"))
	if err != nil {
		errorMessage := `<div class="alert alert-danger" role="alert">Đã có lỗi xảy ra khi thêm xét nghiệm.</div>`
		w.Write([]byte(errorMessage))
	}

	test := models.NewTest(
		r.FormValue("test_name"),
		testPrice,
		r.FormValue("test_normal_value"),
		r.FormValue("test_unit"),
		math.Round(testLowerBound*100)/100,
		math.Round(testUpperBound*100)/100,
	)

	_, err = h.Store.InsertTest(r.Context(), test)
	if err != nil {
		errorMessage := `<div class="alert alert-danger" role="alert">Đã có lỗi xảy ra khi thêm xét nghiệm.</div>`
		w.Write([]byte(errorMessage))
	}

	SetFlashCookie(w, "test:create:success")
	HTMXRedirect(w, "/danh-muc-xet-nghiem")
	return nil
}

func (h *Handler) DeleteTest(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")
	err := h.Store.DeleteTestById(r.Context(), id)
	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusOK)
	return nil
}

// ListTestsV1 handles GET /api/v1/tests
func (h *Handler) ListTestsV1(w http.ResponseWriter, r *http.Request) error {
	log := logger.FromCtx(r.Context())
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		page = 1
	}
	pageSize, err := strconv.Atoi(r.URL.Query().Get("page_size"))
	if err != nil {
		pageSize = 10
	}
	keyword := r.URL.Query().Get("q")
	log.Info("Listing tests", zap.String("keyword", keyword), zap.Int("page", page), zap.Int("pageSize", pageSize))
	tests, pagination, err := h.Store.ListTests(r.Context(), models.TestQueryOptions{Keyword: keyword}, models.GenericQueryOptions{Page: page, PageSize: pageSize})
	if err != nil {
		return err
	}
	RespondJSONWithPagination(w, http.StatusOK, tests, pagination)
	return nil
}

// CreateTestV1 handles POST /api/v1/tests
func (h *Handler) CreateTestV1(w http.ResponseWriter, r *http.Request) error {
	log := logger.FromCtx(r.Context())
	var req models.CreateTestRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error("Failed to decode request body", zap.Error(err))
		return BadRequestError("invalid request body")
	}

	log.Debug("Creating test", zap.Float64("lower_bound", req.LowerBound))

	// // Parse price
	// price, err := strconv.Atoi(req.Price)
	// if err != nil {
	// 	return BadRequestError("invalid price value")
	// }

	// // Parse lower bound
	// lowerBound, err := strconv.ParseFloat(req.LowerBound, 64)
	// if err != nil {
	// 	return BadRequestError("invalid lower_bound value")
	// }

	// // Parse upper bound
	// upperBound, err := strconv.ParseFloat(req.UpperBound, 64)
	// if err != nil {
	// 	return BadRequestError("invalid upper_bound value")
	// }

	if req.Unit == "" {
		req.Unit = "."
	}

	test := models.NewTest(
		req.Name,
		req.Price,
		req.NormalValue,
		req.Unit,
		math.Round(req.LowerBound*1000)/1000,
		math.Round(req.UpperBound*1000)/1000,
	)
	newTest, err := h.Store.InsertTest(r.Context(), test)
	if err != nil {
		return err
	}
	RespondJSON(w, http.StatusCreated, newTest)
	return nil
}

// GetTestV1 handles GET /api/v1/tests/{id}
func (h *Handler) GetTestV1(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")
	test, err := h.Store.GetTestById(r.Context(), id)
	if err != nil {
		return NotFoundError("test not found")
	}
	RespondJSON(w, http.StatusOK, test)
	return nil
}

// UpdateTestV1 handles PUT /api/v1/tests/{id}
func (h *Handler) UpdateTestV1(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")
	var req struct {
		Name        *string  `json:"name"`
		Price       *int     `json:"price"`
		NormalValue *string  `json:"normal_value"`
		Unit        *string  `json:"unit"`
		LowerBound  *float64 `json:"lower_bound"`
		UpperBound  *float64 `json:"upper_bound"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return BadRequestError("invalid request body")
	}
	update := make(map[string]interface{})
	if req.Name != nil {
		update["name"] = *req.Name
	}
	if req.Price != nil {
		update["price"] = *req.Price
	}
	if req.NormalValue != nil {
		update["normal_value"] = *req.NormalValue
	}
	if req.Unit != nil {
		update["unit"] = *req.Unit
	}
	if req.LowerBound != nil {
		update["lower_bound"] = *req.LowerBound
	}
	if req.UpperBound != nil {
		update["upper_bound"] = *req.UpperBound
	}
	if len(update) == 0 {
		RespondJSON(w, http.StatusNotModified, map[string]string{"message": "no fields to update"})
		return nil
	}
	if err := h.Store.UpdateTestById(r.Context(), id, update); err != nil {
		return err
	}
	test, err := h.Store.GetTestById(r.Context(), id)
	if err != nil {
		return NotFoundError("test not found")
	}
	RespondJSON(w, http.StatusOK, test)
	return nil
}

// DeleteTestV1 handles DELETE /api/v1/tests/{id}
func (h *Handler) DeleteTestV1(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")
	if err := h.Store.DeleteTestById(r.Context(), id); err != nil {
		return err
	}
	RespondJSON(w, http.StatusNoContent, map[string]string{"result": "deleted"})
	return nil
}
