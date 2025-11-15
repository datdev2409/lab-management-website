package main

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type MongoTest struct {
	ID            string     `json:"_id"`
	Name          string     `json:"name"`
	Price         int        `json:"price"`
	ImportedPrice int        `json:"imported_price"`
	NormalValue   string     `json:"normal_value"`
	Unit          string     `json:"unit"`
	LowerBound    *float64   `json:"lower_bound"`
	UpperBound    *float64   `json:"upper_bound"`
	CreatedAt     *time.Time `json:"created_at"`
	UpdatedAt     *time.Time `json:"updated_at"`
}

// ConvertTestIDToUUID converts MongoDB test_* ID format to UUID
// Uses SHA-1 deterministic conversion for reproducible results
func ConvertTestIDToUUID(mongoID string) (uuid.UUID, error) {
	// Remove "test_" prefix if present
	if strings.HasPrefix(mongoID, "test_") {
		mongoID = strings.TrimPrefix(mongoID, "test_")
	}

	// Try to parse as UUID directly (if it's already a UUID)
	if id, err := uuid.Parse(mongoID); err == nil {
		return id, nil
	}

	// If not a valid UUID, create a deterministic UUID from the string
	// Use UUID v5 (SHA-1 namespace) for deterministic conversion
	namespace := uuid.MustParse("6ba7b810-9dad-11d1-80b4-00c04fd430c8") // DNS namespace
	return uuid.NewSHA1(namespace, []byte(mongoID)), nil
}

// func main() {
// 	// Read MongoDB test data from JSON file
// 	jsonFile := os.Args[1]
// 	if jsonFile == "" {
// 		log.Fatal("Usage: go run migrate_tests.go <path-to-labadmin.tests.json>")
// 	}

// 	data, err := os.ReadFile(jsonFile)
// 	if err != nil {
// 		log.Fatalf("Failed to read JSON file: %v", err)
// 	}

// 	var tests []MongoTest
// 	if err := json.Unmarshal(data, &tests); err != nil {
// 		log.Fatalf("Failed to parse JSON: %v", err)
// 	}

// 	fmt.Printf("Found %d tests to migrate\n\n", len(tests))

// 	ctx := context.Background()

// 	// Connect to PostgreSQL
// 	pgPool, err := pgxpool.New(ctx, os.Getenv("DATABASE_URL"))
// 	if err != nil {
// 		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
// 	}
// 	defer pgPool.Close()

// 	// Start transaction
// 	tx, err := pgPool.Begin(ctx)
// 	if err != nil {
// 		log.Fatalf("Failed to begin transaction: %v", err)
// 	}
// 	defer tx.Rollback(ctx)

// 	successCount := 0
// 	failCount := 0
// 	idMappings := make(map[string]string) // Map old ID to new UUID

// 	for i, test := range tests {
// 		// Convert MongoDB ID to UUID
// 		testUUID, err := ConvertTestIDToUUID(test.ID)
// 		if err != nil {
// 			log.Printf("Failed to convert ID %s: %v", test.ID, err)
// 			failCount++
// 			continue
// 		}

// 		// Set default timestamps if not provided
// 		now := time.Now()
// 		createdAt := test.CreatedAt
// 		if createdAt == nil {
// 			createdAt = &now
// 		}
// 		updatedAt := test.UpdatedAt
// 		if updatedAt == nil {
// 			updatedAt = &now
// 		}

// 		// Insert test into PostgreSQL
// 		_, err = tx.Exec(ctx,
// 			`INSERT INTO tests (id, name, price, imported_price, normal_value, unit, lower_bound, upper_bound, created_at, updated_at)
//              VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
// 			testUUID,
// 			test.Name,
// 			test.Price,
// 			test.ImportedPrice,
// 			test.NormalValue,
// 			test.Unit,
// 			test.LowerBound,
// 			test.UpperBound,
// 			createdAt,
// 			updatedAt,
// 		)
// 		if err != nil {
// 			log.Printf("Failed to insert test %d (%s): %v", i, test.ID, err)
// 			failCount++
// 			continue
// 		}

// 		// Store mapping for later use in combo migration
// 		idMappings[test.ID] = testUUID.String()

// 		successCount++
// 		if (i+1)%50 == 0 {
// 			log.Printf("✓ Processed %d tests...", i+1)
// 		}
// 	}

// 	// Commit transaction
// 	if err := tx.Commit(ctx); err != nil {
// 		log.Fatalf("Failed to commit transaction: %v", err)
// 	}

// 	fmt.Printf("\n=== Migration Summary ===\n")
// 	fmt.Printf("Total tests: %d\n", len(tests))
// 	fmt.Printf("Successful: %d\n", successCount)
// 	fmt.Printf("Failed: %d\n", failCount)

// 	// Verify migration
// 	var count int
// 	err = pgPool.QueryRow(ctx, "SELECT COUNT(*) FROM tests").Scan(&count)
// 	if err != nil {
// 		log.Fatalf("Failed to verify: %v", err)
// 	}
// 	fmt.Printf("✓ Verified: %d tests in PostgreSQL\n", count)

// 	// Display first few ID mappings for reference
// 	fmt.Println("\n=== Sample ID Mappings (for combo migration) ===")
// 	i := 0
// 	for oldID, newUUID := range idMappings {
// 		if i >= 5 {
// 			break
// 		}
// 		fmt.Printf("%s → %s\n", oldID, newUUID)
// 		i++
// 	}
// 	fmt.Printf("... and %d more mappings\n", len(idMappings)-5)
// }
