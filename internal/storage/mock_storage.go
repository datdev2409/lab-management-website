package storage

import (
	"context"

	"github.com/datdev2409/lab-admin-go/internal/models"
	"github.com/stretchr/testify/mock"
)

// MockStorage is a mock implementation of the Storage interface for testing
type MockStorage struct {
	mock.Mock
}

// User methods
func (m *MockStorage) CreateUser(ctx context.Context, user *models.User) (string, error) {
	args := m.Called(ctx, user)
	return args.String(0), args.Error(1)
}

func (m *MockStorage) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	args := m.Called(ctx, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

// Patient methods
func (m *MockStorage) InsertPatient(ctx context.Context, patient *models.Patient) (string, error) {
	args := m.Called(ctx, patient)
	return args.String(0), args.Error(1)
}

func (m *MockStorage) GetPatientById(ctx context.Context, id string) (*models.Patient, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Patient), args.Error(1)
}

func (m *MockStorage) UpdatePatientById(ctx context.Context, id string, update models.PatientUpdate) error {
	args := m.Called(ctx, id, update)
	return args.Error(0)
}

func (m *MockStorage) DeletePatientById(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockStorage) SearchPatientByNameOrPhone(ctx context.Context, filterOpts models.PatientQueryOptions, opts models.GenericQueryOptions) ([]*models.Patient, *models.PaginationResponse, error) {
	args := m.Called(ctx, filterOpts, opts)
	if args.Get(0) == nil {
		return nil, nil, args.Error(2)
	}
	return args.Get(0).([]*models.Patient), args.Get(1).(*models.PaginationResponse), args.Error(2)
}

func (m *MockStorage) FindPatientByNameAndPhone(ctx context.Context, name, phone string) (*models.Patient, error) {
	args := m.Called(ctx, name, phone)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Patient), args.Error(1)
}

// Doctor methods
func (m *MockStorage) InsertDoctor(ctx context.Context, doctor *models.Doctor) (string, error) {
	args := m.Called(ctx, doctor)
	return args.String(0), args.Error(1)
}

func (m *MockStorage) GetDoctorById(ctx context.Context, id string) (*models.Doctor, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Doctor), args.Error(1)
}

func (m *MockStorage) UpdateDoctorById(ctx context.Context, id string, update models.DoctorUpdate) error {
	args := m.Called(ctx, id, update)
	return args.Error(0)
}

func (m *MockStorage) DeleteDoctorById(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockStorage) SearchDoctorByNameOrPhone(ctx context.Context, filterOpts models.DoctorQueryOptions, opts models.GenericQueryOptions) ([]*models.Doctor, *models.PaginationResponse, error) {
	args := m.Called(ctx, filterOpts, opts)
	if args.Get(0) == nil {
		return nil, nil, args.Error(2)
	}
	return args.Get(0).([]*models.Doctor), args.Get(1).(*models.PaginationResponse), args.Error(2)
}

func (m *MockStorage) FindDoctorByNameAndPhone(ctx context.Context, name, phone string) (*models.Doctor, error) {
	args := m.Called(ctx, name, phone)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Doctor), args.Error(1)
}

// Test methods
func (m *MockStorage) InsertTest(ctx context.Context, test *models.Test) (string, error) {
	args := m.Called(ctx, test)
	return args.String(0), args.Error(1)
}

func (m *MockStorage) ListTests(ctx context.Context, filterOpts models.TestQueryOptions, opts models.GenericQueryOptions) ([]*models.Test, *models.PaginationResponse, error) {
	args := m.Called(ctx, filterOpts, opts)
	if args.Get(0) == nil {
		return nil, nil, args.Error(2)
	}
	return args.Get(0).([]*models.Test), args.Get(1).(*models.PaginationResponse), args.Error(2)
}

func (m *MockStorage) GetTestById(ctx context.Context, id string) (*models.Test, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Test), args.Error(1)
}

func (m *MockStorage) UpdateTestById(ctx context.Context, id string, update map[string]interface{}) error {
	args := m.Called(ctx, id, update)
	return args.Error(0)
}

func (m *MockStorage) DeleteTestById(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// Combo methods
func (m *MockStorage) InsertCombo(ctx context.Context, combo *models.Combo) (string, error) {
	args := m.Called(ctx, combo)
	return args.String(0), args.Error(1)
}

func (m *MockStorage) ListCombos(ctx context.Context, filterOpts models.ComboQueryOptions, opts models.GenericQueryOptions) ([]*models.Combo, *models.PaginationResponse, error) {
	args := m.Called(ctx, filterOpts, opts)
	if args.Get(0) == nil {
		return nil, nil, args.Error(2)
	}
	return args.Get(0).([]*models.Combo), args.Get(1).(*models.PaginationResponse), args.Error(2)
}

func (m *MockStorage) GetComboById(ctx context.Context, id string) (*models.Combo, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Combo), args.Error(1)
}

func (m *MockStorage) UpdateComboById(ctx context.Context, id string, update map[string]interface{}) error {
	args := m.Called(ctx, id, update)
	return args.Error(0)
}

func (m *MockStorage) UpdateComboByIdAndReturn(ctx context.Context, id string, update map[string]interface{}) (*models.Combo, error) {
	args := m.Called(ctx, id, update)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Combo), args.Error(1)
}

func (m *MockStorage) DeleteComboById(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockStorage) GetTestsInCombo(ctx context.Context, comboId string) (*models.Combo, []*models.Test, error) {
	args := m.Called(ctx, comboId)
	if args.Get(0) == nil {
		return nil, nil, args.Error(2)
	}
	return args.Get(0).(*models.Combo), args.Get(1).([]*models.Test), args.Error(2)
}

func (m *MockStorage) GetTestsByComboId(ctx context.Context, comboId string) ([]*models.Test, error) {
	args := m.Called(ctx, comboId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Test), args.Error(1)
}

// Record methods
func (m *MockStorage) InsertRecord(ctx context.Context, record *models.Record) (string, error) {
	args := m.Called(ctx, record)
	return args.String(0), args.Error(1)
}

func (m *MockStorage) ListRecords(ctx context.Context, filters models.RecordQueryOptions, opts models.GenericQueryOptions) ([]*models.Record, *models.PaginationResponse, error) {
	args := m.Called(ctx, filters, opts)
	if args.Get(0) == nil {
		return nil, nil, args.Error(2)
	}
	return args.Get(0).([]*models.Record), args.Get(1).(*models.PaginationResponse), args.Error(2)
}

func (m *MockStorage) GetRecordById(ctx context.Context, id string) (*models.Record, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Record), args.Error(1)
}

func (m *MockStorage) GetRecordsByIds(ctx context.Context, ids []string) ([]*models.Record, error) {
	args := m.Called(ctx, ids)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Record), args.Error(1)
}

func (m *MockStorage) GetRecordsByPatientId(ctx context.Context, patientId string) ([]*models.Record, error) {
	args := m.Called(ctx, patientId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Record), args.Error(1)
}

func (m *MockStorage) UpdateRecord(ctx context.Context, recordId string, updateRequest models.UpdateRecordRequest) error {
	args := m.Called(ctx, recordId, updateRequest)
	return args.Error(0)
}

func (m *MockStorage) DeleteRecord(ctx context.Context, recordId string) error {
	args := m.Called(ctx, recordId)
	return args.Error(0)
}

func (m *MockStorage) GetRecordsWithRevenue(ctx context.Context, filters models.RecordQueryOptions, opts models.GenericQueryOptions) (*models.ReportResponse, error) {
	args := m.Called(ctx, filters, opts)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.ReportResponse), args.Error(1)
}

// Tracking methods
func (m *MockStorage) InsertTracking(ctx context.Context, tracking *models.Tracking) (string, error) {
	args := m.Called(ctx, tracking)
	return args.String(0), args.Error(1)
}

func (m *MockStorage) ListTrackings(ctx context.Context, filterOpts models.TrackingQueryOptions, opts models.GenericQueryOptions) ([]*models.Tracking, *models.PaginationResponse, error) {
	args := m.Called(ctx, filterOpts, opts)
	if args.Get(0) == nil {
		return nil, nil, args.Error(2)
	}
	return args.Get(0).([]*models.Tracking), args.Get(1).(*models.PaginationResponse), args.Error(2)
}

func (m *MockStorage) GetTrackingById(ctx context.Context, id string) (*models.Tracking, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Tracking), args.Error(1)
}

func (m *MockStorage) DeleteTrackingById(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
