# Playwright Test Structure Diagram

## Visual Overview

```
Lab Management E2E Tests (Playwright)
│
├── Configuration
│   ├── playwright.config.js          # Main Playwright configuration
│   └── package.json                  # Dependencies and scripts
│
├── Documentation
│   ├── README_PLAYWRIGHT.md          # Complete setup & usage guide
│   ├── TEST_SUMMARY.md               # Test statistics & coverage
│   ├── TEST_STRUCTURE.md             # This file - visual overview
│   └── playwright/EXAMPLES.js        # 10+ test pattern examples
│
├── Helper Functions (Reusable Operations)
│   ├── helpers/auth.js               # Authentication (login, register, logout)
│   ├── helpers/patients.js           # Patient CRUD operations
│   ├── helpers/doctors.js            # Doctor CRUD operations
│   ├── helpers/tests.js              # Test definition CRUD
│   ├── helpers/combos.js             # Combo (package) CRUD
│   └── helpers/records.js            # Lab record operations
│
└── Test Specifications (44 tests total)
    ├── specs/01-auth.spec.js         # Authentication (7 tests)
    ├── specs/02-patients.spec.js     # Patient Management (6 tests)
    ├── specs/03-tests.spec.js        # Test Management (7 tests)
    ├── specs/04-combos.spec.js       # Combo Management (6 tests)
    ├── specs/05-records.spec.js      # Record Management (9 tests)
    ├── specs/06-integration.spec.js  # Integration Tests (3 tests)
    └── specs/07-doctors.spec.js      # Doctor Management (6 tests)
```

## Test Flow Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                   Test Execution Flow                        │
└─────────────────────────────────────────────────────────────┘
                             │
                             ▼
                    ┌────────────────┐
                    │  Test Runner   │
                    │  (Playwright)  │
                    └────────────────┘
                             │
            ┌────────────────┼────────────────┐
            ▼                ▼                ▼
    ┌──────────────┐ ┌──────────────┐ ┌──────────────┐
    │ Unit Tests   │ │ Integration  │ │   Helper     │
    │ (Individual  │ │    Tests     │ │  Functions   │
    │  Features)   │ │  (Complete   │ │  (Reusable   │
    │              │ │   Workflow)  │ │   Actions)   │
    └──────────────┘ └──────────────┘ └──────────────┘
            │                │                │
            └────────────────┼────────────────┘
                             ▼
                    ┌────────────────┐
                    │  Application   │
                    │ (Lab Management│
                    │    System)     │
                    └────────────────┘
```

## Test Dependencies Map

```
Authentication Flow (01-auth.spec.js)
│
├─► helpers/auth.js
│   ├─► login()
│   ├─► register()
│   └─► logout()
│
└─► All other tests depend on login()

Patient Management (02-patients.spec.js)
│
├─► helpers/auth.js → login()
└─► helpers/patients.js
    ├─► goToPatients()
    ├─► createPatient()
    ├─► searchPatient()
    ├─► editPatient()
    └─► deletePatient()

Doctor Management (07-doctors.spec.js)
│
├─► helpers/auth.js → login()
└─► helpers/doctors.js
    ├─► goToDoctors()
    ├─► createDoctor()
    ├─► searchDoctor()
    ├─► editDoctor()
    └─► deleteDoctor()

Test Management (03-tests.spec.js)
│
├─► helpers/auth.js → login()
└─► helpers/tests.js
    ├─► goToTests()
    ├─► createTest()
    ├─► searchTest()
    └─► deleteTest()

Combo Management (04-combos.spec.js)
│
├─► helpers/auth.js → login()
├─► helpers/tests.js → createTest() (setup)
└─► helpers/combos.js
    ├─► goToCombos()
    ├─► createCombo()
    ├─► searchCombo()
    └─► deleteCombo()

Record Management (05-records.spec.js)
│
├─► helpers/auth.js → login()
├─► helpers/patients.js → createPatient() (setup)
├─► helpers/tests.js → createTest() (setup)
├─► helpers/combos.js → createCombo() (setup)
└─► helpers/records.js
    ├─► goToRecords()
    ├─► createRecord()
    ├─► viewRecordDetails()
    └─► deleteRecord()

Integration Tests (06-integration.spec.js)
│
└─► All helpers combined
    ├─► helpers/auth.js → login()
    ├─► helpers/patients.js → createPatient()
    ├─► helpers/tests.js → createTest()
    ├─► helpers/combos.js → createCombo()
    └─► helpers/records.js → createRecord()
```

## Test Data Flow

```
┌────────────────────────────────────────────────────────────┐
│                  Test Data Generation                       │
└────────────────────────────────────────────────────────────┘
                             │
        ┌────────────────────┼────────────────────┐
        ▼                    ▼                    ▼
┌──────────────┐     ┌──────────────┐     ┌──────────────┐
│  Timestamp   │     │   Random     │     │   Unique     │
│  Based IDs   │     │   Phone      │     │   Names      │
│              │     │   Numbers    │     │              │
│ Date.now()   │     │ 0xxxxxxxxx   │     │ "Test" +     │
│              │     │              │     │  timestamp   │
└──────────────┘     └──────────────┘     └──────────────┘
        │                    │                    │
        └────────────────────┼────────────────────┘
                             ▼
                    ┌────────────────┐
                    │  Test Entity   │
                    │  (Patient,     │
                    │   Doctor,      │
                    │   Test, etc)   │
                    └────────────────┘
                             │
                             ▼
                    ┌────────────────┐
                    │  Application   │
                    │   Database     │
                    └────────────────┘
```

## Execution Sequence

### Sequential Test Execution (1 Worker)

```
Step 1: Initialize Browser
    ↓
Step 2: Run Test File 01-auth.spec.js
    ├─► Test 1: Display login page
    ├─► Test 2: Login with valid credentials
    ├─► Test 3: Show error with invalid credentials
    ├─► Test 4: Display register page
    ├─► Test 5: Navigate between pages
    ├─► Test 6: Logout successfully
    └─► Test 7: Redirect to login
    ↓
Step 3: Run Test File 02-patients.spec.js
    ├─► Test 1: Display patient page
    ├─► Test 2: Create patient
    ├─► Test 3: Search patients
    ├─► Test 4: Edit patient
    ├─► Test 5: Delete patient
    └─► Test 6: Validate fields
    ↓
Step 4: Run Test File 03-tests.spec.js
    └─► ... (7 tests)
    ↓
Step 5: Run Test File 04-combos.spec.js
    └─► ... (6 tests)
    ↓
Step 6: Run Test File 05-records.spec.js
    └─► ... (9 tests)
    ↓
Step 7: Run Test File 06-integration.spec.js
    └─► ... (3 tests)
    ↓
Step 8: Run Test File 07-doctors.spec.js
    └─► ... (6 tests)
    ↓
Step 9: Generate Reports
    ├─► HTML Report
    ├─► Screenshots (on failure)
    └─► Test Results Summary
```

## Component Interaction

```
┌──────────────────────────────────────────────────────────┐
│                     Browser (Chromium)                    │
│  ┌────────────────────────────────────────────────────┐  │
│  │           Application Under Test                   │  │
│  │  ┌──────────────────────────────────────────────┐  │  │
│  │  │  Frontend (HTMX + Alpine.js + Bootstrap)     │  │  │
│  │  └──────────────────┬───────────────────────────┘  │  │
│  │                     │ HTTP/API Calls               │  │
│  │  ┌──────────────────▼───────────────────────────┐  │  │
│  │  │  Backend (Go + Chi Router)                   │  │  │
│  │  └──────────────────┬───────────────────────────┘  │  │
│  │                     │ Database Queries             │  │
│  │  ┌──────────────────▼───────────────────────────┐  │  │
│  │  │  Database (MongoDB)                          │  │  │
│  │  └──────────────────────────────────────────────┘  │  │
│  └────────────────────────────────────────────────────┘  │
└──────────────────┬───────────────────────────────────────┘
                   │
                   │ Playwright API
                   │ (page.goto, page.fill, page.click, etc)
                   │
┌──────────────────▼───────────────────────────────────────┐
│              Playwright Test Runner                       │
│  ┌────────────────────────────────────────────────────┐  │
│  │  Test Specs (01-07)                                │  │
│  └────────────────┬───────────────────────────────────┘  │
│                   │ Uses                                  │
│  ┌────────────────▼───────────────────────────────────┐  │
│  │  Helper Functions (auth, patients, etc)            │  │
│  └────────────────────────────────────────────────────┘  │
└──────────────────────────────────────────────────────────┘
```

## Test Coverage Matrix

| Feature          | Unit Tests | Integration | Search | CRUD | Validation |
|------------------|------------|-------------|--------|------|------------|
| Authentication   | ✅ 7       | ✅          | N/A    | N/A  | ✅         |
| Patients         | ✅ 6       | ✅          | ✅     | ✅   | ✅         |
| Doctors          | ✅ 6       | ✅          | ✅     | ✅   | ✅         |
| Tests            | ✅ 7       | ✅          | ✅     | ✅   | ✅         |
| Combos           | ✅ 6       | ✅          | ✅     | ✅   | ✅         |
| Records          | ✅ 9       | ✅          | ✅     | ✅   | ✅         |
| Full Workflow    | N/A        | ✅ 3        | N/A    | N/A  | N/A        |
| **Total**        | **44**     | **✅**      | **✅** | **✅**| **✅**    |

## File Size & Complexity

```
Test Files (by size):
├── 06-integration.spec.js     7,715 bytes  (Complex workflows)
├── 05-records.spec.js         7,430 bytes  (Most test cases)
├── 03-tests.spec.js           4,489 bytes  (Standard CRUD)
├── 04-combos.spec.js          4,385 bytes  (With dependencies)
├── 02-patients.spec.js        4,299 bytes  (Standard CRUD)
├── 07-doctors.spec.js         3,840 bytes  (Standard CRUD)
└── 01-auth.spec.js            2,701 bytes  (Authentication only)

Helper Files (by size):
├── records.js                 3,432 bytes  (Complex operations)
├── patients.js                3,172 bytes  (Full CRUD)
├── doctors.js                 2,985 bytes  (Full CRUD)
├── tests.js                   2,267 bytes  (Standard CRUD)
├── combos.js                  2,169 bytes  (With test selection)
└── auth.js                    1,421 bytes  (Basic auth)

Documentation (by size):
├── EXAMPLES.js                9,373 bytes  (10 examples + tips)
├── TEST_SUMMARY.md            7,500 bytes  (Complete overview)
├── README_PLAYWRIGHT.md       7,387 bytes  (Setup guide)
└── TEST_STRUCTURE.md          ~8,000 bytes (This file)
```

## Quick Reference

### Run Commands
```bash
npm test                    # All tests
npm run test:ui            # Interactive UI
npm run test:headed        # See browser
npm run test:debug         # Debug mode
npx playwright show-report # View HTML report
```

### File Locations
```
tests/
├── playwright.config.js           # Configuration
├── package.json                   # Dependencies
├── README_PLAYWRIGHT.md           # Documentation
├── TEST_SUMMARY.md                # Statistics
├── TEST_STRUCTURE.md              # This file
└── playwright/
    ├── EXAMPLES.js                # Examples
    ├── helpers/                   # Reusable functions
    └── specs/                     # Test files
```

---

**Legend:**
- ✅ Fully Implemented
- 📊 Statistics Available
- 🎯 Primary Focus
- 📁 File/Directory
- ▼ Flow Direction
- │ Connection
- └─► Dependency/Usage

**Last Updated**: November 2025
