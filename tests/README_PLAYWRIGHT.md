# Playwright E2E Tests for Lab Management System

This directory contains end-to-end (E2E) tests using Playwright for the Lab Management System.

## Structure

```
tests/
├── playwright/
│   ├── helpers/           # Helper functions for test operations
│   │   ├── auth.js       # Authentication helpers (login, register, logout)
│   │   ├── patients.js   # Patient management helpers
│   │   ├── tests.js      # Test management helpers
│   │   ├── combos.js     # Combo management helpers
│   │   └── records.js    # Record management helpers
│   └── specs/            # Test specifications
│       ├── 01-auth.spec.js           # Authentication flow tests
│       ├── 02-patients.spec.js       # Patient management tests
│       ├── 03-tests.spec.js          # Test management tests
│       ├── 04-combos.spec.js         # Combo management tests
│       ├── 05-records.spec.js        # Record management tests
│       └── 06-integration.spec.js    # Complete integration tests
├── playwright.config.js   # Playwright configuration
└── package.json          # NPM dependencies
```

## Prerequisites

1. **Node.js** (v16 or higher)
2. **NPM** or **Yarn**
3. **Running Application** - The application must be running on `http://127.0.0.1:7331`

## Installation

### 1. Install Dependencies

```bash
cd tests
npm install
```

### 2. Install Playwright Browsers

```bash
npx playwright install chromium
```

Or install all browser dependencies:

```bash
npx playwright install-deps
npx playwright install
```

## Running Tests

### Run All Tests

```bash
npm test
```

### Run Tests in Headed Mode (See Browser)

```bash
npm run test:headed
```

### Run Tests with UI Mode (Interactive)

```bash
npm run test:ui
```

### Run Tests in Debug Mode

```bash
npm run test:debug
```

### Run Specific Test File

```bash
npx playwright test playwright/specs/01-auth.spec.js
```

### Run Tests Matching a Pattern

```bash
npx playwright test --grep "should create a new patient"
```

## Test Coverage

The test suite covers the following application flows:

### 1. Authentication Flow (`01-auth.spec.js`)
- Display login page
- Login with valid credentials
- Show error with invalid credentials
- Display register page
- Navigate between login and register
- Logout functionality
- Redirect to login for protected pages

### 2. Patient Management (`02-patients.spec.js`)
- Display patient management page
- Create new patient
- Search for patients
- Edit patient information
- Delete patient
- Validate required fields

### 3. Test Management (`03-tests.spec.js`)
- Display test management page
- Create new test
- Create multiple tests with different parameters
- Search for tests
- Delete tests
- Display test details
- Validate required fields

### 4. Combo Management (`04-combos.spec.js`)
- Display combo management page
- Create combo with multiple tests
- Search for combos
- View combo details
- Delete combos
- Validate required fields

### 5. Record Management (`05-records.spec.js`)
- Display records page
- Navigate to create record page
- Create new record with combo
- View record details
- Filter records by date
- Search records by patient name
- Display pagination controls
- Create patient from record creation page
- Validate abnormal test results

### 6. Integration Tests (`06-integration.spec.js`)
- Complete workflow: patient → tests → combo → record
- Handle workflow with abnormal test results
- Navigate through all main pages

## Configuration

The Playwright configuration is in `playwright.config.js`:

- **Base URL**: `http://127.0.0.1:7331`
- **Timeout**: Default Playwright timeout
- **Retries**: 2 retries on CI, 0 locally
- **Workers**: 1 worker (sequential execution)
- **Browser**: Chromium
- **Screenshots**: Captured on failure
- **Traces**: Captured on first retry

### Modifying Configuration

To change the base URL:

```javascript
// In playwright.config.js
use: {
  baseURL: 'http://localhost:8080',  // Change to your URL
}
```

## Helper Functions

Helper functions are organized by feature area and provide reusable operations:

### Authentication (`helpers/auth.js`)
- `login(page, username, password)` - Login to the application
- `register(page, username, password)` - Register a new user
- `logout(page)` - Logout from the application

### Patients (`helpers/patients.js`)
- `goToPatients(page)` - Navigate to patients page
- `createPatient(page, patientData)` - Create a new patient
- `searchPatient(page, searchTerm)` - Search for a patient
- `editPatient(page, patientName, newData)` - Edit patient information
- `deletePatient(page, patientName)` - Delete a patient

### Tests (`helpers/tests.js`)
- `goToTests(page)` - Navigate to tests page
- `createTest(page, testData)` - Create a new test
- `searchTest(page, searchTerm)` - Search for a test
- `deleteTest(page, testName)` - Delete a test

### Combos (`helpers/combos.js`)
- `goToCombos(page)` - Navigate to combos page
- `createCombo(page, comboData)` - Create a new combo
- `searchCombo(page, searchTerm)` - Search for a combo
- `deleteCombo(page, comboName)` - Delete a combo

### Records (`helpers/records.js`)
- `goToRecords(page)` - Navigate to records page
- `createRecord(page, recordData)` - Create a new record
- `viewRecordDetails(page, recordIdentifier)` - View record details
- `generateReport(page, reportType)` - Generate a report
- `deleteRecord(page, recordIdentifier)` - Delete a record

## Best Practices

1. **Test Independence**: Each test should be independent and not rely on other tests
2. **Data Cleanup**: Tests create unique data using timestamps to avoid conflicts
3. **Wait Strategies**: Use `waitForLoadState('networkidle')` for page transitions
4. **Selectors**: Use text-based selectors for Vietnamese UI elements
5. **Error Handling**: Tests handle dialogs and confirmations appropriately

## Debugging Tests

### View Test Report

After running tests, view the HTML report:

```bash
npx playwright show-report
```

### Debug a Specific Test

```bash
npx playwright test --debug playwright/specs/01-auth.spec.js
```

### Record a New Test

```bash
npx playwright codegen http://127.0.0.1:7331
```

## Continuous Integration

The tests are configured to run in CI environments:

- **Retries**: 2 retries on failure in CI
- **Workers**: 1 worker to avoid race conditions
- **Headless**: Tests run in headless mode by default
- **Screenshots**: Captured on failure for debugging

### Environment Variables

- `CI`: Set to `true` in CI environments to enable CI-specific settings

## Troubleshooting

### Application Not Running

Make sure the application is running on `http://127.0.0.1:7331`:

```bash
# In the project root
make live
```

### Browser Installation Issues

If browsers fail to install, try:

```bash
npx playwright install-deps
npx playwright install chromium --force
```

### Test Timeouts

Increase timeout in `playwright.config.js`:

```javascript
use: {
  timeout: 60000,  // 60 seconds
}
```

### Network Issues

If tests fail due to network issues, check:
1. Application is running and accessible
2. No firewall blocking local connections
3. Correct base URL in configuration

## Contributing

When adding new tests:

1. Create helper functions for reusable operations
2. Use descriptive test names
3. Add appropriate assertions
4. Follow existing code structure
5. Document complex test scenarios

## License

Same as the main project.
