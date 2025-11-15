package service

import (
	"context"
	"errors"

	"github.com/datdev2409/lab-admin-go/internal/logger"
	"github.com/datdev2409/lab-admin-go/internal/models"
	"github.com/datdev2409/lab-admin-go/internal/repository"
	"go.uber.org/zap"
)

const comboNotFoundMsg = "combo not found"

var (
	ErrComboAlreadyExists = errors.New("combo already exists")
	ErrComboNotFound      = errors.New(comboNotFoundMsg)
	ErrInvalidComboTests  = errors.New("combo must contain at least one test")
)

type ComboService struct {
	comboRepository repository.ComboRepository
}

func NewComboService(comboRepo repository.ComboRepository) *ComboService {
	return &ComboService{
		comboRepository: comboRepo,
	}
}

func (s *ComboService) CreateCombo(ctx context.Context, input *models.CreateComboInput) (string, error) {
	log := logger.FromCtx(ctx)

	// Validate test IDs
	if len(input.TestIDs) == 0 {
		log.Warn("Invalid combo input: no tests provided", zap.String("name", input.Name))
		return "", ErrInvalidComboTests
	}

	// Check if combo name already exists
	exists, err := s.comboRepository.IsComboNameExists(ctx, input.Name)
	if err != nil {
		return "", err
	}
	if exists {
		log.Warn("Combo already exists", zap.String("name", input.Name))
		return "", ErrComboAlreadyExists
	}

	// InsertCombo handles transaction internally and returns combo ID
	comboId, err := s.comboRepository.InsertCombo(ctx, input.Name, input.TestIDs)
	if err != nil {
		log.Error("Failed to create combo", zap.String("name", input.Name), zap.Error(err))
		return "", err
	}

	return comboId, nil
}

func (s *ComboService) SearchCombosByName(ctx context.Context, keyword string, page, pageSize int) ([]*models.Combo, *models.PaginationResponse, error) {
	return s.comboRepository.SearchCombosByName(ctx, keyword, page, pageSize)
}

func (s *ComboService) GetComboById(ctx context.Context, id string) (*models.Combo, error) {
	return s.getComboByIdWithErrorHandling(ctx, id)
}

func (s *ComboService) UpdateComboById(ctx context.Context, id string, update models.ComboUpdate) (*models.Combo, error) {
	log := logger.FromCtx(ctx)

	updated, err := s.comboRepository.UpdateComboById(ctx, id, update)
	if err != nil {
		if err.Error() == comboNotFoundMsg {
			log.Warn("Combo not found", zap.String("id", id))
			return nil, ErrComboNotFound
		}
		return nil, err
	}

	return updated, nil
}

func (s *ComboService) DeleteComboById(ctx context.Context, id string) error {
	// Check if combo exists before deleting using the helper
	_, err := s.getComboByIdWithErrorHandling(ctx, id)
	if err != nil {
		return err
	}

	return s.comboRepository.DeleteComboById(ctx, id)
}

func (s *ComboService) ListAllCombos(ctx context.Context) ([]*models.Combo, error) {
	return s.comboRepository.ListAllCombos(ctx)
}

func (s *ComboService) GetComboTests(ctx context.Context, comboId string) ([]*models.Test, error) {
	// Verify combo exists first
	_, err := s.getComboByIdWithErrorHandling(ctx, comboId)
	if err != nil {
		return nil, err
	}

	return s.comboRepository.GetComboTests(ctx, comboId)
}

// testToTestInCombo converts a full Test model into a TestInCombo (excludes timestamps)
func (s *ComboService) testToTestInCombo(t *models.Test) *models.TestInCombo {
	if t == nil {
		return nil
	}
	return &models.TestInCombo{
		ID:            t.ID,
		Name:          t.Name,
		Price:         t.Price,
		ImportedPrice: t.ImportedPrice,
		NormalValue:   t.NormalValue,
		Unit:          t.Unit,
		LowerBound:    t.LowerBound,
		UpperBound:    t.UpperBound,
	}
}

// GetComboDetails returns combo information together with its tests (tests without timestamps)
func (s *ComboService) GetComboDetails(ctx context.Context, comboId string) (*models.ComboDetailsResponse, error) {
	// Fetch combo and handle not-found translation
	combo, err := s.getComboByIdWithErrorHandling(ctx, comboId)
	if err != nil {
		return nil, err
	}

	// Fetch tests for combo
	tests, err := s.comboRepository.GetComboTests(ctx, comboId)
	if err != nil {
		if err.Error() == comboNotFoundMsg {
			return nil, ErrComboNotFound
		}
		return nil, err
	}

	// Convert tests to TestInCombo
	resultTests := make([]*models.TestInCombo, 0, len(tests))
	for _, t := range tests {
		resultTests = append(resultTests, s.testToTestInCombo(t))
	}

	return &models.ComboDetailsResponse{
		ID:    combo.ID,
		Name:  combo.Name,
		Tests: resultTests,
	}, nil
}

// getComboByIdWithErrorHandling retrieves a combo by ID and converts errors to ErrComboNotFound
func (s *ComboService) getComboByIdWithErrorHandling(ctx context.Context, id string) (*models.Combo, error) {
	combo, err := s.comboRepository.GetComboById(ctx, id)
	if err != nil {
		if err.Error() == comboNotFoundMsg {
			log := logger.FromCtx(ctx)
			log.Warn("Combo not found", zap.String("id", id))
			return nil, ErrComboNotFound
		}
		return nil, err
	}
	return combo, nil
}
