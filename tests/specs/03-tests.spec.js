// Test Management E2E tests
const { test, expect } = require('@playwright/test');
const { login } = require('../helpers/auth');
const { goToTests, createTest, searchTest, deleteTest } = require('../helpers/tests');

test.describe('Test Management Flow', () => {
  test.beforeEach(async ({ page }) => {
    // Login before each test
    await login(page);
  });

  test('should display test management page', async ({ page }) => {
    await goToTests(page);

    await expect(page.getByRole('heading', { name: 'Danh mục xét nghiệm' })).toBeVisible();
    await expect(page.getByRole('button', { name: 'Thêm xét nghiệm' })).toBeVisible();
    await expect(page.getByPlaceholder('Tên xét nghiệm')).toBeVisible();
  });

  test('should create a new test', async ({ page }) => {
    await goToTests(page);
    
    const testData = {
      name: `Glucose Test ${Date.now()}`,
      unit: 'mmol/L',
      price: 50000,
      lowerBound: 3.9,
      upperBound: 6.1,
      normalValue: '3.9-6.1',
    };
    
    await createTest(page, testData);

    await searchTest(page, testData.name);
    
    // Verify test appears in the list
    await expect(page.locator(`text=${testData.name}`)).toBeVisible();
    await expect(page.locator(`text=${testData.unit}`)).toBeVisible();
  });

  test('should create multiple tests with different parameters', async ({ page }) => {
    await goToTests(page);
    
    const tests = [
      {
        name: `Hemoglobin ${Date.now()}`,
        unit: 'g/dL',
        price: 30000,
        lowerBound: 12.0,
        upperBound: 16.0,
        normalValue: '12.0-16.0',
      },
      {
        name: `WBC Count ${Date.now()}`,
        unit: '10^9/L',
        price: 40000,
        lowerBound: 4.0,
        upperBound: 11.0,
        normalValue: '4.0-11.0',
      },
    ];
    
    for (const testData of tests) {
      await createTest(page, testData);
      await searchTest(page, testData.name);
      await expect(page.locator(`text=${testData.name}`)).toBeVisible();
    }
  });

  test('should search for tests', async ({ page }) => {
    await goToTests(page);
    
    // Create a test first
    const testData = {
      name: `Cholesterol Search ${Date.now()}`,
      unit: 'mmol/L',
      price: 60000,
      lowerBound: 3.0,
      upperBound: 5.2,
      normalValue: '3.0-5.2',
    };
    
    await createTest(page, testData);
    
    // Search for the test
    await searchTest(page, testData.name);
    
    // Verify search results
    await expect(page.locator(`text=${testData.name}`)).toBeVisible();
    
    // Search for non-existent test
    await searchTest(page, 'NonExistentTest99999');
    await expect(page.locator('text=NonExistentTest99999')).not.toBeVisible();
  });

  test('should delete a test', async ({ page }) => {
    await goToTests(page);
    
    // Create a test to delete
    const testData = {
      name: `Delete Test ${Date.now()}`,
      unit: 'mg/dL',
      price: 45000,
      lowerBound: 70,
      upperBound: 100,
      normalValue: '70-100',
    };
    
    await createTest(page, testData);

    // Verify test exists
    await searchTest(page, testData.name);
    await expect(page.locator(`text=${testData.name}`)).toBeVisible();
    
    // Delete the test
    await deleteTest(page, testData.name);
    
    // Verify test is removed
    await searchTest(page, testData.name);
    await expect(page.locator(`text=${testData.name}`)).not.toBeVisible();
  });

  test('should display test details correctly', async ({ page }) => {
    await goToTests(page);
    
    const testData = {
      name: `Detail Test ${Date.now()}`,
      unit: 'IU/L',
      price: 75000,
      lowerBound: 10,
      upperBound: 40,
      normalValue: '10-40',
    };
    
    await createTest(page, testData);
    
    // Search for the test
    await searchTest(page, testData.name);
    
    // Verify all details are displayed
    const row = page.locator('tr', { hasText: testData.name }).first();
    await expect(row).toContainText(testData.unit);
    await expect(row).toContainText('75,000'); // Price formatted
    await expect(row).toContainText(testData.normalValue);
  });

  test('should validate required fields when creating test', async ({ page }) => {
    await goToTests(page);
    
    await page.click('text=Thêm xét nghiệm');
    
    // Try to submit without filling required fields
    await page.click('button[type="submit"]:has-text("Thêm xét nghiệm")');
    
    // Check for HTML5 validation
    const nameInput = page.getByLabel('Tên xét nghiệm');
    const isInvalid = await nameInput.evaluate(el => !el.checkValidity());
    expect(isInvalid).toBeTruthy();
  });
});
