package repository

import "github.com/jackc/pgx/v5/pgtype"

// Helper to convert pgtype.Float8 to *float64
func fromFloat8Optional(f pgtype.Float8) *float64 {
	if !f.Valid {
		return nil
	}
	val := f.Float64
	return &val
}

// Helper to convert *float64 to pgtype.Float8
func toFloat8Optional(f *float64) pgtype.Float8 {
	if f == nil {
		return pgtype.Float8{Valid: false}
	}
	return pgtype.Float8{Float64: *f, Valid: true}
}
