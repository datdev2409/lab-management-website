package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/datdev2409/lab-admin-go/internal/logger"
	"github.com/datdev2409/lab-admin-go/internal/models"
	"github.com/datdev2409/lab-admin-go/internal/repository"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

var (
	ErrTestAlreadyExists = errors.New("test already exists")
	ErrInvalidBounds     = errors.New("lower_bound must be less than or equal to upper_bound")
	ErrTestNotFound      = errors.New("test not found")
)

type TestService struct {
	testRepository repository.TestRepository
}

func NewTestService(testRepo repository.TestRepository) *TestService {
	return &TestService{
		testRepository: testRepo,
	}
}

func (s *TestService) CreateTest(ctx context.Context, input *models.CreateTestInput) (*models.Test, error) {
	log := logger.FromCtx(ctx)

	// Validate bounds
	if err := s.validateBoundsOptional(input.LowerBound, input.UpperBound); err != nil {
		log.Warn("Invalid bounds", zap.Error(err))
		return nil, err
	}

	// Check if test name already exists
	exists, err := s.testRepository.IsTestNameExists(ctx, input.Name)
	if err != nil {
		return nil, err
	}
	if exists {
		log.Warn("Test already exists", zap.String("name", input.Name))
		return nil, ErrTestAlreadyExists
	}

	return s.testRepository.InsertTest(ctx, input)
}

// BulkCreateTests creates multiple tests in one operation with resilience
// Returns successful and failed counts with error details
func (s *TestService) BulkCreateTests(ctx context.Context, inputs []*models.CreateTestInput) (*models.BulkCreateTestResult, error) {
	log := logger.FromCtx(ctx)

	if len(inputs) == 0 {
		return &models.BulkCreateTestResult{
			Success: 0,
			Failure: 0,
			Errors:  []models.BulkCreateTestError{},
		}, nil
	}

	result := &models.BulkCreateTestResult{
		Errors: []models.BulkCreateTestError{},
	}

	// Validate all inputs first - collect validation errors
	validationErrors := make(map[int]string)
	for i, input := range inputs {
		if err := s.validateBoundsOptional(input.LowerBound, input.UpperBound); err != nil {
			validationErrors[i] = fmt.Sprintf("invalid bounds: %v", err)
		}
	}

	// Check for duplicate names within the input
	nameIndexMap := make(map[string]int)
	for i, input := range inputs {
		if idx, exists := nameIndexMap[input.Name]; exists && validationErrors[i] == "" {
			validationErrors[i] = fmt.Sprintf("duplicate name with entry at index %d", idx)
		}
		nameIndexMap[input.Name] = i
	}

	// Check if any test names already exist in the database
	for i, input := range inputs {
		if validationErrors[i] != "" {
			continue // Skip if already has validation error
		}
		exists, err := s.testRepository.IsTestNameExists(ctx, input.Name)
		if err != nil {
			log.Error("Failed to check test existence", zap.String("name", input.Name), zap.Error(err))
			validationErrors[i] = fmt.Sprintf("database error: %v", err)
		} else if exists {
			validationErrors[i] = "test with this name already exists"
		}
	}

	// Now insert valid tests one by one and collect results
	for i, input := range inputs {
		if errMsg, exists := validationErrors[i]; exists {
			result.Failure++
			result.Errors = append(result.Errors, models.BulkCreateTestError{
				Index:   i,
				Name:    input.Name,
				Message: errMsg,
			})
			continue
		}

		// Try to insert individual test
		_, err := s.testRepository.InsertTest(ctx, input)
		if err != nil {
			result.Failure++
			result.Errors = append(result.Errors, models.BulkCreateTestError{
				Index:   i,
				Name:    input.Name,
				Message: fmt.Sprintf("insert failed: %v", err),
			})
			log.Warn("Failed to insert test in bulk operation", zap.Int("index", i), zap.String("name", input.Name), zap.Error(err))
		} else {
			result.Success++
		}
	}

	log.Info("Bulk create tests completed", zap.Int("success", result.Success), zap.Int("failure", result.Failure))
	return result, nil
}

func (s *TestService) SearchTestsByName(ctx context.Context, keyword string, page, pageSize int) ([]*models.Test, *models.PaginationResponse, error) {
	return s.testRepository.SearchTestsByName(ctx, keyword, page, pageSize)
}

func (s *TestService) GetTestById(ctx context.Context, id string) (*models.Test, error) {
	return s.getTestByIdWithErrorHandling(ctx, id)
}

func (s *TestService) UpdateTestById(ctx context.Context, id string, update models.TestUpdate) (*models.Test, error) {
	log := logger.FromCtx(ctx)

	// Validate bounds if both are provided or if they're being updated
	if update.LowerBound != nil && update.UpperBound != nil {
		if err := s.validateBoundsOptional(update.LowerBound, update.UpperBound); err != nil {
			log.Warn("Invalid bounds", zap.Error(err))
			return nil, err
		}
	} else if update.LowerBound != nil {
		// Only lower bound is being updated, fetch current upper bound to validate
		existing, err := s.getTestByIdWithErrorHandling(ctx, id)
		if err != nil {
			return nil, err
		}
		if err := s.validateBoundsOptional(update.LowerBound, existing.UpperBound); err != nil {
			log.Warn("Invalid bounds", zap.Error(err))
			return nil, err
		}
	} else if update.UpperBound != nil {
		// Only upper bound is being updated, fetch current lower bound to validate
		existing, err := s.getTestByIdWithErrorHandling(ctx, id)
		if err != nil {
			return nil, err
		}
		if err := s.validateBoundsOptional(existing.LowerBound, update.UpperBound); err != nil {
			log.Warn("Invalid bounds", zap.Error(err))
			return nil, err
		}
	}

	updated, err := s.testRepository.UpdateTestById(ctx, id, update)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Warn("Test not found", zap.String("id", id))
			return nil, ErrTestNotFound
		}
		return nil, err
	}

	return updated, nil
}

func (s *TestService) DeleteTestById(ctx context.Context, id string) error {
	// Check if test exists before deleting using the helper
	_, err := s.getTestByIdWithErrorHandling(ctx, id)
	if err != nil {
		return err
	}

	return s.testRepository.DeleteTestById(ctx, id)
}

func (s *TestService) ListAllTests(ctx context.Context) ([]*models.Test, error) {
	return s.testRepository.ListAllTests(ctx)
}

// getTestByIdWithErrorHandling retrieves a test by ID and converts pgx.ErrNoRows to ErrTestNotFound
func (s *TestService) getTestByIdWithErrorHandling(ctx context.Context, id string) (*models.Test, error) {
	test, err := s.testRepository.GetTestById(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log := logger.FromCtx(ctx)
			log.Warn("Test not found", zap.String("id", id))
			return nil, ErrTestNotFound
		}
		return nil, err
	}
	return test, nil
}

// validateBoundsOptional checks that lower_bound <= upper_bound
// Pointers can be nil and validation passes
func (s *TestService) validateBoundsOptional(lowerBound, upperBound *float64) error {
	// Both nil is allowed
	if lowerBound == nil && upperBound == nil {
		return nil
	}

	// Only one is set is allowed
	if lowerBound == nil || upperBound == nil {
		return nil
	}

	// Both are set, validate the relationship
	if *lowerBound > *upperBound {
		return fmt.Errorf("%w: %f > %f", ErrInvalidBounds, *lowerBound, *upperBound)
	}

	return nil
}
