package service

import (
	"context"
	"errors"

	"github.com/datdev2409/lab-admin-go/internal/logger"
	"github.com/datdev2409/lab-admin-go/internal/models"
	"github.com/datdev2409/lab-admin-go/internal/repository"
	"go.uber.org/zap"
)

var ErrPatientAlreadyExists = errors.New("patient already exists")

type PatientService struct {
	patientRepository repository.PatientRepository
}

func NewPatientService(patientRepo repository.PatientRepository) *PatientService {
	return &PatientService{
		patientRepository: patientRepo,
	}
}

func (s *PatientService) CreatePatient(ctx context.Context, patient *models.CreatePatientInput) (*models.Patient, error) {
	log := logger.FromCtx(ctx)
	exists, err := s.patientRepository.IsPatientExists(ctx, patient.Name, patient.Phone)
	if err != nil {
		return nil, err
	}
	if exists {
		log.Warn("Patient already exists", zap.String("name", patient.Name), zap.String("phone", patient.Phone))
		return nil, ErrPatientAlreadyExists
	}
	return s.patientRepository.InsertPatient(ctx, patient)
}

func (s *PatientService) SearchPatientsByKeyword(ctx context.Context, keyword string, page, pageSize int) ([]*models.Patient, *models.PaginationResponse, error) {
	return s.patientRepository.SearchPatientsByKeyword(ctx, keyword, page, pageSize)
}

func (s *PatientService) GetPatientById(ctx context.Context, id string) (*models.Patient, error) {
	return s.patientRepository.GetPatientById(ctx, id)
}

func (s *PatientService) UpdatePatientById(ctx context.Context, id string, update models.PatientUpdate) (*models.Patient, error) {
	return s.patientRepository.UpdatePatientById(ctx, id, update)
}

func (s *PatientService) DeletePatientById(ctx context.Context, id string) error {
	return s.patientRepository.DeletePatientById(ctx, id)
}
