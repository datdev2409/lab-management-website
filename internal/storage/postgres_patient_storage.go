package storage

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/datdev2409/lab-admin-go/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// InsertPatient creates a new patient in PostgreSQL
func (p *PostgresStorage) InsertPatient(ctx context.Context, patient *models.Patient) (string, error) {
	query := `
		INSERT INTO patients (id, name, yob, gender, address, phone, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
	`

	id := uuid.New()
	now := time.Now()

	rows, _ := p.pool.Query(ctx, query,
		id,
		patient.Name,
		patient.YOB,
		patient.Gender,
		patient.Address,
		patient.Phone,
		now,
		now,
	)
	defer rows.Close()

	returnedID, err := pgx.CollectExactlyOneRow(rows, pgx.RowTo[uuid.UUID])
	if err != nil {
		return "", fmt.Errorf("failed to insert patient: %w", err)
	}

	return returnedID.String(), nil
}

// GetPatientById retrieves a patient by UUID
func (p *PostgresStorage) GetPatientById(ctx context.Context, id string) (*models.Patient, error) {
	query := `
		SELECT id, name, yob, gender, address, phone, created_at, updated_at
		FROM patients
		WHERE id = $1
	`

	patientUUID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid patient ID: %w", err)
	}

	rows, _ := p.pool.Query(ctx, query, patientUUID)
	defer rows.Close()

	patient, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[models.Patient])
	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("patient not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get patient: %w", err)
	}

	return &patient, nil
}

// UpdatePatientById updates a patient's information
func (p *PostgresStorage) UpdatePatientById(ctx context.Context, id string, update models.PatientUpdate) error {
	patientUUID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid patient ID: %w", err)
	}

	// Build dynamic update query with named args
	setClauses := []string{}
	args := pgx.NamedArgs{
		"id":         patientUUID,
		"updated_at": time.Now(),
	}

	if update.Name != nil {
		setClauses = append(setClauses, "name = @name")
		args["name"] = *update.Name
	}
	if update.YOB != nil {
		setClauses = append(setClauses, "yob = @yob")
		args["yob"] = *update.YOB
	}
	if update.Gender != nil {
		setClauses = append(setClauses, "gender = @gender")
		args["gender"] = *update.Gender
	}
	if update.Address != nil {
		setClauses = append(setClauses, "address = @address")
		args["address"] = *update.Address
	}
	if update.Phone != nil {
		setClauses = append(setClauses, "phone = @phone")
		args["phone"] = *update.Phone
	}

	if len(setClauses) == 0 {
		return fmt.Errorf("no fields to update")
	}

	// Always update updated_at
	setClauses = append(setClauses, "updated_at = @updated_at")

	query := fmt.Sprintf(`
		UPDATE patients
		SET %s
		WHERE id = @id
	`, strings.Join(setClauses, ", "))

	result, err := p.pool.Exec(ctx, query, args)
	if err != nil {
		return fmt.Errorf("failed to update patient: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("patient not found")
	}

	return nil
}

// DeletePatientById deletes a patient by UUID
func (p *PostgresStorage) DeletePatientById(ctx context.Context, id string) error {
	patientUUID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid patient ID: %w", err)
	}

	query := `DELETE FROM patients WHERE id = $1`

	result, err := p.pool.Exec(ctx, query, patientUUID)
	if err != nil {
		return fmt.Errorf("failed to delete patient: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("patient not found")
	}

	return nil
}

// SearchPatientByNameOrPhone searches patients by name or phone with pagination using similarity scoring
func (p *PostgresStorage) SearchPatientByNameOrPhone(ctx context.Context, filterOpts models.PatientQueryOptions, opts models.GenericQueryOptions) ([]*models.Patient, *models.PaginationResponse, error) {
	pageSize := opts.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}
	offset := (opts.Page - 1) * pageSize
	if offset < 0 {
		offset = 0
	}

	keyword := strings.ToLower(filterOpts.Keyword)
	args := pgx.NamedArgs{
		"keyword": keyword,
	}

	whereClauses := ""
	if keyword != "" {
		whereClauses = "WHERE LOWER(name) LIKE %@keyword% OR phone LIKE %@keyword%"
	}

	// Count total
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM patients %s", whereClauses)
	var total int64
	if err := p.pool.QueryRow(ctx, countQuery, args).Scan(&total); err != nil {
		return nil, nil, fmt.Errorf("failed to count patients: %w", err)
	}

	// Select with pagination
	args["limit"] = pageSize
	args["offset"] = offset
	selectQuery := fmt.Sprintf(`
		SELECT id, name, yob, gender, address, phone, created_at, updated_at
		FROM patients
		%s
		ORDER BY created_at DESC
		LIMIT @limit OFFSET @offset
	`, whereClauses)

	rows, err := p.pool.Query(ctx, selectQuery, args)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to search patients: %w", err)
	}
	defer rows.Close()

	patients, err := pgx.CollectRows(rows, pgx.RowToStructByName[*models.Patient])
	if err != nil {
		return nil, nil, fmt.Errorf("failed to collect patients: %w", err)
	}

	pagination := &models.PaginationResponse{
		Total:     int(total),
		Page:      opts.Page,
		PageSize:  pageSize,
		TotalPage: int((total + int64(pageSize) - 1) / int64(pageSize)),
	}

	return patients, pagination, nil
}

// IsPatientExists checks if a patient exists by name and phone
func (p *PostgresStorage) IsPatientExists(ctx context.Context, name, phone string) (bool, error) {
	query := `SELECT 1 FROM patients WHERE LOWER(name) = LOWER($1) AND phone = $2 LIMIT 1`

	rows, _ := p.pool.Query(ctx, query, name, phone)
	defer rows.Close()

	_, err := pgx.CollectExactlyOneRow(rows, pgx.RowTo[int])
	if err == pgx.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("failed to check patient existence: %w", err)
	}

	return true, nil
}
