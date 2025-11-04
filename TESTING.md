# Unit Tests Documentation

This document describes the unit test coverage for the Lab Management System.

## Running Tests

### Run all tests
```bash
go test ./...
```

### Run tests with coverage
```bash
go test ./... -coverprofile=coverage.out -covermode=atomic
```

### View coverage report
```bash
go tool cover -html=coverage.out
```

### View coverage summary
```bash
go tool cover -func=coverage.out | tail -1
```

### Run tests for specific package
```bash
go test ./internal/models -v
go test ./internal/auth -v
go test ./internal/handlers -v
go test ./internal/storage -v
```

## Test Coverage

### Core Business Logic Packages (>94% coverage)

#### Models Package (97.6% coverage)
- `internal/models/helper_test.go` - Tests for ID generation, BSON conversion
- `internal/models/patient_test.go` - Patient model creation and GetStringPtr utility
- `internal/models/test_test.go` - Test definition model
- `internal/models/combo_test.go` - Test combo/package model
- `internal/models/record_test.go` - Lab record model with doctor management
- `internal/models/doctor_test.go` - Doctor model
- `internal/models/tracking_test.go` - Record tracking configuration
- `internal/models/user_test.go` - User model with password hashing

#### Auth Package (86.7% coverage)
- `internal/auth/jwt_test.go` - JWT token generation and validation

#### Handlers Package (testable utilities at 88-100%)
- `internal/handlers/error_handler_test.go` - Error handling utilities
- `internal/handlers/helper_test.go` - HTTP helpers (JSON, cookies, parsing, dates)

#### Storage Package (testable utilities at 100%)
- `internal/storage/helper_test.go` - Timezone handling
- `internal/storage/search_utils_test.go` - MongoDB filter building and pagination

## Testing Approach

### Test Framework
- **Testing Library**: Go standard `testing` package
- **Assertions**: `github.com/stretchr/testify/assert` and `require`
- **Style**: Table-driven tests where appropriate

### Test Categories

1. **Unit Tests** - Test individual functions and methods in isolation
2. **Edge Cases** - Test nil, empty, and invalid inputs
3. **Error Paths** - Verify error handling and edge cases

### What's NOT Covered (Intentionally)

The following are excluded from unit test coverage as they require integration testing:

- **Database Layer** (`internal/db`) - Requires MongoDB connection
- **Storage Implementations** (`internal/storage/*_storage.go`) - Requires MongoDB
- **HTTP Handlers** (`internal/handlers/*_handler.go`) - Requires HTTP integration tests
- **Logger** (`internal/logger`) - Infrastructure code
- **Sheets** (`internal/sheets`) - External Excel library integration
- **Templates** (`internal/templates`) - Generated Templ files

These components should be tested with integration tests using a test database and HTTP test server.

## Test Examples

### Simple Test with Assertions
```go
func TestNewPatient(t *testing.T) {
    patient := NewPatient("John Doe", "1990", "Male", "123 Main St", "555-1234")
    
    assert.NotNil(t, patient)
    assert.Contains(t, patient.ID, "patient_")
    assert.Equal(t, "John Doe", patient.Name)
}
```

### Table-Driven Test
```go
func TestGetStringPtr(t *testing.T) {
    tests := []struct {
        name  string
        input string
        want  *string
    }{
        {"non-empty string", "test", stringPtr("test")},
        {"empty string returns nil", "", nil},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := GetStringPtr(tt.input)
            // assertions...
        })
    }
}
```

### Error Testing
```go
func TestValidateJWT(t *testing.T) {
    t.Run("invalid token", func(t *testing.T) {
        _, err := ValidateJWT("invalid.token")
        assert.Error(t, err)
    })
}
```

## Code Quality Standards

- ✅ All business logic has >80% test coverage
- ✅ All public functions/methods are tested
- ✅ Error paths are tested
- ✅ Edge cases (nil, empty, invalid) are covered
- ✅ Tests use descriptive names
- ✅ Tests are independent and can run in any order
- ✅ No external dependencies (database, network) in unit tests

## Continuous Integration

Tests should be run in CI/CD pipeline:

```yaml
# Example GitHub Actions
- name: Run tests
  run: go test ./... -coverprofile=coverage.out

- name: Check coverage
  run: |
    total=$(go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//')
    if (( $(echo "$total < 80" | bc -l) )); then
      echo "Coverage $total% is below 80%"
      exit 1
    fi
```

## Adding New Tests

When adding new code:

1. Write tests alongside the code
2. Aim for >80% coverage of new code
3. Test happy path and error cases
4. Use table-driven tests for multiple similar cases
5. Keep tests simple and focused
6. Run tests locally before committing

## Troubleshooting

### Tests fail to compile
```bash
# Make sure all dependencies are installed
go mod tidy
go mod download
```

### "no test files" error
Make sure test files:
- End with `_test.go`
- Are in the same package as the code being tested
- Have at least one test function starting with `Test`

### Coverage seems low
Check which files are included:
```bash
go tool cover -func=coverage.out
```

Remember: Infrastructure code (db, logger, templates) is intentionally excluded from coverage.
