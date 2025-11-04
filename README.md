## Introduction

Business: An Admin website to manage the Anh Quan laboratory.

Techstack: Go + HTMX + Alpine + MongoDB

## Features

- Implement the record status.
- Tracking page: Provide user the comparision between multiple records

## Bug Fixes

- When create new patient in Create Record page -> should stay in the create record page and select the newly created patient instead of redirecting to patient management page

## Improvements

- Tracking unsaved changes and asking user to confirm before they leave or reload the page

## Testing

This project includes comprehensive E2E testing with Playwright covering all major application flows.

### Running E2E Tests

```bash
# Navigate to tests directory
cd tests

# Install dependencies (first time only)
npm install
npx playwright install chromium

# Run all tests
npm test

# Run tests with UI (interactive mode)
npm run test:ui

# Run tests in headed mode (see browser)
npm run test:headed
```

For more details, see [tests/README_PLAYWRIGHT.md](tests/README_PLAYWRIGHT.md).

### Test Coverage

- ✅ Authentication (login/register/logout)
- ✅ Patient Management (CRUD operations)
- ✅ Doctor Management (CRUD operations)
- ✅ Test Management (CRUD operations)
- ✅ Combo Management (test packages)
- ✅ Record Management (lab test records)
- ✅ Complete integration workflows
- ✅ Abnormal test result validation
