# E2E Test Suite Documentation

## Overview

The E2E test suite has been reorganized into **3 main flows** to eliminate duplication and improve test coverage. Each flow tests a specific aspect of the lab management system while sharing reusable test data.

## Test Structure

### Main Test Flows

#### 1. Flow 1: Record CRUD (`02-record-flow.spec.js`)
**Main business flow - Patient comes to the lab**

This flow covers the complete end-to-end process of managing a lab test record:
- **Step 1**: Admin creates a patient in Patient page
- **Step 2**: Select existing combo (seeded data) → Validate combo name and tests are populated
- **Step 3**: Add and remove tests to the record
- **Step 4**: Test back button unsaved warning functionality
- **Step 5**: Input test results
- **Step 6a**: Test automatic abnormal detection (values outside bounds)
- **Step 6b**: Test manual abnormal override
- **Step 7**: Save record and search for it
- **Step 8**: View record details and validate all information

**Key Features Tested:**
- Patient selection via autocomplete
- Combo selection and test population
- Dynamic test addition/removal
- Unsaved changes warning
- Abnormal value detection (automatic & manual)
- Record persistence and retrieval

#### 2. Flow 2: Admin CRUD (`03-admin-crud-flow.spec.js`)
**Entity management - Patients, Tests, Combos**

Consolidated CRUD operations for all administrative entities (excluding records):

**Patient Management:**
- Display patient management page
- Create new patients
- Search for patients
- Edit patient information
- Delete patients
- Validate required fields

**Test Management:**
- Display test management page
- Create single and multiple tests
- Search for tests
- Display test details
- Delete tests
- Validate required fields

**Combo Management:**
- Display combo management page
- Create combos with multiple tests
- Search for combos
- View combo details
- Delete combos
- Display test count
- Navigate back from creation form

**Navigation:**
- Verify all admin pages are accessible

#### 3. Flow 3: Report Export (`04-report-export-flow.spec.js`)
**Report generation and validation using Playwright + SheetJS**

Tests all report types with actual file downloads and content validation:

**Report Types Tested:**
- `phieu_thu` - Billing report (Excel)
- `phieu_ket_qua` - Results report (Excel)
- `phieu_ket_qua_chu_ky` - Results with signature (Excel)
- `phieu_ket_qua_chu_ky_pdf` - PDF report
- `phieu_theo_doi` - Tracking/comparison report (Excel)

**Validation Process:**
1. Navigate to record details page
2. Trigger report download
3. Save file to `/tmp`
4. Parse Excel with SheetJS (or verify PDF magic bytes)
5. Validate patient information and test data in report
6. Clean up downloaded files

### Supporting Test Files

#### Auth Flow (`01-auth.spec.js`)
Tests authentication functionality:
- Display login page
- Login with valid credentials
- Show error with invalid credentials
- Display register page
- Navigate between login/register
- Logout functionality
- Redirect to login for protected pages

#### Doctors Management (`07-doctors.spec.js`)
Tests doctor-specific functionality (preserved from original)

## Test Helpers

### Helper Modules (`tests/helpers/`)

- **`auth.js`** - Login, register, logout operations
- **`patients.js`** - Patient CRUD operations
- **`tests.js`** - Test definition CRUD operations
- **`combos.js`** - Combo CRUD operations
- **`records.js`** - Record CRUD operations
- **`doctors.js`** - Doctor management operations
- **`seed.js`** - **NEW** - Reusable test data seeding

### Seed Helper

The `seed.js` helper provides `seedBasicTestData(page)` which creates:
- 4 standard test definitions (Glucose, Hemoglobin, WBC Count, Cholesterol)
- 3 standard combos (Basic Health Check, Complete Blood Count, Full Panel)

This eliminates duplication by allowing multiple test flows to share the same seeded data.

## Key Improvements

### ✅ Eliminated Duplication
- Record tests no longer create their own patients, tests, and combos
- Data is seeded once and reused across test flows
- Helper functions shared across all tests

### ✅ Improved Coverage
- Comprehensive record CRUD flow (8 steps)
- Both automatic and manual abnormal detection
- All report types validated with actual file parsing
- Unsaved changes warning tested

### ✅ Better Organization
- Tests grouped by business flow (not entity type)
- Clear separation of concerns
- Easy to understand test progression

### ✅ Enhanced Validation
- SheetJS library for Excel content validation
- PDF magic byte verification
- Patient and test data validation in reports

## Running Tests

```bash
# Run all E2E tests
npm run test:ci

# Run tests with UI
npm run test:ui

# Run specific flow
npx playwright test tests/specs/02-record-flow.spec.js

# Run specific test
npx playwright test tests/specs/02-record-flow.spec.js -g "Step 1"
```

## Configuration

- **Base URL**: `http://localhost:9000` (configured in `playwright.config.js`)
- **Browser**: Chromium (Desktop Chrome)
- **Parallel**: Enabled (2 workers on CI)
- **Retries**: 0
- **Timeout**: Default 30s (120s for integration tests)

## Wait Strategy

All tests use explicit `waitForTimeout()` instead of `waitForLoadState('networkidle')` for more predictable timing:
- Page navigation: `1000ms`
- Autocomplete results: `600ms`
- Form interactions: `300ms`
- Element updates: `200-500ms`

## Dependencies

- `@playwright/test` - E2E testing framework
- `xlsx` - SheetJS library for Excel parsing and validation

## Migration Notes

### Old Files (Archived)
The following files have been archived with `.old` extension:
- `02-patients.spec.js.old`
- `03-tests.spec.js.old`
- `04-combos.spec.js.old`
- `05-records.spec.js.old`
- `06-integration.spec.js.old`

These are ignored by git (`.gitignore`) and can be safely deleted once the new tests are verified.

### Breaking Changes
None - all test functionality has been preserved and consolidated.

## Troubleshooting

### Test Failures

1. **Authentication failures**: Ensure the test user is seeded in MongoDB
   ```bash
   ./scripts/seed-test-data.sh
   ```

2. **Timeout errors**: Increase wait times if server is slow
   ```javascript
   await page.waitForTimeout(2000); // Increase from 1000
   ```

3. **Element not found**: Check for recent UI changes in selectors

4. **Report download failures**: Ensure report generation endpoints are working

## Future Enhancements

- [ ] Add test data cleanup after test runs
- [ ] Implement test database snapshots for faster resets
- [ ] Add visual regression testing for reports
- [ ] Parallel test execution optimization
- [ ] Add performance metrics tracking
