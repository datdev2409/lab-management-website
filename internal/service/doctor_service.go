package service

import (
	"context"
	"errors"

	"github.com/datdev2409/lab-admin-go/internal/logger"
	"github.com/datdev2409/lab-admin-go/internal/models"
	"github.com/datdev2409/lab-admin-go/internal/repository"
	"go.uber.org/zap"
)

var ErrDoctorAlreadyExists = errors.New("doctor already exists")

type DoctorService struct {
	doctorRepository repository.DoctorRepository
}

func NewDoctorService(doctorRepo repository.DoctorRepository) *DoctorService {
	return &DoctorService{
		doctorRepository: doctorRepo,
	}
}

func (s *DoctorService) CreateDoctor(ctx context.Context, input *models.CreateDoctorInput) (*models.Doctor, error) {
	log := logger.FromCtx(ctx)
	exists, err := s.doctorRepository.IsDoctorExists(ctx, input.Name, input.Phone)
	if err != nil {
		return nil, err
	}
	if exists {
		log.Warn("Doctor already exists", zap.String("name", input.Name), zap.String("phone", input.Phone))
		return nil, ErrDoctorAlreadyExists
	}
	return s.doctorRepository.InsertDoctor(ctx, input)
}

func (s *DoctorService) SearchDoctorsByKeyword(ctx context.Context, keyword string, page, pageSize int) ([]*models.Doctor, *models.PaginationResponse, error) {
	return s.doctorRepository.SearchDoctorsByKeyword(ctx, keyword, page, pageSize)
}

func (s *DoctorService) GetDoctorById(ctx context.Context, id string) (*models.Doctor, error) {
	return s.doctorRepository.GetDoctorById(ctx, id)
}

func (s *DoctorService) UpdateDoctorById(ctx context.Context, id string, update models.DoctorUpdate) (*models.Doctor, error) {
	return s.doctorRepository.UpdateDoctorById(ctx, id, update)
}

func (s *DoctorService) DeleteDoctorById(ctx context.Context, id string) error {
	return s.doctorRepository.DeleteDoctorById(ctx, id)
}
