package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/datdev2409/lab-admin-go/internal/models"
	"github.com/datdev2409/lab-admin-go/internal/service"
	"github.com/datdev2409/lab-admin-go/internal/templates/pages"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type TestHandler struct {
	testService *service.TestService
	validator   *validator.Validate
}

func NewTestHandler(testService *service.TestService, validator *validator.Validate) *TestHandler {
	return &TestHandler{
		testService: testService,
		validator:   validator,
	}
}

// HandleTestPage handles GET /danh-muc-xet-nghiem
func (h *TestHandler) HandleTestPage(w http.ResponseWriter, r *http.Request) error {
	return Render(r.Context(), w, pages.TestPage())
}

// CreateTest handles POST /api/v1/tests
func (h *TestHandler) CreateTest(w http.ResponseWriter, r *http.Request) error {
	var input models.CreateTestInput

	if err := BindAndValidate(r, h.validator, &input); err != nil {
		return err
	}

	test, err := h.testService.CreateTest(r.Context(), &input)
	if err != nil {
		if errors.Is(err, service.ErrTestAlreadyExists) {
			return &AppError{StatusCode: http.StatusConflict, Message: TEST_ALREADY_EXISTS}
		}
		if errors.Is(err, service.ErrInvalidBounds) {
			return &AppError{StatusCode: http.StatusBadRequest, Message: INVALID_TEST_BOUNDS}
		}
		return err
	}

	RespondJSON(w, http.StatusCreated, test)
	return nil
}

// SearchTestsByName handles GET /api/v1/tests?q=keyword&page=1&page_size=10
func (h *TestHandler) SearchTestsByName(w http.ResponseWriter, r *http.Request) error {
	queryOpts := ParseListParams(r, 10) // default page size 10

	keyword := r.URL.Query().Get("q")
	tests, pagination, err := h.testService.SearchTestsByName(
		r.Context(),
		keyword,
		queryOpts.Page,
		queryOpts.PageSize,
	)
	if err != nil {
		return err
	}

	RespondJSONWithPagination(w, http.StatusOK, tests, pagination)
	return nil
}

// GetTest handles GET /api/v1/tests/{id}
func (h *TestHandler) GetTest(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")

	test, err := h.testService.GetTestById(r.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrTestNotFound) {
			return &AppError{StatusCode: http.StatusNotFound, Message: TEST_NOT_FOUND}
		}
		return err
	}

	RespondJSON(w, http.StatusOK, test)
	return nil
}

// UpdateTest handles PATCH /api/v1/tests/{id}
func (h *TestHandler) UpdateTest(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")

	var update models.TestUpdate
	if err := BindAndValidate(r, h.validator, &update); err != nil {
		return err
	}

	test, err := h.testService.UpdateTestById(r.Context(), id, update)
	if err != nil {
		if errors.Is(err, service.ErrTestNotFound) {
			return &AppError{StatusCode: http.StatusNotFound, Message: TEST_NOT_FOUND}
		}
		if errors.Is(err, service.ErrInvalidBounds) {
			return &AppError{StatusCode: http.StatusBadRequest, Message: INVALID_TEST_BOUNDS}
		}
		return err
	}

	RespondJSON(w, http.StatusOK, test)
	return nil
}

// DeleteTest handles DELETE /api/v1/tests/{id}
func (h *TestHandler) DeleteTest(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")

	err := h.testService.DeleteTestById(r.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrTestNotFound) {
			return &AppError{StatusCode: http.StatusNotFound, Message: TEST_NOT_FOUND}
		}
		return err
	}

	RespondJSON(w, http.StatusOK, map[string]string{"result": "deleted"})
	return nil
}

// BulkCreateTests handles POST /api/v1/tests/bulk
func (h *TestHandler) BulkCreateTests(w http.ResponseWriter, r *http.Request) error {
	var inputs []*models.CreateTestInput

	// Decode JSON - inputs must be a pointer for decoding
	if err := json.NewDecoder(r.Body).Decode(&inputs); err != nil {
		return &AppError{StatusCode: http.StatusBadRequest, Message: "Invalid JSON: " + err.Error()}
	}

	// Validate each input
	for i, input := range inputs {
		if err := h.validator.Struct(input); err != nil {
			return &AppError{StatusCode: http.StatusBadRequest, Message: fmt.Sprintf("Validation error at index %d: %v", i, err)}
		}
	}

	if len(inputs) == 0 {
		return &AppError{StatusCode: http.StatusBadRequest, Message: "Request body must contain at least one test"}
	}

	result, err := h.testService.BulkCreateTests(r.Context(), inputs)
	if err != nil {
		return err
	}

	RespondJSON(w, http.StatusCreated, result)
	return nil
}

// ListAllTests handles GET /api/v1/tests/all (get all tests without pagination)
func (h *TestHandler) ListAllTests(w http.ResponseWriter, r *http.Request) error {
	tests, err := h.testService.ListAllTests(r.Context())
	if err != nil {
		return err
	}

	RespondJSON(w, http.StatusOK, tests)
	return nil
}
