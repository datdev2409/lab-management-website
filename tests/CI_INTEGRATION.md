# CI/CD Integration for E2E Tests

This document describes the CI/CD integration for Playwright E2E tests.

## Overview

The E2E tests are automatically run in GitHub Actions CI pipeline on every pull request and push to `main` or `develop` branches.

## CI Workflow

The workflow (`.github/workflows/e2e-tests.yml`) performs the following steps:

### 1. Environment Setup
- Checks out code
- Sets up Go 1.24 and Node.js 20
- Installs Go dependencies
- Installs and runs templ generator
- Builds the Go application

### 2. Database & Services
- Starts MongoDB and Gotenberg using Docker Compose
- Waits for MongoDB to be healthy
- Seeds test database with initial data (admin user)

### 3. Application Server
- Starts the application server on port 7331
- Waits for server health check to pass
- Runs in background during tests

### 4. Test Execution
- Installs Playwright and browser dependencies
- Runs all 44 Playwright tests
- Captures screenshots on failure
- Generates HTML test report

### 5. Cleanup
- Stops application server
- Stops and removes Docker containers

## Test Data

The CI pipeline seeds the database with:
- **Admin User**: 
  - Username: `admin`
  - Password: `admin123`

Tests create additional data dynamically using unique timestamps to avoid conflicts.

## Environment Variables

The following environment variables are configured in CI:

```bash
MONGODB_URI=mongodb://root:password123@localhost:27017/labadmin?authSource=admin
SERVER_PORT=7331
ENV=test
CI=true
```

## Running Tests Locally with CI Setup

To replicate the CI environment locally:

```bash
# 1. Start services
docker-compose -f docker-compose.test.yml up -d

# 2. Wait for MongoDB
timeout 60 bash -c 'until docker exec $(docker ps -qf "name=mongo_db") mongosh --eval "db.adminCommand(\"ping\")" > /dev/null 2>&1; do sleep 2; done'

# 3. Seed test data
MONGODB_URI="mongodb://root:password123@localhost:27017" tests/scripts/seed-test-data.sh

# 4. Build and start application
go build -o bin/server cmd/api/main.go
MONGODB_URI="mongodb://root:password123@localhost:27017/labadmin?authSource=admin" \
SERVER_PORT=7331 \
ENV=test \
./bin/server &

# 5. Run tests
cd tests
npm test

# 6. Cleanup
kill %1  # Stop server
docker-compose -f docker-compose.test.yml down -v
```

## CI Artifacts

When tests run in CI, the following artifacts are uploaded:

1. **playwright-report**: HTML test report (always uploaded)
2. **playwright-screenshots**: Screenshots of failed tests (uploaded on failure)

These artifacts are retained for 7 days and can be downloaded from the GitHub Actions workflow run page.

## Test Execution Time

Typical test execution times:
- Setup (build, services, seed): ~2-3 minutes
- Test execution: ~3-5 minutes
- **Total**: ~5-8 minutes

## Monitoring & Debugging

### Viewing Test Results

1. Go to GitHub Actions tab in the repository
2. Select the "E2E Tests" workflow
3. Click on a workflow run
4. View the job logs or download artifacts

### Common Issues

**Issue**: MongoDB connection timeout
- **Solution**: Check MongoDB health check status in logs

**Issue**: Server health check timeout
- **Solution**: Verify server logs, check port 7331 availability

**Issue**: Tests failing due to missing data
- **Solution**: Verify seed script executed successfully

**Issue**: Browser not found
- **Solution**: Ensure `playwright install --with-deps` ran successfully

## Configuration Files

- `.github/workflows/e2e-tests.yml` - CI workflow definition
- `docker-compose.test.yml` - Docker services for testing
- `tests/scripts/seed-test-data.sh` - Database seeding script
- `tests/playwright.config.js` - Playwright configuration

## Adding New Tests

New tests added to `tests/playwright/specs/` will automatically be included in CI runs. No changes to the CI configuration are needed.

## Performance Optimization

The CI pipeline is optimized for:
- **Fast startup**: Health checks with short intervals
- **Parallel operations**: Where possible (e.g., building while services start)
- **Efficient cleanup**: Removes volumes to ensure clean state
- **Artifact management**: 7-day retention, only uploads on failure for screenshots

## Security Considerations

- MongoDB credentials are hardcoded for test environment only
- Test database is isolated and destroyed after each run
- No production credentials are used
- All services run in isolated Docker network

## Future Enhancements

Potential improvements:
- [ ] Add test result comments to PRs
- [ ] Implement test result trends/history
- [ ] Add performance benchmarks
- [ ] Support multiple browsers (Firefox, WebKit)
- [ ] Parallel test execution (when tests become independent)
- [ ] Visual regression testing
