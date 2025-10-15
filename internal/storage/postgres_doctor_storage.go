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

// InsertDoctor creates a new doctor in PostgreSQL
func (p *PostgresStorage) InsertDoctor(ctx context.Context, doctor *models.Doctor) (string, error) {
	query := `
		INSERT INTO doctors (id, name, phone, address, created_at, updated_at)
		VALUES (@id, @name, @phone, @address, @created_at, @updated_at)
		RETURNING id
	`

	args := StructToNamedArgs(*doctor)

	rows, _ := p.pool.Query(ctx, query, args)
	defer rows.Close()

	returnedID, err := pgx.CollectExactlyOneRow(rows, pgx.RowTo[uuid.UUID])
	if err != nil {
		return "", fmt.Errorf("failed to insert doctor: %w", err)
	}

	return returnedID.String(), nil
}

// GetDoctorById retrieves a doctor by UUID
func (p *PostgresStorage) GetDoctorById(ctx context.Context, id string) (*models.Doctor, error) {
	query := `
		SELECT id, name, phone, address, created_at, updated_at
		FROM doctors
		WHERE id = $1
	`

	doctorUUID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid doctor ID: %w", err)
	}

	rows, _ := p.pool.Query(ctx, query, doctorUUID)
	defer rows.Close()

	doctor, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[models.Doctor])
	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("doctor not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get doctor: %w", err)
	}

	return &doctor, nil
}

// UpdateDoctorById updates a doctor's information
func (p *PostgresStorage) UpdateDoctorById(ctx context.Context, id string, update models.DoctorUpdate) error {
	doctorUUID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid doctor ID: %w", err)
	}

	// Build dynamic update query with named args
	setClauses := []string{}
	args := pgx.NamedArgs{
		"id":         doctorUUID,
		"updated_at": time.Now(),
	}

	setStmt := "updated_at = @updated_at"

	if update.Name != nil {
		setStmt += ", name = @name"
		args["name"] = *update.Name
	}
	if update.Phone != nil {
		setStmt += ", phone = @phone"
		args["phone"] = *update.Phone
	}
	if update.Address != nil {
		setStmt += ", address = @address"
		args["address"] = *update.Address
	}

	if len(setClauses) == 0 {
		return fmt.Errorf("no fields to update")
	}

	query := fmt.Sprintf(`
		UPDATE doctors
		SET %s
		WHERE id = @id
	`, setStmt)

	result, err := p.pool.Exec(ctx, query, args)
	if err != nil {
		return fmt.Errorf("failed to update doctor: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("doctor not found")
	}

	return nil
}

// DeleteDoctorById deletes a doctor by UUID
func (p *PostgresStorage) DeleteDoctorById(ctx context.Context, id string) error {
	doctorUUID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid doctor ID: %w", err)
	}

	query := `DELETE FROM doctors WHERE id = $1`

	result, err := p.pool.Exec(ctx, query, doctorUUID)
	if err != nil {
		return fmt.Errorf("failed to delete doctor: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("doctor not found")
	}

	return nil
}

// SearchDoctorByNameOrPhone searches doctors by name or phone with pagination using named arguments
func (p *PostgresStorage) SearchDoctorByNameOrPhone(ctx context.Context, filterOpts models.DoctorQueryOptions, opts models.GenericQueryOptions) ([]models.Doctor, *models.PaginationResponse, error) {
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
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM doctors %s", whereClauses)
	var total int64
	if err := p.pool.QueryRow(ctx, countQuery, args).Scan(&total); err != nil {
		return nil, nil, fmt.Errorf("failed to count doctors: %w", err)
	}

	// Select with pagination
	args["limit"] = pageSize
	args["offset"] = offset
	selectQuery := fmt.Sprintf(`
		SELECT id, name, phone, address, created_at, updated_at
		FROM doctors
		%s
		ORDER BY created_at DESC
		LIMIT @limit OFFSET @offset
	`, whereClauses)

	rows, err := p.pool.Query(ctx, selectQuery, args)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to search doctors: %w", err)
	}
	defer rows.Close()

	doctors, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.Doctor])
	if err != nil {
		return nil, nil, fmt.Errorf("failed to collect doctors: %w", err)
	}

	pagination := &models.PaginationResponse{
		Total:     int(total),
		Page:      opts.Page,
		PageSize:  pageSize,
		TotalPage: int((total + int64(pageSize) - 1) / int64(pageSize)),
	}

	return doctors, pagination, nil
}
