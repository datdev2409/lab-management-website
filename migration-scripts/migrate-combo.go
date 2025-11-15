package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// MongoDate represents MongoDB's date format
type MongoDate struct {
	Date time.Time `json:"$date"`
}

// UnmarshalJSON handles MongoDB date format parsing
func (md *MongoDate) UnmarshalJSON(data []byte) error {
	// Handle both formats: { "$date": "ISO-8601" } and plain ISO-8601 string
	var obj map[string]interface{}
	if err := json.Unmarshal(data, &obj); err == nil {
		if dateStr, ok := obj["$date"].(string); ok {
			// Parse ISO-8601 format
			t, err := time.Parse(time.RFC3339Nano, dateStr)
			if err != nil {
				return fmt.Errorf("failed to parse date %q: %w", dateStr, err)
			}
			md.Date = t
			return nil
		}
	}

	// Fallback: try parsing as direct string
	var dateStr string
	if err := json.Unmarshal(data, &dateStr); err == nil {
		t, err := time.Parse(time.RFC3339Nano, dateStr)
		if err != nil {
			return fmt.Errorf("failed to parse date string %q: %w", dateStr, err)
		}
		md.Date = t
		return nil
	}

	return fmt.Errorf("invalid date format")
}

type MongoCombo struct {
	ID        string     `json:"_id"`
	Name      string     `json:"name"`
	TestIDs   []string   `json:"test_ids"`
	CreatedAt *MongoDate `json:"created_at"`
	UpdatedAt *MongoDate `json:"updated_at"`
}

// ConvertTestIDToUUID converts MongoDB test_* ID format to UUID
func ConvertTestIDToUUID(mongoID string) (uuid.UUID, error) {
	if strings.HasPrefix(mongoID, "test_") {
		mongoID = strings.TrimPrefix(mongoID, "test_")
	}

	if id, err := uuid.Parse(mongoID); err == nil {
		return id, nil
	}

	namespace := uuid.MustParse("6ba7b810-9dad-11d1-80b4-00c04fd430c8")
	return uuid.NewSHA1(namespace, []byte(mongoID)), nil
}

func main() {
	// Read combo data from JSON file
	jsonFile := os.Args[1]
	if jsonFile == "" {
		log.Fatal("Usage: go run migrate_combos.go <path-to-combos.json>")
	}

	data, err := os.ReadFile(jsonFile)
	if err != nil {
		log.Fatalf("Failed to read JSON file: %v", err)
	}

	var combos []MongoCombo
	if err := json.Unmarshal(data, &combos); err != nil {
		log.Fatalf("Failed to parse JSON: %v", err)
	}

	fmt.Printf("Found %d combos to migrate\n\n", len(combos))

	ctx := context.Background()

	// Connect to PostgreSQL
	pgPool, err := pgxpool.New(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}
	defer pgPool.Close()

	// Start transaction
	tx, err := pgPool.Begin(ctx)
	if err != nil {
		log.Fatalf("Failed to begin transaction: %v", err)
	}
	defer tx.Rollback(ctx)

	successCount := 0
	failCount := 0
	skippedTests := 0

	for i, combo := range combos {
		// Generate UUID for combo
		comboUUID := uuid.New()

		// Set default timestamps if not provided
		now := time.Now().UTC()
		createdAt := now
		updatedAt := now

		if combo.CreatedAt != nil {
			createdAt = combo.CreatedAt.Date
		}
		if combo.UpdatedAt != nil {
			updatedAt = combo.UpdatedAt.Date
		}

		// Insert combo
		_, err := tx.Exec(ctx,
			`INSERT INTO combos (id, name, created_at, updated_at)
             VALUES ($1, $2, $3, $4)`,
			comboUUID,
			combo.Name,
			createdAt,
			updatedAt,
		)
		if err != nil {
			log.Printf("Failed to insert combo %d (%s): %v", i, combo.ID, err)
			failCount++
			continue
		}

		// Add tests to combo using bulk insert with UNNEST
		if len(combo.TestIDs) > 0 {
			testUUIDs := make([]uuid.UUID, 0, len(combo.TestIDs))

			for _, testID := range combo.TestIDs {
				testUUID, err := ConvertTestIDToUUID(testID)
				if err != nil {
					log.Printf("Failed to convert test ID %s for combo %s: %v", testID, combo.Name, err)
					skippedTests++
					continue
				}
				testUUIDs = append(testUUIDs, testUUID)
			}

			// Bulk insert tests using UNNEST WITH ORDINALITY
			if len(testUUIDs) > 0 {
				_, err := tx.Exec(ctx,
					`INSERT INTO combo_tests (combo_id, test_id, test_order)
                     SELECT $1::uuid as combo_id, test_id, row_number
                     FROM UNNEST($2::uuid[]) WITH ORDINALITY AS t(test_id, row_number)
                     ON CONFLICT DO NOTHING`,
					comboUUID,
					testUUIDs,
				)
				if err != nil {
					log.Printf("Failed to add tests for combo %s: %v", combo.Name, err)
					failCount++
					continue
				}
			}
		}

		successCount++
		if (i+1)%5 == 0 {
			log.Printf("✓ Processed %d combos...", i+1)
		}
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		log.Fatalf("Failed to commit transaction: %v", err)
	}

	fmt.Printf("\n=== Migration Summary ===\n")
	fmt.Printf("Total combos: %d\n", len(combos))
	fmt.Printf("Successful: %d\n", successCount)
	fmt.Printf("Failed: %d\n", failCount)
	fmt.Printf("Skipped tests: %d\n", skippedTests)

	// Verify migration
	var comboCount int
	err = pgPool.QueryRow(ctx, "SELECT COUNT(*) FROM combos").Scan(&comboCount)
	if err != nil {
		log.Fatalf("Failed to verify combos: %v", err)
	}

	var testCount int
	err = pgPool.QueryRow(ctx, "SELECT COUNT(*) FROM combo_tests").Scan(&testCount)
	if err != nil {
		log.Fatalf("Failed to verify combo_tests: %v", err)
	}

	fmt.Printf("✓ Verified: %d combos and %d combo_tests in PostgreSQL\n", comboCount, testCount)
}
