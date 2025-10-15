package storage

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/datdev2409/lab-admin-go/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Storage interface {
	// User
	CreateUser(ctx context.Context, user *models.User) (string, error)
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
	IsUserExists(ctx context.Context, username string) (bool, error)
	// Patient
	InsertPatient(ctx context.Context, patient *models.Patient) (string, error)
	GetPatientById(ctx context.Context, id string) (*models.Patient, error)
	UpdatePatientById(ctx context.Context, id string, update models.PatientUpdate) error
	DeletePatientById(ctx context.Context, id string) error
	SearchPatientByNameOrPhone(ctx context.Context, filterOpts models.PatientQueryOptions, opts models.GenericQueryOptions) ([]*models.Patient, *models.PaginationResponse, error)
	IsPatientExists(ctx context.Context, name, phone string) (bool, error)
	// Doctor
	InsertDoctor(ctx context.Context, doctor *models.Doctor) (string, error)
	GetDoctorById(ctx context.Context, id string) (*models.Doctor, error)
	UpdateDoctorById(ctx context.Context, id string, update models.DoctorUpdate) error
	DeleteDoctorById(ctx context.Context, id string) error
	SearchDoctorByNameOrPhone(ctx context.Context, filterOpts models.DoctorQueryOptions, opts models.GenericQueryOptions) ([]models.Doctor, *models.PaginationResponse, error)
	// Test
	InsertTest(ctx context.Context, test *models.Test) (string, error)
	ListTests(ctx context.Context, filterOpts models.TestQueryOptions, opts models.GenericQueryOptions) ([]*models.Test, *models.PaginationResponse, error)
	GetTestById(ctx context.Context, id string) (*models.Test, error)
	UpdateTestById(ctx context.Context, id string, update map[string]interface{}) error
	DeleteTestById(ctx context.Context, id string) error
	// Combo
	InsertCombo(ctx context.Context, combo *models.Combo) (string, error)
	ListCombos(ctx context.Context, filterOpts models.ComboQueryOptions, opts models.GenericQueryOptions) ([]*models.Combo, *models.PaginationResponse, error)
	GetComboById(ctx context.Context, id string) (*models.Combo, error)
	UpdateComboById(ctx context.Context, id string, update map[string]interface{}) error
	UpdateComboByIdAndReturn(ctx context.Context, id string, update map[string]interface{}) (*models.Combo, error)
	DeleteComboById(ctx context.Context, id string) error
	GetTestsInCombo(ctx context.Context, comboId string) (*models.Combo, []*models.Test, error)
	GetTestsByComboId(ctx context.Context, comboId string) ([]*models.Test, error)
	// Record
	InsertRecord(ctx context.Context, record *models.Record) (string, error)
	ListRecords(ctx context.Context, filters models.RecordQueryOptions, opts models.GenericQueryOptions) ([]*models.Record, *models.PaginationResponse, error)
	GetRecordById(ctx context.Context, id string) (*models.Record, error)
	GetRecordsByIds(ctx context.Context, ids []string) ([]*models.Record, error)
	GetRecordsByPatientId(ctx context.Context, patientId string) ([]*models.Record, error)
	UpdateRecord(ctx context.Context, recordId string, updateRequest models.UpdateRecordRequest) error
	DeleteRecord(ctx context.Context, recordId string) error
	GetRecordsWithRevenue(ctx context.Context, filters models.RecordQueryOptions, opts models.GenericQueryOptions) (*models.ReportResponse, error)
	// Tracking
	InsertTracking(ctx context.Context, tracking *models.Tracking) (string, error)
	ListTrackings(ctx context.Context, filterOpts models.TrackingQueryOptions, opts models.GenericQueryOptions) ([]*models.Tracking, *models.PaginationResponse, error)
	GetTrackingById(ctx context.Context, id string) (*models.Tracking, error)
	DeleteTrackingById(ctx context.Context, id string) error
}

type MongoStorage struct {
	db *mongo.Database
}

func NewMongoStorage(dbClient *mongo.Client) *MongoStorage {
	db := dbClient.Database("labadmin")
	return &MongoStorage{db: db}
}

type PostgresStorage struct {
	pool *pgxpool.Pool
}

func NewPostgresStorage(pool *pgxpool.Pool) *PostgresStorage {
	return &PostgresStorage{pool: pool}
}

var ErrNotImplemented = errors.New("not implemented yet")

// CreateUser creates a new user in PostgreSQL
func (p *PostgresStorage) CreateUser(ctx context.Context, user *models.User) (string, error) {
	query := `
		INSERT INTO users (id, username, password, role, active, created_at, updated_at)
		VALUES (@id, @username, @password, @role, @active, @created_at, @updated_at)
		RETURNING id
	`

	now := time.Now()

	args := pgx.NamedArgs{
		"id":         uuid.New(),
		"username":   user.Username,
		"password":   user.Password,
		"role":       user.Role,
		"active":     user.Active,
		"created_at": now,
		"updated_at": now,
	}

	var returnedID uuid.UUID
	err := p.pool.QueryRow(ctx, query, args).Scan(&returnedID)

	if err != nil {
		return "", fmt.Errorf("failed to create user: %w", err)
	}

	return returnedID.String(), nil
}

// GetUserByUsername retrieves a user by username
func (p *PostgresStorage) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	query := `
		SELECT id, username, password, role, active, created_at, updated_at
		FROM users
		WHERE username = $1 AND active = true
		LIMIT 1
	`

	row, _ := p.pool.Query(ctx, query, username)
	user, err := pgx.CollectExactlyOneRow(row, pgx.RowToStructByName[models.User])
	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("user not found or inactive")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

// IsUserExists checks if a user exists by username
func (p *PostgresStorage) IsUserExists(ctx context.Context, username string) (bool, error) {
	query := `SELECT 1 FROM users WHERE username = $1 LIMIT 1`

	var exists int
	err := p.pool.QueryRow(ctx, query, username).Scan(&exists)

	if err == pgx.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("failed to check if user exists: %w", err)
	}

	return true, nil
}

// Combo methods
func (p *PostgresStorage) InsertCombo(ctx context.Context, combo *models.Combo) (string, error) {
	return "", ErrNotImplemented
}

func (p *PostgresStorage) ListCombos(ctx context.Context, filterOpts models.ComboQueryOptions, opts models.GenericQueryOptions) ([]*models.Combo, *models.PaginationResponse, error) {
	return nil, nil, ErrNotImplemented
}

func (p *PostgresStorage) GetComboById(ctx context.Context, id string) (*models.Combo, error) {
	return nil, ErrNotImplemented
}

func (p *PostgresStorage) UpdateComboById(ctx context.Context, id string, update map[string]interface{}) error {
	return ErrNotImplemented
}

func (p *PostgresStorage) UpdateComboByIdAndReturn(ctx context.Context, id string, update map[string]interface{}) (*models.Combo, error) {
	return nil, ErrNotImplemented
}

func (p *PostgresStorage) DeleteComboById(ctx context.Context, id string) error {
	return ErrNotImplemented
}

func (p *PostgresStorage) GetTestsInCombo(ctx context.Context, comboId string) (*models.Combo, []*models.Test, error) {
	return nil, nil, ErrNotImplemented
}

func (p *PostgresStorage) GetTestsByComboId(ctx context.Context, comboId string) ([]*models.Test, error) {
	return nil, ErrNotImplemented
}

// Record methods
func (p *PostgresStorage) InsertRecord(ctx context.Context, record *models.Record) (string, error) {
	return "", ErrNotImplemented
}

func (p *PostgresStorage) ListRecords(ctx context.Context, filters models.RecordQueryOptions, opts models.GenericQueryOptions) ([]*models.Record, *models.PaginationResponse, error) {
	return nil, nil, ErrNotImplemented
}

func (p *PostgresStorage) GetRecordById(ctx context.Context, id string) (*models.Record, error) {
	return nil, ErrNotImplemented
}

func (p *PostgresStorage) GetRecordsByIds(ctx context.Context, ids []string) ([]*models.Record, error) {
	return nil, ErrNotImplemented
}

func (p *PostgresStorage) GetRecordsByPatientId(ctx context.Context, patientId string) ([]*models.Record, error) {
	return nil, ErrNotImplemented
}

func (p *PostgresStorage) UpdateRecord(ctx context.Context, recordId string, updateRequest models.UpdateRecordRequest) error {
	return ErrNotImplemented
}

func (p *PostgresStorage) DeleteRecord(ctx context.Context, recordId string) error {
	return ErrNotImplemented
}

func (p *PostgresStorage) GetRecordsWithRevenue(ctx context.Context, filters models.RecordQueryOptions, opts models.GenericQueryOptions) (*models.ReportResponse, error) {
	return nil, ErrNotImplemented
}

// Tracking methods
func (p *PostgresStorage) InsertTracking(ctx context.Context, tracking *models.Tracking) (string, error) {
	return "", ErrNotImplemented
}

func (p *PostgresStorage) ListTrackings(ctx context.Context, filterOpts models.TrackingQueryOptions, opts models.GenericQueryOptions) ([]*models.Tracking, *models.PaginationResponse, error) {
	return nil, nil, ErrNotImplemented
}

func (p *PostgresStorage) GetTrackingById(ctx context.Context, id string) (*models.Tracking, error) {
	return nil, ErrNotImplemented
}

func (p *PostgresStorage) DeleteTrackingById(ctx context.Context, id string) error {
	return ErrNotImplemented
}
