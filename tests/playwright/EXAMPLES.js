// Example test file showing how to create new Playwright tests
// This file demonstrates various test patterns and best practices

const { test, expect } = require('@playwright/test');
const { login } = require('../helpers/auth');

// ============================================================================
// EXAMPLE 1: Simple Page Navigation Test
// ============================================================================
test.describe('Example: Simple Navigation', () => {
  test('should navigate to a page and verify content', async ({ page }) => {
    // Login first
    await login(page);
    
    // Navigate to a page
    await page.goto('/danh-muc-benh-nhan');
    
    // Wait for page to load
    await page.waitForLoadState('networkidle');
    
    // Verify page content
    await expect(page.locator('h3')).toContainText('Danh mục bệnh nhân');
  });
});

// ============================================================================
// EXAMPLE 2: Form Submission Test
// ============================================================================
test.describe('Example: Form Submission', () => {
  test('should fill and submit a form', async ({ page }) => {
    await login(page);
    await page.goto('/danh-muc-benh-nhan');
    
    // Click button to open form
    await page.click('text=Thêm bệnh nhân');
    
    // Fill form fields
    await page.fill('input[name="patient_name"]', 'Test Name');
    await page.fill('input[name="patient_yob"]', '1990');
    await page.check('input[name="patient_gender"][value="Nam"]');
    await page.fill('input[name="patient_address"]', 'Test Address');
    await page.fill('input[name="patient_phone"]', '0123456789');
    
    // Submit form
    await page.click('button[type="submit"]');
    
    // Wait for operation to complete
    await page.waitForLoadState('networkidle');
    
    // Verify success
    await expect(page.locator('text=Test Name')).toBeVisible();
  });
});

// ============================================================================
// EXAMPLE 3: Search and Filter Test
// ============================================================================
test.describe('Example: Search Functionality', () => {
  test('should search and find results', async ({ page }) => {
    await login(page);
    await page.goto('/danh-muc-xet-nghiem');
    
    // Type in search box
    await page.fill('#test-search', 'Glucose');
    
    // Submit search
    await page.click('#test-search-form button[type="submit"]');
    
    // Wait for results
    await page.waitForLoadState('networkidle');
    
    // Verify results (if they exist)
    const results = page.locator('table tbody tr');
    const count = await results.count();
    
    if (count > 0) {
      await expect(results.first()).toContainText('Glucose');
    }
  });
});

// ============================================================================
// EXAMPLE 4: Dialog/Confirmation Handling
// ============================================================================
test.describe('Example: Dialog Handling', () => {
  test('should handle confirmation dialogs', async ({ page }) => {
    await login(page);
    await page.goto('/danh-muc-benh-nhan');
    
    // Setup dialog handler BEFORE triggering the action
    page.once('dialog', dialog => {
      expect(dialog.type()).toBe('confirm');
      dialog.accept(); // or dialog.dismiss()
    });
    
    // Trigger action that shows dialog
    const deleteButton = page.locator('text=Xoá').first();
    if (await deleteButton.count() > 0) {
      await deleteButton.click();
    }
  });
});

// ============================================================================
// EXAMPLE 5: Waiting for Elements
// ============================================================================
test.describe('Example: Waiting Strategies', () => {
  test('should wait for elements properly', async ({ page }) => {
    await login(page);
    await page.goto('/phieu-xet-nghiem');
    
    // Wait for network to be idle
    await page.waitForLoadState('networkidle');
    
    // Wait for specific element to be visible
    await page.waitForSelector('table', { state: 'visible' });
    
    // Wait for element with timeout
    await page.waitForSelector('h3', { timeout: 5000 });
    
    // Wait for custom condition
    await page.waitForFunction(() => {
      return document.querySelectorAll('table tbody tr').length > 0;
    });
  });
});

// ============================================================================
// EXAMPLE 6: Autocomplete/Suggestions
// ============================================================================
test.describe('Example: Autocomplete', () => {
  test('should interact with autocomplete fields', async ({ page }) => {
    await login(page);
    await page.goto('/phieu-xet-nghiem/new');
    await page.waitForLoadState('networkidle');
    
    // Type to trigger autocomplete
    await page.fill('input[placeholder*="Tìm kiếm bệnh nhân"]', 'Test');
    
    // Wait for suggestions to appear
    await page.waitForTimeout(500);
    
    // Click on a suggestion
    const suggestion = page.locator('text=Test').first();
    if (await suggestion.count() > 0) {
      await suggestion.click();
    }
  });
});

// ============================================================================
// EXAMPLE 7: Data-Driven Tests
// ============================================================================
test.describe('Example: Data-Driven Testing', () => {
  const testData = [
    { name: 'Test 1', value: 'Value 1' },
    { name: 'Test 2', value: 'Value 2' },
    { name: 'Test 3', value: 'Value 3' },
  ];

  testData.forEach(data => {
    test(`should test with ${data.name}`, async ({ page }) => {
      await login(page);
      // Your test logic using data.name and data.value
    });
  });
});

// ============================================================================
// EXAMPLE 8: Setup and Teardown
// ============================================================================
test.describe('Example: Setup and Teardown', () => {
  let testData;

  // Runs once before all tests in this describe block
  test.beforeAll(async ({ browser }) => {
    // Setup code
    testData = { id: Date.now() };
  });

  // Runs before each test
  test.beforeEach(async ({ page }) => {
    await login(page);
  });

  test('should use setup data', async ({ page }) => {
    // Test using testData
  });

  // Runs after each test
  test.afterEach(async ({ page }) => {
    // Cleanup if needed
  });

  // Runs once after all tests
  test.afterAll(async () => {
    // Final cleanup
  });
});

// ============================================================================
// EXAMPLE 9: Taking Screenshots
// ============================================================================
test.describe('Example: Screenshots', () => {
  test('should take screenshots', async ({ page }) => {
    await login(page);
    await page.goto('/phieu-xet-nghiem');
    
    // Take full page screenshot
    await page.screenshot({ path: 'screenshot-full.png', fullPage: true });
    
    // Take element screenshot
    const element = page.locator('table').first();
    await element.screenshot({ path: 'screenshot-table.png' });
  });
});

// ============================================================================
// EXAMPLE 10: API Requests
// ============================================================================
test.describe('Example: API Testing', () => {
  test('should make API requests', async ({ request, page }) => {
    await login(page);
    
    // Make API request
    const response = await request.get('/api/v1/patients');
    expect(response.ok()).toBeTruthy();
    
    // Parse response
    const data = await response.json();
    expect(data).toHaveProperty('data');
  });
});

// ============================================================================
// TIPS AND BEST PRACTICES
// ============================================================================

/*
1. UNIQUE TEST DATA
   - Use timestamps or random values to avoid data conflicts
   - Example: `Test Patient ${Date.now()}`

2. PROPER WAITS
   - Use waitForLoadState('networkidle') after navigation
   - Use waitForTimeout() sparingly, prefer waitForSelector()

3. SELECTORS
   - Prefer text-based selectors for Vietnamese UI: text=Thêm bệnh nhân
   - Use data-testid attributes in production code for stable selectors
   - Fallback to name, id, or CSS selectors

4. ERROR HANDLING
   - Handle dialogs before triggering actions that show them
   - Use count() to check if elements exist before interacting
   - Use conditional logic: if (await element.count() > 0) { ... }

5. TEST INDEPENDENCE
   - Each test should be able to run independently
   - Don't rely on execution order
   - Clean up test data when possible

6. ASSERTIONS
   - Use expect() with appropriate matchers
   - Common matchers: toBeVisible(), toContainText(), toHaveURL()
   - Check both positive and negative cases

7. DEBUGGING
   - Use test.only() to run a single test
   - Use test.skip() to skip tests temporarily
   - Run with --debug flag to step through tests
   - Use page.pause() to pause execution

8. PERFORMANCE
   - Reuse browser contexts when possible
   - Use beforeAll for expensive setup operations
   - Keep tests focused and fast
*/
