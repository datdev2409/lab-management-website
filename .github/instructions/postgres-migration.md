# PostgreSQL Migration Guide - MongoDB to PostgreSQL with Clean Architecture

This guide documents the step-by-step process for migrating entities from MongoDB to PostgreSQL using the **handler → service → repository → sqlc** pattern.

## Architecture Overview

```
HTTP Request
    ↓
Handler (presentation layer - internal/handlers/)
    ↓
Service (business logic layer - internal/service/)
    ↓
Repository (data access interface - internal/repository/)
    ↓
sqlc.Queries (PostgreSQL via pgx/v5 - internal/db/sqlc/)
```

## Prerequisites

- PostgreSQL database running and accessible
- `goose` CLI installed for migrations
- `sqlc` CLI installed for code generation
- `DATABASE_URL` environment variable set

## Reference Implementation

**See the complete working example**:

- `internal/repository/patient_repository.go` - Repository layer
- `internal/service/patient_service.go` - Service layer
- `internal/handlers/patient_handler.go` - Handler layer
- `internal/db/migrations/20251107130607_initial_schema_setup.sql` - Migration script
- `internal/db/queries/patients.sql` - sqlc queries
- `internal/models/patient.go` - Domain models with Input/Update DTOs

## Migration Steps for Each Entity

### Step 1: Create PostgreSQL Migration Script

**Location**: `internal/db/migrations/`

**Create Migration File**:

```bash
# Navigate to project root
# Create migration with descriptive name
goose create entity_name_setup sql

# This creates: YYYYMMDDHHMMSS_entity_name_setup.sql
```

**Key Points**:

- Use UUID primary keys with `gen_random_uuid()`
- Add NOT NULL constraints based on domain model in `internal/models/`
- Create indexes for search fields, date filtering, and foreign keys
- Include `created_at` and `updated_at` TIMESTAMPTZ columns

**Reference**: See `internal/db/migrations/20251107130607_initial_schema_setup.sql` for patients table example

**Commands**:

```bash
# Run migration
goose -dir internal/db/migrations postgres $DATABASE_URL up

# Verify migration
goose -dir internal/db/migrations postgres $DATABASE_URL status
```

### Step 2: Create sqlc Queries

**Location**: `internal/db/queries/`

**Naming Convention**: `entity_name.sql` (e.g., `patients.sql`, `doctors.sql`)

**Common Query Types** (adapt based on entity needs):

- Create (INSERT ... RETURNING \*)
- GetByID (SELECT with UUID)
- Search with pagination (ILIKE with LIMIT/OFFSET)
- Count for pagination
- Existence check (SELECT EXISTS)
- Update (UPDATE ... SET with updated_at = NOW())
- Delete (DELETE by ID)

**Reference**: See `internal/db/queries/patients.sql` for complete query examples

**Commands**:

```bash
# Generate sqlc code
sqlc generate

# Verify generated code
ls internal/db/sqlc/entity_name.sql.go
```

### Step 3: Create Repository Layer

**Location**: `internal/repository/entity_repository.go`

**Structure**:

- Define interface with repository methods
- Implement `PgEntityRepository` struct with `*sqlc.Queries`
- Create constructor `NewPgEntityRepository(queries *sqlc.Queries)`
- Implement each interface method calling sqlc queries
- Add `ToDomainEntity()` mapper to convert sqlc types to domain models

**Key Patterns**:

- Use `parseUUID()` helper for ID validation
- For partial updates: fetch existing, merge changes, update, return fresh record
- Convert UUID to string when returning domain models
- Calculate pagination with `math.Ceil(float64(total) / float64(pageSize))`

**Reference**: See `internal/repository/patient_repository.go` for complete implementation

### Step 4: Create Service Layer

**Location**: `internal/service/entity_service.go`

**Structure**:

- Create service struct with repository dependency
- Implement business logic methods
- Add duplicate checks, validation, or orchestration as needed
- Define custom error types (e.g., `ErrEntityAlreadyExists`)
- Use logger for important business events

**Responsibilities**:

- Business validation and rules
- Coordinate multiple repository calls
- Error translation to domain errors
- Logging business events

**Reference**: See `internal/service/patient_service.go` for implementation example

### Step 5: Migrate Handler Layer

**Location**: `internal/handlers/entity_handler.go`

**Structure**:

- Create handler struct with service dependency and validator
- Implement HTTP handler methods for each endpoint
- Use `BindAndValidate()` for input validation
- Extract URL params with `chi.URLParam(r, "id")`
- Parse query params with `ParseListParams()` for pagination
- Handle service errors and return appropriate HTTP status codes
- Use `RespondJSON()` or `RespondJSONWithPagination()` for responses

**Important Guidelines**:

- Import `github.com/go-chi/chi/v5` (not v4)
- Use PATCH for partial updates (not PUT)
- Check for specific service errors and return `AppError` with status codes
- Ensure all old MongoDB endpoints are migrated

**Reference**: See `internal/handlers/patient_handler.go` for complete implementation

### Step 6: Update Domain Models

**Location**: `internal/models/entity.go`

**Add Input/Update Models**:

- `CreateEntityInput` - For creation requests with validation tags
- `EntityUpdate` - For partial updates with pointer fields (omitempty)
- Update existing model ID field to string type (for UUID)

**Reference**: See `internal/models/patient.go` for examples

### Step 7: Wire Up in Main and Handler

#### Update `cmd/api/main.go`:

```go
// Initialize repository, service, and handler
queries := sqlc.New(pgPool)
entityRepository := repository.NewPgEntityRepository(queries)
entityService := service.NewEntityService(entityRepository)
entityHandler := handlers.NewEntityHandler(entityService, v)

// Register routes
r.Route("/api/v1/entities", func(r chi.Router) {
	r.Get("/", handlers.Make(entityHandler.SearchEntitiesByKeyword))
	r.Post("/", handlers.Make(entityHandler.CreateEntity))
	r.Get("/{id}", handlers.Make(entityHandler.GetEntity))
	r.Patch("/{id}", handlers.Make(entityHandler.UpdateEntity))
	r.Delete("/{id}", handlers.Make(entityHandler.DeleteEntity))
})
```

#### Update `internal/handlers/handler.go`:

```go
type Handler struct {
	Router        http.Handler
	Validator     *validator.Validate
	Store         storage.Storage
	entityHandler *EntityHandler
	// ... other handlers
}

func NewHandler(store storage.Storage, log *zap.Logger, entityHandler *EntityHandler) *Handler {
	// ... initialization
	h := &Handler{
		Router:        r,
		Store:         store,
		Validator:     v,
		entityHandler: entityHandler,
	}

	// Register routes
	r.Route("/api/v1/entities", func(r chi.Router) {
		r.Get("/", Make(h.entityHandler.SearchEntitiesByKeyword))
		r.Post("/", Make(h.entityHandler.CreateEntity))
		r.Get("/{id}", Make(h.entityHandler.GetEntity))
		r.Patch("/{id}", Make(h.entityHandler.UpdateEntity))
		r.Delete("/{id}", Make(h.entityHandler.DeleteEntity))
	})
}
```

### Step 8: Update Error Constants

**Location**: `internal/handlers/error_msg.go`

```go
const (
	// ... existing constants
	ENTITY_ALREADY_EXISTS = "Entity with this field1 and field2 already exists"
	ENTITY_NOT_FOUND      = "Entity not found"
)
```

### Step 9: Create Bruno API Documentation

**Location**: `docs/bruno/lab-admin-web/entity_name/`

Create a folder for the entity with Bruno API test files for all endpoints:

- `folder.bru` - Folder metadata
- `Create Entity.bru` - POST endpoint
- `Search entities.bru` - GET list with pagination
- `Get entity by ID.bru` - GET single entity
- `Update entity by ID.bru` - PATCH endpoint
- `Delete entity by ID.bru` - DELETE endpoint

**Reference**: See `docs/bruno/lab-admin-web/patients/` for complete examples

## Common Patterns and Helpers

### UUID Parsing Helper (in repository files)

```go
func parseUUID(id string) (uuid.UUID, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid UUID: %w", err)
	}
	return uid, nil
}
```

### Pagination Calculation

```go
pagination := &models.PaginationResponse{
	Total:     int(total),
	TotalPage: int(math.Ceil(float64(total) / float64(pageSize))),
	Page:      page,
	PageSize:  pageSize,
}
```

### Input Validation and Binding

```go
if err := BindAndValidate(r, h.validator, &input); err != nil {
	return err
}
```

## Testing Strategy

1. **Repository Tests**: Test sqlc query functions with real PostgreSQL
2. **Service Tests**: Mock repository interface, test business logic
3. **Handler Tests**: Test HTTP endpoints with mocked service
4. **Integration Tests**: End-to-end tests with test database

## Verification Checklist

After migrating each entity, verify:

- [ ] Migration script runs without errors (`goose up`)
- [ ] sqlc generates code without errors (`sqlc generate`)
- [ ] Repository implements all interface methods
- [ ] Service layer includes business logic (if needed)
- [ ] Handler methods use correct HTTP methods (PATCH for updates)
- [ ] All old MongoDB endpoints are migrated
- [ ] Error messages are user-friendly (Vietnamese)
- [ ] chi/v5 is used consistently (not chi/v4)
- [ ] UUID parsing is handled correctly
- [ ] Pagination works correctly
- [ ] Input validation is present
- [ ] Bruno API documentation created with all endpoints
- [ ] Code compiles: `go build ./...`

## Common Pitfalls to Avoid

1. **Chi Version Mismatch**: Always import `github.com/go-chi/chi/v5`
2. **PUT vs PATCH**: Use PATCH for partial updates, not PUT
3. **Missing UUID Conversion**: Always convert UUID to string when returning to handlers
4. **Forgetting Indexes**: Create indexes based on query patterns
5. **Not Handling Partial Updates**: Fetch existing record before updating
6. **Missing Error Translation**: Convert service errors to appropriate HTTP status codes
7. **Hardcoded Pagination**: Use `ParseListParams()` helper
8. **Missing Validation**: Always validate input using `BindAndValidate()`

## Entity-Specific Considerations

### Entities with Relationships

- Add foreign key columns in migration
- Create indexes on foreign keys
- Handle cascade operations in service layer

### Entities with Complex Queries

- Create specialized query methods in sqlc
- Implement custom repository methods
- Consider using PostgreSQL JSON fields for flexible data

### Entities with File Uploads

- Store file metadata in database
- Handle file operations in service layer
- Keep file paths relative for portability

## Migration Order Recommendation

Suggested order to minimize dependencies:

1. **Patients** ✅ (Complete - reference implementation)
2. **Doctors** (Similar to patients)
3. **Tests** (No dependencies)
4. **Combos** (Depends on Tests - has test_ids array)
5. **Records** (Depends on Patients, Doctors, Combos - most complex)
6. **Trackings** (Depends on Tests)
7. **Users** (Authentication - can be done anytime)

Each entity should be fully migrated (steps 1-8) and tested before moving to the next.
