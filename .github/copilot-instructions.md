# Lab Management System - Copilot Instructions

## Project Overview

This is a **Lab Management System** for **Anh Quan Laboratory** built with Go, HTMX, Alpine.js, and PostgreSQL. The system manages laboratory test records, patients, test definitions, test combos, and tracking/comparison features.

## Technology Stack

- **Backend**: Go 1.24.7 with Chi router
- **Frontend**: HTMX + Alpine.js + Bootstrap (server-side rendered with Templ). Please note that we are migrating the HTMX parts to Alpine.js gradually, so prefer using Alpine.js over HTMX.
- **Database**: PostgreSQL with sqlc, pgx/v5, goose migrations
- **Authentication**: JWT with HTTP-only cookies
- **Report Generation**: Excel files using excelize, PDF conversion with Gotenberg
- **Deployment**: Docker, Docker Compose, systemd service
- **Monitoring**: Traefik reverse proxy, structured logging with Zap

## PostgreSQL Migration

⚠️ **Active Migration**: The system is being migrated from MongoDB to PostgreSQL using the **handler → service → repository → sqlc** architecture pattern.

📖 **Complete Migration Guide**: See [instructions/postgres-migration.md](instructions/postgres-migration.md) for detailed step-by-step instructions on migrating each entity.

**Migration Status**:

- ✅ **Patients**: Fully migrated (reference implementation)
- ✅ **Doctors**: Fully migrated
- ⏳ **Tests**: Pending
- ⏳ **Combos**: Pending
- ⏳ **Records**: Pending
- ⏳ **Trackings**: Pending
- ⏳ **Users**: Pending

**New Architecture Pattern**:

```
HTTP Request → Handler → Service → Repository → sqlc.Queries → PostgreSQL
```

## Architecture Overview

### Project Structure

```
├── cmd/api/main.go                 # Application entry point
├── internal/
│   ├── auth/jwt.go                 # JWT authentication
│   ├── db/database.go              # MongoDB connection
│   ├── handlers/                   # HTTP handlers (controllers)
│   │   ├── handler.go              # Main router setup
│   │   ├── auth_handler.go         # Authentication endpoints
│   │   ├── patient_handler.go      # Patient management
│   │   ├── record_handler.go       # Lab record management
│   │   ├── test_handler.go         # Test definition management
│   │   ├── combo_handler.go        # Test combo management
│   │   ├── tracking_handler.go     # Record tracking/comparison
│   │   ├── midlewares.go           # JWT auth, logging middleware
│   │   └── helper.go               # Utility functions
│   ├── models/                     # Data models and DTOs
│   │   ├── patient.go              # Patient model
│   │   ├── record.go               # Lab record model
│   │   ├── test.go                 # Test definition model
│   │   ├── combo.go                # Test combo model
│   │   ├── tracking.go             # Tracking model
│   │   └── user.go                 # User/auth model
│   ├── storage/                    # Data access layer
│   │   ├── storage.go              # Storage interface
│   │   ├── base.go                 # Generic MongoDB operations
│   │   ├── patient_storage.go      # Patient CRUD
│   │   ├── record_storage.go       # Record CRUD
│   │   ├── test_storage.go         # Test CRUD
│   │   ├── combo_storage.go        # Combo CRUD
│   │   └── tracking_storage.go     # Tracking CRUD
│   ├── sheets/                     # Excel/PDF report generation
│   ├── templates/                  # Templ templates
│   │   ├── pages/                  # Full page templates
│   │   └── partials/               # Reusable components
│   └── logger/                     # Structured logging
├── templates/                      # Excel templates
├── reports/                        # Generated reports
└── deploy/                         # Deployment configs
```

### Core Domain Models

#### Patient

- **ID**: Custom string ID (patient\_\*)
- **Fields**: Name, YOB, Gender, Address, Phone
- **Relationships**: Has many Records

#### Test (Test Definition)

- **ID**: Custom string ID (test\_\*)
- **Fields**: Name, Price, NormalValue, Unit, LowerBound, UpperBound
- **Purpose**: Template for test types

#### Combo (Test Package)

- **ID**: Custom string ID (combo\_\*)
- **Fields**: Name, TestIDs[]
- **Purpose**: Groups multiple tests together

#### Record (Lab Test Record)

- **ID**: Custom string ID (record\_\*)
- **Fields**: Patient (embedded), ComboName, TestResults[], Status, timestamps
- **Status**: "pending", "completed"
- **TestResults**: Array of test results with actual values

#### Tracking

- **ID**: Custom string ID (tracking\_\*)
- **Purpose**: Define which tests to compare across multiple records
- **Fields**: Name, Tests[] (test configurations for comparison)

### Database Collections

- **patients**: Patient documents
- **tests**: Test definition documents
- **combos**: Test combo documents
- **records**: Lab record documents
- **trackings**: Tracking configuration documents
- **users**: User authentication documents

## Key Features & Workflows

### 1. Patient Management

- Create, update, delete patients
- Search by name or phone number
- Auto-complete during record creation

### 2. Test Management

- Define test types with normal ranges
- Set pricing and units
- Manage test metadata

### 3. Combo Management

- Group tests into packages
- Reusable test combinations
- Used when creating records

### 4. Record Management (Core Feature)

- Create lab test records for patients
- Select test combo or individual tests
- Enter test results with automatic abnormal detection
- Status tracking (pending → completed)
- Generate multiple report types

### 5. Record Tracking & Comparison

- Create tracking configurations
- Compare test results across multiple records
- Useful for monitoring patient progress over time

### 6. Report Generation

- **Billing Report** (`phieu_thu`): Invoice/payment receipt
- **Results Report** (`phieu_ket_qua`): Lab results
- **Results with Signature** (`phieu_ket_qua_chu_ky`): Signed results
- **Results PDF** (`phieu_ket_qua_chu_ky_pdf`): PDF version
- **Tracking Report** (`phieu_theo_doi`): Comparison report

## API Design Patterns

### Route Structure

```
Web Pages (SSR):
- GET / → Record management page
- GET /phieu-xet-nghiem → Record list
- GET /phieu-xet-nghiem/new → Create new record
- GET /danh-muc-benh-nhan → Patient management
- GET /danh-muc-xet-nghiem → Test management
- GET /danh-muc-goi-xet-nghiem → Combo management
- GET /so-sanh-ket-qua → Tracking comparison
- GET /danh-muc-so-sanh → Tracking management

API Endpoints:
- /api/v1/* → RESTful API with JWT auth
- /api/reports/* → Report generation
- /api/tracking/* → Tracking operations
```

### Authentication Flow

1. **Login**: POST `/api/v1/auth/login` → Sets HTTP-only auth_token cookie
2. **Registration**: POST `/api/v1/auth/register`
3. **Middleware**: `JWTAuthWebEndpoint` for web pages, `JWTAuthAPIEndpoint` for API
4. **Token Validation**: Extracts userID from JWT claims

### Data Flow Patterns

1. **HTMX → Alpine.js Migration**: Frontend gradually migrating from HTMX to Alpine.js - prefer Alpine.js for new components
2. **Templ Rendering**: Server-side HTML generation with Go templates
3. **Alpine.js**: Client-side reactivity and state management
4. **JSON API**: RESTful endpoints for data operations

## Coding Conventions

### Error Handling

```go
// Always use context-aware logging
logger.FromCtx(ctx).Error("Failed to create record", zap.Error(err))

// Use Make wrapper for error handling in handlers
r.Get("/endpoint", Make(h.HandlerFunction))

// Use AppError for business logic errors with HTTP status codes
func (h *Handler) SomeHandler(w http.ResponseWriter, r *http.Request) error {
    if someCondition {
        return &AppError{StatusCode: 404, Message: "Resource not found"}
    }
    return WriteJSON(w, http.StatusOK, data)
}

// Return appropriate HTTP status codes
http.Error(w, "Invalid request", http.StatusBadRequest)

// Use custom error types when needed
var ErrUserNotFound = errors.New("user not found")
```

### Database Operations

```go
// Use generic MongoDB functions
result, err := MongoInsert(ctx, col, entity)
entity, err := MongoGetById[Model](ctx, col, id)
entities, pagination, err := MongoList[Model](ctx, col, filter, opts)

// Always use context
func (m *MongoStorage) CreateEntity(ctx context.Context, entity *Model) (string, error) {
    col := m.getCollection("entities")
    return MongoInsert(ctx, col, entity)
}
```

### MongoDB v2 Indexing Strategy

Based on query pattern analysis, create these indexes for optimal performance:

**Essential Indexes (High Priority):**

```javascript
// Records collection - date filtering & embedded patient search
db.records.createIndex({ created_at: -1 }); // Date filtering & sorting
db.records.createIndex({ "patient.name": 1, "patient.phone": 1 }); // Patient search in records
db.records.createIndex({ "patient._id": 1 }); // Patient lookup
db.records.createIndex({ status: 1 }); // Status filtering

// Patients collection - frequent searches
db.patients.createIndex({ name: 1, phone: 1 }); // Compound search (primary use case)
db.patients.createIndex({ name: "text", phone: "text", address: "text" }); // Text search fallback

// Users collection - authentication
db.users.createIndex({ username: 1 }, { unique: true }); // Login queries + uniqueness

// Tests collection - name-based search
db.tests.createIndex({ name: 1 }); // Exact & prefix matching for autocomplete

// Combos collection - name-based search
db.combos.createIndex({ name: 1 }); // Exact & prefix matching for autocomplete

// Trackings collection - name-based search
db.trackings.createIndex({ name: 1 }); // Name search functionality
```

**Evidence for indexes:**

- **Records patient search**: `ListRecords()` uses `$or` with `$regex` on `patient.name` and `patient.phone` (lines 19-23)
- **Records date filtering**: Date range filtering with `$gte` and `$lte` on `created_at` (lines 34, 37)
- **Patient compound search**: `FindPatientByNameAndPhone()` queries both fields together (line 15)
- **Test/Combo name search**: `ListTests()` and `ListCombos()` use `$regex` on `name` field for autocomplete
- **Authentication**: `GetUserByUsername()` queries by username for login (line 23)

**Index Implementation Guide:**

```javascript
// Connect to your MongoDB instance
use labadmin

// 1. Records Collection Indexes
db.records.createIndex({ "created_at": -1 })
db.records.createIndex({ "patient.name": 1, "patient.phone": 1 })
db.records.createIndex({ "patient._id": 1 })
db.records.createIndex({ "status": 1 })

// 2. Patients Collection Indexes
db.patients.createIndex({ "name": 1, "phone": 1 })
db.patients.createIndex({ "name": "text", "phone": "text", "address": "text" })

// 3. Users Collection Index
db.users.createIndex({ "username": 1 }, { unique: true })

// 4. Tests Collection Index
db.tests.createIndex({ "name": 1 })

// 5. Combos Collection Index
db.combos.createIndex({ "name": 1 })

// 6. Trackings Collection Index
db.trackings.createIndex({ "name": 1 })

// Verify indexes were created
db.records.getIndexes()
db.patients.getIndexes()
db.users.getIndexes()
db.tests.getIndexes()
db.combos.getIndexes()
db.trackings.getIndexes()
```

**Index Monitoring & Maintenance:**

```javascript
// Check index usage statistics
db.records.aggregate([{ $indexStats: {} }]);

// Analyze query performance
db.records
  .find({
    "patient.name": /john/i,
  })
  .explain("executionStats");

// Drop index if needed (be careful!)
// db.collection.dropIndex("indexName")
```

### Handler Patterns

```go
// Use Make wrapper for error handling
r.Get("/endpoint", Make(h.HandlerFunction))

// Templ rendering
func (h *Handler) HandlePage(w http.ResponseWriter, r *http.Request) error {
    return Render(r.Context(), w, pages.PageTemplate())
}

// JSON responses
func (h *Handler) APIEndpoint(w http.ResponseWriter, r *http.Request) error {
    return WriteJSON(w, http.StatusOK, response)
}
```

### ID Generation

- All entities use custom string IDs with prefixes
- Format: `{type}_{randomString}` (e.g., `patient_abc123`)
- Generated via `GenerateRandomID(prefix string)`

### Time Handling

- Use Vietnam timezone: `Asia/Ho_Chi_Minh`
- Helper: `GetCurrentTime()` in storage package
- Auto-update `updated_at` in update operations

### Abnormal Test Result Detection

- **Automatic Detection**: System automatically marks test results as abnormal for numeric values outside the test's `lower_bound` and `upper_bound` range
- **Non-numeric Results**: Ignored by automatic detection (no abnormal marking)
- **Manual Override**: Doctors can manually mark any test result as abnormal/normal using `manual_abnormal_override` flag
- **Priority**: Manual override takes precedence over automatic detection
- **UI Pattern**: Frontend provides toggle between automatic and manual modes with clear visual indicators

## Development Guidelines

### Adding New Features

1. **Models**: Define in `internal/models/`
2. **Storage**: Add CRUD operations in `internal/storage/`
3. **Handlers**: Create HTTP handlers in `internal/handlers/`
4. **Templates**: Add Templ templates in `internal/templates/`
5. **Routes**: Register in `handler.go`

### Testing Strategy

- **Unit tests**: Storage layer functions
- **Integration tests**: HTTP handlers
- **E2E tests**: Cypress in `tests/` directory
- **Pre-commit workflow**: `./scripts/pre-commit.sh` runs:
  - `templ generate` for templates
  - `go fmt` and `go vet` for code formatting
  - `golangci-lint run` for linting
  - Unit tests with Go test runner
- **Test commands**:
  - `go test ./...` for all tests
  - `npm test` in `tests/` for Cypress

### Development Workflow

- **Live Development**: `make live` runs 3 concurrent processes:
  - `make live/server`: Air for Go hot reload
  - `make live/templ`: Templ template generation with proxy
  - `make live/esbuild`: JavaScript bundling and minification
- **Generate Templ**: `templ generate` before building
- **Pre-commit**: `./scripts/pre-commit.sh` runs tests, linting, and formatting
- **Testing**: Cypress E2E tests in `tests/` directory

### Deployment

- **Local**: `make live` (3-way hot reload: Air + Templ + ESBuild)
- **Build**: `make build` creates Linux binary in `bin/main`
- **Production**: Docker multi-stage build → systemd service
- **Environment**: Use `.env` files for configuration
- **Database**: MongoDB connection via `MONGODB_URI`
- **Service Management**: systemd with `goweb.service` file

## Common Operations

### Creating a New Entity Type

1. Define model in `internal/models/new_entity.go`
2. Add to Storage interface in `internal/storage/storage.go`
3. Implement in `internal/storage/new_entity_storage.go`
4. Create handlers in `internal/handlers/new_entity_handler.go`
5. Add routes in `internal/handlers/handler.go`
6. Create templates in `internal/templates/`

### Report Generation

1. Use existing templates in `templates/` directory
2. Extend `internal/sheets/excel.go` for new report types
3. Add report type constant in `models/record.go`
4. Wire up in report handlers

### Frontend Patterns

- **HTMX → Alpine.js Migration**: Gradually migrating from HTMX to Alpine.js - prefer Alpine.js for new features
- **Alpine.js**: For client-side state and interactions
- **Bootstrap**: For styling and layout
- **Templ**: For type-safe HTML generation

## Security Considerations

- JWT tokens in HTTP-only cookies
- Input validation in handlers
- SQL injection prevention via BSON queries
- CORS configuration for API endpoints
- Environment-based configuration secrets

## Performance Notes

- MongoDB aggregation for complex queries
- Pagination for large datasets
- Background processes for report generation
- File cleanup for generated reports
- Connection pooling for database access

This system follows domain-driven design principles with clear separation between web presentation, business logic, and data persistence layers.
